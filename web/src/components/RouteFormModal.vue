<template>
  <Modal
    :show="show"
    :title="route ? 'Edit Route' : 'New Route'"
    size="xl"
    @close="emit('close')"
  >
    <div class="modal-body">
      <div class="form-group">
        <label class="form-label">Name</label>
        <input
          v-model="formName"
          required
          class="form-input"
          placeholder="my-api-route"
        />
      </div>

      <!-- Mode toggle -->
      <div class="tabs mode-tabs">
        <button
          :class="['tab', { active: formMode === 'simple' }]"
          @click="switchMode('simple')"
        >
          Simple
        </button>
        <button
          :class="['tab', { active: formMode === 'advanced' }]"
          @click="switchMode('advanced')"
        >
          Advanced
        </button>
      </div>

      <!-- Simple mode -->
      <template v-if="formMode === 'simple'">
        <div class="form-group">
          <label class="form-label">Path</label>
          <input v-model="simpleForm.path" class="form-input" placeholder="/" />
        </div>
        <div class="form-group">
          <label class="form-label">Target</label>
          <input
            v-model="simpleForm.target"
            class="form-input"
            placeholder="http://backend:8080"
          />
        </div>
        <div class="form-group">
          <label class="form-label" for="hosts"
            >Hosts <span class="form-hint-inline">(comma-separated)</span>
          </label>

          <div class="input-host-container">
            <div class="host-list">
              <span v-for="(host, index) in parsedHosts" :key="index" class="host-chip">
                {{ host }}
                <button
                  type="button"
                  @click="handleRemoveHost(host, $event)"
                  class="host-remove-btn"
                  @touchend.prevent="handleRemoveHost(host, $event)"
                >
                  <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                  >
                    <line x1="18" y1="6" x2="6" y2="18"></line>
                    <line x1="6" y1="6" x2="18" y2="18"></line>
                  </svg>
                </button>
              </span>

              <input
                id="hosts"
                v-model="hostInput"
                class="host-input-field"
                type="text"
                placeholder="api.example.com, api2.example.com"
                @input="handleHostInput"
                @keydown="handleHostInputKeydown"
                @blur="handleHostInputBlur"
              />
            </div>
          </div>

          <p class="form-hint">Use Enter, Space, or Comma to add a host</p>
          <p class="form-tip">
            Tip: Type "api.example.com" then press comma to add quickly
          </p>
        </div>

        <div class="form-group">
          <label class="form-label"
            >Methods
            <span class="form-hint-inline">(optional, defaults to all)</span></label
          >
          <div class="methods-grid">
            <label v-for="m in allMethods" :key="m" class="checkbox-label">
              <input type="checkbox" :value="m" v-model="simpleForm.methods" />
              {{ m }}
            </label>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label"
            >Middlewares <span class="form-hint-inline">(optional)</span></label
          >
          <div v-if="availableMiddlewares.length" class="middleware-select">
            <label
              v-for="mw in availableMiddlewares"
              :key="mw.id"
              class="checkbox-label middleware-option"
            >
              <input type="checkbox" :value="mw.name" v-model="simpleForm.middlewares" />
              <span class="middleware-option-name">{{ mw.name }}</span>
              <span class="badge badge-info middleware-option-type">{{ mw.type }}</span>
            </label>
          </div>
          <p v-else class="form-hint">No middlewares available. Create middlewares first.</p>
        </div>
        <div class="form-group">
          <label class="form-label"
            >Rewrite <span class="form-hint-inline">(optional)</span></label
          >
          <input
            v-model="simpleForm.rewrite"
            class="form-input"
            placeholder="/new-prefix/"
          />
        </div>
        <div class="form-group toggle-row">
          <label class="form-label">Enabled</label>
          <button
            :class="['toggle-btn', { active: simpleForm.enabled }]"
            @click="simpleForm.enabled = !simpleForm.enabled"
          >
            <span class="toggle-slider"></span>
          </button>
        </div>
      </template>

      <!-- Advanced mode -->
      <template v-else>
        <div class="form-warning">
          <svg
            width="16"
            height="16"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 9v2m0 4h.01M10.29 3.86l-8.6 14.86A1 1 0 0 0 2.54 20h18.92a1 1 0 0 0 .85-1.28l-8.6-14.86a1 1 0 0 0-1.42 0z"
            />
          </svg>
          <span>Enter advanced configuration at your own risk!</span>
        </div>
        <div class="form-group">
          <label class="form-label">Configuration (YAML)</label>
          <CodeEditor v-model="yamlContent" language="yaml" min-height="360px" />
        </div>
      </template>

      <div v-if="yamlError" class="form-error yaml-error">{{ yamlError }}</div>

      <div class="modal-footer">
        <button type="button" class="btn btn-secondary" @click="emit('close')">
          Cancel
        </button>
        <button
          type="button"
          class="btn btn-primary"
          :disabled="saving"
          @click="handleSubmit"
        >
          {{ saving ? "Saving..." : route ? "Update" : "Create" }}
        </button>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { routesApi, type Route, type RouteCreateRequest } from '@/api/routes'
import { middlewaresApi, type Middleware } from '@/api/middlewares'
import Modal from '@/components/Modal.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import { parseYaml, toYaml } from '@/utils/yaml'

const props = defineProps<{
  show: boolean
  /** Route being edited, or null to create a new one. */
  route: Route | null
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'saved', route: Route): void
}>()

const allMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']

/* ── Form state ── */
const saving = ref(false)
const formName = ref('')
const yamlContent = ref('')
const yamlError = ref('')
const formMode = ref<'simple' | 'advanced'>('simple')
const availableMiddlewares = ref<Middleware[]>([])
const hostInput = ref('')

const simpleForm = ref({
  path: '/',
  target: '',
  hosts: '',
  methods: [] as string[],
  middlewares: [] as string[],
  rewrite: '',
  enabled: true,
})

const defaultYaml = `path: /api/v1/*
hosts: ["api.example.com"]
target: https://backend.example.com
methods:
  - GET
  - POST
enabled: true`

/* ── Initialize the form whenever the modal opens ── */
watch(
  () => props.show,
  (open) => {
    if (open) initForm()
  },
  { immediate: true }
)

function initForm() {
  yamlError.value = ''
  if (props.route) {
    formName.value = props.route.name
    const config = props.route.config || {}
    // Default to advanced mode if the config has fields the simple form can't represent.
    if (hasAdvancedFields(config)) {
      formMode.value = 'advanced'
      yamlContent.value = toYaml(config)
    } else {
      formMode.value = 'simple'
      configToSimpleForm(config)
      yamlContent.value = toYaml(config)
    }
  } else {
    formName.value = ''
    formMode.value = 'simple'
    resetSimpleForm()
    yamlContent.value = defaultYaml
  }
  fetchAvailableMiddlewares()
}

/* ── Simple ↔ Advanced mode switching ── */
const simpleFieldKeys = new Set(['path', 'target', 'hosts', 'methods', 'middlewares', 'rewrite', 'enabled'])

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
  if (simpleForm.value.middlewares.length) config.middlewares = [...simpleForm.value.middlewares]
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
  const mws = config.middlewares
  simpleForm.value.middlewares = Array.isArray(mws) ? mws.map(String) : []
  simpleForm.value.rewrite = (config.rewrite as string) || ''
  simpleForm.value.enabled = config.enabled !== false
}

function resetSimpleForm() {
  simpleForm.value = { path: '/', target: '', hosts: '', methods: [], middlewares: [], rewrite: '', enabled: true }
}

function switchMode(mode: 'simple' | 'advanced') {
  if (mode === formMode.value) return
  yamlError.value = ''
  if (mode === 'advanced') {
    // Simple → Advanced: serialize form fields to YAML
    yamlContent.value = toYaml(simpleFormToConfig())
  } else {
    // Advanced → Simple: parse YAML into form fields
    try {
      const config = parseYaml(yamlContent.value)
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

/* ── Hosts management ── */
const parsedHosts = computed(() => {
  if (!simpleForm.value.hosts) return []
  return simpleForm.value.hosts
    .split(',')
    .map((t) => t.trim())
    .filter((t) => t)
})

const addHost = () => {
  const newHost = hostInput.value.trim()
  if (!newHost) return

  const currentHosts = parsedHosts.value
  if (currentHosts.includes(newHost)) {
    hostInput.value = ''
    return
  }

  simpleForm.value.hosts = [...currentHosts, newHost].join(', ')
  hostInput.value = ''
}

const removeHost = (hostToRemove: string) => {
  simpleForm.value.hosts = parsedHosts.value.filter((host) => host !== hostToRemove).join(', ')
}

const handleHostInputKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    e.preventDefault()
    addHost()
  } else if (e.key === ' ') {
    // Allow space for multi-word Hosts, but trim and add if there's content
    if (hostInput.value.trim()) {
      e.preventDefault()
      addHost()
    }
  } else if (e.key === 'Backspace' && !hostInput.value && parsedHosts.value.length > 0) {
    // Remove last host when backspace is pressed on empty input
    e.preventDefault()
    removeHost(parsedHosts.value[parsedHosts.value.length - 1])
  }
}

// Handle comma and space input to work on all keyboards/devices
const handleHostInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  const value = target.value

  // Check if the last character is a comma or space
  if (value.endsWith(',') || (value.endsWith(' ') && value.trim().length > 0)) {
    // Remove the comma/space and add the host
    const cleanValue = value.slice(0, -1).trim()
    if (cleanValue) {
      hostInput.value = cleanValue
      addHost()
    } else {
      hostInput.value = ''
    }
  }
}

// Add host on blur (when user clicks/taps outside) for better mobile experience
const handleHostInputBlur = () => {
  if (hostInput.value.trim()) {
    addHost()
  }
}

// Handle host removal with proper touch event support
const handleRemoveHost = (host: string, e: Event) => {
  e.preventDefault()
  e.stopPropagation()
  removeHost(host)
}

/* ── Middleware list for simple mode ── */
async function fetchAvailableMiddlewares() {
  try {
    const res = await middlewaresApi.list(0, 100)
    availableMiddlewares.value = res.data.data || []
  } catch { /* non-critical */ }
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
      config = parseYaml(yamlContent.value)
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
    const res = props.route
      ? await routesApi.update(props.route.id, payload)
      : await routesApi.create(payload)
    emit('saved', res.data)
    emit('close')
  } catch (e: any) {
    yamlError.value = e.response?.data?.message || 'Failed to save route.'
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
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

.middleware-select {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 200px;
  overflow-y: auto;
  padding: 8px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-secondary);
}
.middleware-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 6px;
  border-radius: 4px;
}
.middleware-option:hover {
  background: var(--bg-tertiary);
}
.middleware-option-name {
  font-size: 13px;
  font-weight: 500;
}
.middleware-option-type {
  font-size: 10px;
  margin-left: auto;
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
</style>
