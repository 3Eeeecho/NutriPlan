<template>
  <img
    v-bind="$attrs"
    :class="{ 'recipe-image--placeholder': isPlaceholder }"
    :data-placeholder="isPlaceholder ? 'true' : 'false'"
    :src="currentSrc"
    :alt="resolvedAlt"
    @error="handleError"
  />
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { buildRecipePlaceholderSrc } from '@/utils/recipePlaceholder'

defineOptions({
  inheritAttrs: false
})

const props = defineProps({
  src: {
    type: String,
    default: ''
  },
  name: {
    type: String,
    default: ''
  },
  mealType: {
    type: String,
    default: ''
  },
  ingredients: {
    type: [Array, String],
    default: () => []
  },
  targetUsers: {
    type: [Array, String],
    default: () => []
  },
  alt: {
    type: String,
    default: ''
  }
})

const failed = ref(false)

watch(() => props.src, () => {
  failed.value = false
})

const placeholderSrc = computed(() => buildRecipePlaceholderSrc({
  name: props.name,
  mealType: props.mealType,
  ingredients: props.ingredients,
  targetUsers: props.targetUsers
}))

const currentSrc = computed(() => {
  const source = String(props.src || '').trim()
  if (source && !failed.value) {
    return source
  }
  return placeholderSrc.value
})

const isPlaceholder = computed(() => {
  const source = String(props.src || '').trim()
  return !source || failed.value
})

const resolvedAlt = computed(() => props.alt || props.name || 'recipe image')

function handleError(event) {
  if (!failed.value && String(props.src || '').trim()) {
    failed.value = true
    return
  }
  if (event?.target) {
    event.target.onerror = null
  }
}
</script>
