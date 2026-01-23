<template>
  <div class="stat-card" :class="{ clickable: clickable }" @click="handleClick">
    <div class="stat-icon-wrapper" :style="{ background: iconBg }">
      <span class="stat-icon">{{ icon }}</span>
    </div>
    <div class="stat-content">
      <div class="stat-value">{{ value }}</div>
      <div class="stat-label">{{ label }}</div>
      <div v-if="subtitle" class="stat-subtitle">{{ subtitle }}</div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  icon: {
    type: String,
    required: true,
  },
  iconBg: {
    type: String,
    default: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
  },
  value: {
    type: [String, Number],
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
  subtitle: {
    type: String,
    default: '',
  },
  clickable: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['click'])

const handleClick = () => {
  if (props.clickable) {
    emit('click')
  }
}
</script>

<style scoped>
.stat-card {
  background: rgba(255, 255, 255, 0.4);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  display: flex;
  align-items: center;
  gap: var(--spacing-lg);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all var(--transition-base);
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: var(--color-primary);
  background: rgba(255, 255, 255, 0.6);
}

.stat-icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.6);
}

.stat-icon {
  font-size: 28px;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  line-height: 1.2;
  font-family: var(--font-family-number);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin-top: 4px;
}

.stat-subtitle {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  margin-top: 2px;
}

@media (max-width: 768px) {
  .stat-card {
    padding: var(--spacing-lg);
  }

  .stat-icon-wrapper {
    width: 48px;
    height: 48px;
  }

  .stat-icon {
    font-size: 24px;
  }

  .stat-value {
    font-size: var(--font-size-2xl);
  }
}
</style>
