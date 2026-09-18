import { createApp } from 'vue'
import 'element-plus/theme-chalk/dark/css-vars.css'
import App from './App.vue'
import router from './router'
import i18n from './i18n'

import './assets/css/var.css'
import './assets/css/dark.css'

const app = createApp(App)

app.use(i18n)
app.use(router)

app.mount('#app')
