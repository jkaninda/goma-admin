<template>
  <div class="card">
    <div class="empty-state">
      <svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <h3>{{ title }}</h3>
      <p>{{ description }}</p>
      <div v-if="$slots.action" class="empty-action">
        <slot name="action" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description: string
}>()
</script>

<style scoped>
.empty-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 12px;
  color: var(--text-muted);
}
.empty-action {
  margin-top: 16px;
}
</style>
