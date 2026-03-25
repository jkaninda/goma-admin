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

        <div class="stat-card">
          <div class="stat-header">
            <span class="stat-label">Users</span>
            <span class="stat-icon stat-icon-danger">
              <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
            </span>
          </div>
          <span class="stat-value">{{ data.users }}</span>
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
import { ref, onMounted, watch } from 'vue'
import { dashboardApi, type DashboardStats } from '@/api/dashboard'
import { useInstanceStore } from '@/stores/instance'

const instanceStore = useInstanceStore()
const loading = ref(true)
const data = ref<DashboardStats>({ users: 0, instances: 0, middlewares: 0, routes: 0 })

async function fetchStats() {
  loading.value = true
  try {
    const res = await dashboardApi.getStats()
    data.value = res.data
  } catch {
    // silently handle - stats will show zeros
  } finally {
    loading.value = false
  }
}

watch(() => instanceStore.currentInstanceId, fetchStats)
onMounted(fetchStats)
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
</style>
