<template>
  <Teleport to="body">
    <div v-if="show" class="modal-overlay" @click.self="$emit('close')">
      <div class="modal" :class="sizeClass">
        <div class="modal-header">
          <h3>{{ title }}</h3>
          <button class="btn btn-ghost btn-icon" @click="$emit('close')">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M18 6L6 18" />
              <path d="M6 6l12 12" />
            </svg>
          </button>
        </div>
        <slot />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    show: boolean
    title: string
    size?: 'sm' | 'md' | 'lg' | 'xl'
  }>(),
  {
    size: 'md',
  }
)

defineEmits<{
  close: []
}>()

const sizeClass = computed(() => {
  switch (props.size) {
    case 'lg':
      return 'modal-lg'
    case 'xl':
      return 'modal-xl'
    default:
      return ''
  }
})
</script>
