import { createApp } from 'vue'
import 'element-plus/theme-chalk/dark/css-vars.css'
import App from './App.vue'
import router from './router'

import './assets/css/var.css'
import './assets/css/dark.css'
import './assets/css/tailwind.css'
import './assets/css/element-override.css'

const app = createApp(App)

app.use(router)

app.mount('#app')
