<template>
  <div>
    <div class="page-header">
      <h1>Repositories</h1>
      <button class="btn btn-primary" @click="showCreate = true">
        <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Repository
      </button>
    </div>

    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <div v-else-if="repos.length === 0" class="card">
      <div class="empty-state">
        <svg width="48" height="48" fill="none" stroke="currentColor" viewBox="0 0 24 24" style="margin: 0 auto 12px;">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22" />
        </svg>
        <h3>No repositories</h3>
        <p>Add a git repository to sync gateway configs from.</p>
        <button class="btn btn-primary" @click="showCreate = true">Add Repository</button>
      </div>
    </div>

    <template v-else>
      <div class="card">
        <div class="card-body" style="padding: 0">
          <table class="data-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>URL</th>
                <th>Branch</th>
                <th>Status</th>
                <th>Last Sync</th>
                <th>Commit</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="repo in repos" :key="repo.id">
                <td>
                  <router-link :to="`/repositories/${repo.id}`" class="link-primary">{{ repo.name }}</router-link>
                </td>
                <td class="text-mono text-muted" style="max-width: 300px; overflow: hidden; text-overflow: ellipsis;">{{ repo.url }}</td>
                <td>{{ repo.branch }}</td>
                <td>
                  <span :class="['badge', statusBadge(repo.status)]">{{ repo.status }}</span>
                </td>
                <td class="text-muted">{{ repo.lastSyncedAt ? relativeTime(repo.lastSyncedAt) : 'Never' }}</td>
                <td class="text-mono">{{ repo.lastCommit || '-' }}</td>
                <td>
                  <div style="display: flex; gap: 6px; justify-content: flex-end;">
                    <button class="btn btn-secondary btn-sm" :disabled="syncing === repo.id" @click="syncRepo(repo)">
                      {{ syncing === repo.id ? 'Syncing...' : 'Sync' }}
                    </button>
                    <button class="btn btn-secondary btn-sm" @click="editRepo(repo)">Edit</button>
                    <button class="btn btn-danger btn-sm" @click="deleteRepo(repo)">Delete</button>
                  </div>
                </td>
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
            <h3>{{ editing ? 'Edit Repository' : 'Add Repository' }}</h3>
            <button class="btn-ghost btn-icon" @click="closeModal">
              <svg width="18" height="18" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <form @submit.prevent="handleSubmit">
            <div class="modal-body">
              <div class="form-group">
                <label class="form-label" for="repo-name">Name</label>
                <input id="repo-name" v-model="form.name" class="form-input" required placeholder="gateway-configs" />
              </div>
              <div class="form-group">
                <label class="form-label" for="repo-url">Git URL</label>
                <input id="repo-url" v-model="form.url" class="form-input" required placeholder="https://github.com/org/configs.git" />
              </div>
              <div class="form-group">
                <label class="form-label" for="repo-branch">Branch</label>
                <input id="repo-branch" v-model="form.branch" class="form-input" placeholder="main" />
              </div>
              <div class="form-section-divider">
                <span class="form-section-title">Authentication</span>
              </div>
              <div class="form-grid">
                <div class="form-group">
                  <label class="form-label" for="repo-auth-type">Auth Type</label>
                  <select id="repo-auth-type" v-model="form.authType" class="form-select">
                    <option value="">None (public)</option>
                    <option value="token">Personal Access Token</option>
                    <option value="basic">Basic Auth</option>
                  </select>
                </div>
                <div v-if="form.authType" class="form-group">
                  <label class="form-label" for="repo-auth-value">Credentials</label>
                  <input
                    id="repo-auth-value"
                    v-model="form.authValue"
                    class="form-input"
                    type="password"
                    :placeholder="editing?.hasAuth ? '••••••• (leave empty to keep)' : (form.authType === 'basic' ? 'user:password' : 'token')"
                  />
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
              <button type="submit" class="btn btn-primary" :disabled="saving">
                {{ saving ? 'Saving...' : (editing ? 'Update' : 'Add') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { repositoriesApi, type Repository, type RepositoryCreateRequest } from '@/api/repositories'
import { useConfirm } from '@/composables/useConfirm'

const { confirm } = useConfirm()

const loading = ref(true)
const saving = ref(false)
const syncing = ref<number | null>(null)
const repos = ref<Repository[]>([])
const showCreate = ref(false)
const editing = ref<Repository | null>(null)

const form = reactive<RepositoryCreateRequest>({
  name: '',
  url: '',
  branch: 'main',
  authType: '',
  authValue: '',
})

function resetForm() {
  form.name = ''
  form.url = ''
  form.branch = 'main'
  form.authType = ''
  form.authValue = ''
}

function closeModal() {
  showCreate.value = false
  editing.value = null
  resetForm()
}

function editRepo(repo: Repository) {
  editing.value = repo
  form.name = repo.name
  form.url = repo.url
  form.branch = repo.branch
  form.authType = repo.authType || ''
  form.authValue = ''
}

async function handleSubmit() {
  saving.value = true
  try {
    if (editing.value) {
      await repositoriesApi.update(editing.value.id, form)
    } else {
      await repositoriesApi.create(form)
    }
    closeModal()
    await fetchRepos()
  } catch {
    // handle error
  } finally {
    saving.value = false
  }
}

async function syncRepo(repo: Repository) {
  syncing.value = repo.id
  try {
    await repositoriesApi.sync(repo.id)
    await fetchRepos()
  } catch {
    // handle error
  } finally {
    syncing.value = null
  }
}

async function deleteRepo(repo: Repository) {
  const confirmed = await confirm({
    title: 'Delete Repository',
    message: `Are you sure you want to delete "${repo.name}"? The local clone will be removed.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await repositoriesApi.delete(repo.id)
    await fetchRepos()
  } catch {
    // handle error
  }
}

function statusBadge(status: string): string {
  const map: Record<string, string> = {
    synced: 'badge-success',
    error: 'badge-danger',
    pending: 'badge-warning',
  }
  return map[status] || 'badge-neutral'
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
  return Math.floor(diffHr / 24) + 'd ago'
}

async function fetchRepos() {
  loading.value = true
  try {
    const res = await repositoriesApi.list()
    repos.value = res.data
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
}

onMounted(fetchRepos)
</script>

<style scoped>
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
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
</style>
