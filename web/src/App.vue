<template>
  <router-view v-slot="{ Component, route }">
    <transition name="fade" mode="out-in">
      <component :is="Component" :key="route.path" />
    </transition>
  </router-view>
</template>

<script setup>
// App 根组件，主要负责路由视图的渲染
</script>

<style>
#app {
  min-height: 100vh;
}
</style>

