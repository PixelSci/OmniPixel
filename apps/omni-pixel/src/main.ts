import { createApp } from 'vue'
import App from './App.vue'
import { router } from './routes/route'
import './style.css'

createApp(App).use(router).mount('#app')
