<template>
  <div>
    <div v-if="loading" class="loading-page"><div class="spinner"></div></div>

    <div v-else-if="error" class="detail-error">
      <div class="card card-body" style="text-align: center; padding: 48px 24px">
        <svg width="48" height="48" fill="none" stroke="var(--text-muted)" viewBox="0 0 24 24" style="margin: 0 auto 16px">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <path stroke-linecap="round" stroke-width="1.5" d="M12 8v4m0 4h.01" />
        </svg>
        <h2 style="margin-bottom: 8px">Middleware not found</h2>
        <p class="text-muted" style="margin-bottom: 20px">The middleware you're looking for doesn't exist or has been deleted.</p>
        <router-link to="/middlewares" class="btn btn-primary">Back to Middlewares</router-link>
      </div>
    </div>

    <div v-else-if="middleware">
      <!-- Header -->
      <div class="detail-header">
        <router-link to="/middlewares" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <div class="detail-header-info">
          <h1>{{ middleware.name }}</h1>
          <span :class="['badge', typeBadgeClass]">{{ middleware.type }}</span>
        </div>
        <div class="detail-header-actions">
          <button class="btn btn-secondary btn-sm" @click="handleEdit">Edit</button>
          <button class="btn btn-danger btn-sm" @click="handleDelete">Delete</button>
        </div>
      </div>

      <!-- Type summary banner -->
      <div v-if="typeDescription" class="type-banner">
        <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{{ typeDescription }}</span>
      </div>

      <div class="detail-cards">
        <!-- Overview -->
        <div class="card">
          <div class="card-header"><h2>Overview</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Name</dt>
                <dd class="text-mono">{{ middleware.name }}</dd>
              </div>
              <div class="detail-row">
                <dt>Type</dt>
                <dd>
                  <span :class="['badge', typeBadgeClass]">{{ middleware.type }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Configuration - structured -->
        <div class="card">
          <div class="card-header"><h2>Configuration</h2></div>
          <div class="card-body">
            <div v-if="!middleware.config || configEntries.length === 0" class="empty-config">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>No configuration parameters set</span>
            </div>
            <dl v-else class="detail-grid">
              <div v-for="entry in configEntries" :key="entry.key" class="detail-row">
                <dt>{{ formatKey(entry.key) }}</dt>
                <dd>
                  <!-- Boolean -->
                  <span v-if="typeof entry.value === 'boolean'" :class="['badge', entry.value ? 'badge-success' : 'badge-neutral']">
                    {{ entry.value ? 'Enabled' : 'Disabled' }}
                  </span>
                  <!-- Array of strings -->
                  <div v-else-if="isStringArray(entry.value)" class="config-tags">
                    <span v-for="(item, i) in toArray(entry.value)" :key="i" class="config-tag">{{ item }}</span>
                  </div>
                  <!-- Array of objects -->
                  <div v-else-if="Array.isArray(entry.value)" class="config-object-list">
                    <div v-for="(item, i) in toArray(entry.value)" :key="i" class="config-object-item">
                      <template v-if="typeof item === 'object' && item !== null">
                        <div v-for="k in objectKeys(item)" :key="k" class="config-nested-row">
                          <span class="config-nested-key">{{ k }}:</span>
                          <span class="text-mono">{{ formatNestedValue(item, k) }}</span>
                        </div>
                      </template>
                      <span v-else class="text-mono">{{ item }}</span>
                    </div>
                  </div>
                  <!-- Object -->
                  <div v-else-if="typeof entry.value === 'object' && entry.value !== null" class="config-nested">
                    <div v-for="k in objectKeys(entry.value)" :key="k" class="config-nested-row">
                      <span class="config-nested-key">{{ k }}:</span>
                      <span v-if="isNestedBool(entry.value, k)" :class="['badge', getNestedValue(entry.value, k) ? 'badge-success' : 'badge-neutral']">
                        {{ getNestedValue(entry.value, k) ? 'Enabled' : 'Disabled' }}
                      </span>
                      <span v-else class="text-mono">{{ formatNestedValue(entry.value, k) }}</span>
                    </div>
                  </div>
                  <!-- Number -->
                  <span v-else-if="typeof entry.value === 'number'" class="text-mono config-number">{{ entry.value }}</span>
                  <!-- Duration-like string -->
                  <span v-else-if="isDuration(entry.value)" class="config-duration">{{ entry.value }}</span>
                  <!-- URL-like string -->
                  <span v-else-if="isUrl(entry.value)" class="text-mono">{{ entry.value }}</span>
                  <!-- String / other -->
                  <span v-else class="text-mono">{{ entry.value }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Used by Routes -->
        <div v-if="usedByRoutes.length" class="card">
          <div class="card-header">
            <h2>Used by Routes</h2>
            <span class="badge badge-info">{{ usedByRoutes.length }}</span>
          </div>
          <div class="card-body">
            <div class="used-by-list">
              <router-link
                v-for="r in usedByRoutes"
                :key="r.id"
                :to="`/routes/${r.id}`"
                class="used-by-item"
              >
                <div class="used-by-info">
                  <span class="used-by-name">{{ r.name }}</span>
                  <span class="used-by-path text-mono">{{ r.config?.path || '/' }}</span>
                </div>
                <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
              </router-link>
            </div>
          </div>
        </div>

        <!-- Raw Configuration -->
        <div v-if="configEntries.length > 0" class="card">
          <div class="card-header">
            <h2>Raw Configuration</h2>
            <button class="btn btn-secondary btn-sm" @click="showRaw = !showRaw">
              {{ showRaw ? 'Hide' : 'Show' }}
            </button>
          </div>
          <div v-if="showRaw" class="card-body">
            <pre class="code-block">{{ JSON.stringify(middleware.config, null, 2) }}</pre>
          </div>
        </div>

        <!-- Metadata -->
        <div class="card">
          <div class="card-header"><h2>Metadata</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Middleware ID</dt>
                <dd class="text-mono">{{ middleware.id }}</dd>
              </div>
              <div class="detail-row">
                <dt>Created</dt>
                <dd>{{ formatDate(middleware.createdAt) }}</dd>
              </div>
              <div class="detail-row">
                <dt>Last Updated</dt>
                <dd>{{ formatDate(middleware.updatedAt) }}</dd>
              </div>
            </dl>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { middlewaresApi, type Middleware, type MiddlewareTypeInfo } from '@/api/middlewares'
import type { Route } from '@/api/routes'
import { useConfirm } from '@/composables/useConfirm'
import { useNotificationStore } from '@/stores/notification'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { confirm } = useConfirm()
const notify = useNotificationStore()

const loading = ref(true)
const error = ref(false)
const middleware = ref<Middleware | null>(null)
const showRaw = ref(false)
const usedByRoutes = ref<Route[]>([])
const typesCatalog = ref<MiddlewareTypeInfo[]>([])

// Map category → badge class
const categoryBadgeMap: Record<string, string> = {
  auth: 'badge-danger',
  security: 'badge-danger',
  traffic: 'badge-warning',
  transform: 'badge-info',
  performance: 'badge-success',
  observability: 'badge-neutral',
}

const currentTypeInfo = computed(() => {
  const t = middleware.value?.type || ''
  return typesCatalog.value.find(ti => ti.type === t) || null
})

const typeBadgeClass = computed(() => {
  const cat = currentTypeInfo.value?.category || ''
  return categoryBadgeMap[cat] || 'badge-info'
})

const typeDescription = computed(() => {
  return currentTypeInfo.value?.description || ''
})

const configEntries = computed(() => {
  const config = middleware.value?.config
  if (!config) return []
  return Object.entries(config)
    .filter(([, v]) => v !== undefined && v !== null)
    .map(([key, value]) => ({ key, value }))
})

function isStringArray(value: unknown): boolean {
  return Array.isArray(value) && value.length > 0 && value.every(v => typeof v === 'string' || typeof v === 'number')
}

function isDuration(value: unknown): boolean {
  return typeof value === 'string' && /^\d+[smhd]$/.test(value)
}

function isUrl(value: unknown): boolean {
  return typeof value === 'string' && /^https?:\/\//.test(value)
}

function toArray(value: unknown): unknown[] {
  return Array.isArray(value) ? value : []
}

function objectKeys(value: unknown): string[] {
  return value && typeof value === 'object' ? Object.keys(value) : []
}

function isNestedBool(obj: unknown, key: string): boolean {
  if (!obj || typeof obj !== 'object') return false
  return typeof (obj as Record<string, unknown>)[key] === 'boolean'
}

function getNestedValue(obj: unknown, key: string): unknown {
  if (!obj || typeof obj !== 'object') return undefined
  return (obj as Record<string, unknown>)[key]
}

function formatKey(key: string): string {
  return key.replace(/([A-Z])/g, ' $1').replace(/^./, s => s.toUpperCase()).trim()
}

function formatNestedValue(obj: unknown, key: string): string {
  const rec = obj as Record<string, unknown>
  const v = rec[key]
  return typeof v === 'string' ? v : JSON.stringify(v)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function handleEdit() {
  router.push(`/middlewares?edit=${middleware.value?.id}`)
}

async function handleDelete() {
  if (!middleware.value) return
  const confirmed = await confirm({
    title: 'Delete Middleware',
    message: usedByRoutes.value.length
      ? `"${middleware.value.name}" is used by ${usedByRoutes.value.length} route(s). Are you sure you want to delete it?`
      : `Are you sure you want to delete "${middleware.value.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await middlewaresApi.delete(middleware.value.id)
    notify.success(`Middleware "${middleware.value.name}" deleted`)
    router.push('/middlewares')
  } catch {
    notify.error('Failed to delete middleware')
  }
}

async function loadUsedByRoutes() {
  if (!middleware.value) return
  try {
    const res = await middlewaresApi.usage(middleware.value.id)
    usedByRoutes.value = res.data || []
  } catch {
    // non-critical
  }
}

onMounted(async () => {
  try {
    const [mwRes, typesRes] = await Promise.all([
      middlewaresApi.get(Number(props.id)),
      middlewaresApi.types(),
    ])
    middleware.value = mwRes.data
    typesCatalog.value = typesRes.data || []
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
  await loadUsedByRoutes()
})
</script>

<style scoped>
/* Header */
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}
.detail-header-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}
.detail-header h1 {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.detail-header-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}
.back-link {
  display: flex;
  align-items: center;
  color: var(--text-muted);
  transition: color var(--transition);
}
.back-link:hover { color: var(--text-secondary); }

/* Type banner */
.type-banner {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  margin-bottom: 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
}
.type-banner svg {
  flex-shrink: 0;
  margin-top: 1px;
  color: var(--text-muted);
}

/* Cards */
.detail-cards {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

/* Detail grid */
.detail-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0;
  margin: 0;
}
.detail-row {
  display: flex;
  align-items: baseline;
  padding: 10px 0;
  border-bottom: 1px solid var(--border-secondary);
}
.detail-row:last-child {
  border-bottom: none;
}
.detail-row dt {
  width: 160px;
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}
.detail-row dd {
  margin: 0;
  font-size: 14px;
  color: var(--text-primary);
}

.text-mono {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
}
.text-muted {
  color: var(--text-muted);
  font-style: italic;
}

/* Config values */
.config-number {
  color: var(--primary-600);
  font-weight: 600;
}
.config-duration {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
  color: var(--warning-700, #a16207);
  font-weight: 500;
}

.config-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
.config-tag {
  display: inline-block;
  padding: 2px 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: 4px;
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 12px;
  color: var(--text-secondary);
}

.config-nested {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.config-nested-row {
  display: flex;
  align-items: baseline;
  gap: 6px;
  font-size: 13px;
}
.config-nested-key {
  font-weight: 500;
  color: var(--text-tertiary);
}

.config-object-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.config-object-item {
  padding: 8px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
}

.empty-config {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 0;
  color: var(--text-muted);
  font-size: 13px;
}

/* Used by routes */
.used-by-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.used-by-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
  text-decoration: none;
  color: inherit;
  transition: border-color 0.15s, background 0.15s;
}
.used-by-item:hover {
  border-color: var(--primary-300);
  background: var(--primary-50, rgba(147, 51, 234, 0.04));
}
.used-by-item svg {
  flex-shrink: 0;
  color: var(--text-muted);
}
.used-by-info {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}
.used-by-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.used-by-path {
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
}
</style>
