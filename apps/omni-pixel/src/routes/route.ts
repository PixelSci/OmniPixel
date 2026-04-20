import { createMemoryHistory, createRouter } from 'vue-router'

const routes = [
    { path: '/', component: () => import('@/pages/index.vue') },
]

export const router = createRouter({
    history: createMemoryHistory(),
    routes,
})
