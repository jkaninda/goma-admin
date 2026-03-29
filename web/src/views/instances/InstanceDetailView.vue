<template>
  <div>
    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <div v-else-if="loadError" class="error-page">
      <div class="card">
        <div class="empty-state" style="padding: 48px 24px">
          <svg width="48" height="48" fill="none" stroke="currentColor" viewBox="0 0 24 24" style="color: var(--danger-500); margin-bottom: 12px;">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          <h3>Failed to load instance</h3>
          <p>{{ loadError }}</p>
          <div style="margin-top: 16px; display: flex; gap: 8px;">
            <router-link to="/instances" class="btn btn-secondary">Back to Instances</router-link>
            <button class="btn btn-primary" @click="loadInstance">Retry</button>
          </div>
        </div>
      </div>
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
        <div class="header-actions">
          <button v-if="!instance.builtIn" class="btn btn-secondary btn-sm" :disabled="exporting" @click="exportConfig">
            <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Export
          </button>
          <button class="btn btn-secondary btn-sm" @click="openEdit">
            <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            Edit
          </button>
          <button v-if="!instance.builtIn" class="btn btn-danger btn-sm" @click="handleDelete">
            <svg width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            Delete
          </button>
        </div>
      </div>

      <div class="detail-grid">
        <!-- Sidebar column -->
        <div class="detail-sidebar">
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

          <!-- Metadata Card -->
          <div class="card">
            <div class="card-header">
              <h2>Metadata</h2>
            </div>
            <div class="card-body">
              <dl class="detail-list">
                <div class="detail-item">
                  <dt>ID</dt>
                  <dd class="text-mono">{{ instance.id }}</dd>
                </div>
                <div class="detail-item">
                  <dt>Created</dt>
                  <dd>{{ formatDate(instance.createdAt) }}</dd>
                </div>
                <div class="detail-item">
                  <dt>Updated</dt>
                  <dd>{{ formatDate(instance.updatedAt) }}</dd>
                </div>
              </dl>
            </div>
          </div>
        </div>

        <!-- Main column -->
        <div class="detail-main">
          <!-- Metrics Card -->
          <div v-if="!instance.builtIn && instance.enableMetrics" class="card">
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
                <router-link
                  v-for="route in instance.routes"
                  :key="route.id"
                  :to="`/routes/${route.id}`"
                  class="route-item route-item-link"
                >
                  <div class="route-info">
                    <span class="route-name">{{ route.name }}</span>
                    <span class="route-path text-mono">{{ route.path || '-' }}</span>
                  </div>
                  <div class="route-meta">
                    <span
                      v-for="method in (route.methods || []).slice(0, 3)"
                      :key="method"
                      class="badge badge-info"
                    >{{ method }}</span>
                    <span :class="['badge', route.enabled !== false ? 'badge-success' : 'badge-neutral']">
                      {{ route.enabled !== false ? 'enabled' : 'disabled' }}
                    </span>
                  </div>
                </router-link>
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
                <router-link
                  v-for="mw in instance.middlewares"
                  :key="mw.id"
                  :to="`/middlewares/${mw.id}`"
                  class="route-item route-item-link"
                >
                  <div class="route-info">
                    <span class="route-name">{{ mw.name }}</span>
                  </div>
                  <span class="badge badge-info">{{ mw.type }}</span>
                </router-link>
              </div>
            </div>
          </div>

        </div>
      </div>

      <!-- Edit Modal -->
      <Modal :show="showEdit" title="Edit Instance" @close="closeEdit">
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
      </Modal>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { instancesApi, type Instance, type InstanceCreateRequest } from '@/api/instances'
import { dockerApi, type DockerStatus } from '@/api/docker'
import { metricsApi, type InstanceMetrics } from '@/api/metrics'
import { repositoriesApi, type Repository } from '@/api/repositories'
import { useNotificationStore } from '@/stores/notification'
import { useConfirm } from '@/composables/useConfirm'
import Modal from '@/components/Modal.vue'

const props = defineProps<{ id: string }>()
const numericId = computed(() => Number(props.id))
const router = useRouter()
const notify = useNotificationStore()
const { confirm } = useConfirm()

const loading = ref(true)
const loadError = ref('')
const syncing = ref(false)
const saving = ref(false)
const checking = ref(false)
const exporting = ref(false)
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

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString()
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
    notify.success('Instance updated')
    if (instance.value.enableMetrics) {
      fetchMetrics()
    } else {
      metrics.value = null
    }
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Failed to update instance')
  } finally {
    editSaving.value = false
  }
}

async function handleDelete() {
  if (!instance.value) return
  const ok = await confirm({
    title: 'Delete Instance',
    message: `Are you sure you want to delete "${instance.value.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!ok) return
  try {
    await instancesApi.delete(instance.value.id)
    notify.success('Instance deleted')
    router.push('/instances')
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Failed to delete instance')
  }
}

async function exportConfig() {
  if (!instance.value) return
  exporting.value = true
  try {
    const res = await instancesApi.exportConfig(instance.value.id)
    const blob = new Blob([res.data as string], { type: 'application/x-yaml' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `${instance.value.name}.yaml`
    a.click()
    URL.revokeObjectURL(url)
    notify.success('Config exported')
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Failed to export config')
  } finally {
    exporting.value = false
  }
}

async function syncFromRepo() {
  if (!instance.value) return
  repoSyncing.value = true
  try {
    await instancesApi.syncRepo(instance.value.id)
    const res = await instancesApi.get(instance.value.id)
    instance.value = res.data
    notify.success('Synced from repository')
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Failed to sync from repository')
  } finally {
    repoSyncing.value = false
  }
}

async function checkHealth() {
  if (!instance.value) return
  checking.value = true
  try {
    await instancesApi.checkHealth(instance.value.id)
    const res = await instancesApi.get(instance.value.id)
    instance.value = res.data
    notify.success(`Health: ${instance.value.status}`)
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Health check failed')
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
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Failed to update setting')
  } finally {
    saving.value = false
  }
}

async function triggerSync() {
  syncing.value = true
  try {
    await dockerApi.sync()
    const [statusRes, instRes] = await Promise.all([
      dockerApi.status(),
      instancesApi.get(numericId.value),
    ])
    dockerStatus.value = statusRes.data
    instance.value = instRes.data
    notify.success('Docker sync complete')
  } catch (err: any) {
    notify.error(err.response?.data?.message || 'Docker sync failed')
  } finally {
    syncing.value = false
  }
}

async function loadInstance() {
  loading.value = true
  loadError.value = ''
  try {
    const res = await instancesApi.get(numericId.value)
    instance.value = res.data

    if (instance.value?.builtIn) {
      try {
        const statusRes = await dockerApi.status()
        dockerStatus.value = statusRes.data
      } catch {
        // Docker status not available
      }
    }

    if (instance.value && !instance.value.builtIn && instance.value.enableMetrics) {
      fetchMetrics()
    }

    repositoriesApi.list().then(res => { repositories.value = res.data }).catch(() => {})
  } catch (err: any) {
    loadError.value = err.response?.data?.message || 'Could not load instance'
  } finally {
    loading.value = false
  }
}

onMounted(loadInstance)
</script>

<style scoped>
/* ─── Header ─── */
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
  flex-wrap: wrap;
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

/* ─── Grid Layout ─── */
.detail-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 24px;
  align-items: start;
}

.detail-sidebar {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-main {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

@media (max-width: 768px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

/* ─── Detail List ─── */
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
  display: flex;
  align-items: center;
}

.tags-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.toggle-hint {
  margin-left: 10px;
  font-size: 12px;
  color: var(--text-muted);
}

.last-seen-text {
  margin-left: 8px;
  font-size: 12px;
  color: var(--text-muted);
}

/* ─── Route / Middleware List ─── */
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

.route-item-link {
  text-decoration: none;
  cursor: pointer;
}

.route-item-link:hover {
  background: var(--bg-hover);
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

.route-meta {
  display: flex;
  align-items: center;
  gap: 6px;
}

/* ─── Metrics ─── */
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

/* ─── Modal Form ─── */
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

/* ─── Error Page ─── */
.error-page {
  max-width: 480px;
  margin: 48px auto;
}
</style>
