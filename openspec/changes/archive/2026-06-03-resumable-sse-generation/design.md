## Context

当前 Chat usecase 中 AI 生成与 SSE 连接耦合：`WriteToken` 失败 → `return err` → 消息不写入 DB、goroutine 泄漏。需要将生成生命周期与连接生命周期解耦。

## Goals / Non-Goals

**Goals:**
- SSE 客户端断开后，AI 生成继续执行，完成后写入 DB
- `GET /conversations/:id` 返回 `generating: bool` 标识生成状态
- `POST /conversation` 传入已有 `conversation_id` 时，若活跃生成存在则走 resume（回放 + 直播），否则走正常新消息
- 聊天列表能感知生成中的对话

**Non-Goals:**
- 不新建 API 路由
- 不持久化部分生成内容到 DB（只存内存，服务器重启丢失——可接受）
- 不处理多用户同时 resume 的并发问题（先单用户）

## Decisions

### GenerationManager

```go
// internal/generation/manager.go

type Generation struct {
    ConversationID uuid.UUID
    Buffer         []string          // 已累积的全部 token
    subscribers    map[chan string]struct{}
    Content        string            // 完整内容
    Done           bool
    mu             sync.RWMutex
}

type Manager struct {
    mu    sync.RWMutex
    items map[uuid.UUID]*Generation
}
```

全局单例，在 bootstrap 中初始化并注入 usecase。

**方法：**
- `Start(convID)` — 创建 Generation
- `Append(convID, token)` — 追加 token，广播给所有 subscriber
- `Subscribe(convID)` → `(chan string, []string)` — 返回 channel + 已累积 buffer（用于 replay）
- `Unsubscribe(convID, ch)` — 移除 subscriber
- `Finish(convID, content)` → 标记 done，插入 DB，清理
- `Get(convID)` → `*Generation, bool`

### UseCase Chat 双路径

```
Chat(userID, request, writer):
    │
    ├─ request.ConversationID 存在 && Manager.Has(convID)?
    │    │
    │    ├─ Yes (Resume 路径):
    │    │    ch, buffer := Manager.Subscribe(convID)
    │    │    // 1. replay: 遍历 buffer，逐条 WriteToken
    │    │    for _, token := range buffer { writer.WriteToken(token) }
    │    │    // 2. live: 从 ch 读新 token，WriteToken
    │    │    for token := range ch { writer.WriteToken(token) }
    │    │    // 3. Manager 的 goroutine 在 done 时 close(ch)
    │    │
    │    └─ No (新消息路径):
    │         // 与当前逻辑相同，但改为：
    │         Manager.Start(convID)
    │         go 调 AI → 写入 Manager
    │         ch, _ := Manager.Subscribe(convID)
    │         从 ch 读 token → WriteToken（失败忽略，继续消费 ch）
    │         Manager.Finish → 插入 DB
```

关键变化：新消息路径中，SSE write 失败不再 return，而是继续消费 channel 直到 done，由 Manager 的 goroutine 完成 DB 写入。

### GetConversation 返回 generating 状态

```go
func (u *ConversationUseCase) GetConversation(userID, conversationID uuid.UUID) (*domain.ConversationDetailResponse, error) {
    // ... 现有逻辑不变 ...
    
    gen, active := u.manager.Get(conversationID)
    resp.Generating = active
    
    // 如果有活跃生成且部分内容，拼接到 messages 末尾
    if active && gen.Content != "" {
        resp.Messages = append(resp.Messages, domain.Message{
            Type: 1,  // assistant
            Content: gen.Content,
            Generating: true,
        })
    }
    
    return &resp, nil
}
```

### 前端流恢复

```typescript
// pages/index.vue
async function handleSelect(id: string) {
    setActive(id)
    const detail = await getConversation(id)
    messages.value = detail.messages ?? []
    
    if (detail.generating) {
        // 恢复流式——复用 streamChat，只传 conversation_id，不传 message
        isLoading.value = true
        streamBuffer.value = detail.messages.filter(m => m.type === 1 && m.generating)[0]?.content ?? ''
        // 继续流式接收
        await streamChat(
            { conversation_id: id },
            { onToken, onDone, onError },
            abortController.signal,
        )
    }
}
```

`ChatRequest.message` 在 resume 时可为空，后端通过 Manager 判断是否走 resume 分支。

### 生成中的对话列表展示

`ListConversations` 不包含 `generating` 字段（避免每个都查 Manager）。前端通过已打开的对话的 `GetConversation` 来判断。如需在列表中显示"生成中"标识，可在后续迭代中加。

## Risks / Trade-offs

- 部分内容只在内存中，服务器重启丢失（需用户重新发送消息触发生成）
- Manager 的 subscriber map 需要处理客户端中途断开时的清理
- 当前 `ChatStream` 返回 unbuffered channel，Manager 的 goroutine 写入 channel 时需另一个 goroutine 消费，避免阻塞
