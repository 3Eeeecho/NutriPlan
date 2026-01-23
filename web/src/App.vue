<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-loading-bar-provider>
        <n-dialog-provider>
          <router-view v-slot="{ Component, route }">
            <transition name="fade" mode="out-in">
              <keep-alive include="RecipeRecommend">
                <component :is="Component" :key="route.path" />
              </keep-alive>
            </transition>
          </router-view>
        </n-dialog-provider>
      </n-loading-bar-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup>
import { NConfigProvider, NMessageProvider, NLoadingBarProvider, NDialogProvider } from 'naive-ui'
import { themeOverrides } from './styles/theme'
</script>

<style>
#app {
  min-height: 100vh;
}
</style>

