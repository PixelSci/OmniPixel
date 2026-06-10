## Context

`/image` 是图片功能入口页，无侧边栏，纯 Hub 布局。子功能页先只做路由占位，后续再实现。

## Goals / Non-Goals

**Goals:**
- `/image` 页面：居中偏上的输入框 + 下方功能卡片网格
- 输入框带 model 选择器，输入后跳转 `/image/text-to-image`
- 6 个功能卡片，点击跳转对应子路由
- 6 个子路由占位页面
- 无侧边栏

**Non-Goals:**
- 不实现子功能页的具体功能
- 不做图片生成历史

## Layout

```
┌──────────────────────────────────────────────┐
│  MenuBar                                     │
├──────────────────────────────────────────────┤
│                                              │
│         ┌─────────────────────────┐          │
│         │  Describe the image     │          │
│         │  you want to create...  │  输入框   │
│         │              [model ▼][→]│          │
│         └─────────────────────────┘          │
│                                              │
│    ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│    │  Text to │ │ Image to │ │  Remove  │   │
│    │  Image   │ │  Image   │ │    BG    │   │
│    └──────────┘ └──────────┘ └──────────┘   │
│                                              │
│    ┌──────────┐ ┌──────────┐ ┌──────────┐   │
│    │  Enhance │ │  Style   │ │ Inpaint  │   │
│    │          │ │ Transfer │ │          │   │
│    └──────────┘ └──────────┘ └──────────┘   │
│                                              │
└──────────────────────────────────────────────┘
```

## Decisions

### 路由结构

```
/image                    → pages/image.vue (Hub)
  /image/text-to-image    → pages/image/TextToImage.vue
  /image/image-to-image   → pages/image/ImageToImage.vue
  /image/remove-bg        → pages/image/RemoveBg.vue
  /image/enhance          → pages/image/Enhance.vue
  /image/style-transfer   → pages/image/StyleTransfer.vue
  /image/inpaint          → pages/image/Inpaint.vue
```

使用 Vue Router 的嵌套 `children`，Hub 页通过 `<RouterView>` 渲染子页面（子路由匹配时 Hub 内容被替换）。

### 输入框

复用现有 `HigPromptInput` 组件，带 model 选择。子组件通过 emit `@send` 携带 prompt + model，Hub 页用 `router.push` 跳转。

### 功能卡片

使用现有 `Card` 组件（shadcn/ui 风格），2 行 × 3 列网格。卡片内容：
- 图标（lucide-vue-next 中选取）
- 英文标题
- 简短描述

卡片点击用 `router.push` 跳转到对应子路由。

### 子页占位

每个子页只做居中显示标题的占位，和当前 video/audio 页一致风格。

### Model 选择

`HigPromptInput` 已内置 model 选择逻辑，直接复用，无需额外开发。

## Risks

- `HigPromptInput` 的 model 选择引入后，确认其 `@send` 的 payload 结构，确保 model 字段正确拼接到跳转参数
