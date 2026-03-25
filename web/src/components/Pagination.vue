<script setup lang="ts">
import { computed } from 'vue'
import type { Pageable } from '@/api/types'

const props = defineProps<{
  pageable: Pageable
}>()

const emit = defineEmits<{
  (e: 'page', page: number): void
}>()

const currentPage = computed(() => props.pageable.current_page)
const totalPages = computed(() => props.pageable.total_pages)

const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  if (total <= 7) return Array.from({ length: total }, (_, i) => i)

  const pages: (number | null)[] = [0]
  const start = Math.max(1, current - 1)
  const end = Math.min(total - 2, current + 1)

  if (start > 1) pages.push(null)
  for (let i = start; i <= end; i++) pages.push(i)
  if (end < total - 2) pages.push(null)
  pages.push(total - 1)

  return pages
})
</script>

<template>
  <div v-if="totalPages > 1" class="pagination">
    <span class="pagination-info">
      Page {{ currentPage + 1 }} of {{ totalPages }}
      ({{ pageable.total_elements }} total)
    </span>
    <div class="pagination-buttons">
      <button
        class="btn btn-secondary btn-sm"
        :disabled="currentPage === 0"
        @click="emit('page', currentPage - 1)"
      >
        Previous
      </button>
      <template v-for="(p, idx) in visiblePages">
        <span v-if="p === null" :key="'e' + idx" class="pagination-ellipsis">...</span>
        <button
          v-else
          :key="p"
          class="btn btn-sm"
          :class="p === currentPage ? 'btn-primary' : 'btn-secondary'"
          @click="emit('page', p)"
        >
          {{ p + 1 }}
        </button>
      </template>
      <button
        class="btn btn-secondary btn-sm"
        :disabled="currentPage >= totalPages - 1"
        @click="emit('page', currentPage + 1)"
      >
        Next
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination-ellipsis {
  display: inline-flex;
  align-items: center;
  padding: 0 4px;
  color: var(--text-muted);
  font-size: 13px;
}
</style>
