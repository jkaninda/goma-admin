<template>
  <div>
    <div v-if="loading" class="loading-page"><div class="spinner"></div></div>

    <div v-else-if="error" class="detail-error">
      <div class="card card-body" style="text-align: center; padding: 48px 24px">
        <svg width="48" height="48" fill="none" stroke="var(--text-muted)" viewBox="0 0 24 24" style="margin: 0 auto 16px">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <path stroke-linecap="round" stroke-width="1.5" d="M12 8v4m0 4h.01" />
        </svg>
        <h2 style="margin-bottom: 8px">Route not found</h2>
        <p class="text-muted" style="margin-bottom: 20px">The route you're looking for doesn't exist or has been deleted.</p>
        <router-link to="/routes" class="btn btn-primary">Back to Routes</router-link>
      </div>
    </div>

    <div v-else-if="route">
      <!-- Header -->
      <div class="detail-header">
        <router-link to="/routes" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <div class="detail-header-info">
          <h1>{{ route.name }}</h1>
          <span :class="['status-indicator', config.enabled !== false ? 'status-online' : 'status-offline']">
            <span class="status-dot"></span>
            {{ config.enabled !== false ? 'Online' : 'Offline' }}
          </span>
        </div>
        <div class="detail-header-actions">
          <button class="btn btn-secondary btn-sm" @click="handleEdit">Edit</button>
          <button class="btn btn-danger btn-sm" @click="handleDelete">Delete</button>
        </div>
      </div>

      <!-- Route Flow -->
      <div class="route-flow card">
        <div class="flow-node">
          <div class="flow-label">Client</div>
          <div class="flow-value">
            <template v-if="hostsList.length">
              <a v-for="h in hostsList" :key="h" :href="hostUrl(h)" target="_blank" rel="noopener" class="badge badge-info flow-badge flow-host-link">{{ h }}</a>
            </template>
            <span v-else class="text-muted">Any host</span>
          </div>
        </div>
        <div class="flow-arrow">
          <svg width="24" height="24" fill="none" stroke="var(--text-muted)" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
        <div class="flow-node flow-node-highlight">
          <div class="flow-label">Path</div>
          <div class="flow-value text-mono">{{ config.path || '/' }}</div>
        </div>
        <div class="flow-arrow">
          <svg width="24" height="24" fill="none" stroke="var(--text-muted)" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
          <div v-if="config.rewrite" class="flow-rewrite">{{ config.rewrite }}</div>
        </div>
        <!-- Single target -->
        <div v-if="!backendsList.length" class="flow-node">
          <div class="flow-label">Target</div>
          <div class="flow-value text-mono">{{ config.target || '-' }}</div>
        </div>
        <!-- Load-balanced backends -->
        <div v-else class="flow-backends">
          <div v-for="(b, i) in backendsList" :key="i" class="flow-node flow-backend">
            <div class="flow-label">Backend {{ i + 1 }}</div>
            <div class="flow-value text-mono">{{ b.endpoint }}</div>
            <div v-if="isWeighted" class="flow-weight">
              <div class="flow-weight-bar" :style="{ width: b.weight + '%' }"></div>
              <span class="flow-weight-label">{{ b.weight }}%</span>
            </div>
          </div>
        </div>
      </div>

      <div class="detail-cards">
        <!-- Core Configuration -->
        <div class="card">
          <div class="card-header"><h2>Configuration</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Path</dt>
                <dd class="text-mono">{{ config.path || '-' }}</dd>
              </div>
              <div v-if="!backendsList.length" class="detail-row">
                <dt>Target</dt>
                <dd class="text-mono">{{ config.target || '-' }}</dd>
              </div>
              <div v-else class="detail-row">
                <dt>Routing</dt>
                <dd>
                  <span v-if="isWeighted" class="badge badge-warning">Canary</span>
                  <span v-else class="badge badge-info">Round Robin</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Rewrite</dt>
                <dd class="text-mono">{{ config.rewrite || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Hosts</dt>
                <dd>
                  <template v-if="hostsList.length">
                    <span v-for="h in hostsList" :key="h" class="badge badge-info detail-badge">{{ h }}</span>
                  </template>
                  <span v-else class="text-muted">Any</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Methods</dt>
                <dd>
                  <template v-if="methodsList.length">
                    <span v-for="m in methodsList" :key="m" :class="['badge', 'detail-badge', methodBadge(m)]">{{ m }}</span>
                  </template>
                  <span v-else class="text-muted">All</span>
                </dd>
              </div>
              <div v-if="config.priority" class="detail-row">
                <dt>Priority</dt>
                <dd>{{ config.priority }}</dd>
              </div>
              <div class="detail-row">
                <dt>Metrics</dt>
                <dd>
                  <span :class="['badge', config.disableMetrics ? 'badge-neutral' : 'badge-success']">
                    {{ config.disableMetrics ? 'Disabled' : 'Enabled' }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Backends (Load Balancing) -->
        <div v-if="backendsList.length" class="card">
          <div class="card-header">
            <h2>Backends</h2>
            <span v-if="isWeighted" class="badge badge-warning">Canary</span>
            <span v-else class="badge badge-info">Round Robin</span>
            <span class="badge badge-info" style="margin-left: 4px">{{ backendsList.length }} endpoints</span>
          </div>
          <div class="card-body">
            <div class="backends-list">
              <div v-for="(b, i) in backendsList" :key="i" class="backend-item">
                <div class="backend-info">
                  <span class="backend-endpoint text-mono">{{ b.endpoint }}</span>
                  <span v-if="isWeighted" class="backend-weight-label">{{ b.weight }}%</span>
                </div>
                <div v-if="isWeighted" class="backend-bar-track">
                  <div class="backend-bar-fill" :style="{ width: b.weight + '%' }"></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Middlewares -->
        <div v-if="middlewaresList.length" class="card">
          <div class="card-header">
            <h2>Middlewares</h2>
            <span class="badge badge-info">{{ middlewaresList.length }}</span>
          </div>
          <div class="card-body">
            <div class="middleware-list">
              <div v-for="(mw, i) in middlewaresList" :key="mw" class="middleware-item">
                <span class="middleware-order">{{ i + 1 }}</span>
                <span class="middleware-name text-mono">{{ mw }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Health Check -->
        <div v-if="healthCheck" class="card">
          <div class="card-header">
            <h2>Health Check</h2>
            <span class="badge badge-success">Active</span>
          </div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Path</dt>
                <dd class="text-mono">{{ healthCheck.path || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Interval</dt>
                <dd>{{ healthCheck.interval || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Timeout</dt>
                <dd>{{ healthCheck.timeout || '-' }}</dd>
              </div>
              <div v-if="healthCheck.healthyStatuses?.length" class="detail-row">
                <dt>Healthy Statuses</dt>
                <dd>
                  <span v-for="s in healthCheck.healthyStatuses" :key="s" class="badge badge-success detail-badge">{{ s }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Security -->
        <div v-if="security" class="card">
          <div class="card-header"><h2>Security</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Forward Host Headers</dt>
                <dd>
                  <span :class="['badge', security.forwardHostHeaders !== false ? 'badge-success' : 'badge-neutral']">
                    {{ security.forwardHostHeaders !== false ? 'Yes' : 'No' }}
                  </span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Exploit Protection</dt>
                <dd>
                  <span :class="['badge', security.enableExploitProtection ? 'badge-success' : 'badge-neutral']">
                    {{ security.enableExploitProtection ? 'Enabled' : 'Disabled' }}
                  </span>
                </dd>
              </div>
              <div v-if="security.tls" class="detail-row">
                <dt>TLS Skip Verify</dt>
                <dd>
                  <span :class="['badge', security.tls.insecureSkipVerify ? 'badge-warning' : 'badge-success']">
                    {{ security.tls.insecureSkipVerify ? 'Yes (insecure)' : 'No' }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Extra fields (catch-all for anything unknown) -->
        <div v-if="hasExtraFields" class="card">
          <div class="card-header"><h2>Additional Configuration</h2></div>
          <div class="card-body">
            <pre class="code-block">{{ formatExtraFields() }}</pre>
          </div>
        </div>

        <!-- Metadata -->
        <div class="card">
          <div class="card-header"><h2>Metadata</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Route ID</dt>
                <dd class="text-mono">{{ route.id }}</dd>
              </div>
              <div class="detail-row">
                <dt>Created</dt>
                <dd>{{ formatDate(route.createdAt) }}</dd>
              </div>
              <div class="detail-row">
                <dt>Last Updated</dt>
                <dd>{{ formatDate(route.updatedAt) }}</dd>
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
import { routesApi, type Route } from '@/api/routes'
import { useConfirm } from '@/composables/useConfirm'
import { useNotificationStore } from '@/stores/notification'

const props = defineProps<{ id: string }>()
const router = useRouter()
const { confirm } = useConfirm()
const notify = useNotificationStore()

const loading = ref(true)
const error = ref(false)
const route = ref<Route | null>(null)

const knownFieldKeys = new Set([
  'path', 'target', 'hosts', 'methods', 'rewrite', 'enabled',
  'priority', 'disableMetrics', 'middlewares', 'healthCheck', 'security',
  'backends',
])

interface Backend {
  endpoint: string
  weight?: number
}

const config = computed(() => route.value?.config || {})

const hostsList = computed(() => {
  const hosts = config.value.hosts
  return Array.isArray(hosts) ? hosts.map(String) : []
})

const methodsList = computed(() => {
  const methods = config.value.methods
  return Array.isArray(methods) ? methods.map(String) : []
})

const backendsList = computed((): Backend[] => {
  const backends = config.value.backends
  if (!Array.isArray(backends)) return []
  return backends.map((b: unknown) => {
    if (typeof b === 'string') return { endpoint: b }
    if (b && typeof b === 'object') {
      const obj = b as Record<string, unknown>
      return {
        endpoint: String(obj.endpoint || obj.url || obj.target || ''),
        weight: typeof obj.weight === 'number' ? obj.weight : undefined,
      }
    }
    return { endpoint: String(b) }
  })
})

const isWeighted = computed(() =>
  backendsList.value.some(b => b.weight != null && b.weight > 0)
)

const middlewaresList = computed(() => {
  const mw = config.value.middlewares
  return Array.isArray(mw) ? mw.map(String) : []
})

interface HealthCheckConfig {
  path?: string
  interval?: string
  timeout?: string
  healthyStatuses?: number[]
}

interface SecurityConfig {
  forwardHostHeaders?: boolean
  enableExploitProtection?: boolean
  tls?: { insecureSkipVerify?: boolean }
}

const healthCheck = computed((): HealthCheckConfig | null => {
  const hc = config.value.healthCheck
  return hc && typeof hc === 'object' ? hc as HealthCheckConfig : null
})

const security = computed((): SecurityConfig | null => {
  const sec = config.value.security
  return sec && typeof sec === 'object' ? sec as SecurityConfig : null
})

const extraFields = computed(() => {
  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(config.value)) {
    if (!knownFieldKeys.has(key)) {
      result[key] = value
    }
  }
  return result
})

const hasExtraFields = computed(() => Object.keys(extraFields.value).length > 0)

function hostUrl(host: string): string {
  const base = (host.startsWith('http://') || host.startsWith('https://')) ? host : `http://${host}`
  const path = (config.value.path as string) || ''
  const cleanPath = path.replace(/\/?\*$/, '')
  return base.replace(/\/$/, '') + cleanPath
}

function methodBadge(method: string): string {
  switch (method) {
    case 'GET': return 'badge-success'
    case 'POST': return 'badge-primary'
    case 'PUT': case 'PATCH': return 'badge-warning'
    case 'DELETE': return 'badge-danger'
    default: return 'badge-info'
  }
}

function formatExtraFields(): string {
  return JSON.stringify(extraFields.value, null, 2)
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
  router.push(`/routes?edit=${route.value?.id}`)
}

async function handleDelete() {
  if (!route.value) return
  const confirmed = await confirm({
    title: 'Delete Route',
    message: `Are you sure you want to delete "${route.value.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await routesApi.delete(route.value.id)
    notify.success(`Route "${route.value.name}" deleted`)
    router.push('/routes')
  } catch {
    notify.error('Failed to delete route')
  }
}

onMounted(async () => {
  try {
    const res = await routesApi.get(Number(props.id))
    route.value = res.data
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
/* Header */
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
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

/* Status indicator */
.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 9999px;
  flex-shrink: 0;
}
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-online {
  color: var(--success-600);
  background: var(--success-50);
}
.status-online .status-dot {
  background: var(--success-500);
  animation: pulse-green 2s infinite;
}
.status-offline {
  color: var(--danger-500);
  background: var(--danger-50);
}
.status-offline .status-dot {
  background: var(--danger-500);
  animation: pulse-red 2s infinite;
}
@keyframes pulse-green {
  0% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.6); }
  70% { box-shadow: 0 0 0 6px rgba(34, 197, 94, 0); }
  100% { box-shadow: 0 0 0 0 rgba(34, 197, 94, 0); }
}
@keyframes pulse-red {
  0% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6); }
  70% { box-shadow: 0 0 0 6px rgba(239, 68, 68, 0); }
  100% { box-shadow: 0 0 0 0 rgba(239, 68, 68, 0); }
}

/* Route flow visualization */
.route-flow {
  display: flex;
  align-items: center;
  gap: 0;
  padding: 20px 24px;
  margin-bottom: 20px;
  overflow-x: auto;
}
.flow-node {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px 16px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
  min-width: 120px;
}
.flow-node-highlight {
  border-color: var(--primary-300);
  background: var(--primary-50, rgba(147, 51, 234, 0.04));
}
.flow-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-muted);
}
.flow-value {
  font-size: 13px;
  color: var(--text-primary);
  word-break: break-all;
}
.flow-badge {
  margin-right: 4px;
  margin-bottom: 2px;
}
.flow-host-link {
  text-decoration: none;
  cursor: pointer;
  transition: opacity 0.15s;
}
.flow-host-link:hover {
  opacity: 0.8;
  text-decoration: underline;
}
.flow-arrow {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0 8px;
  flex-shrink: 0;
}
.flow-rewrite {
  font-size: 10px;
  color: var(--primary-600);
  font-family: 'SF Mono', 'Fira Code', Menlo, Consolas, monospace;
  white-space: nowrap;
  margin-top: 2px;
}

/* Cards layout */
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
.detail-badge {
  margin-right: 4px;
  margin-bottom: 2px;
}
.text-mono {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
}
.text-muted {
  color: var(--text-muted);
  font-style: italic;
}

/* Middlewares */
.middleware-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.middleware-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
}
.middleware-order {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--primary-100, rgba(147, 51, 234, 0.1));
  color: var(--primary-700);
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}
.middleware-name {
  font-size: 14px;
  color: var(--text-primary);
}

/* Backends (load balancing) */
.flow-backends {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.flow-backend {
  min-width: 180px;
}
.flow-weight {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}
.flow-weight-bar {
  height: 4px;
  border-radius: 2px;
  background: var(--primary-500);
  min-width: 4px;
}
.flow-weight-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  white-space: nowrap;
}

.backends-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.backend-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 10px 14px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
}
.backend-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.backend-endpoint {
  font-size: 14px;
  color: var(--text-primary);
}
.backend-weight-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  flex-shrink: 0;
}
.backend-bar-track {
  height: 6px;
  border-radius: 3px;
  background: var(--border-primary);
  overflow: hidden;
}
.backend-bar-fill {
  height: 100%;
  border-radius: 3px;
  background: var(--primary-500);
  transition: width 0.3s ease;
  min-width: 4px;
}
</style>
