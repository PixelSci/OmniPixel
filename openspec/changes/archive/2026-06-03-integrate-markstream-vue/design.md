## Context

`markstream-vue@1.0.1-beta.3` 提供 `<MarkdownRender>` 组件，核心依赖 `markstream-core` + `stream-markdown-parser`。可选 peer deps：`stream-markdown`（Shiki 代码高亮）、`katex`、`mermaid`。

当前 `HigAssistantMessage` 只接收 `content` string，渲染为 plain text。需要改为 markdown 渲染，区分 streaming 和 completed 模式。

## Goals / Non-Goals

**Goals:**
- 安装 `markstream-vue` + `stream-markdown` + `katex` + `mermaid`
- `HigAssistantMessage` 使用 `<MarkdownRender>` 渲染内容
- Streaming 模式：增量批处理 (`max-live-nodes="0"`, `:final="false"`)
- Completed 模式：虚拟化窗口 (`max-live-nodes="320"`, `:final="true"`)
- 暗色模式跟随系统

**Non-Goals:**
- 不集成 `stream-monaco`（代码编辑器，非聊天场景需要）
- 不集成 `@antv/infographic`、`@terrastruct/d2`
- 不修改 SSE streaming 逻辑（token 累积方式不变）
- 不添加 SSR 预解析（纯客户端渲染）

## Decisions

### 依赖选择

| 包 | 用途 | 安装 |
|---|---|---|
| `markstream-vue` | 核心渲染组件 | ✓ |
| `stream-markdown` | Shiki 代码高亮 | ✓ |
| `katex` | 数学公式渲染 | ✓ |
| `mermaid` | 图表渲染 | ✓ |
| `stream-monaco` | 代码编辑器 | ✗（不需要） |

### 双模式渲染

```
streaming=true (正在接收 token):
  <MarkdownRender
    :content="content"
    :final="false"
    :max-live-nodes="0"
    :batch-rendering="true"
    :render-batch-size="16"
    :render-batch-delay="8"
    :typewriter="true"
    :is-dark="isDark"
  />

streaming=false (已完成的消息):
  <MarkdownRender
    :content="content"
    :final="true"
    :max-live-nodes="320"
    :is-dark="isDark"
  />
```

### 暗色模式

`@vueuse/core` 的 `useDark()` 已是全局单例，直接在 `HigAssistantMessage` 中调用即可获取响应式暗色状态。

### 启用可选功能

在组件 setup 中调用：
```ts
import { enableKatex, enableMermaid } from 'markstream-vue'
import 'katex/dist/katex.min.css'

enableKatex()
enableMermaid()
```

### 滚动行为

当前 `onToken` 回调中调用 `scrollToBottom()`。`markstream-vue` 的增量批处理可能在 token 到达后延迟渲染。改为用 `ResizeObserver` 监听内容容器高度变化自动滚底。

### `HigAssistantMessage` 接口

```typescript
interface Props {
    content: string
    streaming?: boolean   // 新增：标记是否为流式输出中
    model?: string
}
```

## Risks / Trade-offs

- `katex` + `mermaid` + `stream-markdown`（含 shiki）会增加约 5-10MB bundle。首屏加载时间略有增加，但聊天应用的核心体验值得
- `stream-markdown` 使用 Shiki，需要下载语法文件。首次渲染代码块时可能略有延迟
- `mermaid` 渲染大型图表时可能阻塞主线程，可后续迁移到 CDN worker
