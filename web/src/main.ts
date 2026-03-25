import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/styles.css'
// Import instance store to register the axios interceptor
import '@/stores/instance'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
