<template>
  <div>
    <div class="page-header">
      <div class="page-header-left">
        <h1>Audit Log</h1>
      </div>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <EmptyState
      v-else-if="snapshots.length === 0"
      title="No config changes recorded"
      description="Changes to routes and middlewares will appear here."
    />

    <div v-else class="card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Time</th>
              <th>Action</th>
              <th>Resource</th>
              <th>Name</th>
              <th>User</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <template v-for="snap in snapshots" :key="snap.id">
              <tr>
                <td>{{ formatDate(snap.createdAt) }}</td>
                <td>
                  <span :class="['badge', actionBadge(snap.action)]">{{ formatAction(snap.action) }}</span>
                </td>
                <td>{{ snap.resource }}</td>
                <td class="cell-name">{{ snap.name }}</td>
                <td>{{ snap.userId || '-' }}</td>
                <td class="text-right actions-cell">
                  <button class="btn btn-ghost btn-sm" @click="toggleExpand(snap.id)">
                    {{ expanded === snap.id ? 'Hide' : 'Details' }}
                  </button>
                  <button
                    v-if="snap.before"
                    class="btn btn-ghost btn-sm action-rollback"
                    @click="confirmRollback(snap)"
                  >Rollback</button>
                </td>
              </tr>
              <tr v-if="expanded === snap.id" class="detail-row">
                <td colspan="6">
                  <div class="diff-container">
                    <div v-if="snap.before" class="diff-panel">
                      <h4 class="diff-title diff-title-before">Before</h4>
                      <pre class="diff-content">{{ formatJSON(snap.before) }}</pre>
                    </div>
                    <div v-if="snap.after" class="diff-panel">
                      <h4 class="diff-title diff-title-after">After</h4>
                      <pre class="diff-content">{{ formatJSON(snap.after) }}</pre>
                    </div>
                    <div v-if="!snap.before && !snap.after" class="diff-panel">
                      <p class="text-muted">No diff data available.</p>
                    </div>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination">
        <button class="btn btn-ghost btn-sm" :disabled="page === 0" @click="goToPage(page - 1)">
          Previous
        </button>
        <span class="pagination-info">Page {{ page + 1 }} of {{ totalPages }}</span>
        <button class="btn btn-ghost btn-sm" :disabled="page >= totalPages - 1" @click="goToPage(page + 1)">
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { auditApi, type ConfigSnapshot } from '@/api/audit'
import { useConfirm } from '@/composables/useConfirm'
import EmptyState from '@/components/EmptyState.vue'

const { confirm } = useConfirm()

/* -- State -- */
const loading = ref(true)
const snapshots = ref<ConfigSnapshot[]>([])
const expanded = ref<number | null>(null)
const page = ref(0)
const pageSize = 20
const totalPages = ref(0)

/* -- Helpers -- */
function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleString()
}

function formatAction(action: string): string {
  return action.replace(/_/g, ' ')
}

function actionBadge(action: string): string {
  if (action.includes('deleted')) return 'badge-danger'
  if (action.includes('created')) return 'badge-success'
  if (action.includes('updated')) return 'badge-warning'
  return 'badge-info'
}

function formatJSON(obj: Record<string, unknown> | undefined): string {
  if (!obj) return ''
  return JSON.stringify(obj, null, 2)
}

function toggleExpand(id: number) {
  expanded.value = expanded.value === id ? null : id
}

/* -- Pagination -- */
function goToPage(p: number) {
  page.value = p
  fetchSnapshots()
}

/* -- Rollback -- */
async function confirmRollback(snap: ConfigSnapshot) {
  const confirmed = await confirm({
    title: 'Rollback Config',
    message: `Are you sure you want to rollback "${snap.name}" to its previous state? This will overwrite the current configuration.`,
    confirmText: 'Rollback',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await auditApi.rollback(snap.id)
    await fetchSnapshots()
  } catch {
    // Error handled by API interceptor
  }
}

/* -- Fetch -- */
async function fetchSnapshots() {
  loading.value = true
  try {
    const res = await auditApi.list(page.value, pageSize)
    snapshots.value = res.data.data
    totalPages.value = res.data.pageable.total_pages
  } catch {
    // Error handled by API interceptor
  } finally {
    loading.value = false
  }
}

onMounted(fetchSnapshots)
</script>

<style scoped>
.page-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cell-name {
  font-weight: 600;
  color: var(--text-primary);
}

.actions-cell {
  display: flex;
  gap: 4px;
  justify-content: flex-end;
}

.action-rollback {
  color: var(--warning-600);
}
.action-rollback:hover {
  color: var(--warning-700);
}

.detail-row td {
  padding: 0 !important;
  background: var(--bg-secondary);
}

.diff-container {
  display: flex;
  gap: 16px;
  padding: 16px;
}

.diff-panel {
  flex: 1;
  min-width: 0;
}

.diff-title {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 8px;
}

.diff-title-before {
  color: var(--danger-500);
}

.diff-title-after {
  color: var(--success-500, #22c55e);
}

.diff-content {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: 6px;
  padding: 12px;
  font-size: 0.8rem;
  line-height: 1.5;
  overflow-x: auto;
  max-height: 400px;
  overflow-y: auto;
  white-space: pre-wrap;
  word-break: break-word;
  margin: 0;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 16px;
  border-top: 1px solid var(--border-primary);
}

.pagination-info {
  font-size: 0.875rem;
  color: var(--text-secondary);
}
</style>
