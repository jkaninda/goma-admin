<template>
  <div>
    <div class="page-header">
      <h1>Instances</h1>
      <button class="btn btn-primary" @click="showCreate = true">
        <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New Instance
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="instances.length === 0" class="card">
      <div class="empty-state">
        <svg width="48" height="48" fill="none" stroke="currentColor" viewBox="0 0 24 24" style="margin: 0 auto 12px;">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2" />
        </svg>
        <h3>No instances</h3>
        <p>Get started by creating a new gateway instance.</p>
        <button class="btn btn-primary" @click="showCreate = true">Create Instance</button>
      </div>
    </div>

    <!-- Search + Table -->
    <template v-else>
      <div class="search-bar">
        <svg class="search-icon" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8" stroke-width="2" /><path stroke-linecap="round" stroke-width="2" d="m21 21-4.35-4.35" />
        </svg>
        <input
          v-model="search"
          class="form-input search-input"
          placeholder="Search instances..."
        />
      </div>

      <div class="card">
        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Environment</th>
                <th>Endpoint</th>
                <th>Status</th>
                <th>Routes</th>
                <th>Middlewares</th>
                <th class="text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="inst in filteredInstances" :key="inst.id">
                <td>
                  <router-link :to="`/instances/${inst.id}`" class="instance-name-link">
                    {{ inst.name }}
                  </router-link>
                  <span v-if="inst.builtIn" class="badge badge-info" style="margin-left: 6px; font-size: 10px;">Built-in</span>
                </td>
                <td>
                  <span :class="['badge', envBadge(inst.environment)]">{{ inst.environment }}</span>
                </td>
                <td class="text-mono cell-endpoint truncate">{{ inst.endpoint }}</td>
                <td>
                  <span :class="['badge', statusBadge(inst.status)]">{{ inst.status }}</span>
                </td>
                <td>
                  <span class="count-badge">{{ inst.routes?.length || 0 }}</span>
                </td>
                <td>
                  <span class="count-badge">{{ inst.middlewares?.length || 0 }}</span>
                </td>
                <td class="text-right">
                  <div style="display: flex; gap: 6px; justify-content: flex-end">
                    <button class="btn btn-secondary btn-sm" @click="editInstance(inst)">Edit</button>
                    <button v-if="!inst.builtIn" class="btn btn-danger btn-sm" @click="deleteInstance(inst)">Delete</button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredInstances.length === 0">
                <td colspan="7" class="text-center text-muted" style="padding: 32px">No matching instances</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Create/Edit Modal -->
    <Teleport to="body">
      <div v-if="showCreate || !!editing" class="modal-overlay" @click.self="closeModal">
        <div class="modal">
          <div class="modal-header">
            <h3>{{ editing ? 'Edit Instance' : 'New Instance' }}</h3>
            <button class="btn-ghost btn-icon" @click="closeModal">
              <svg width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <form @submit.prevent="handleSubmit">
            <div class="modal-body">
              <div class="form-group">
                <label class="form-label" for="inst-name">Name</label>
                <input id="inst-name" v-model="form.name" class="form-input" required placeholder="my-gateway" :disabled="editing?.builtIn" />
              </div>
              <div class="form-grid">
                <div class="form-group">
                  <label class="form-label" for="inst-env">Environment</label>
                  <select id="inst-env" v-model="form.environment" class="form-select">
                    <option value="development">Development</option>
                    <option value="staging">Staging</option>
                    <option value="production">Production</option>
                    <option value="testing">Testing</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="form-label" for="inst-region">Region</label>
                  <input id="inst-region" v-model="form.region" class="form-input" placeholder="us-east-1" />
                </div>
              </div>
              <div class="form-group">
                <label class="form-label" for="inst-endpoint">Endpoint</label>
                <input id="inst-endpoint" v-model="form.endpoint" class="form-input" required placeholder="https://gateway.example.com" />
              </div>
              <div class="form-group">
                <label class="form-label" for="inst-desc">Description</label>
                <input id="inst-desc" v-model="form.description" class="form-input" placeholder="Optional description" />
              </div>

              <!-- Metrics Section -->
              <div class="form-section-divider">
                <span class="form-section-title">Metrics</span>
              </div>
              <div class="form-group">
                <label class="form-label form-label-inline" for="inst-enable-metrics">
                  <input id="inst-enable-metrics" type="checkbox" v-model="form.enableMetrics" class="form-checkbox" />
                  Enable Metrics
                </label>
                <span class="form-hint">Fetch Prometheus metrics from this instance</span>
              </div>
              <template v-if="form.enableMetrics">
                <div class="form-group">
                  <label class="form-label" for="inst-metrics-endpoint">Metrics Endpoint</label>
                  <input id="inst-metrics-endpoint" v-model="form.metricsEndpoint" class="form-input" placeholder="Defaults to endpoint/metrics" />
                  <span class="form-hint">Leave empty to use {endpoint}/metrics</span>
                </div>
                <div class="form-grid">
                  <div class="form-group">
                    <label class="form-label" for="inst-metrics-auth-type">Auth Type</label>
                    <select id="inst-metrics-auth-type" v-model="form.metricsAuthType" class="form-select">
                      <option value="">None</option>
                      <option value="basic">Basic Auth</option>
                      <option value="bearer">Bearer Token</option>
                      <option value="header">Custom Header</option>
                    </select>
                  </div>
                  <div v-if="form.metricsAuthType" class="form-group">
                    <label class="form-label" for="inst-metrics-auth-value">Auth Value</label>
                    <input
                      id="inst-metrics-auth-value"
                      v-model="form.metricsAuthValue"
                      class="form-input"
                      type="password"
                      :placeholder="editing?.hasMetricsAuth ? '••••••• (leave empty to keep)' : (form.metricsAuthType === 'basic' ? 'user:password' : 'token')"
                    />
                  </div>
                </div>
              </template>

              <!-- Repository Section -->
              <div class="form-section-divider">
                <span class="form-section-title">Git Repository</span>
              </div>
              <div class="form-group">
                <label class="form-label" for="inst-repo">Repository</label>
                <select id="inst-repo" v-model="form.repositoryId" class="form-select">
                  <option :value="null">None</option>
                  <option v-for="r in repositories" :key="r.id" :value="r.id">{{ r.name }} ({{ r.branch }})</option>
                </select>
              </div>
              <template v-if="form.repositoryId">
                <div class="form-group">
                  <label class="form-label" for="inst-repo-path">Config Path</label>
                  <input id="inst-repo-path" v-model="form.repositoryPath" class="form-input" placeholder="e.g. production/gateway-1" />
                  <span class="form-hint">Path within the repository where YAML configs are stored</span>
                </div>
                <div class="form-group">
                  <label class="form-label form-label-inline" for="inst-auto-sync">
                    <input id="inst-auto-sync" type="checkbox" v-model="form.autoSync" class="form-checkbox" />
                    Auto-sync on push
                  </label>
                  <span class="form-hint">Automatically import configs when the repository receives a push</span>
                </div>
              </template>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="saving">
                {{ saving ? 'Saving...' : (editing ? 'Update' : 'Create') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { instancesApi, type Instance, type InstanceCreateRequest } from '@/api/instances'
import { repositoriesApi, type Repository } from '@/api/repositories'
import { useConfirm } from '@/composables/useConfirm'

const { confirm } = useConfirm()

const loading = ref(true)
const saving = ref(false)
const instances = ref<Instance[]>([])
const search = ref('')
const showCreate = ref(false)
const editing = ref<Instance | null>(null)
const repositories = ref<Repository[]>([])

const form = reactive<InstanceCreateRequest>({
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

function resetForm() {
  form.name = ''
  form.environment = 'development'
  form.description = ''
  form.endpoint = ''
  form.region = ''
  form.enableMetrics = false
  form.metricsEndpoint = ''
  form.metricsAuthType = ''
  form.metricsAuthValue = ''
  form.repositoryId = null
  form.repositoryPath = ''
  form.autoSync = false
}

function closeModal() {
  showCreate.value = false
  editing.value = null
  resetForm()
}

function editInstance(inst: Instance) {
  editing.value = inst
  form.name = inst.name
  form.environment = inst.environment
  form.description = inst.description
  form.endpoint = inst.endpoint
  form.region = inst.region
  form.enableMetrics = inst.enableMetrics
  form.metricsEndpoint = inst.metricsEndpoint || ''
  form.metricsAuthType = inst.metricsAuthType || ''
  form.metricsAuthValue = ''
  form.repositoryId = inst.repositoryId || null
  form.repositoryPath = inst.repositoryPath || ''
  form.autoSync = inst.autoSync
}

async function handleSubmit() {
  saving.value = true
  try {
    if (editing.value) {
      await instancesApi.update(editing.value.id, form)
    } else {
      await instancesApi.create(form)
    }
    closeModal()
    await fetchInstances()
  } catch {
    // handle error
  } finally {
    saving.value = false
  }
}

async function deleteInstance(inst: Instance) {
  const confirmed = await confirm({
    title: 'Delete Instance',
    message: `Are you sure you want to delete "${inst.name}"? All routes and middlewares belonging to this instance will also be deleted.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await instancesApi.delete(inst.id)
    await fetchInstances()
  } catch {
    // handle error
  }
}

async function fetchInstances() {
  loading.value = true
  try {
    const res = await instancesApi.list()
    instances.value = res.data
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

const filteredInstances = computed(() => {
  if (!search.value.trim()) return instances.value
  const q = search.value.toLowerCase()
  return instances.value.filter(
    (i) => i.name.toLowerCase().includes(q) || i.environment.toLowerCase().includes(q) || i.endpoint.toLowerCase().includes(q)
  )
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

onMounted(() => {
  fetchInstances()
  repositoriesApi.list().then(res => { repositories.value = res.data }).catch(() => {})
})
</script>

<style scoped>
.instance-name-link {
  font-weight: 500;
  color: var(--primary-600);
}
.instance-name-link:hover {
  text-decoration: underline;
}

.action-delete:hover {
  color: var(--danger-600);
}

.search-bar {
  position: relative;
  margin-bottom: 16px;
}
.search-icon {
  position: absolute;
  left: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-muted);
  pointer-events: none;
}
.search-input {
  padding-left: 38px;
}

.cell-endpoint {
  max-width: 240px;
}

.count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  padding: 2px 8px;
  background: var(--bg-tertiary);
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
}

.form-section-divider {
  padding-top: 8px;
  margin-bottom: -4px;
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
