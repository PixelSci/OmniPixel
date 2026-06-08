import { createMemoryHistory, createRouter } from 'vue-router'

const routes = [
    { path: '/', redirect: '/text' },
    { path: '/text', component: () => import('@/pages/index.vue') },
]

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
})
