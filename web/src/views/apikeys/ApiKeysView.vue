<template>
  <div>
    <div class="page-header">
      <div class="page-header-left">
        <h1>API Keys</h1>
        <span v-if="instanceStore.currentInstance" class="badge badge-info">
          Scoped to: {{ instanceStore.currentInstance.name }}
        </span>
      </div>
      <button class="btn btn-primary" @click="openCreate">Create Key</button>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <EmptyState
      v-else-if="keys.length === 0"
      title="No API keys"
      description="Create an API key to authenticate programmatic access."
    >
      <template #action>
        <button class="btn btn-primary" @click="openCreate">Create Key</button>
      </template>
    </EmptyState>

    <div v-else class="card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Key Prefix</th>
              <th>Created</th>
              <th>Last Used</th>
              <th>Expires</th>
              <th>Allowed IPs</th>
              <th>Status</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="key in keys" :key="key.id">
              <td class="cell-name">{{ key.name }}</td>
              <td><code class="text-mono">{{ key.key_prefix }}</code></td>
              <td>{{ formatDate(key.created_at) }}</td>
              <td>{{ key.last_used_at ? formatDate(key.last_used_at) : 'Never' }}</td>
              <td>{{ key.expires_at ? formatDate(key.expires_at) : 'Never' }}</td>
              <td>
                <span v-if="key.allowed_ips && key.allowed_ips.length" class="text-sm text-mono">
                  {{ key.allowed_ips.join(', ') }}
                </span>
                <span v-else class="text-muted text-sm">Any</span>
              </td>
              <td>
                <span :class="['badge', statusBadge(key)]">{{ statusLabel(key) }}</span>
              </td>
              <td class="text-right">
                <button
                  v-if="isActive(key)"
                  class="btn btn-ghost btn-sm action-revoke"
                  @click="confirmRevoke(key)"
                >Revoke</button>
                <button
                  v-else
                  class="btn btn-ghost btn-sm action-delete"
                  @click="confirmDelete(key)"
                >Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create Modal -->
    <Modal
      :show="createModalOpen"
      title="Create API Key"
      @close="closeCreateModal"
    >
      <div class="modal-body">
        <form @submit.prevent="handleCreate">
          <div class="form-group">
            <label class="form-label">Name</label>
            <input
              v-model="form.name"
              required
              class="form-input"
              placeholder="e.g. CI/CD Pipeline"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Expiration</label>
            <select v-model="form.expiration" class="form-select">
              <option value="90">Default (90 days)</option>
              <option value="30">30 days</option>
              <option value="60">60 days</option>
              <option value="90">90 days</option>
              <option value="180">180 days</option>
              <option value="365">365 days</option>
              <option value="never">Never</option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Allowed IPs</label>
            <textarea
              v-model="form.allowedIPs"
              class="form-textarea"
              rows="3"
              placeholder="192.168.1.0/24, 10.0.0.1"
            ></textarea>
            <span class="form-hint">
              Optional. Comma or newline separated. Leave empty to allow any IP.
            </span>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeCreateModal">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? 'Creating...' : 'Create' }}
            </button>
          </div>
        </form>
      </div>
    </Modal>

    <!-- Key Created Modal -->
    <Modal
      :show="keyCreatedModalOpen"
      title="Key Created"
      @close="closeKeyCreatedModal"
    >
      <div class="modal-body">
        <p class="key-warning">
          Save this key securely. It will not be shown again.
        </p>
        <div class="code-block">
          <code class="text-mono">{{ createdKeyRaw }}</code>
          <button class="btn btn-ghost btn-sm copy-btn" @click="copyKey">
            {{ copied ? 'Copied!' : 'Copy' }}
          </button>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-primary" @click="closeKeyCreatedModal">Done</button>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { apiKeysApi, type ApiKey, type ApiKeyCreateResponse } from '@/api/apikeys'
import { useConfirm } from '@/composables/useConfirm'
import Modal from '@/components/Modal.vue'
import EmptyState from '@/components/EmptyState.vue'
import { useInstanceStore } from '@/stores/instance'

const { confirm } = useConfirm()

/* -- Stores -- */
const instanceStore = useInstanceStore()

/* -- State -- */
const loading = ref(true)
const saving = ref(false)
const keys = ref<ApiKey[]>([])
const createModalOpen = ref(false)
const keyCreatedModalOpen = ref(false)
const createdKeyRaw = ref('')
const copied = ref(false)

/* -- Form -- */
const form = reactive({
  name: '',
  expiration: '90',
  allowedIPs: '',
})

/* -- Helpers -- */
function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString()
}

function isExpired(key: ApiKey): boolean {
  if (!key.expires_at) return false
  return new Date(key.expires_at) < new Date()
}

function isActive(key: ApiKey): boolean {
  return !key.revoked && !isExpired(key)
}

function statusLabel(key: ApiKey): string {
  if (key.revoked) return 'Revoked'
  if (isExpired(key)) return 'Expired'
  return 'Active'
}

function statusBadge(key: ApiKey): string {
  if (key.revoked) return 'badge-danger'
  if (isExpired(key)) return 'badge-warning'
  return 'badge-success'
}

function parseAllowedIPs(raw: string): string[] {
  return raw
    .split(/[,\n]/)
    .map(s => s.trim())
    .filter(Boolean)
}

/* -- Modal open / close -- */
function resetForm() {
  form.name = ''
  form.expiration = '90'
  form.allowedIPs = ''
}

function openCreate() {
  resetForm()
  createModalOpen.value = true
}

function closeCreateModal() {
  createModalOpen.value = false
  resetForm()
}

function closeKeyCreatedModal() {
  keyCreatedModalOpen.value = false
  createdKeyRaw.value = ''
  copied.value = false
}

/* -- Copy key -- */
async function copyKey() {
  try {
    await navigator.clipboard.writeText(createdKeyRaw.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // Fallback: select text for manual copy
  }
}

/* -- Create -- */
async function handleCreate() {
  saving.value = true
  const allowedIPs = parseAllowedIPs(form.allowedIPs)
  const expiresInDays = form.expiration === 'never' ? undefined : parseInt(form.expiration, 10)

  try {
    const res = await apiKeysApi.create(form.name, allowedIPs.length ? allowedIPs : undefined, expiresInDays)
    const created: ApiKeyCreateResponse = res.data
    closeCreateModal()
    createdKeyRaw.value = created.key
    keyCreatedModalOpen.value = true
    await fetchKeys()
  } catch {
    // Error handled by API interceptor
  } finally {
    saving.value = false
  }
}

/* -- Revoke -- */
async function confirmRevoke(key: ApiKey) {
  const confirmed = await confirm({
    title: 'Revoke API Key',
    message: `Are you sure you want to revoke "${key.name}"? This key will immediately stop working and cannot be reactivated.`,
    confirmText: 'Revoke',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await apiKeysApi.revoke(key.id)
    await fetchKeys()
  } catch {
    // Error handled by API interceptor
  }
}

/* -- Delete -- */
async function confirmDelete(key: ApiKey) {
  const confirmed = await confirm({
    title: 'Delete API Key',
    message: `Are you sure you want to permanently delete "${key.name}"? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await apiKeysApi.delete(key.id)
    await fetchKeys()
  } catch {
    // Error handled by API interceptor
  }
}

/* -- Fetch -- */
async function fetchKeys() {
  loading.value = true
  try {
    const res = await apiKeysApi.list()
    keys.value = res.data.data
  } catch {
    // Error handled by API interceptor
  } finally {
    loading.value = false
  }
}

onMounted(fetchKeys)
</script>

<style scoped>
.page-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cell-name {
  font-weight: 600;
  color: var(--text-primary);
}

.action-revoke {
  color: var(--warning-600);
}
.action-revoke:hover {
  color: var(--warning-700);
}

.action-delete {
  color: var(--danger-500);
}
.action-delete:hover {
  color: var(--danger-700);
}

.key-warning {
  color: var(--warning-600);
  font-weight: 600;
  margin-bottom: 12px;
}

.code-block {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-primary);
  border-radius: 6px;
  word-break: break-all;
}

.copy-btn {
  flex-shrink: 0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding-top: 20px;
}
</style>
