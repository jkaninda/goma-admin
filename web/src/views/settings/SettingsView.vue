<template>
  <div>
    <div class="page-header">
      <h1>Settings</h1>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <template v-else>
      <!-- Export Section -->
      <div class="card section-card">
        <div class="card-header">
          <h2>Export Config</h2>
        </div>
        <div class="card-body">
          <p class="section-desc">
            Export all routes and middlewares for an instance as a YAML file.
            The exported file can be imported into another Goma Admin instance.
          </p>
          <div class="form-grid">
            <div class="form-group">
              <label class="form-label">Instance</label>
              <select v-model="exportInstanceId" class="form-select">
                <option :value="null" disabled>Select instance</option>
                <option v-for="inst in instances" :key="inst.id" :value="inst.id">
                  {{ inst.name }} ({{ inst.environment }})
                </option>
              </select>
            </div>
            <div class="form-group form-actions-cell">
              <button
                class="btn btn-primary"
                :disabled="!exportInstanceId || exporting"
                @click="handleExport"
              >
                {{ exporting ? 'Exporting...' : 'Export YAML' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Import Section -->
      <div class="card section-card">
        <div class="card-header">
          <h2>Import Config</h2>
        </div>
        <div class="card-body">
          <p class="section-desc">
            Import routes and middlewares from a YAML file into an instance.
            Existing items with the same name will be updated.
          </p>
          <div class="form-grid">
            <div class="form-group">
              <label class="form-label">Target Instance</label>
              <select v-model="importInstanceId" class="form-select">
                <option :value="null" disabled>Select instance</option>
                <option v-for="inst in instances" :key="inst.id" :value="inst.id">
                  {{ inst.name }} ({{ inst.environment }})
                </option>
              </select>
            </div>
            <div class="form-group form-actions-cell">
              <label class="btn btn-secondary" :class="{ disabled: !importInstanceId }">
                Load YAML file
                <input
                  type="file"
                  accept=".yaml,.yml"
                  hidden
                  :disabled="!importInstanceId"
                  @change="handleImportFile"
                />
              </label>
            </div>
          </div>

          <div v-if="importYaml" class="form-group">
            <label class="form-label">Preview</label>
            <CodeEditor v-model="importYaml" language="yaml" min-height="260px" />
          </div>

          <div v-if="importError" class="form-error">{{ importError }}</div>

          <div v-if="importResult" class="import-result">
            <span v-if="importResult.created" class="badge badge-success">{{ importResult.created }} created</span>
            <span v-if="importResult.updated" class="badge badge-info">{{ importResult.updated }} updated</span>
            <div v-if="importResult.errors?.length" class="import-errors">
              <div v-for="(err, i) in importResult.errors" :key="i" class="form-error">{{ err }}</div>
            </div>
          </div>

          <div v-if="importYaml" class="form-actions">
            <button class="btn btn-secondary" @click="importYaml = ''; importError = ''; importResult = null">
              Clear
            </button>
            <button class="btn btn-primary" :disabled="importing" @click="handleImport">
              {{ importing ? 'Importing...' : 'Import' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Copy Section -->
      <div class="card section-card">
        <div class="card-header">
          <h2>Copy Instance Config</h2>
        </div>
        <div class="card-body">
          <p class="section-desc">
            Copy all routes and middlewares from one instance to another.
            For example, copy dev config into production.
            Existing items with the same name on the target will be updated.
          </p>
          <div class="form-grid-3">
            <div class="form-group">
              <label class="form-label">Source Instance</label>
              <select v-model="copySourceId" class="form-select">
                <option :value="null" disabled>Select source</option>
                <option v-for="inst in instances" :key="inst.id" :value="inst.id">
                  {{ inst.name }} ({{ inst.environment }})
                </option>
              </select>
            </div>
            <div class="copy-arrow">
              <svg width="24" height="24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
            </div>
            <div class="form-group">
              <label class="form-label">Target Instance</label>
              <select v-model="copyTargetId" class="form-select">
                <option :value="null" disabled>Select target</option>
                <option
                  v-for="inst in instances"
                  :key="inst.id"
                  :value="inst.id"
                  :disabled="inst.id === copySourceId"
                >
                  {{ inst.name }} ({{ inst.environment }})
                </option>
              </select>
            </div>
          </div>

          <div v-if="copyError" class="form-error">{{ copyError }}</div>

          <div v-if="copyResult" class="import-result">
            <span v-if="copyResult.created" class="badge badge-success">{{ copyResult.created }} created</span>
            <span v-if="copyResult.updated" class="badge badge-info">{{ copyResult.updated }} updated</span>
            <div v-if="copyResult.errors?.length" class="import-errors">
              <div v-for="(err, i) in copyResult.errors" :key="i" class="form-error">{{ err }}</div>
            </div>
          </div>

          <div class="form-actions">
            <button
              class="btn btn-primary"
              :disabled="!copySourceId || !copyTargetId || copying"
              @click="handleCopy"
            >
              {{ copying ? 'Copying...' : 'Copy Config' }}
            </button>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { instancesApi, type Instance, type ImportResult } from '@/api/instances'
import { useConfirm } from '@/composables/useConfirm'
import { useNotificationStore } from '@/stores/notification'
import CodeEditor from '@/components/CodeEditor.vue'

const { confirm } = useConfirm()
const notify = useNotificationStore()

/* ── State ── */
const loading = ref(true)
const instances = ref<Instance[]>([])

/* Export */
const exportInstanceId = ref<number | null>(null)
const exporting = ref(false)

/* Import */
const importInstanceId = ref<number | null>(null)
const importing = ref(false)
const importYaml = ref('')
const importError = ref('')
const importResult = ref<ImportResult | null>(null)

/* Copy */
const copySourceId = ref<number | null>(null)
const copyTargetId = ref<number | null>(null)
const copying = ref(false)
const copyError = ref('')
const copyResult = ref<ImportResult | null>(null)

/* ── Export ── */
async function handleExport() {
  if (!exportInstanceId.value) return
  exporting.value = true
  try {
    const res = await instancesApi.exportConfig(exportInstanceId.value)
    const yamlContent = typeof res.data === 'string' ? res.data : JSON.stringify(res.data)

    const inst = instances.value.find((i) => i.id === exportInstanceId.value)
    const filename = `goma-config-${inst?.name || 'instance'}.yaml`

    const blob = new Blob([yamlContent], { type: 'application/x-yaml' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    URL.revokeObjectURL(url)

    notify.success('Config exported successfully')
  } catch {
    notify.error('Failed to export config')
  } finally {
    exporting.value = false
  }
}

/* ── Import ── */
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
  if (!importInstanceId.value || !importYaml.value.trim()) return
  importError.value = ''
  importResult.value = null

  importing.value = true
  try {
    const res = await instancesApi.importConfig(importInstanceId.value, importYaml.value)
    importResult.value = res.data
    const r = res.data
    if (r.created || r.updated) {
      notify.success(`Import complete: ${r.created} created, ${r.updated} updated`)
    }
  } catch (e: any) {
    importError.value = e.response?.data?.message || 'Import failed. Please try again.'
  } finally {
    importing.value = false
  }
}

/* ── Copy ── */
async function handleCopy() {
  if (!copySourceId.value || !copyTargetId.value) return
  copyError.value = ''
  copyResult.value = null

  const source = instances.value.find((i) => i.id === copySourceId.value)
  const target = instances.value.find((i) => i.id === copyTargetId.value)

  const confirmed = await confirm({
    title: 'Copy Instance Config',
    message: `Copy all routes and middlewares from "${source?.name}" to "${target?.name}"? Existing items with the same name on the target will be updated.`,
    confirmText: 'Copy',
    variant: 'warning',
  })
  if (!confirmed) return

  copying.value = true
  try {
    const res = await instancesApi.copyTo(copySourceId.value, copyTargetId.value)
    copyResult.value = res.data
    const r = res.data
    if (r.created || r.updated) {
      notify.success(`Copy complete: ${r.created} created, ${r.updated} updated`)
    }
  } catch (e: any) {
    copyError.value = e.response?.data?.message || 'Copy failed. Please try again.'
  } finally {
    copying.value = false
  }
}

/* ── Fetch ── */
onMounted(async () => {
  try {
    const res = await instancesApi.list()
    instances.value = res.data || []
  } catch {
    // handled by interceptor
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.section-card {
  margin-bottom: 24px;
}

.section-desc {
  color: var(--text-secondary);
  font-size: 14px;
  margin-bottom: 16px;
  line-height: 1.5;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 16px;
  align-items: end;
}

.form-grid-3 {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 16px;
  align-items: end;
}

.form-actions-cell {
  display: flex;
  align-items: flex-end;
  padding-bottom: 2px;
}

.form-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 16px;
}

.copy-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
  justify-self: center;
  align-self: center;
}

.import-result {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 12px;
}

.import-errors {
  width: 100%;
  margin-top: 4px;
}

.disabled {
  opacity: 0.5;
  pointer-events: none;
}

@media (max-width: 768px) {
  .form-grid,
  .form-grid-3 {
    grid-template-columns: 1fr;
  }
  .copy-arrow {
    transform: rotate(90deg);
    padding: 0;
  }
}
</style>
