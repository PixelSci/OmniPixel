## Why

`/image` 页面当前为空占位。需要创建一个功能入口 Hub 页——居中靠上的输入框用于快速进入文生图，下方是功能卡片网格用于进入其它图片功能子页。

## What Changes

- **新增 `/image` Hub 页** — 居中偏上布局：输入框（带 model 选择器）+ 6 个功能卡片
- **功能卡片** — Text to Image、Image to Image、Remove BG、Enhance、Style Transfer、Inpaint
- **路由扩展** — 新增 6 个三级占位路由 `/image/text-to-image` 等，创建对应占位页面
- **输入框交互** — 输入 prompt → 跳转 `/image/text-to-image?prompt=...&model=...`
- **卡片交互** — 点击 → 跳转对应子路由

## Impact

- `apps/omni-pixel/src/pages/image.vue` — 重写为 Hub 页
- `apps/omni-pixel/src/routes/route.ts` — 新增嵌套路由
- `apps/omni-pixel/src/pages/image/` — 新建目录，存放子页占位
