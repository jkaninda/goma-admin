<template>
  <div>
    <div class="page-header">
      <h1>Middlewares</h1>
      <div class="page-header-actions">
        <button class="btn btn-secondary" @click="openImport">Import YAML</button>
        <button class="btn btn-primary" @click="openCreate">New Middleware</button>
      </div>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <EmptyState
      v-else-if="middlewares.length === 0"
      title="No middlewares"
      description="Create your first middleware for rate limiting, authentication, and more."
    >
      <template #action>
        <button class="btn btn-primary" @click="openCreate">Create Middleware</button>
      </template>
    </EmptyState>

    <template v-else>
      <div class="search-bar">
        <svg class="search-icon" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="11" cy="11" r="8" stroke-width="2" /><path stroke-linecap="round" stroke-width="2" d="m21 21-4.35-4.35" />
        </svg>
        <input v-model="search" class="form-input search-input" placeholder="Search middlewares..." />
      </div>

      <div class="card">
        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Config</th>
                <th class="text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mw in middlewares" :key="mw.id">
                <td>
                  <router-link :to="`/middlewares/${mw.id}`" class="cell-name-link">{{ mw.name }}</router-link>
                </td>
                <td><span class="badge badge-info">{{ mw.type }}</span></td>
                <td class="cell-config truncate">{{ configPreview(mw.config) }}</td>
                <td class="text-right">
                  <div style="display: flex; gap: 6px; justify-content: flex-end">
                    <button class="btn btn-secondary btn-sm" @click="openEdit(mw)">Edit</button>
                    <button class="btn btn-danger btn-sm" @click="confirmDelete(mw)">Delete</button>
                  </div>
                </td>
              </tr>
              <tr v-if="middlewares.length === 0">
                <td colspan="4" class="text-center text-muted" style="padding: 32px">No matching middlewares</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <!-- Create / Edit Modal -->
    <Modal
      :show="modalOpen"
      :title="editing ? 'Edit Middleware' : 'New Middleware'"
      size="xl"
      @close="closeModal"
    >
      <div class="modal-body">
        <form @submit.prevent="handleSubmit">
          <div class="form-grid">
            <div class="form-group">
              <label class="form-label">Name</label>
              <input v-model="form.name" required class="form-input" placeholder="rate-limit" />
            </div>
            <div class="form-group">
              <label class="form-label">Type</label>
              <input v-model="form.type" required class="form-input" placeholder="rateLimit" />
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">Config (YAML)</label>
            <CodeEditor
              v-model="configYaml"
              language="yaml"
              min-height="260px"
            />
            <div v-if="configError" class="form-error">{{ configError }}</div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? 'Saving...' : (editing ? 'Update' : 'Create') }}
            </button>
          </div>
        </form>
      </div>
    </Modal>
    <!-- Import Modal -->
    <Modal
      :show="importModalOpen"
      title="Import Middlewares from YAML"
      size="xl"
      @close="closeImportModal"
    >
      <div class="modal-body">
        <p class="form-hint import-hint">
          Paste a YAML file containing a <code>middlewares</code> list. Existing middlewares with the same name will be updated.
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
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { middlewaresApi, type Middleware, type MiddlewareCreateRequest, type ImportResult } from '@/api/middlewares'
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
const middlewares = ref<Middleware[]>([])
const modalOpen = ref(false)
const editing = ref<Middleware | null>(null)
const editingId = ref<number | null>(null)
const configYaml = ref('')
const configError = ref('')
const search = ref('')
let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

watch(search, () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    fetchMiddlewares()
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

const form = reactive({
  name: '',
  type: '',
})

/* ── YAML helpers ── */

function objectToYaml(obj: unknown, indent = 0): string {
  const pad = '  '.repeat(indent)
  const lines: string[] = []

  if (Array.isArray(obj)) {
    for (const item of obj) {
      if (typeof item === 'object' && item !== null && !Array.isArray(item)) {
        const nested = objectToYaml(item, indent + 1)
        const nestedLines = nested.split('\n').filter(Boolean)
        if (nestedLines.length > 0) {
          lines.push(`${pad}- ${nestedLines[0].replace(/^\s+/, '')}`)
          for (let i = 1; i < nestedLines.length; i++) {
            lines.push(`${pad}  ${nestedLines[i].replace(/^\s+/, '')}`)
          }
        }
      } else {
        lines.push(`${pad}- ${formatScalar(item)}`)
      }
    }
    return lines.join('\n')
  }

  if (typeof obj === 'object' && obj !== null) {
    for (const [key, value] of Object.entries(obj)) {
      if (value === undefined) continue
      if (value === null) {
        lines.push(`${pad}${key}: null`)
      } else if (Array.isArray(value)) {
        if (value.length === 0) {
          lines.push(`${pad}${key}: []`)
        } else {
          lines.push(`${pad}${key}:`)
          lines.push(objectToYaml(value, indent + 1))
        }
      } else if (typeof value === 'object') {
        if (Object.keys(value).length === 0) {
          lines.push(`${pad}${key}: {}`)
        } else {
          lines.push(`${pad}${key}:`)
          lines.push(objectToYaml(value, indent + 1))
        }
      } else {
        lines.push(`${pad}${key}: ${formatScalar(value)}`)
      }
    }
    return lines.join('\n')
  }

  return `${pad}${formatScalar(obj)}`
}

function formatScalar(val: unknown): string {
  if (val === null) return 'null'
  if (val === undefined) return ''
  if (typeof val === 'boolean') return val ? 'true' : 'false'
  if (typeof val === 'number') return String(val)
  const s = String(val)
  if (s === '' || s === 'true' || s === 'false' || s === 'null' || s === '~' || /^[\d.-]/.test(s) || /[:#{}[\],&*?|>!%@`]/.test(s)) {
    return `"${s.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`
  }
  return s
}

function parseScalar(val: string): unknown {
  if (val === 'true') return true
  if (val === 'false') return false
  if (val === 'null' || val === '~') return null
  if (val === '[]') return []
  if (val === '{}') return {}
  if (/^-?\d+$/.test(val)) return parseInt(val, 10)
  if (/^-?\d+\.\d+$/.test(val)) return parseFloat(val)
  return val.replace(/^["']|["']$/g, '')
}

function yamlToObject(yamlStr: string): Record<string, unknown> {
  const lines = yamlStr.split('\n')
  return parseBlock(lines, 0, 0).value as Record<string, unknown>
}

function getIndent(line: string): number {
  const m = line.match(/^(\s*)/)
  return m ? m[1].length : 0
}

function parseBlock(lines: string[], start: number, baseIndent: number): { value: unknown; nextIndex: number } {
  const result: Record<string, unknown> = {}
  let i = start

  while (i < lines.length) {
    const raw = lines[i]
    const stripped = raw.replace(/\s+$/, '')

    // Skip empty / comment lines
    if (stripped === '' || stripped.replace(/^\s+/, '').startsWith('#')) {
      i++
      continue
    }

    const indent = getIndent(stripped)
    if (indent < baseIndent) break

    const content = stripped.replace(/^\s+/, '')

    // Array item at this level
    if (content.startsWith('- ') || content === '-') {
      // This block is actually an array — parse as array
      return parseArray(lines, start, baseIndent)
    }

    // Key-value line
    const kvMatch = content.match(/^([\w][\w.-]*):\s*(.*)$/)
    if (kvMatch) {
      const key = kvMatch[1]
      const val = kvMatch[2].trim()

      if (val === '' || val === '|' || val === '>') {
        // Check if next non-empty line is indented further
        const nextNonEmpty = findNextNonEmpty(lines, i + 1)
        if (nextNonEmpty !== -1 && getIndent(lines[nextNonEmpty]) > indent) {
          const childIndent = getIndent(lines[nextNonEmpty])
          const child = parseBlock(lines, i + 1, childIndent)
          result[key] = child.value
          i = child.nextIndex
        } else {
          result[key] = null
          i++
        }
      } else {
        result[key] = parseScalar(val)
        i++
      }
    } else {
      i++
    }
  }

  return { value: result, nextIndex: i }
}

function parseArray(lines: string[], start: number, baseIndent: number): { value: unknown; nextIndex: number } {
  const result: unknown[] = []
  let i = start

  while (i < lines.length) {
    const raw = lines[i]
    const stripped = raw.replace(/\s+$/, '')

    if (stripped === '' || stripped.replace(/^\s+/, '').startsWith('#')) {
      i++
      continue
    }

    const indent = getIndent(stripped)
    if (indent < baseIndent) break

    const content = stripped.replace(/^\s+/, '')

    if (content.startsWith('- ')) {
      const itemContent = content.slice(2).trim()

      // Check if the item itself contains a key: value (inline map item)
      const inlineKv = itemContent.match(/^([\w][\w.-]*):\s*(.*)$/)

      if (inlineKv) {
        // Could be a mapping entry — check for further nested lines
        const nextNonEmpty = findNextNonEmpty(lines, i + 1)
        const itemIndent = indent + 2 // items inside "- " are indented 2 more

        if (nextNonEmpty !== -1 && getIndent(lines[nextNonEmpty]) >= itemIndent) {
          // Multi-line mapping item: reconstruct lines with adjusted indentation
          const subLines: string[] = ['  '.repeat(0) + itemContent]
          let j = i + 1
          while (j < lines.length) {
            const subRaw = lines[j]
            const subStripped = subRaw.replace(/\s+$/, '')
            if (subStripped === '' || subStripped.replace(/^\s+/, '').startsWith('#')) {
              j++
              continue
            }
            if (getIndent(subStripped) < itemIndent) break
            // Re-indent relative to itemIndent
            subLines.push(subStripped.slice(itemIndent))
            j++
          }
          const child = parseBlock(subLines, 0, 0)
          result.push(child.value)
          i = j
        } else {
          // Single key-value — treat as small object
          const obj: Record<string, unknown> = {}
          obj[inlineKv[1]] = inlineKv[2].trim() === '' ? null : parseScalar(inlineKv[2].trim())
          result.push(obj)
          i++
        }
      } else {
        result.push(parseScalar(itemContent))
        i++
      }
    } else if (content === '-') {
      // Bare dash, next indented block is the item
      const nextNonEmpty = findNextNonEmpty(lines, i + 1)
      if (nextNonEmpty !== -1 && getIndent(lines[nextNonEmpty]) > indent) {
        const childIndent = getIndent(lines[nextNonEmpty])
        const child = parseBlock(lines, i + 1, childIndent)
        result.push(child.value)
        i = child.nextIndex
      } else {
        result.push(null)
        i++
      }
    } else {
      break
    }
  }

  return { value: result, nextIndex: i }
}

function findNextNonEmpty(lines: string[], start: number): number {
  for (let i = start; i < lines.length; i++) {
    const stripped = lines[i].replace(/\s+$/, '')
    if (stripped !== '' && !stripped.replace(/^\s+/, '').startsWith('#')) return i
  }
  return -1
}

/* ── Config preview for table ── */
function configPreview(config: Record<string, unknown>): string {
  if (!config || typeof config !== 'object') return '-'
  const keys = Object.keys(config)
  if (keys.length === 0) return '-'

  const parts: string[] = []
  for (const key of keys) {
    const val = config[key]
    if (Array.isArray(val)) {
      if (val.length <= 2) {
        parts.push(`${key}: ${val.join(', ')}`)
      } else {
        parts.push(`${key}: ${val.slice(0, 2).join(', ')} +${val.length - 2}`)
      }
    } else if (typeof val === 'object' && val !== null) {
      const count = Object.keys(val).length
      parts.push(`${key}: ${count} key${count !== 1 ? 's' : ''}`)
    } else {
      parts.push(`${key}: ${val}`)
    }
  }
  const summary = parts.join(', ')
  return summary.length > 80 ? summary.slice(0, 77) + '...' : summary
}

/* ── Default config template ── */
const defaultConfigYaml = `paths:
  - /api/*
rule:
  requestsPerUnit: 100
  unit: minute
  burst: 10`

/* ── Modal open / close ── */
function resetForm() {
  form.name = ''
  form.type = ''
  configYaml.value = ''
  configError.value = ''
}

function openCreate() {
  resetForm()
  editing.value = null
  editingId.value = null
  configYaml.value = defaultConfigYaml
  modalOpen.value = true
}

function openEdit(mw: Middleware) {
  editing.value = mw
  editingId.value = mw.id
  form.name = mw.name
  form.type = mw.type
  configYaml.value = mw.config && Object.keys(mw.config).length > 0
    ? objectToYaml(mw.config)
    : ''
  configError.value = ''
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editing.value = null
  editingId.value = null
  resetForm()
}

/* ── Submit ── */
function buildPayload(): MiddlewareCreateRequest | null {
  configError.value = ''

  let config: Record<string, unknown> = {}
  try {
    if (configYaml.value.trim()) {
      const parsed = yamlToObject(configYaml.value)
      config = parsed
    }
  } catch {
    configError.value = 'Invalid YAML in config. Please check your syntax.'
    return null
  }

  if (!form.name.trim()) return null
  if (!form.type.trim()) return null

  return {
    name: form.name,
    type: form.type,
    config,
  }
}

async function handleSubmit() {
  const data = buildPayload()
  if (!data) return

  saving.value = true
  try {
    if (editing.value && editingId.value) {
      await middlewaresApi.update(editingId.value, data)
    } else {
      await middlewaresApi.create(data)
    }
    closeModal()
    await fetchMiddlewares()
  } catch {
    // Error handled by API interceptor
  } finally {
    saving.value = false
  }
}

/* ── Delete ── */
async function confirmDelete(mw: Middleware) {
  const confirmed = await confirm({
    title: 'Delete Middleware',
    message: `Are you sure you want to delete "${mw.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await middlewaresApi.delete(mw.id)
    await fetchMiddlewares()
  } catch {
    // Error handled by API interceptor
  }
}

/* ── Import ── */
const defaultImportYaml = `middlewares:
  - name: rate-limit
    type: rateLimit
    rule:
      unit: minute
      requestsPerUnit: 60
      burst: 100`

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
    const res = await middlewaresApi.importMiddlewares(importYaml.value)
    importResult.value = res.data
    const r = res.data
    if (r.created || r.updated) {
      notify.success(`Import complete: ${r.created} created, ${r.updated} updated`)
      await fetchMiddlewares()
    }
  } catch (e: any) {
    importError.value = e.response?.data?.message || 'Import request failed. Please try again.'
  } finally {
    importing.value = false
  }
}

/* ── Fetch ── */
async function fetchMiddlewares() {
  loading.value = true
  try {
    const res = await middlewaresApi.list(0, 20, search.value)
    middlewares.value = res.data.data || []
  } catch {
    // Error handled by API interceptor
  } finally {
    loading.value = false
  }
}

onMounted(fetchMiddlewares)
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

.cell-config {
  max-width: 320px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 12px;
  color: var(--text-tertiary);
}

.action-delete {
  color: var(--danger-500);
}
.action-delete:hover {
  color: var(--danger-700);
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

.yaml-error {
  margin-top: 8px;
}
</style>
