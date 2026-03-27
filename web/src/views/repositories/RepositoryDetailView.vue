<template>
  <div>
    <div v-if="loading" class="loading-page">
      <div class="spinner"></div>
    </div>

    <div v-else-if="repo">
      <div class="detail-header">
        <router-link to="/repositories" class="back-link">
          <svg width="20" height="20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </router-link>
        <h1>{{ repo.name }}</h1>
        <span :class="['badge', statusBadge(repo.status)]">{{ repo.status }}</span>
        <button class="btn btn-secondary btn-sm" :disabled="syncing" @click="syncRepo" style="margin-left: auto;">
          {{ syncing ? 'Syncing...' : 'Sync Now' }}
        </button>
      </div>

      <div class="detail-grid">
        <!-- Details Card -->
        <div class="card">
          <div class="card-header"><h2>Details</h2></div>
          <div class="card-body">
            <dl class="detail-list">
              <div class="detail-item">
                <dt>URL</dt>
                <dd class="text-mono">{{ repo.url }}</dd>
              </div>
              <div class="detail-item">
                <dt>Branch</dt>
                <dd>{{ repo.branch }}</dd>
              </div>
              <div class="detail-item">
                <dt>Authentication</dt>
                <dd>{{ repo.authType || 'None' }}{{ repo.hasAuth ? ' (configured)' : '' }}</dd>
              </div>
              <div v-if="repo.lastCommit" class="detail-item">
                <dt>Last Commit</dt>
                <dd class="text-mono">{{ repo.lastCommit }}</dd>
              </div>
              <div v-if="repo.lastSyncedAt" class="detail-item">
                <dt>Last Synced</dt>
                <dd>{{ new Date(repo.lastSyncedAt).toLocaleString() }}</dd>
              </div>
              <div v-if="repo.statusMessage" class="detail-item">
                <dt>Error</dt>
                <dd class="text-danger">{{ repo.statusMessage }}</dd>
              </div>
              <div class="detail-item">
                <dt>Webhook URL</dt>
                <dd class="text-mono" style="font-size: 12px; word-break: break-all;">
                  {{ webhookUrl }}
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- File Browser Card -->
        <div class="card">
          <div class="card-header">
            <h2>Browse Files</h2>
          </div>
          <div class="card-body">
            <div v-if="repo.status !== 'synced'" class="empty-state" style="padding: 24px">
              <p>Repository must be synced before browsing.</p>
            </div>
            <template v-else>
              <!-- Breadcrumb -->
              <div class="browse-breadcrumb">
                <span class="browse-crumb" :class="{ active: !currentPath }" @click="browseTo('')">/</span>
                <template v-for="(segment, i) in pathSegments" :key="i">
                  <span class="browse-sep">/</span>
                  <span
                    class="browse-crumb"
                    :class="{ active: i === pathSegments.length - 1 }"
                    @click="browseTo(pathSegments.slice(0, i + 1).join('/'))"
                  >{{ segment }}</span>
                </template>
              </div>

              <div v-if="browseLoading" class="empty-state" style="padding: 24px">
                <div class="spinner" style="width: 20px; height: 20px;"></div>
              </div>
              <div v-else class="file-list">
                <div v-if="currentPath" class="file-item" @click="browseTo(parentPath)">
                  <svg width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                  </svg>
                  <span>..</span>
                </div>
                <div v-for="entry in entries" :key="entry.path" class="file-item" @click="entry.isDir ? browseTo(entry.path) : null" :class="{ clickable: entry.isDir }">
                  <svg v-if="entry.isDir" width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                  </svg>
                  <svg v-else width="16" height="16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z" /><polyline points="14 2 14 8 20 8" />
                  </svg>
                  <span>{{ entry.name }}</span>
                </div>
                <div v-if="entries.length === 0" class="empty-state" style="padding: 16px">
                  <p>Empty directory</p>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { repositoriesApi, type Repository, type BrowseEntry } from '@/api/repositories'

const props = defineProps<{ id: string }>()

const loading = ref(true)
const syncing = ref(false)
const browseLoading = ref(false)
const repo = ref<Repository | null>(null)
const entries = ref<BrowseEntry[]>([])
const currentPath = ref('')

const webhookUrl = computed(() => {
  if (!repo.value) return ''
  return `${window.location.origin}/api/v1/repositories/${repo.value.id}/webhook`
})

const pathSegments = computed(() => {
  if (!currentPath.value) return []
  return currentPath.value.split('/').filter(Boolean)
})

const parentPath = computed(() => {
  const segs = pathSegments.value
  if (segs.length <= 1) return ''
  return segs.slice(0, -1).join('/')
})

function statusBadge(status: string): string {
  const map: Record<string, string> = {
    synced: 'badge-success',
    error: 'badge-danger',
    pending: 'badge-warning',
  }
  return map[status] || 'badge-neutral'
}

async function browseTo(path: string) {
  if (!repo.value) return
  currentPath.value = path
  browseLoading.value = true
  try {
    const res = await repositoriesApi.browse(repo.value.id, path || undefined)
    entries.value = res.data
  } catch {
    entries.value = []
  } finally {
    browseLoading.value = false
  }
}

async function syncRepo() {
  if (!repo.value) return
  syncing.value = true
  try {
    await repositoriesApi.sync(repo.value.id)
    const res = await repositoriesApi.get(repo.value.id)
    repo.value = res.data
    await browseTo(currentPath.value)
  } catch {
    // handle error
  } finally {
    syncing.value = false
  }
}

onMounted(async () => {
  try {
    const res = await repositoriesApi.get(props.id)
    repo.value = res.data
    if (repo.value?.status === 'synced') {
      await browseTo('')
    }
  } catch {
    // handle error
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.detail-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
}
.detail-header h1 {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
}
.back-link {
  display: flex;
  align-items: center;
  color: var(--text-muted);
  transition: color var(--transition);
}
.back-link:hover { color: var(--text-secondary); }

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 24px;
}
@media (max-width: 768px) {
  .detail-grid { grid-template-columns: 1fr; }
}

.detail-list { display: flex; flex-direction: column; gap: 16px; }
.detail-item dt { font-size: 12px; color: var(--text-muted); margin-bottom: 2px; }
.detail-item dd { font-size: 14px; font-weight: 500; color: var(--text-primary); }
.text-danger { color: var(--color-danger, #ef4444); }

.browse-breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
  font-size: 13px;
  font-family: var(--font-mono, monospace);
}
.browse-crumb {
  color: var(--primary-600);
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
}
.browse-crumb:hover { background: var(--bg-hover); }
.browse-crumb.active { color: var(--text-primary); cursor: default; }
.browse-crumb.active:hover { background: none; }
.browse-sep { color: var(--text-muted); }

.file-list { display: flex; flex-direction: column; }
.file-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: var(--radius-sm, 4px);
  font-size: 13px;
  color: var(--text-primary);
  transition: background var(--transition);
}
.file-item.clickable { cursor: pointer; }
.file-item.clickable:hover, .file-item:first-child:hover { background: var(--bg-hover); cursor: pointer; }
.file-item svg { color: var(--text-muted); flex-shrink: 0; }
</style>
