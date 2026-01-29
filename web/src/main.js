import { createApp } from 'vue'
import { createPinia } from 'pinia'

// 样式
import './style.css'
import './styles/variables.css'
import './styles/mixins.css'

// 主题配置
import { themeOverrides } from './styles/theme'

import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// 全局配置 Naive UI 主题
app.provide('themeOverrides', themeOverrides)

app.mount('#app')
