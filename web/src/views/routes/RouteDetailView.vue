<template>
  <div>
    <div v-if="loading" class="loading-page"><div class="spinner"></div></div>

    <div v-else-if="route">
      <div class="detail-header">
        <router-link to="/routes" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <h1>{{ route.name }}</h1>
        <span :class="['badge', route.config?.enabled !== false ? 'badge-success' : 'badge-neutral']">
          {{ route.config?.enabled !== false ? 'enabled' : 'disabled' }}
        </span>
      </div>

      <div class="detail-cards">
        <!-- Core Configuration -->
        <div class="card">
          <div class="card-header"><h2>Route Configuration</h2></div>
          <div class="card-body">
            <dl class="detail-grid">
              <div class="detail-row">
                <dt>Path</dt>
                <dd class="text-mono">{{ config.path || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Target</dt>
                <dd class="text-mono">{{ config.target || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Hosts</dt>
                <dd>
                  <template v-if="hostsList.length">
                    <a v-for="h in hostsList" :key="h" :href="`http://${h}`" target="_blank" class="badge badge-info detail-badge">{{ h }}</a>
                  </template>
                  <span v-else class="text-muted">None</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Methods</dt>
                <dd>
                  <template v-if="methodsList.length">
                    <span v-for="m in methodsList" :key="m" class="badge badge-info detail-badge">{{ m }}</span>
                  </template>
                  <span v-else class="text-muted">All</span>
                </dd>
              </div>
              <div class="detail-row">
                <dt>Rewrite</dt>
                <dd class="text-mono">{{ config.rewrite || '-' }}</dd>
              </div>
              <div class="detail-row">
                <dt>Status</dt>
                <dd>
                  <span :class="['badge', config.enabled !== false ? 'badge-success' : 'badge-neutral']">
                    {{ config.enabled !== false ? 'Enabled' : 'Disabled' }}
                  </span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Advanced Configuration (only if extra fields exist) -->
        <div v-if="hasAdvancedFields" class="card">
          <div class="card-header"><h2>Advanced Configuration</h2></div>
          <div class="card-body">
            <pre class="code-block">{{ formatAdvancedConfig() }}</pre>
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
import { routesApi, type Route } from '@/api/routes'

const props = defineProps<{ id: string }>()

const loading = ref(true)
const route = ref<Route | null>(null)

const simpleFieldKeys = new Set(['path', 'target', 'hosts', 'methods', 'rewrite', 'enabled'])

const config = computed(() => route.value?.config || {})

const hostsList = computed(() => {
  const hosts = config.value.hosts
  return Array.isArray(hosts) ? hosts.map(String) : []
})

const methodsList = computed(() => {
  const methods = config.value.methods
  return Array.isArray(methods) ? methods.map(String) : []
})

const advancedFields = computed(() => {
  const result: Record<string, unknown> = {}
  for (const [key, value] of Object.entries(config.value)) {
    if (!simpleFieldKeys.has(key)) {
      result[key] = value
    }
  }
  return result
})

const hasAdvancedFields = computed(() => Object.keys(advancedFields.value).length > 0)

function formatAdvancedConfig(): string {
  const obj = advancedFields.value
  const lines: string[] = []
  for (const [key, value] of Object.entries(obj)) {
    if (value === undefined || value === null) continue
    if (Array.isArray(value)) {
      lines.push(`${key}:`)
      for (const item of value) {
        if (typeof item === 'object') {
          lines.push(`  - ${JSON.stringify(item)}`)
        } else {
          lines.push(`  - ${item}`)
        }
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
    const res = await routesApi.get(Number(props.id))
    route.value = res.data
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
</style>
