import { createApp } from 'vue'
import { createPinia } from 'pinia'

// Element Plus (逐步移除)
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

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

// 注册所有 Element Plus 图标（逐步移除）
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(pinia)
app.use(router)
app.use(ElementPlus, { locale: zhCn })

// 全局配置 Naive UI 主题
app.provide('themeOverrides', themeOverrides)

app.mount('#app')


