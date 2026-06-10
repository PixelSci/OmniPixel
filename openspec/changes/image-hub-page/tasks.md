## 1. Routes

- [x] 1.1 `routes/route.ts` — 将 `/image` 改为嵌套路由，添加 6 个子路由（均指向占位组件）

## 2. Sub-page placeholders

- [x] 2.1 `pages/image/TextToImage.vue` — 占位页
- [x] 2.2 `pages/image/ImageToImage.vue` — 占位页
- [x] 2.3 `pages/image/RemoveBg.vue` — 占位页
- [x] 2.4 `pages/image/Enhance.vue` — 占位页
- [x] 2.5 `pages/image/StyleTransfer.vue` — 占位页
- [x] 2.6 `pages/image/Inpaint.vue` — 占位页

## 3. Hub page

- [x] 3.1 `pages/image.vue` — 重写为 Hub 页：输入框 + 功能卡片网格，挂接跳转逻辑

## 4. Verify

- [x] 4.1 验证 `/image` 展示 Hub 页，输入框 + 6 张卡片
- [x] 4.2 验证输入 prompt 后跳转到 `/image/text-to-image?prompt=...&model=...`
- [x] 4.3 验证点击卡片跳转到对应子路由
- [x] 4.4 验证子路由不显示侧边栏
- [x] 4.5 `pnpm build` 编译通过
