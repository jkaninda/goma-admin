<template>
  <div>
    <div v-if="loading" class="loading-page"><div class="spinner"></div></div>

    <div v-else-if="middleware">
      <div class="detail-header">
        <router-link to="/middlewares" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <h1>{{ middleware.name }}</h1>
        <span :class="['badge', typeBadgeClass]">{{ middleware.type }}</span>
      </div>

      <div class="detail-cards">
        <!-- Overview -->
        <div class="card">
          <div class="card-header"><h2>Overview</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Name</dt>
                <dd>{{ middleware.name }}</dd>
              </div>
              <div class="detail-row">
                <dt>Type</dt>
                <dd>
                  <span :class="['badge', typeBadgeClass]">{{ middleware.type }}</span>
                  <span class="type-description">{{ typeDescription }}</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Instance ID</dt>
                <dd class="text-mono">{{ middleware.instanceId }}</dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Configuration -->
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
                <dt>{{ entry.key }}</dt>
                <dd>
                  <!-- Boolean -->
                  <span v-if="typeof entry.value === 'boolean'" :class="['badge', entry.value ? 'badge-success' : 'badge-neutral']">
                    {{ entry.value ? 'true' : 'false' }}
                  </span>
                  <!-- Array -->
                  <div v-else-if="Array.isArray(entry.value)" class="config-list">
                    <span v-for="(item, i) in entry.value" :key="i" class="config-list-item">
                      <template v-if="typeof item === 'object'">{{ JSON.stringify(item) }}</template>
                      <template v-else>{{ item }}</template>
                    </span>
                  </div>
                  <!-- Object -->
                  <div v-else-if="typeof entry.value === 'object' && entry.value !== null" class="config-nested">
                    <div v-for="k in Object.keys(entry.value)" :key="k" class="config-nested-row">
                      <span class="config-nested-key">{{ k }}:</span>
                      <span class="text-mono">{{ formatNestedValue(entry.value, k) }}</span>
                    </div>
                  </div>
                  <!-- Number -->
                  <span v-else-if="typeof entry.value === 'number'" class="text-mono config-number">{{ entry.value }}</span>
                  <!-- String / other -->
                  <span v-else class="text-mono">{{ entry.value }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Raw Configuration -->
        <div v-if="configEntries.length > 0" class="card">
          <div class="card-header">
            <h2>Raw Configuration</h2>
            <button class="btn btn-secondary btn-sm" @click="toggleRaw">
              {{ showRaw ? 'Hide' : 'Show' }}
            </button>
          </div>
          <div v-if="showRaw" class="card-body">
            <pre class="code-block">{{ formatConfig(middleware.config) }}</pre>
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
import { middlewaresApi, type Middleware } from '@/api/middlewares'

const props = defineProps<{ id: string }>()

const loading = ref(true)
const middleware = ref<Middleware | null>(null)
const showRaw = ref(false)

const typeDescriptions: Record<string, string> = {
  rateLimit: 'Limits the number of requests per time window',
  rateLimiter: 'Limits the number of requests per time window',
  cors: 'Configures Cross-Origin Resource Sharing headers',
  auth: 'Handles authentication and authorization',
  basicAuth: 'HTTP Basic Authentication',
  jwtAuth: 'JSON Web Token authentication',
  jwt: 'JSON Web Token authentication',
  oauth: 'OAuth 2.0 authentication flow',
  apiKey: 'API key-based authentication',
  headers: 'Adds or modifies request/response headers',
  addHeaders: 'Adds custom headers to requests or responses',
  rewrite: 'Rewrites the request path',
  redirect: 'Redirects requests to a different URL',
  retry: 'Retries failed requests',
  timeout: 'Sets request timeout limits',
  circuitBreaker: 'Circuit breaker for fault tolerance',
  cache: 'Caches responses for improved performance',
  compress: 'Compresses response bodies',
  ipWhitelist: 'Restricts access by IP address',
  accessLog: 'Logs request access details',
}

const typeBadgeClasses: Record<string, string> = {
  rateLimit: 'badge-warning',
  rateLimiter: 'badge-warning',
  cors: 'badge-info',
  auth: 'badge-danger',
  basicAuth: 'badge-danger',
  jwtAuth: 'badge-danger',
  jwt: 'badge-danger',
  oauth: 'badge-danger',
  apiKey: 'badge-danger',
  headers: 'badge-info',
  addHeaders: 'badge-info',
  rewrite: 'badge-neutral',
  redirect: 'badge-neutral',
  retry: 'badge-warning',
  timeout: 'badge-warning',
  circuitBreaker: 'badge-danger',
  cache: 'badge-success',
  compress: 'badge-success',
  ipWhitelist: 'badge-danger',
  accessLog: 'badge-neutral',
}

const typeBadgeClass = computed(() => {
  const t = middleware.value?.type || ''
  return typeBadgeClasses[t] || 'badge-info'
})

const typeDescription = computed(() => {
  const t = middleware.value?.type || ''
  return typeDescriptions[t] || ''
})

const configEntries = computed(() => {
  const config = middleware.value?.config
  if (!config) return []
  return Object.entries(config)
    .filter(([, v]) => v !== undefined && v !== null)
    .map(([key, value]) => ({ key, value }))
})

function formatNestedValue(obj: unknown, key: string): string {
  const rec = obj as Record<string, unknown>
  const v = rec[key]
  return typeof v === 'string' ? v : JSON.stringify(v)
}

function toggleRaw() {
  showRaw.value = !showRaw.value
}

function formatConfig(config: Record<string, unknown>): string {
  if (!config) return ''
  const lines: string[] = []
  for (const [key, value] of Object.entries(config)) {
    if (value === undefined || value === null) continue
    if (Array.isArray(value)) {
      lines.push(`${key}:`)
      for (const item of value) {
        lines.push(`  - ${typeof item === 'object' ? JSON.stringify(item) : item}`)
      }
    } else if (typeof value === 'object') {
      lines.push(`${key}:`)
      for (const [k, v] of Object.entries(value as Record<string, unknown>)) {
        lines.push(`  ${k}: ${typeof v === 'string' ? v : JSON.stringify(v)}`)
      }
    } else {
      lines.push(`${key}: ${value}`)
    }
  }
  return lines.join('\n')
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

onMounted(async () => {
  try {
    const res = await middlewaresApi.get(Number(props.id))
    middleware.value = res.data
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
}
.detail-header h1 {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}
.back-link {
  display: flex;
  align-items: center;
  color: var(--text-muted);
  transition: color var(--transition);
}
.back-link:hover { color: var(--text-secondary); }
.detail-cards {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

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
  width: 140px;
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

.type-description {
  margin-left: 8px;
  font-size: 13px;
  color: var(--text-muted);
}

.text-mono {
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  font-size: 13px;
}

.config-number {
  color: var(--primary-600);
  font-weight: 600;
}

.config-list {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.config-list-item {
  display: inline-block;
  padding: 2px 8px;
  background: var(--bg-tertiary);
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

.empty-config {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 0;
  color: var(--text-muted);
  font-size: 13px;
}
</style>
