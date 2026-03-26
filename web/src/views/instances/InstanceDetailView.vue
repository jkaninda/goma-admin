<template>
  <div>
    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <div v-else-if="instance">
      <div class="detail-header">
        <router-link to="/instances" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <h1>{{ instance.name }}</h1>
        <span v-if="instance.builtIn" class="badge badge-info">Built-in</span>
        <span :class="['badge', statusBadge(instance.status)]">{{ instance.status }}</span>
      </div>

      <div class="detail-grid">
        <!-- Details Card -->
        <div class="card">
          <div class="card-header">
            <h2>Details</h2>
          </div>
          <div class="card-body">
            <dl class="detail-list">
              <div class="detail-item">
                <dt>Environment</dt>
                <dd>
                  <span :class="['badge', envBadge(instance.environment)]">{{ instance.environment }}</span>
                </dd>
              </div>
              <div class="detail-item">
                <dt>Health Status</dt>
                <dd>
                  <span :class="['badge', statusBadge(instance.status)]">{{ instance.status }}</span>
                  <span v-if="instance.lastSeen" class="last-seen-text">{{ relativeTime(instance.lastSeen) }}</span>
                </dd>
              </div>
              <div v-if="!instance.builtIn && instance.healthEndpoint" class="detail-item">
                <dt>Health Check</dt>
                <dd>
                  <button class="btn btn-secondary btn-sm" :disabled="checking" @click="checkHealth">
                    {{ checking ? 'Checking...' : 'Check Health' }}
                  </button>
                </dd>
              </div>
              <div class="detail-item">
                <dt>Endpoint</dt>
                <dd class="text-mono">{{ instance.endpoint }}</dd>
              </div>
              <div v-if="instance.region" class="detail-item">
                <dt>Region</dt>
                <dd>{{ instance.region }}</dd>
              </div>
              <div v-if="instance.version" class="detail-item">
                <dt>Version</dt>
                <dd>{{ instance.version }}</dd>
              </div>
              <div v-if="instance.description" class="detail-item">
                <dt>Description</dt>
                <dd>{{ instance.description }}</dd>
              </div>
              <div v-if="instance.tags?.length" class="detail-item">
                <dt>Tags</dt>
                <dd class="tags-list">
                  <span v-for="tag in instance.tags" :key="tag" class="badge badge-info">{{ tag }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Provider Settings Card -->
        <div v-if="!instance.builtIn" class="card">
          <div class="card-header">
            <h2>Provider Settings</h2>
          </div>
          <div class="card-body">
            <dl class="detail-list">
              <div class="detail-item">
                <dt>Write Config to Disk</dt>
                <dd>
                  <button
                    :class="['toggle-btn', { active: instance.writeConfig }]"
                    @click="toggleField('writeConfig', !instance.writeConfig)"
                    :disabled="saving"
                  >
                    <span class="toggle-slider"></span>
                  </button>
                  <span class="toggle-hint">{{ instance.writeConfig ? 'Enabled' : 'Disabled' }}</span>
                </dd>
              </div>
              <div class="detail-item">
                <dt>Include Docker Routes</dt>
                <dd>
                  <button
                    :class="['toggle-btn', { active: instance.includeDockerRoutes }]"
                    @click="toggleField('includeDockerRoutes', !instance.includeDockerRoutes)"
                    :disabled="saving"
                  >
                    <span class="toggle-slider"></span>
                  </button>
                  <span class="toggle-hint">{{ instance.includeDockerRoutes ? 'Writes docker-provider.yaml' : 'Disabled' }}</span>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Docker Status Card (only for built-in Docker provider) -->
        <div v-if="instance.builtIn && dockerStatus" class="card">
          <div class="card-header">
            <h2>Docker Provider</h2>
            <button class="btn btn-secondary btn-sm" :disabled="syncing" @click="triggerSync">
              {{ syncing ? 'Syncing...' : 'Sync Now' }}
            </button>
          </div>
          <div class="card-body">
            <dl class="detail-list">
              <div class="detail-item">
                <dt>Connection</dt>
                <dd>
                  <span :class="['badge', dockerStatus.connected ? 'badge-success' : 'badge-danger']">
                    {{ dockerStatus.connected ? 'Connected' : 'Disconnected' }}
                  </span>
                </dd>
              </div>
              <div class="detail-item">
                <dt>Swarm Mode</dt>
                <dd>{{ dockerStatus.swarmMode ? 'Yes' : 'No' }}</dd>
              </div>
              <div class="detail-item">
                <dt>Discovered Routes</dt>
                <dd>{{ dockerStatus.routeCount }}</dd>
              </div>
              <div v-if="dockerStatus.lastSync" class="detail-item">
                <dt>Last Sync</dt>
                <dd>{{ new Date(dockerStatus.lastSync).toLocaleString() }}</dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Config Cards -->
        <div class="config-stack">
          <!-- Routes Card -->
          <div class="card">
            <div class="card-header">
              <h2>Routes ({{ instance.routes?.length || 0 }})</h2>
            </div>
            <div class="card-body">
              <div v-if="!instance.routes?.length" class="empty-state" style="padding: 24px">
                <p>No routes attached</p>
              </div>
              <div v-else class="route-list">
                <div v-for="route in instance.routes" :key="route.id" class="route-item">
                  <div class="route-info">
                    <span class="route-name">{{ route.name }}</span>
                    <span class="route-path text-mono">{{ getRouteConfig(route, 'path') }}</span>
                  </div>
                  <div class="route-meta">
                    <span
                      v-for="method in getRouteMethods(route)"
                      :key="method"
                      class="badge badge-info"
                    >{{ method }}</span>
                    <span :class="['badge', isRouteEnabled(route) ? 'badge-success' : 'badge-neutral']">
                      {{ isRouteEnabled(route) ? 'enabled' : 'disabled' }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Middlewares Card -->
          <div class="card">
            <div class="card-header">
              <h2>Middlewares ({{ instance.middlewares?.length || 0 }})</h2>
            </div>
            <div class="card-body">
              <div v-if="!instance.middlewares?.length" class="empty-state" style="padding: 24px">
                <p>No middlewares attached</p>
              </div>
              <div v-else class="route-list">
                <div v-for="mw in instance.middlewares" :key="mw.id" class="route-item">
                  <div class="route-info">
                    <span class="route-name">{{ mw.name }}</span>
                  </div>
                  <span class="badge badge-info">{{ mw.type }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { instancesApi, type Instance } from '@/api/instances'
import { dockerApi, type DockerStatus } from '@/api/docker'

const props = defineProps<{ id: string }>()

const loading = ref(true)
const syncing = ref(false)
const saving = ref(false)
const checking = ref(false)
const instance = ref<Instance | null>(null)
const dockerStatus = ref<DockerStatus | null>(null)

function statusBadge(status: string): string {
  const map: Record<string, string> = {
    active: 'badge-success',
    inactive: 'badge-neutral',
    unhealthy: 'badge-danger',
    unknown: 'badge-warning',
  }
  return map[status] || 'badge-neutral'
}

function envBadge(env: string): string {
  const map: Record<string, string> = {
    production: 'badge-danger',
    staging: 'badge-warning',
    development: 'badge-info',
    testing: 'badge-neutral',
  }
  return map[env] || 'badge-neutral'
}

function getRouteConfig(route: any, key: string): string {
  return route.config?.[key] || route[key] || '-'
}

function getRouteMethods(route: any): string[] {
  const methods = route.config?.methods || route.methods || []
  return Array.isArray(methods) ? methods.slice(0, 3) : []
}

function isRouteEnabled(route: any): boolean {
  return route.config?.enabled ?? route.enabled ?? true
}

function relativeTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  if (diffSec < 60) return `${diffSec}s ago`
  const diffMin = Math.floor(diffSec / 60)
  if (diffMin < 60) return `${diffMin}m ago`
  const diffHr = Math.floor(diffMin / 60)
  if (diffHr < 24) return `${diffHr}h ago`
  const diffDay = Math.floor(diffHr / 24)
  return `${diffDay}d ago`
}

async function checkHealth() {
  if (!instance.value) return
  checking.value = true
  try {
    await instancesApi.checkHealth(instance.value.id)
    // Refresh instance data to get updated status and lastSeen
    const res = await instancesApi.get(instance.value.id)
    instance.value = res.data
  } catch {
    // handle error
  } finally {
    checking.value = false
  }
}

async function toggleField(field: 'writeConfig' | 'includeDockerRoutes', value: boolean) {
  if (!instance.value) return
  saving.value = true
  try {
    const res = await instancesApi.patch(instance.value.id, { [field]: value })
    instance.value = res.data
  } catch {
    // handle error
  } finally {
    saving.value = false
  }
}

async function triggerSync() {
  syncing.value = true
  try {
    await dockerApi.sync()
    // Refresh both docker status and instance data
    const [statusRes, instRes] = await Promise.all([
      dockerApi.status(),
      instancesApi.get(props.id),
    ])
    dockerStatus.value = statusRes.data
    instance.value = instRes.data
  } catch {
    // handle error
  } finally {
    syncing.value = false
  }
}

onMounted(async () => {
  try {
    const res = await instancesApi.get(props.id)
    instance.value = res.data

    // Fetch Docker status if this is a built-in instance
    if (instance.value?.builtIn) {
      try {
        const statusRes = await dockerApi.status()
        dockerStatus.value = statusRes.data
      } catch {
        // Docker status not available
      }
    }
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

.back-link:hover {
  color: var(--text-secondary);
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 24px;
}

@media (max-width: 768px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

.detail-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-item dt {
  font-size: 12px;
  color: var(--text-muted);
  margin-bottom: 2px;
}

.detail-item dd {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.route-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.route-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  border-radius: var(--radius);
  background: var(--bg-tertiary);
  transition: background var(--transition);
}

.route-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.route-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.route-path {
  font-size: 12px;
  color: var(--text-muted);
}

.config-stack {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.route-meta {
  display: flex;
  align-items: center;
  gap: 6px;
}

.toggle-hint {
  margin-left: 10px;
  font-size: 12px;
  color: var(--text-muted);
}

.detail-item dd {
  display: flex;
  align-items: center;
}

.last-seen-text {
  margin-left: 8px;
  font-size: 12px;
  color: var(--text-muted);
}
</style>
