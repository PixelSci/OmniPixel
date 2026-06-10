import { createMemoryHistory, createRouter } from 'vue-router'

const routes = [
    { path: '/', redirect: '/text' },
    { path: '/text', component: () => import('@/pages/index.vue') },
    {
        path: '/image',
        component: () => import('@/pages/image.vue'),
        children: [
            { path: 'text-to-image', component: () => import('@/pages/image/TextToImage.vue') },
            { path: 'image-to-image', component: () => import('@/pages/image/ImageToImage.vue') },
            { path: 'remove-bg', component: () => import('@/pages/image/RemoveBg.vue') },
            { path: 'enhance', component: () => import('@/pages/image/Enhance.vue') },
            { path: 'style-transfer', component: () => import('@/pages/image/StyleTransfer.vue') },
            { path: 'inpaint', component: () => import('@/pages/image/Inpaint.vue') },
        ],
    },
    { path: '/audio', component: () => import('@/pages/audio.vue') },
    { path: '/video', component: () => import('@/pages/video.vue') },
]

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
})
