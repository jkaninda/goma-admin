<template>
  <div>
    <div class="page-header">
      <h1>Users</h1>
      <button class="btn btn-primary" @click="openCreate">Create User</button>
    </div>

    <!-- Filters -->
    <div class="card card-body filters-bar">
      <div class="filters-row">
        <input
          v-model="searchQuery"
          type="text"
          class="form-input"
          placeholder="Search by name or email..."
          @input="debouncedFetch"
        />
        <select v-model="roleFilter" class="form-select" @change="fetchUsers">
          <option value="">All roles</option>
          <option value="viewer">Viewer</option>
          <option value="user">User</option>
          <option value="admin">Admin</option>
          <option value="superadmin">Super Admin</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <EmptyState
      v-else-if="users.length === 0"
      title="No users found"
      description="No users match your filters, or create a new user."
    >
      <template #action>
        <button class="btn btn-primary" @click="openCreate">Create User</button>
      </template>
    </EmptyState>

    <div v-else class="card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Email</th>
              <th>Role</th>
              <th>Auth</th>
              <th>2FA</th>
              <th>Status</th>
              <th>Last Login</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="u in users" :key="u.id">
              <td class="cell-name">
                <div class="user-cell">
                  <div v-if="u.avatar" class="user-avatar">
                    <img :src="u.avatar" :alt="u.name" />
                  </div>
                  <div v-else class="user-avatar user-avatar-fallback">
                    {{ u.name?.charAt(0)?.toUpperCase() || '?' }}
                  </div>
                  <span>{{ u.name }}</span>
                </div>
              </td>
              <td>{{ u.email }}</td>
              <td><span :class="['badge', roleBadge(u.role)]">{{ u.role }}</span></td>
              <td>
                <span v-if="u.oauth_provider" class="badge badge-info">{{ u.oauth_provider }}</span>
                <span v-else class="text-muted text-sm">Local</span>
              </td>
              <td>
                <span :class="['badge', u.two_factor_enabled ? 'badge-success' : 'badge-secondary']">
                  {{ u.two_factor_enabled ? 'Enabled' : 'Off' }}
                </span>
              </td>
              <td>
                <span :class="['badge', u.active ? 'badge-success' : 'badge-danger']">
                  {{ u.active ? 'Active' : 'Disabled' }}
                </span>
              </td>
              <td>{{ u.last_login_at ? formatDate(u.last_login_at) : 'Never' }}</td>
              <td class="text-right">
                <button class="btn btn-ghost btn-sm" @click="openEdit(u)">Edit</button>
                <button
                  v-if="u.two_factor_enabled && u.id !== currentUserId"
                  class="btn btn-ghost btn-sm action-delete"
                  @click="confirmDisable2FA(u)"
                >Disable 2FA</button>
                <button
                  v-if="u.id !== currentUserId"
                  class="btn btn-ghost btn-sm action-delete"
                  @click="confirmDelete(u)"
                >Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <Pagination
        v-if="pageable.total_pages > 1"
        :current-page="pageable.current_page"
        :total-pages="pageable.total_pages"
        :total-elements="pageable.total_elements"
        @page-change="goToPage"
      />
    </div>

    <!-- Create / Edit Modal -->
    <Modal
      :show="modalOpen"
      :title="editingUser ? 'Edit User' : 'Create User'"
      @close="closeModal"
    >
      <div class="modal-body">
        <form @submit.prevent="handleSave">
          <div class="form-group">
            <label class="form-label">Name</label>
            <input v-model="form.name" required class="form-input" placeholder="Full name" />
          </div>
          <div class="form-group">
            <label class="form-label">Email</label>
            <input v-model="form.email" type="email" required class="form-input" placeholder="email@example.com" />
          </div>
          <div v-if="!editingUser" class="form-group">
            <label class="form-label">Password</label>
            <input v-model="form.password" type="password" required minlength="6" class="form-input" placeholder="Min 6 characters" />
          </div>
          <div class="form-group">
            <label class="form-label">Role</label>
            <select v-model="form.role" class="form-select">
              <option value="viewer">Viewer</option>
              <option value="user">User</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div v-if="editingUser" class="form-group">
            <label class="checkbox-label">
              <input v-model="form.active" type="checkbox" />
              Active
            </label>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" @click="closeModal">Cancel</button>
            <button type="submit" class="btn btn-primary" :disabled="saving">
              {{ saving ? 'Saving...' : (editingUser ? 'Save' : 'Create') }}
            </button>
          </div>
        </form>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { usersApi, type UserDetail } from '@/api/users'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { useConfirm } from '@/composables/useConfirm'
import Modal from '@/components/Modal.vue'
import Pagination from '@/components/Pagination.vue'
import EmptyState from '@/components/EmptyState.vue'

const { confirm } = useConfirm()
const authStore = useAuthStore()
const notifications = useNotificationStore()
const currentUserId = computed(() => authStore.user?.id)

const loading = ref(true)
const saving = ref(false)
const users = ref<UserDetail[]>([])
const searchQuery = ref('')
const roleFilter = ref('')
const page = ref(1)
const pageable = ref({ current_page: 1, total_pages: 1, total_elements: 0, size: 20, empty: true })

const modalOpen = ref(false)
const editingUser = ref<UserDetail | null>(null)
const form = reactive({ name: '', email: '', password: '', role: 'user', active: true })

let debounceTimer: ReturnType<typeof setTimeout>
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchUsers(), 300)
}

function formatDate(d: string) { return new Date(d).toLocaleDateString() }

function roleBadge(role: string) {
  switch (role) {
    case 'superadmin': return 'badge-danger'
    case 'admin': return 'badge-warning'
    case 'user': return 'badge-primary'
    default: return 'badge-info'
  }
}

function openCreate() {
  editingUser.value = null
  form.name = ''
  form.email = ''
  form.password = ''
  form.role = 'user'
  form.active = true
  modalOpen.value = true
}

function openEdit(u: UserDetail) {
  editingUser.value = u
  form.name = u.name
  form.email = u.email
  form.password = ''
  form.role = u.role
  form.active = u.active
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editingUser.value = null
}

async function handleSave() {
  saving.value = true
  try {
    if (editingUser.value) {
      await usersApi.update(editingUser.value.id, {
        name: form.name,
        email: form.email,
        role: form.role,
        active: form.active,
      })
    } else {
      await usersApi.create({
        name: form.name,
        email: form.email,
        password: form.password,
        role: form.role,
      })
    }
    closeModal()
    await fetchUsers()
  } catch {
    // handled by interceptor
  } finally {
    saving.value = false
  }
}

async function confirmDisable2FA(u: UserDetail) {
  const confirmed = await confirm({
    title: 'Disable 2FA',
    message: `Are you sure you want to disable two-factor authentication for "${u.name}" (${u.email})? They will need to set it up again.`,
    confirmText: 'Disable 2FA',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await usersApi.disable2FA(u.id)
    notifications.success(`2FA disabled for ${u.name}`)
    await fetchUsers()
  } catch {
    notifications.error('Failed to disable 2FA')
  }
}

async function confirmDelete(u: UserDetail) {
  const confirmed = await confirm({
    title: 'Delete User',
    message: `Are you sure you want to delete "${u.name}" (${u.email})? This action cannot be undone.`,
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await usersApi.delete(u.id)
    await fetchUsers()
  } catch { /* handled */ }
}

function goToPage(p: number) {
  page.value = p
  fetchUsers()
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await usersApi.list(page.value, 20, roleFilter.value || undefined, searchQuery.value || undefined)
    users.value = res.data.data ?? []
    pageable.value = res.data.pageable
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

onMounted(fetchUsers)
</script>

<style scoped>
.filters-bar {
  margin-bottom: 16px;
}
.filters-row {
  display: flex;
  gap: 12px;
}
.filters-row .form-input {
  flex: 1;
}
.filters-row .form-select {
  width: 180px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}
.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}
.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.user-avatar-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-100);
  color: var(--primary-700);
  font-weight: 600;
  font-size: 14px;
}

.cell-name {
  font-weight: 600;
  color: var(--text-primary);
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
</style>
