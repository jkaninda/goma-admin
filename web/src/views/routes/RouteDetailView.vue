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
        <div class="card">
          <div class="card-header"><h2>Configuration</h2></div>
          <div class="card-body">
            <pre class="code-block">{{ formatConfig(route.config) }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { routesApi, type Route } from '@/api/routes'

const props = defineProps<{ id: string }>()

const loading = ref(true)
const route = ref<Route | null>(null)

function formatConfig(config: Record<string, unknown>): string {
  if (!config) return ''
  const lines: string[] = []
  for (const [key, value] of Object.entries(config)) {
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
  max-width: 800px;
}
</style>
