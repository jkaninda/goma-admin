<template>
  <div>
    <div class="page-header">
      <div>
        <h1>Dashboard</h1>
        <p class="page-subtitle">{{ instanceStore.contextLabel }}</p>
      </div>
    </div>

    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <div class="stats-grid">
        <router-link to="/instances" class="stat-card stat-card-link">
          <div class="stat-header">
            <span class="stat-label">Instances</span>
            <span class="stat-icon stat-icon-primary">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
              </svg>
            </span>
          </div>
          <span class="stat-value">{{ data.instances }}</span>
        </router-link>

        <router-link to="/routes" class="stat-card stat-card-link">
          <div class="stat-header">
            <span class="stat-label">Routes</span>
            <span class="stat-icon stat-icon-success">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </span>
          </div>
          <span class="stat-value">{{ data.routes }}</span>
        </router-link>

        <router-link to="/middlewares" class="stat-card stat-card-link">
          <div class="stat-header">
            <span class="stat-label">Middlewares</span>
            <span class="stat-icon stat-icon-warning">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
              </svg>
            </span>
          </div>
          <span class="stat-value">{{ data.middlewares }}</span>
        </router-link>

        <router-link to="/users" class="stat-card stat-card-link">
          <div class="stat-header">
            <span class="stat-label">Users</span>
            <span class="stat-icon stat-icon-danger">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </span>
          </div>
          <span class="stat-value">{{ data.users }}</span>
        </router-link>
      
      </div>

      <div class="dashboard-panels">
        <!-- Docker Provider Card -->
        <div v-if="dockerStatus" class="card docker-card">
          <div class="card-header">
            <div class="docker-header-left">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z" />
              </svg>
              <h2>Docker Provider</h2>
            </div>
            <div class="docker-header-right">
              <span :class="['badge', dockerStatus.connected ? 'badge-success' : 'badge-danger']">
                {{ dockerStatus.connected ? 'Connected' : 'Disconnected' }}
              </span>
              <button class="btn btn-secondary btn-sm" :disabled="syncing" @click="triggerSync">
                <svg v-if="syncing" width="14" height="14" class="spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <svg v-else width="14" height="14" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                {{ syncing ? 'Syncing...' : 'Sync Now' }}
              </button>
            </div>
          </div>
          <div class="card-body">
            <div class="docker-stats">
              <div class="docker-stat">
                <span class="docker-stat-value">{{ dockerStatus.routeCount }}</span>
                <span class="docker-stat-label">Discovered Routes</span>
              </div>
              <div class="docker-stat">
                <span class="docker-stat-value">{{ dockerStatus.swarmMode ? 'Swarm' : 'Standalone' }}</span>
                <span class="docker-stat-label">Mode</span>
              </div>
              <div class="docker-stat">
                <span class="docker-stat-value">{{ formatLastSync(dockerStatus.lastSync) }}</span>
                <span class="docker-stat-label">Last Sync</span>
              </div>
            </div>
            <div v-if="dockerEvents.length" class="docker-events">
              <div class="docker-events-header">
                <span class="docker-events-title">Live Activity</span>
                <span :class="['sse-indicator', sseConnected ? 'sse-connected' : 'sse-disconnected']"></span>
              </div>
              <div class="docker-event-list">
                <div v-for="(evt, i) in dockerEvents" :key="i" class="docker-event-item">
                  <span :class="['docker-event-dot', `event-${evt.type}`]"></span>
                  <span class="docker-event-msg">{{ evt.message }}</span>
                  <span class="docker-event-time">{{ formatEventTime(evt.timestamp) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity (Config SSE) -->
        <div v-if="configEvents.length" class="card activity-card">
          <div class="card-header">
            <div class="docker-header-left">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <h2>Recent Activity</h2>
            </div>
            <span :class="['sse-indicator', configSseConnected ? 'sse-connected' : 'sse-disconnected']"></span>
          </div>
          <div class="card-body">
            <div class="docker-event-list">
              <div v-for="(evt, i) in configEvents" :key="i" class="docker-event-item">
                <span :class="['docker-event-dot', `event-${evt.resource}`]"></span>
                <span class="docker-event-msg">{{ evt.message }}</span>
                <span class="docker-event-time">{{ formatEventTime(evt.timestamp) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Instances Health -->
        <div v-if="instances.length" class="card">
          <div class="card-header">
            <h2>Instances</h2>
            <router-link to="/instances" class="btn btn-secondary btn-sm">View All</router-link>
          </div>
          <div class="card-body instance-list-body">
            <div v-for="inst in instances" :key="inst.id" class="instance-row">
              <div class="instance-row-left">
                <span :class="['instance-dot', `dot-${inst.status}`]"></span>
                <router-link :to="`/instances/${inst.id}`" class="instance-name">{{ inst.name }}</router-link>
                <span v-if="inst.builtIn" class="badge badge-info">Built-in</span>
              </div>
              <div class="instance-row-right">
                <span class="instance-meta">{{ inst.routes?.length || 0 }} routes</span>
                <span :class="['badge', envBadge(inst.environment)]">{{ inst.environment }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Gateway Metrics -->
      <div v-if="metricsInstance" class="card gateway-metrics-card">
        <div class="card-header">
          <div class="docker-header-left">
            <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <h2>Gateway Metrics</h2>
          </div>
          <div class="docker-header-right">
            <span class="instance-meta">{{ metricsInstance.name }}</span>
            <button class="btn btn-secondary btn-sm" :disabled="metricsLoading" @click="fetchDashboardMetrics">
              {{ metricsLoading ? 'Loading...' : 'Refresh' }}
            </button>
          </div>
        </div>
        <div class="card-body">
          <div v-if="metricsLoading && !dashboardMetrics" class="empty-state" style="padding: 24px; text-align: center;">
            <div class="spinner" style="width: 24px; height: 24px; margin: 0 auto;"></div>
          </div>
          <div v-else-if="metricsError" class="empty-state" style="padding: 24px; text-align: center;">
            <p style="color: var(--text-muted); font-size: 13px;">{{ metricsError }}</p>
          </div>
          <div v-else-if="dashboardMetrics">
            <div class="gw-metrics-summary">
              <div class="gw-metric-stat">
                <span class="gw-metric-value">{{ formatMetricNumber(dashboardMetrics.totalRequests) }}</span>
                <span class="gw-metric-label">Total Requests</span>
              </div>
              <div class="gw-metric-stat">
                <span class="gw-metric-value" :class="{ 'gw-metric-danger': dashboardMetrics.errorRate > 5 }">{{ dashboardMetrics.errorRate }}%</span>
                <span class="gw-metric-label">Error Rate</span>
              </div>
              <div class="gw-metric-stat">
                <span class="gw-metric-value">{{ dashboardMetrics.avgLatencyMs }}ms</span>
                <span class="gw-metric-label">Avg Latency</span>
              </div>
              <div class="gw-metric-stat">
                <span class="gw-metric-value">{{ dashboardMetrics.realtimeVisitors }}</span>
                <span class="gw-metric-label">Active Visitors</span>
              </div>
              <div class="gw-metric-stat">
                <span class="gw-metric-value">{{ dashboardMetrics.routesCount }}</span>
                <span class="gw-metric-label">Routes</span>
              </div>
              <div class="gw-metric-stat">
                <span class="gw-metric-value">{{ formatMetricUptime(dashboardMetrics.uptimeSeconds) }}</span>
                <span class="gw-metric-label">Uptime</span>
              </div>
            </div>
            <div v-if="dashboardMetrics.routeMetrics?.length" class="gw-route-table-wrap">
              <div class="gw-route-table-header">Top Routes</div>
              <table class="gw-route-table">
                <thead>
                  <tr>
                    <th>Route</th>
                    <th class="text-right">Requests</th>
                    <th class="text-right">Error Rate</th>
                    <th class="text-right">Avg Latency</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="rm in dashboardMetrics.routeMetrics.slice(0, 10)" :key="rm.routeName">
                    <td class="text-mono">{{ rm.routeName }}</td>
                    <td class="text-right">{{ formatMetricNumber(rm.totalRequests) }}</td>
                    <td class="text-right" :class="{ 'gw-metric-danger': rm.errorRate > 5 }">{{ rm.errorRate }}%</td>
                    <td class="text-right">{{ rm.avgLatencyMs }}ms</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- Quick actions -->
      <div class="card quick-actions-card">
        <div class="card-header">
          <h2>Quick Actions</h2>
        </div>
        <div class="card-body">
          <div class="quick-actions">
            <router-link to="/routes" class="quick-action">
              <span class="quick-action-icon stat-icon-success">
                <svg width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
              </span>
              <div>
                <div class="quick-action-title">Manage Routes</div>
                <div class="quick-action-desc">Create, edit, or import routes</div>
              </div>
            </router-link>
            <router-link to="/middlewares" class="quick-action">
              <span class="quick-action-icon stat-icon-warning">
                <svg width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
              </span>
              <div>
                <div class="quick-action-title">Manage Middlewares</div>
                <div class="quick-action-desc">Configure rate limits, auth, and more</div>
              </div>
            </router-link>
            <router-link to="/settings" class="quick-action">
              <span class="quick-action-icon stat-icon-primary">
                <svg width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </span>
              <div>
                <div class="quick-action-title">Import / Export</div>
                <div class="quick-action-desc">Transfer config between instances</div>
              </div>
            </router-link>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { dashboardApi, type DashboardStats } from '@/api/dashboard'
import { dockerApi, type DockerStatus, type DockerEvent } from '@/api/docker'
import { instancesApi, type Instance } from '@/api/instances'
import { metricsApi, type InstanceMetrics } from '@/api/metrics'
import { connectConfigSSE, type ConfigEvent } from '@/api/events'
import { useInstanceStore } from '@/stores/instance'

const MAX_EVENTS = 5

const instanceStore = useInstanceStore()
const loading = ref(true)
const syncing = ref(false)
const data = ref<DashboardStats>({ users: 0, instances: 0, middlewares: 0, routes: 0 })
const dockerStatus = ref<DockerStatus | null>(null)
const instances = ref<Instance[]>([])
const dockerEvents = ref<DockerEvent[]>([])
const sseConnected = ref(false)
const configEvents = ref<ConfigEvent[]>([])
const configSseConnected = ref(false)
let eventSource: EventSource | null = null
let configEventSource: EventSource | null = null
const dashboardMetrics = ref<InstanceMetrics | null>(null)
const metricsLoading = ref(false)
const metricsError = ref('')
const metricsInstance = ref<Instance | null>(null)

function envBadge(env: string): string {
  const map: Record<string, string> = {
    production: 'badge-danger',
    staging: 'badge-warning',
    development: 'badge-info',
    testing: 'badge-neutral',
  }
  return map[env] || 'badge-neutral'
}

function formatLastSync(lastSync: string): string {
  if (!lastSync || lastSync === '0001-01-01T00:00:00Z') return 'Never'
  const date = new Date(lastSync)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  if (diffSec < 60) return `${diffSec}s ago`
  const diffMin = Math.floor(diffSec / 60)
  if (diffMin < 60) return `${diffMin}m ago`
  const diffHr = Math.floor(diffMin / 60)
  if (diffHr < 24) return `${diffHr}h ago`
  return date.toLocaleDateString()
}

function formatEventTime(ts: string): string {
  if (!ts) return ''
  return new Date(ts).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

function formatMetricNumber(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1) + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1) + 'K'
  return String(n)
}

function formatMetricUptime(seconds: number): string {
  if (seconds < 60) return Math.floor(seconds) + 's'
  if (seconds < 3600) return Math.floor(seconds / 60) + 'm'
  if (seconds < 86400) return Math.floor(seconds / 3600) + 'h'
  return Math.floor(seconds / 86400) + 'd'
}

async function fetchDashboardMetrics() {
  if (!metricsInstance.value) return
  metricsLoading.value = true
  metricsError.value = ''
  try {
    const res = await metricsApi.getMetrics(metricsInstance.value.id)
    dashboardMetrics.value = res.data
  } catch (err: any) {
    metricsError.value = err.response?.data?.message || 'Failed to fetch metrics'
    dashboardMetrics.value = null
  } finally {
    metricsLoading.value = false
  }
}

function connectSSE() {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }

  if (!dockerStatus.value?.enabled) return

  eventSource = dockerApi.events(
    (evt: DockerEvent) => {
      dockerEvents.value = [evt, ...dockerEvents.value].slice(0, MAX_EVENTS)

      // Update route count live when routes change
      if (evt.type === 'routes_changed' || evt.type === 'sync_completed') {
        if (evt.routeCount !== undefined && dockerStatus.value) {
          dockerStatus.value = { ...dockerStatus.value, routeCount: evt.routeCount, lastSync: evt.timestamp }
        }
        // Refresh dashboard stats when routes change
        if (evt.type === 'routes_changed') {
          dashboardApi.getStats().then(res => { data.value = res.data }).catch(() => {})
        }
      }
    },
    (status: DockerStatus) => {
      dockerStatus.value = status
    },
  )

  eventSource.onopen = () => { sseConnected.value = true }
  eventSource.onerror = () => { sseConnected.value = false }
}

function connectConfigEvents() {
  if (configEventSource) {
    configEventSource.close()
    configEventSource = null
  }

  configEventSource = connectConfigSSE((evt: ConfigEvent) => {
    configEvents.value = [evt, ...configEvents.value].slice(0, MAX_EVENTS)

    // Auto-refresh dashboard stats on any mutation event
    if (evt.resource === 'route' || evt.resource === 'middleware' || evt.resource === 'instance') {
      dashboardApi.getStats().then(res => { data.value = res.data }).catch(() => {})
    }
  })

  configEventSource.onopen = () => { configSseConnected.value = true }
  configEventSource.onerror = () => { configSseConnected.value = false }
}

async function triggerSync() {
  syncing.value = true
  try {
    await dockerApi.sync()
    const [statusRes, statsRes] = await Promise.all([
      dockerApi.status(),
      dashboardApi.getStats(),
    ])
    dockerStatus.value = statusRes.data
    data.value = statsRes.data
  } catch {
    // handle error
  } finally {
    syncing.value = false
  }
}

async function fetchStats() {
  loading.value = true
  try {
    const [statsRes, instancesRes] = await Promise.all([
      dashboardApi.getStats(),
      instancesApi.list(),
    ])
    data.value = statsRes.data
    instances.value = instancesRes.data

    // Connect config SSE for live activity feed
    connectConfigEvents()

    // Pick the first active non-built-in instance for metrics
    // If a specific instance is selected, use that
    const candidateId = instanceStore.currentInstanceId
    const candidate = candidateId
      ? instances.value.find(i => i.id === candidateId && !i.builtIn && i.enableMetrics)
      : instances.value.find(i => !i.builtIn && i.enableMetrics)
    if (candidate) {
      metricsInstance.value = candidate
      fetchDashboardMetrics()
    }

    // Fetch Docker status if any built-in instance exists
    try {
      const dockerRes = await dockerApi.status()
      if (dockerRes.data.enabled) {
        dockerStatus.value = dockerRes.data
        connectSSE()
      }
    } catch {
      // Docker provider not available
    }
  } catch {
    // silently handle - stats will show zeros
  } finally {
    loading.value = false
  }
}

watch(() => instanceStore.currentInstanceId, fetchStats)
onMounted(fetchStats)
onUnmounted(() => {
  if (eventSource) {
    eventSource.close()
    eventSource = null
  }
  if (configEventSource) {
    configEventSource.close()
    configEventSource = null
  }
})
</script>

<style scoped>
.page-subtitle {
  font-size: 13px;
  color: var(--text-muted);
  margin-top: 2px;
}

.stat-card-link {
  text-decoration: none;
  color: inherit;
}
.stat-card-link:hover {
  color: inherit;
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

/* Dashboard panels */
.dashboard-panels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18px;
  margin-bottom: 18px;
}

@media (max-width: 900px) {
  .dashboard-panels {
    grid-template-columns: 1fr;
  }
}

/* Docker card */
.docker-header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-primary);
}

.docker-header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.docker-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.docker-stat {
  text-align: center;
  padding: 12px 8px;
  background: var(--bg-tertiary);
  border-radius: var(--radius);
}

.docker-stat-value {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.docker-stat-label {
  display: block;
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-weight: 500;
}

.spin {
  animation: spin 0.8s linear infinite;
}

/* Docker live events */
.docker-events {
  margin-top: 16px;
  border-top: 1px solid var(--border-secondary);
  padding-top: 14px;
}

.docker-events-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.docker-events-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.sse-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.sse-connected {
  background: var(--success-500);
  box-shadow: 0 0 6px var(--success-500);
}

.sse-disconnected {
  background: var(--text-muted);
}

.docker-event-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.docker-event-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: var(--radius-sm);
  background: var(--bg-tertiary);
  font-size: 12px;
  animation: fadeSlideIn 200ms ease;
}

.docker-event-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.event-sync_started { background: var(--primary-500); }
.event-sync_completed { background: var(--success-500); }
.event-routes_changed { background: var(--warning-500); }
.event-sync_error { background: var(--danger-500); }
.event-connected { background: var(--success-500); }

/* Config event resource-type dots */
.event-route { background: var(--success-500); }
.event-middleware { background: var(--warning-500); }
.event-instance { background: var(--primary-500); }

.docker-event-msg {
  flex: 1;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.docker-event-time {
  color: var(--text-muted);
  font-size: 11px;
  flex-shrink: 0;
}

@keyframes fadeSlideIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Instance list */
.instance-list-body {
  padding: 0 !important;
}

.instance-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 24px;
  border-bottom: 1px solid var(--border-secondary);
  transition: background var(--transition);
}

.instance-row:last-child {
  border-bottom: none;
}

.instance-row:hover {
  background: var(--bg-hover);
}

.instance-row-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.instance-row-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.instance-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dot-active {
  background: var(--success-500);
  box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.6);
  animation: pulse-green 2s infinite;
}
.dot-inactive {
  background: var(--text-muted);
}
.dot-unhealthy {
  background: var(--danger-500);
  box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6);
  animation: pulse-red 2s infinite;
}
.dot-unknown {
  background: var(--warning-500);
  box-shadow: 0 0 0 0 rgba(245, 158, 11, 0.6);
  animation: pulse-warning 2s infinite;
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

@keyframes pulse-warning {
  0% { box-shadow: 0 0 0 0 rgba(245, 158, 11, 0.6); }
  70% { box-shadow: 0 0 0 6px rgba(245, 158, 11, 0); }
  100% { box-shadow: 0 0 0 0 rgba(245, 158, 11, 0); }
}

.instance-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  text-decoration: none;
}

.instance-name:hover {
  color: var(--primary-600);
}

.instance-meta {
  font-size: 12px;
  color: var(--text-muted);
}

/* Quick actions */
.quick-actions-card {
  margin-top: 4px;
}

.quick-actions {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 12px;
}

.quick-action {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  border-radius: var(--radius);
  background: var(--bg-tertiary);
  text-decoration: none;
  color: inherit;
  transition: background var(--transition), transform var(--transition);
}
.quick-action:hover {
  background: var(--bg-hover);
  color: inherit;
  transform: translateX(2px);
}

.quick-action-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--radius);
  flex-shrink: 0;
}

.quick-action-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.quick-action-desc {
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 1px;
}

/* ─── Gateway Metrics ─── */
.gateway-metrics-card {
  margin-bottom: 18px;
}

.gw-metrics-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.gw-metric-stat {
  text-align: center;
  padding: 12px 8px;
  background: var(--bg-tertiary);
  border-radius: var(--radius);
}

.gw-metric-value {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}

.gw-metric-label {
  display: block;
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 2px;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  font-weight: 500;
}

.gw-metric-danger {
  color: var(--color-danger, #ef4444);
}

.gw-route-table-wrap {
  overflow-x: auto;
}

.gw-route-table-header {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-bottom: 8px;
}

.gw-route-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.gw-route-table th {
  text-align: left;
  padding: 8px 12px;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
  border-bottom: 1px solid var(--border-primary);
}

.gw-route-table td {
  padding: 8px 12px;
  color: var(--text-primary);
  border-bottom: 1px solid var(--border-secondary);
}

.gw-route-table tbody tr:last-child td {
  border-bottom: none;
}

.gw-route-table .text-right {
  text-align: right;
}

.gw-route-table .text-mono {
  font-family: var(--font-mono, monospace);
  font-size: 12px;
}
</style>
