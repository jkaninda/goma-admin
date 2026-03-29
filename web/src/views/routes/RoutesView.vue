<template>
  <div class="flex-container">
    <div class="page-header">
      <h1>Routes</h1>
      <div class="page-header-actions">
        <button class="btn btn-secondary" @click="openImport">Import YAML</button>
        <button class="btn btn-primary" @click="openCreate">New Route</button>
      </div>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <EmptyState
      v-else-if="routes.length === 0"
      title="No routes"
      description="Create your first route to start proxying traffic."
    >
      <template #action>
        <button class="btn btn-primary" @click="openCreate">Create Route</button>
      </template>
    </EmptyState>

    <template v-else>
      <div class="search-bar">
        <svg class="search-icon" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8" stroke-width="2" /><path stroke-linecap="round" stroke-width="2" d="m21 21-4.35-4.35" />
        </svg>
        <input v-model="search" class="form-input search-input" placeholder="Search routes..." />
      </div>

      <div class="card">
        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Path</th>
                <th>Methods</th>
                <th>Target</th>
                <th>Status</th>
                <th class="text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="route in routes" :key="route.id">
                <td>
                  <router-link :to="`/routes/${route.id}`" class="cell-name-link">
                    {{ route.name }}
                  </router-link>
                </td>
                <td class="text-mono">{{ route.config?.path || '-' }}</td>
                <td>
                  <span
                    v-for="m in getMethodsList(route)"
                    :key="m"
                    class="badge badge-info method-badge"
                  >{{ m }}</span>
                </td>
                <td class="text-mono truncate cell-target">
                  {{ route.config?.target || '-' }}
                </td>
                <td>
                  <span :class="['status-indicator', route.config?.enabled ? 'status-online' : 'status-offline']">
                    <span class="status-dot"></span>
                    {{ route.config?.enabled ? 'Online' : 'Offline' }}
                  </span>
                </td>
                <td class="text-right">
                  <div style="display: flex; gap: 6px; justify-content: flex-end">
                    <button class="btn btn-secondary btn-sm" @click="openEdit(route)">Edit</button>
                    <button class="btn btn-danger btn-sm" @click="confirmDelete(route)">Delete</button>
                  </div>
                </td>
              </tr>
              <tr v-if="routes.length === 0">
                <td colspan="6" class="text-center text-muted" style="padding: 32px">No matching routes</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Create / Edit Modal -->
    <Modal
      :show="modalOpen"
      :title="editing ? 'Edit Route' : 'New Route'"
      size="xl"
      @close="closeModal"
    >
      <div class="modal-body">
        <div class="form-group">
          <label class="form-label">Name</label>
          <input v-model="formName" required class="form-input" placeholder="my-api-route" />
        </div>

        <!-- Mode toggle -->
        <div class="tabs mode-tabs">
          <button :class="['tab', { active: formMode === 'simple' }]" @click="switchMode('simple')">Simple</button>
          <button :class="['tab', { active: formMode === 'advanced' }]" @click="switchMode('advanced')">Advanced</button>
        </div>

        <!-- Simple mode -->
        <template v-if="formMode === 'simple'">
          <div class="form-group">
            <label class="form-label">Path</label>
            <input v-model="simpleForm.path" class="form-input" placeholder="/" />
          </div>
          <div class="form-group">
            <label class="form-label">Target</label>
            <input v-model="simpleForm.target" class="form-input" placeholder="http://backend:8080" />
          </div>
          <div class="form-group">
            <label class="form-label">Hosts <span class="form-hint-inline">(comma-separated)</span></label>
            <input v-model="simpleForm.hosts" class="form-input" placeholder="api.example.com, api2.example.com" />
          </div>
          <div class="form-group">
            <label class="form-label">Methods <span class="form-hint-inline">(optional, defaults to all)</span></label>
            <div class="methods-grid">
              <label v-for="m in allMethods" :key="m" class="checkbox-label">
                <input type="checkbox" :value="m" v-model="simpleForm.methods" />
                {{ m }}
              </label>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">Rewrite <span class="form-hint-inline">(optional)</span></label>
            <input v-model="simpleForm.rewrite" class="form-input" placeholder="/new-prefix/" />
          </div>
          <div class="form-group toggle-row">
            <label class="form-label">Enabled</label>
            <button :class="['toggle-btn', { active: simpleForm.enabled }]" @click="simpleForm.enabled = !simpleForm.enabled">
              <span class="toggle-slider"></span>
            </button>
          </div>
        </template>

        <!-- Advanced mode -->
        <template v-else>
          <div class="form-warning">
            <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86l-8.6 14.86A1 1 0 0 0 2.54 20h18.92a1 1 0 0 0 .85-1.28l-8.6-14.86a1 1 0 0 0-1.42 0z" />
            </svg>
            <span>Enter advanced configuration at your own risk!</span>
          </div>
          <div class="form-group">
            <label class="form-label">Configuration (YAML)</label>
            <CodeEditor
              v-model="yamlContent"
              language="yaml"
              min-height="360px"
            />
          </div>
        </template>

        <div v-if="yamlError" class="form-error yaml-error">{{ yamlError }}</div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
          <button type="button" class="btn btn-primary" :disabled="saving" @click="handleSubmit">
            {{ saving ? 'Saving...' : (editing ? 'Update' : 'Create') }}
          </button>
        </div>
      </div>
    </Modal>
    <!-- Import Modal -->
    <Modal
      :show="importModalOpen"
      title="Import Routes from YAML"
      size="xl"
      @close="closeImportModal"
    >
      <div class="modal-body">
        <p class="form-hint import-hint">
          Paste a YAML file containing a <code>routes</code> list. Existing routes with the same name will be updated.
        </p>

        <div class="form-group">
          <div class="import-file-row">
            <label class="form-label">YAML Content</label>
            <label class="btn btn-secondary btn-sm import-file-btn">
              Load file
              <input type="file" accept=".yaml,.yml" hidden @change="handleImportFile" />
            </label>
          </div>
          <CodeEditor
            v-model="importYaml"
            language="yaml"
            min-height="360px"
          />
        </div>

        <div v-if="importError" class="form-error yaml-error">{{ importError }}</div>

        <div v-if="importResult" class="import-result">
          <span v-if="importResult.created" class="badge badge-success">{{ importResult.created }} created</span>
          <span v-if="importResult.updated" class="badge badge-info">{{ importResult.updated }} updated</span>
          <div v-if="importResult.errors?.length" class="import-errors">
            <div v-for="(err, i) in importResult.errors" :key="i" class="form-error">{{ err }}</div>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeImportModal">Cancel</button>
          <button type="button" class="btn btn-primary" :disabled="importing" @click="handleImport">
            {{ importing ? 'Importing...' : 'Import' }}
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { routesApi, type Route, type RouteCreateRequest, type ImportResult } from '@/api/routes'
import { useConfirm } from '@/composables/useConfirm'
import { useNotificationStore } from '@/stores/notification'
import Modal from '@/components/Modal.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import EmptyState from '@/components/EmptyState.vue'

const { confirm } = useConfirm()
const notify = useNotificationStore()

/* ── State ── */
const loading = ref(true)
const saving = ref(false)
const routes = ref<Route[]>([])
const modalOpen = ref(false)
const editing = ref<Route | null>(null)
const formName = ref('')
const yamlContent = ref('')
const yamlError = ref('')
const search = ref('')
const formMode = ref<'simple' | 'advanced'>('simple')

const allMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']

const simpleForm = ref({
  path: '/',
  target: '',
  hosts: '',
  methods: [] as string[],
  rewrite: '',
  enabled: true,
})
let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

watch(search, () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    fetchRoutes()
  }, 300)
})

onUnmounted(() => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
})

/* ── Import State ── */
const importModalOpen = ref(false)
const importing = ref(false)
const importYaml = ref('')
const importError = ref('')
const importResult = ref<ImportResult | null>(null)

/* ── Template helpers ── */
function getMethodsList(route: Route): string[] {
  const methods = route.config?.methods
  return Array.isArray(methods) ? methods : []
}

/* ── YAML helpers ── */
function objectToYaml(obj: Record<string, unknown>): string {
  const lines: string[] = []
  for (const [key, value] of Object.entries(obj)) {
    if (value === undefined || value === null) continue
    if (Array.isArray(value)) {
      if (value.length === 0) continue
      lines.push(`${key}:`)
      for (const item of value) {
        lines.push(`  - ${item}`)
      }
    } else if (typeof value === 'object') {
      lines.push(`${key}:`)
      for (const [k, v] of Object.entries(value as Record<string, unknown>)) {
        lines.push(`  ${k}: ${JSON.stringify(v)}`)
      }
    } else if (typeof value === 'string') {
      lines.push(`${key}: ${value}`)
    } else {
      lines.push(`${key}: ${value}`)
    }
  }
  return lines.join('\n')
}

function parseValue(val: string): unknown {
  if (val === 'true') return true
  if (val === 'false') return false
  if (val === 'null' || val === '~') return null
  if (/^-?\d+$/.test(val)) return parseInt(val, 10)
  if (/^-?\d+\.\d+$/.test(val)) return parseFloat(val)
  return val.replace(/^["']|["']$/g, '')
}

function yamlToObject(yamlStr: string): Record<string, unknown> {
  const obj: Record<string, unknown> = {}
  const lines = yamlStr.split('\n')
  let currentKey = ''
  let isCollectingArray = false
  let isCollectingObject = false
  let collectedArray: unknown[] = []
  let collectedObject: Record<string, unknown> = {}

  function flush() {
    if (!currentKey) return
    if (isCollectingArray && collectedArray.length > 0) {
      obj[currentKey] = collectedArray
    } else if (isCollectingObject && Object.keys(collectedObject).length > 0) {
      obj[currentKey] = collectedObject
    }
    isCollectingArray = false
    isCollectingObject = false
    collectedArray = []
    collectedObject = {}
  }

  for (const rawLine of lines) {
    const line = rawLine.replace(/\s+$/, '')
    if (!line || line.startsWith('#')) continue

    if (/^\s+/.test(line) && currentKey) {
      const arrayMatch = line.match(/^\s+-\s+(.*)$/)
      if (arrayMatch) {
        isCollectingArray = true
        collectedArray.push(parseValue(arrayMatch[1].trim()))
        continue
      }
      const nestedMatch = line.match(/^\s+(\w[\w-]*):\s*(.*)$/)
      if (nestedMatch) {
        isCollectingObject = true
        const val = nestedMatch[2].trim()
        collectedObject[nestedMatch[1]] = val === '' ? null : parseValue(val)
        continue
      }
      continue
    }

    const kvMatch = line.match(/^(\w[\w-]*):\s*(.*)$/)
    if (kvMatch) {
      flush()
      const key = kvMatch[1]
      const val = kvMatch[2].trim()
      currentKey = key
      if (val === '') continue
      obj[key] = parseValue(val)
      currentKey = ''
    }
  }

  flush()
  return obj
}

/* ── Default YAML template ── */
const defaultYaml = `path: /api/v1/*
target: https://backend.example.com
methods:
  - GET
  - POST
enabled: true`

/* ── Simple ↔ Advanced mode switching ── */
const simpleFieldKeys = new Set(['path', 'target', 'hosts', 'methods', 'rewrite', 'enabled'])

function hasAdvancedFields(config: Record<string, unknown>): boolean {
  return Object.keys(config).some(k => !simpleFieldKeys.has(k))
}

function simpleFormToConfig(): Record<string, unknown> {
  const config: Record<string, unknown> = {
    path: simpleForm.value.path || '/',
    target: simpleForm.value.target,
    enabled: simpleForm.value.enabled,
  }
  const hosts = simpleForm.value.hosts
    .split(',')
    .map(h => h.trim())
    .filter(Boolean)
  if (hosts.length) config.hosts = hosts
  if (simpleForm.value.methods.length) config.methods = [...simpleForm.value.methods]
  if (simpleForm.value.rewrite) config.rewrite = simpleForm.value.rewrite
  return config
}

function configToSimpleForm(config: Record<string, unknown>) {
  simpleForm.value.path = (config.path as string) || '/'
  simpleForm.value.target = (config.target as string) || ''
  const hosts = config.hosts
  simpleForm.value.hosts = Array.isArray(hosts) ? hosts.join(', ') : ''
  const methods = config.methods
  simpleForm.value.methods = Array.isArray(methods) ? methods.map(String) : []
  simpleForm.value.rewrite = (config.rewrite as string) || ''
  simpleForm.value.enabled = config.enabled !== false
}

function resetSimpleForm() {
  simpleForm.value = { path: '/', target: '', hosts: '', methods: [], rewrite: '', enabled: true }
}

function switchMode(mode: 'simple' | 'advanced') {
  if (mode === formMode.value) return
  yamlError.value = ''
  if (mode === 'advanced') {
    // Simple → Advanced: serialize form fields to YAML
    yamlContent.value = objectToYaml(simpleFormToConfig())
  } else {
    // Advanced → Simple: parse YAML into form fields
    try {
      const config = yamlToObject(yamlContent.value)
      if (hasAdvancedFields(config)) {
        yamlError.value = 'This config has advanced fields that will be lost in Simple mode. Switch anyway or stay in Advanced.'
      }
      configToSimpleForm(config)
    } catch {
      yamlError.value = 'Could not parse YAML into simple form.'
      return
    }
  }
  formMode.value = mode
}

/* ── Modal open / close ── */
function openCreate() {
  editing.value = null
  formName.value = ''
  formMode.value = 'simple'
  resetSimpleForm()
  yamlContent.value = defaultYaml
  yamlError.value = ''
  modalOpen.value = true
}

function openEdit(route: Route) {
  editing.value = route
  formName.value = route.name
  yamlError.value = ''
  const config = route.config || {}
  // Default to advanced mode if config has fields the simple form can't represent
  if (hasAdvancedFields(config)) {
    formMode.value = 'advanced'
    yamlContent.value = objectToYaml(config)
  } else {
    formMode.value = 'simple'
    configToSimpleForm(config)
    yamlContent.value = objectToYaml(config)
  }
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editing.value = null
  formName.value = ''
  yamlContent.value = ''
  yamlError.value = ''
  formMode.value = 'simple'
  resetSimpleForm()
}

/* ── Submit ── */
async function handleSubmit() {
  yamlError.value = ''

  if (!formName.value.trim()) {
    yamlError.value = 'Name is required.'
    return
  }

  let config: Record<string, unknown>
  if (formMode.value === 'simple') {
    if (!simpleForm.value.path.trim()) {
      yamlError.value = 'Path is required.'
      return
    }
    if (!simpleForm.value.target.trim()) {
      yamlError.value = 'Target is required.'
      return
    }
    if (!simpleForm.value.hosts.trim()) {
      yamlError.value = 'Hosts is required.'
      return
    }
    config = simpleFormToConfig()
  } else {
    try {
      config = yamlToObject(yamlContent.value)
    } catch {
      yamlError.value = 'Failed to parse YAML. Please check your syntax.'
      return
    }
  }

  const payload: RouteCreateRequest = {
    name: formName.value.trim(),
    config,
  }

  saving.value = true
  try {
    if (editing.value) {
      await routesApi.update(editing.value.id, payload)
    } else {
      await routesApi.create(payload)
    }
    closeModal()
    await fetchRoutes()
  } catch (e: any) {
    yamlError.value = e.response?.data?.message || 'Failed to save route.'
  } finally {
    saving.value = false
  }
}

/* ── Delete ── */
async function confirmDelete(route: Route) {
  const confirmed = await confirm({
    title: 'Delete Route',
    message: `Are you sure you want to delete "${route.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await routesApi.delete(route.id)
    await fetchRoutes()
  } catch {
    // Error handled by API interceptor
  }
}

/* ── Import ── */
const defaultImportYaml = `routes:
  - name: my-route
    path: /api/v1
    target: http://backend:8080
    methods:
      - GET
      - POST
    enabled: true`

function openImport() {
  importYaml.value = defaultImportYaml
  importError.value = ''
  importResult.value = null
  importModalOpen.value = true
}

function closeImportModal() {
  importModalOpen.value = false
  importYaml.value = ''
  importError.value = ''
  importResult.value = null
}

function handleImportFile(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  const reader = new FileReader()
  reader.onload = (e) => {
    importYaml.value = e.target?.result as string
    importError.value = ''
    importResult.value = null
  }
  reader.readAsText(file)
  // Reset the input so the same file can be re-selected
  ;(event.target as HTMLInputElement).value = ''
}

async function handleImport() {
  importError.value = ''
  importResult.value = null

  if (!importYaml.value.trim()) {
    importError.value = 'Please provide YAML content.'
    return
  }

  importing.value = true
  try {
    const res = await routesApi.importRoutes(importYaml.value)
    importResult.value = res.data
    const r = res.data
    if (r.created || r.updated) {
      notify.success(`Import complete: ${r.created} created, ${r.updated} updated`)
      await fetchRoutes()
    }
  } catch (e: any) {
    importError.value = e.response?.data?.message || 'Import request failed. Please try again.'
  } finally {
    importing.value = false
  }
}

/* ── Fetch ── */
async function fetchRoutes() {
  loading.value = true
  try {
    const res = await routesApi.list(0, 20, search.value)
    routes.value = res.data.data || []
  } catch {
    // Error handled by API interceptor
  } finally {
    loading.value = false
  }
}

onMounted(fetchRoutes)
</script>

<style scoped>
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

.cell-name-link {
  font-weight: 600;
  color: var(--primary-600);
  text-decoration: none;
}
.cell-name-link:hover {
  text-decoration: underline;
}

.cell-target {
  max-width: 220px;
}

.method-badge {
  margin-right: 4px;
  margin-bottom: 2px;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
  padding: 4px 10px;
  border-radius: 9999px;
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
  box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.6);
  animation: pulse-green 2s infinite;
}

.status-offline {
  color: var(--danger-500);
  background: var(--danger-50);
}

.status-offline .status-dot {
  background: var(--danger-500);
  box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6);
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

.action-delete {
  color: var(--danger-500);
}
.action-delete:hover {
  color: var(--danger-700);
}

.yaml-error {
  margin-top: 8px;
}

.form-warning {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 14px;
  margin-bottom: 16px;
  background: var(--warning-50, #fffbeb);
  border: 1px solid var(--warning-300, #fcd34d);
  border-radius: 6px;
  color: var(--warning-800, #92400e);
  font-size: 13px;
  line-height: 1.4;
}

.form-warning svg {
  flex-shrink: 0;
  margin-top: 1px;
  color: var(--warning-500, #f59e0b);
}

.mode-tabs {
  margin-bottom: 16px;
}

.methods-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 18px;
}

.toggle-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toggle-row .form-label {
  margin-bottom: 0;
}

.form-hint-inline {
  font-weight: 400;
  color: var(--text-muted);
  font-size: 12px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 20px;
}

.page-header-actions {
  display: flex;
  gap: 8px;
}

.import-hint {
  margin-bottom: 12px;
}
.import-hint code {
  background: var(--bg-tertiary);
  padding: 1px 5px;
  border-radius: 3px;
  font-size: 13px;
}

.import-file-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.import-result {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.import-errors {
  width: 100%;
  margin-top: 4px;
}
</style>
