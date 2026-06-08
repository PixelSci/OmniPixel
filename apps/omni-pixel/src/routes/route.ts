import { createMemoryHistory, createRouter } from 'vue-router'

const routes = [
    { path: '/', redirect: '/text' },
    { path: '/text', component: () => import('@/pages/index.vue') },
    { path: '/image', component: () => import('@/pages/image.vue') },
    { path: '/audio', component: () => import('@/pages/audio.vue') },
    { path: '/video', component: () => import('@/pages/video.vue') },
]

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
})
