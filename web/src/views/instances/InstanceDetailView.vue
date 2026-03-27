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
        <button class="btn btn-secondary btn-sm" @click="openEdit" style="margin-left: auto;">
          <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
          Edit
        </button>
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
              <div v-if="!instance.builtIn && instance.endpoint" class="detail-item">
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

        <!-- Repository Card -->
        <div v-if="!instance.builtIn && instance.repositoryId" class="card">
          <div class="card-header">
            <h2>Git Repository</h2>
            <button class="btn btn-secondary btn-sm" :disabled="repoSyncing" @click="syncFromRepo">
              {{ repoSyncing ? 'Syncing...' : 'Sync from Repo' }}
            </button>
          </div>
          <div class="card-body">
            <dl class="detail-list">
              <div class="detail-item">
                <dt>Repository</dt>
                <dd>
                  <router-link :to="`/repositories/${instance.repositoryId}`" class="link-primary">
                    {{ repoName || `#${instance.repositoryId}` }}
                  </router-link>
                </dd>
              </div>
              <div v-if="instance.repositoryPath" class="detail-item">
                <dt>Config Path</dt>
                <dd class="text-mono">{{ instance.repositoryPath }}</dd>
              </div>
              <div class="detail-item">
                <dt>Auto-sync</dt>
                <dd>{{ instance.autoSync ? 'Enabled' : 'Disabled' }}</dd>
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

        <!-- Metrics Card (non-built-in instances with endpoint) -->
        <div v-if="!instance.builtIn && instance.enableMetrics" class="card metrics-card">
          <div class="card-header">
            <h2>Metrics</h2>
            <button class="btn btn-secondary btn-sm" :disabled="metricsLoading" @click="fetchMetrics">
              {{ metricsLoading ? 'Loading...' : 'Refresh' }}
            </button>
          </div>
          <div class="card-body">
            <div v-if="metricsLoading && !metrics" class="empty-state" style="padding: 24px">
              <div class="spinner" style="width: 24px; height: 24px;"></div>
            </div>
            <div v-else-if="metricsError" class="empty-state" style="padding: 24px">
              <p>{{ metricsError }}</p>
            </div>
            <div v-else-if="metrics">
              <!-- Summary stats -->
              <div class="metrics-summary">
                <div class="metric-stat">
                  <span class="metric-value">{{ formatNumber(metrics.totalRequests) }}</span>
                  <span class="metric-label">Requests</span>
                </div>
                <div class="metric-stat">
                  <span class="metric-value" :class="{ 'metric-danger': metrics.errorRate > 5 }">{{ metrics.errorRate }}%</span>
                  <span class="metric-label">Error Rate</span>
                </div>
                <div class="metric-stat">
                  <span class="metric-value">{{ metrics.avgLatencyMs }}ms</span>
                  <span class="metric-label">Avg Latency</span>
                </div>
                <div class="metric-stat">
                  <span class="metric-value">{{ metrics.realtimeVisitors }}</span>
                  <span class="metric-label">Active Visitors</span>
                </div>
                <div class="metric-stat">
                  <span class="metric-value">{{ formatUptime(metrics.uptimeSeconds) }}</span>
                  <span class="metric-label">Uptime</span>
                </div>
              </div>

              <!-- Route metrics table -->
              <div v-if="metrics.routeMetrics?.length" class="metrics-table-wrap">
                <table class="metrics-table">
                  <thead>
                    <tr>
                      <th>Route</th>
                      <th class="text-right">Requests</th>
                      <th class="text-right">Errors</th>
                      <th class="text-right">Error Rate</th>
                      <th class="text-right">Avg Latency</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="rm in metrics.routeMetrics" :key="rm.routeName">
                      <td class="text-mono">{{ rm.routeName }}</td>
                      <td class="text-right">{{ formatNumber(rm.totalRequests) }}</td>
                      <td class="text-right">{{ formatNumber(rm.errorCount) }}</td>
                      <td class="text-right" :class="{ 'metric-danger': rm.errorRate > 5 }">{{ rm.errorRate }}%</td>
                      <td class="text-right">{{ rm.avgLatencyMs }}ms</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
            <div v-else class="empty-state" style="padding: 24px">
              <p>No metrics available</p>
            </div>
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

      <!-- Edit Modal -->
      <Teleport to="body">
        <div v-if="showEdit" class="modal-overlay" @click.self="closeEdit">
          <div class="modal">
            <div class="modal-header">
              <h3>Edit Instance</h3>
              <button class="btn-ghost btn-icon" @click="closeEdit">
                <svg width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <form @submit.prevent="handleEditSubmit">
              <div class="modal-body">
                <div class="form-group">
                  <label class="form-label" for="edit-name">Name</label>
                  <input id="edit-name" v-model="editForm.name" class="form-input" required placeholder="my-gateway" :disabled="instance.builtIn" />
                </div>
                <div class="form-grid">
                  <div class="form-group">
                    <label class="form-label" for="edit-env">Environment</label>
                    <select id="edit-env" v-model="editForm.environment" class="form-select">
                      <option value="development">Development</option>
                      <option value="staging">Staging</option>
                      <option value="production">Production</option>
                      <option value="testing">Testing</option>
                    </select>
                  </div>
                  <div class="form-group">
                    <label class="form-label" for="edit-region">Region</label>
                    <input id="edit-region" v-model="editForm.region" class="form-input" placeholder="us-east-1" />
                  </div>
                </div>
                <div class="form-group">
                  <label class="form-label" for="edit-endpoint">Endpoint</label>
                  <input id="edit-endpoint" v-model="editForm.endpoint" class="form-input" required placeholder="https://gateway.example.com" />
                </div>
                <div class="form-group">
                  <label class="form-label" for="edit-desc">Description</label>
                  <input id="edit-desc" v-model="editForm.description" class="form-input" placeholder="Optional description" />
                </div>

                <!-- Metrics Section -->
                <div class="form-section-divider">
                  <span class="form-section-title">Metrics</span>
                </div>
                <div class="form-group">
                  <label class="form-label form-label-inline" for="edit-enable-metrics">
                    <input id="edit-enable-metrics" type="checkbox" v-model="editForm.enableMetrics" class="form-checkbox" />
                    Enable Metrics
                  </label>
                  <span class="form-hint">Fetch Prometheus metrics from this instance</span>
                </div>
                <template v-if="editForm.enableMetrics">
                  <div class="form-group">
                    <label class="form-label" for="edit-metrics-endpoint">Metrics Endpoint</label>
                    <input id="edit-metrics-endpoint" v-model="editForm.metricsEndpoint" class="form-input" placeholder="Defaults to endpoint/metrics" />
                    <span class="form-hint">Leave empty to use {endpoint}/metrics</span>
                  </div>
                  <div class="form-grid">
                    <div class="form-group">
                      <label class="form-label" for="edit-metrics-auth-type">Auth Type</label>
                      <select id="edit-metrics-auth-type" v-model="editForm.metricsAuthType" class="form-select">
                        <option value="">None</option>
                        <option value="basic">Basic Auth</option>
                        <option value="bearer">Bearer Token</option>
                        <option value="header">Custom Header</option>
                      </select>
                    </div>
                    <div v-if="editForm.metricsAuthType" class="form-group">
                      <label class="form-label" for="edit-metrics-auth-value">Auth Value</label>
                      <input
                        id="edit-metrics-auth-value"
                        v-model="editForm.metricsAuthValue"
                        class="form-input"
                        type="password"
                        :placeholder="instance?.hasMetricsAuth ? '••••••• (leave empty to keep)' : (editForm.metricsAuthType === 'basic' ? 'user:password' : 'token')"
                      />
                    </div>
                  </div>
                </template>

                <!-- Repository Section -->
                <div class="form-section-divider">
                  <span class="form-section-title">Git Repository</span>
                </div>
                <div class="form-group">
                  <label class="form-label" for="edit-repo">Repository</label>
                  <select id="edit-repo" v-model="editForm.repositoryId" class="form-select">
                    <option :value="null">None</option>
                    <option v-for="r in repositories" :key="r.id" :value="r.id">{{ r.name }} ({{ r.branch }})</option>
                  </select>
                </div>
                <template v-if="editForm.repositoryId">
                  <div class="form-group">
                    <label class="form-label" for="edit-repo-path">Config Path</label>
                    <input id="edit-repo-path" v-model="editForm.repositoryPath" class="form-input" placeholder="e.g. production/gateway-1" />
                    <span class="form-hint">Path within the repository where YAML configs are stored</span>
                  </div>
                  <div class="form-group">
                    <label class="form-label form-label-inline" for="edit-auto-sync">
                      <input id="edit-auto-sync" type="checkbox" v-model="editForm.autoSync" class="form-checkbox" />
                      Auto-sync on push
                    </label>
                  </div>
                </template>
              </div>
              <div class="modal-footer">
                <button type="button" class="btn btn-secondary" @click="closeEdit">Cancel</button>
                <button type="submit" class="btn btn-primary" :disabled="editSaving">
                  {{ editSaving ? 'Saving...' : 'Update' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </Teleport>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { instancesApi, type Instance, type InstanceCreateRequest } from '@/api/instances'
import { dockerApi, type DockerStatus } from '@/api/docker'
import { metricsApi, type InstanceMetrics } from '@/api/metrics'
import { repositoriesApi, type Repository } from '@/api/repositories'

const props = defineProps<{ id: string }>()

const loading = ref(true)
const syncing = ref(false)
const saving = ref(false)
const checking = ref(false)
const instance = ref<Instance | null>(null)
const dockerStatus = ref<DockerStatus | null>(null)
const metrics = ref<InstanceMetrics | null>(null)
const metricsLoading = ref(false)
const metricsError = ref('')
const showEdit = ref(false)
const editSaving = ref(false)
const repoSyncing = ref(false)
const repositories = ref<Repository[]>([])

const repoName = computed(() => {
  if (!instance.value?.repositoryId) return ''
  const r = repositories.value.find(r => r.id === instance.value?.repositoryId)
  return r?.name || ''
})
const editForm = reactive<InstanceCreateRequest>({
  name: '',
  environment: 'development',
  description: '',
  endpoint: '',
  region: '',
  enableMetrics: false,
  metricsEndpoint: '',
  metricsAuthType: '',
  metricsAuthValue: '',
  repositoryId: null as number | null,
  repositoryPath: '',
  autoSync: false,
})

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

function formatNumber(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K'
  return String(n)
}

function formatUptime(seconds: number): string {
  if (seconds < 60) return Math.floor(seconds) + 's'
  if (seconds < 3600) return Math.floor(seconds / 60) + 'm'
  if (seconds < 86400) return Math.floor(seconds / 3600) + 'h'
  return Math.floor(seconds / 86400) + 'd'
}

async function fetchMetrics() {
  if (!instance.value) return
  metricsLoading.value = true
  metricsError.value = ''
  try {
    const res = await metricsApi.getMetrics(instance.value.id)
    metrics.value = res.data
  } catch (err: any) {
    metricsError.value = err.response?.data?.message || 'Failed to fetch metrics'
    metrics.value = null
  } finally {
    metricsLoading.value = false
  }
}

function openEdit() {
  if (!instance.value) return
  const inst = instance.value
  editForm.name = inst.name
  editForm.environment = inst.environment
  editForm.description = inst.description || ''
  editForm.endpoint = inst.endpoint
  editForm.region = inst.region || ''
  editForm.enableMetrics = inst.enableMetrics
  editForm.metricsEndpoint = inst.metricsEndpoint || ''
  editForm.metricsAuthType = inst.metricsAuthType || ''
  editForm.metricsAuthValue = ''
  editForm.repositoryId = inst.repositoryId || null
  editForm.repositoryPath = inst.repositoryPath || ''
  editForm.autoSync = inst.autoSync
  showEdit.value = true
}

function closeEdit() {
  showEdit.value = false
}

async function handleEditSubmit() {
  if (!instance.value) return
  editSaving.value = true
  try {
    const res = await instancesApi.update(instance.value.id, editForm)
    instance.value = res.data
    closeEdit()
    // Refresh metrics if enabled
    if (instance.value.enableMetrics) {
      fetchMetrics()
    } else {
      metrics.value = null
    }
  } catch {
    // handle error
  } finally {
    editSaving.value = false
  }
}

async function syncFromRepo() {
  if (!instance.value) return
  repoSyncing.value = true
  try {
    await instancesApi.syncRepo(instance.value.id)
    const res = await instancesApi.get(instance.value.id)
    instance.value = res.data
  } catch {
    // handle error
  } finally {
    repoSyncing.value = false
  }
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

    // Fetch metrics for non-built-in instances with metrics enabled
    if (instance.value && !instance.value.builtIn && instance.value.enableMetrics) {
      fetchMetrics()
    }

    // Fetch repositories for the edit form
    repositoriesApi.list().then(res => { repositories.value = res.data }).catch(() => {})
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

/* ─── Metrics ─── */
.metrics-card {
  grid-column: 1 / -1;
}

.metrics-summary {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.metric-stat {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 100px;
}

.metric-value {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.metric-label {
  font-size: 12px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.metric-danger {
  color: var(--color-danger, #ef4444);
}

.metrics-table-wrap {
  overflow-x: auto;
}

.metrics-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.metrics-table th {
  text-align: left;
  padding: 8px 12px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  border-bottom: 1px solid var(--border-primary);
}

.metrics-table td {
  padding: 10px 12px;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-primary);
}

.metrics-table tbody tr:last-child td {
  border-bottom: none;
}

.metrics-table .text-right {
  text-align: right;
}

.metrics-table .text-mono {
  font-family: var(--font-mono, monospace);
  font-size: 12px;
}

/* ─── Edit Modal ─── */
.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(2px);
}

.modal {
  background: var(--bg-primary);
  border-radius: var(--radius-lg, 12px);
  box-shadow: var(--shadow-lg);
  width: 100%;
  max-width: 520px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px 0;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.modal-body {
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 24px;
  border-top: 1px solid var(--border-primary);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-section-divider {
  padding-top: 8px;
  border-top: 1px solid var(--border-secondary);
}

.form-section-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.form-label-inline {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.form-checkbox {
  width: 16px;
  height: 16px;
  accent-color: var(--primary-600);
  cursor: pointer;
}

.form-hint {
  display: block;
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
}
</style>
