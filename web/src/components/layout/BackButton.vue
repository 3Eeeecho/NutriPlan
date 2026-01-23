<template>
  <n-button
    text
    @click="goBack"
    class="back-button-component"
  >
    <template #icon>
      <n-icon size="20">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
          <path fill="currentColor" d="M20 11H7.83l5.59-5.59L12 4l-8 8 8 8 1.41-1.41L7.83 13H20v-2z"/>
        </svg>
      </n-icon>
    </template>
    {{ text }}
  </n-button>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { NButton, NIcon } from 'naive-ui'

const props = defineProps({
  text: {
    type: String,
    default: '返回',
  },
  to: {
    type: String,
    default: null,
  },
})

const router = useRouter()

const goBack = () => {
  if (props.to) {
    router.push(props.to)
  } else {
    router.back()
  }
}
</script>

<style scoped>
.back-button-component {
  color: var(--color-text-secondary);
  font-weight: var(--font-weight-medium);
}

.back-button-component:hover {
  color: var(--color-primary);
}
</style>
