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
        <svg
          class="search-icon"
          width="16"
          height="16"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <circle cx="11" cy="11" r="8" stroke-width="2" />
          <path stroke-linecap="round" stroke-width="2" d="m21 21-4.35-4.35" />
        </svg>
        <input
          v-model="search"
          class="form-input search-input"
          placeholder="Search routes..."
        />
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
              <tr
                v-for="route in routes"
                :key="route.id"
                class="cursor-pointer"
                @click="$router.push(`/routes/${route.id}`)"
              >
                <td>
                  <router-link :to="`/routes/${route.id}`" class="cell-name-link">
                    {{ route.name }}
                  </router-link>
                </td>
                <td class="text-mono">{{ route.config?.path || "-" }}</td>
                <td>
                  <span
                    v-for="m in getMethodsList(route)"
                    :key="m"
                    class="badge badge-info method-badge"
                    >{{ m }}</span
                  >
                </td>
                <td class="text-mono truncate cell-target">
                  {{ route.config?.target || "-" }}
                </td>
                <td>
                  <span
                    :class="[
                      'status-indicator',
                      route.config?.enabled ? 'status-online' : 'status-offline',
                    ]"
                  >
                    <span class="status-dot"></span>
                    {{ route.config?.enabled ? "Online" : "Offline" }}
                  </span>
                </td>
                <td class="text-right">
                  <div style="display: flex; gap: 6px; justify-content: flex-end">
                    <button
                      class="btn btn-secondary btn-sm"
                      @click.stop="openEdit(route)"
                    >
                      Edit
                    </button>
                    <button
                      class="btn btn-danger btn-sm"
                      @click.stop="confirmDelete(route)"
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="routes.length === 0">
                <td colspan="6" class="text-center text-muted" style="padding: 32px">
                  No matching routes
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <Pagination
          :pageable="pageable"
          @page="goToPage"
        />
      </div>
    </template>

    <!-- Create / Edit Modal -->
    <RouteFormModal
      :show="modalOpen"
      :route="editing"
      @close="closeModal"
      @saved="onSaved"
    />
    <!-- Import Modal -->
    <Modal
      :show="importModalOpen"
      title="Import Routes from YAML"
      size="xl"
      @close="closeImportModal"
    >
      <div class="modal-body">
        <p class="form-hint import-hint">
          Paste a YAML file containing a <code>routes</code> list. Existing routes with
          the same name will be updated.
        </p>

        <div class="form-group">
          <div class="import-file-row">
            <label class="form-label">YAML Content</label>
            <label class="btn btn-secondary btn-sm import-file-btn">
              Load file
              <input type="file" accept=".yaml,.yml" hidden @change="handleImportFile" />
            </label>
          </div>
          <CodeEditor v-model="importYaml" language="yaml" min-height="360px" />
        </div>

        <div v-if="importError" class="form-error yaml-error">{{ importError }}</div>

        <div v-if="importResult" class="import-result">
          <span v-if="importResult.created" class="badge badge-success"
            >{{ importResult.created }} created</span
          >
          <span v-if="importResult.updated" class="badge badge-info"
            >{{ importResult.updated }} updated</span
          >
          <div v-if="importResult.errors?.length" class="import-errors">
            <div v-for="(err, i) in importResult.errors" :key="i" class="form-error">
              {{ err }}
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeImportModal">
            Cancel
          </button>
          <button
            type="button"
            class="btn btn-primary"
            :disabled="importing"
            @click="handleImport"
          >
            {{ importing ? "Importing..." : "Import" }}
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { routesApi, type Route, type ImportResult } from '@/api/routes'
import { useConfirm } from '@/composables/useConfirm'
import { useNotificationStore } from '@/stores/notification'
import Modal from '@/components/Modal.vue'
import Pagination from '@/components/Pagination.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import EmptyState from '@/components/EmptyState.vue'
import RouteFormModal from '@/components/RouteFormModal.vue'

const { confirm } = useConfirm()
const notify = useNotificationStore()

const currentRoute = useRoute()
const router = useRouter()

// Restore the page from the URL (1-based) so the list survives refresh / back-navigation.
function initialPage(): number {
  const p = currentRoute.query.page
  if (typeof p === 'string') {
    const n = parseInt(p, 10)
    if (!isNaN(n) && n > 0) return n - 1
  }
  return 0
}

/* ── State ── */
const loading = ref(true)
const page = ref(initialPage())
const pageable = ref({ current_page: 0, total_pages: 1, total_elements: 0, size: 20, empty: true })
const routes = ref<Route[]>([])
const modalOpen = ref(false)
const editing = ref<Route | null>(null)
const search = ref(typeof currentRoute.query.q === 'string' ? currentRoute.query.q : '')
let searchDebounceTimer: ReturnType<typeof setTimeout> | null = null

watch(search, () => {
  if (searchDebounceTimer) clearTimeout(searchDebounceTimer)
  searchDebounceTimer = setTimeout(() => {
    page.value = 0
    syncQuery()
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

/* ── Modal open / close ── */
function openCreate() {
  editing.value = null
  modalOpen.value = true
}

function openEdit(route: Route) {
  editing.value = route
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editing.value = null
}

// Refresh the list after the form modal saves a route.
function onSaved() {
  fetchRoutes()
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
    hosts: ["api.example.com"]
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
// Mirror the current page (1-based) and search term into the URL query.
function syncQuery() {
  const query: Record<string, string> = {}
  if (page.value > 0) query.page = String(page.value + 1)
  if (search.value) query.q = search.value
  router.replace({ query })
}

function goToPage(p: number) {
  page.value = p
  syncQuery()
  fetchRoutes()
}

async function fetchRoutes() {
  loading.value = true
  try {
    const res = await routesApi.list(page.value, 20, search.value)
    routes.value = res.data.data || []
    pageable.value = res.data.pageable
  } catch {
    // Error handled by API interceptor
  } finally {
    loading.value = false
  }
}

async function openEditFromQuery() {
  const editId = currentRoute.query.edit
  if (!editId) return
  try {
    const res = await routesApi.get(Number(editId))
    openEdit(res.data)
  } catch {
    notify.error('Failed to load route for editing')
  }
  // Drop the edit param (keeping page/search) so closing/refreshing doesn't re-open the modal.
  syncQuery()
}

onMounted(async () => {
  await fetchRoutes()
  await openEditFromQuery()
})
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
  0% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.6);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(34, 197, 94, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);
  }
}

@keyframes pulse-red {
  0% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.6);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(239, 68, 68, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0);
  }
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
