<template>
  <div>
    <div class="page-header">
      <h1>OAuth Provider</h1>
    </div>

    <div v-if="loading" class="card card-body loading-page">
      <div class="spinner"></div>
    </div>

    <div v-else class="card card-body">
      <p v-if="!provider && !editing" class="text-muted" style="margin-bottom: 16px;">
        No OAuth provider configured. Users can only sign in with email and password.
      </p>

      <!-- Display mode -->
      <div v-if="provider && !editing" class="provider-display">
        <div class="provider-header">
          <div>
            <h3 class="provider-name">{{ provider.display_name || provider.name }}</h3>
            <span :class="['badge', provider.enabled ? 'badge-success' : 'badge-danger']">
              {{ provider.enabled ? 'Enabled' : 'Disabled' }}
            </span>
          </div>
          <div class="provider-actions">
            <button class="btn btn-secondary btn-sm" @click="startEdit">Edit</button>
            <button class="btn btn-ghost btn-sm action-delete" @click="confirmDelete">Delete</button>
          </div>
        </div>

        <div class="detail-grid">
          <div class="detail-item">
            <span class="detail-label">Name</span>
            <span class="detail-value">{{ provider.name }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Client ID</span>
            <code class="detail-value text-mono">{{ provider.client_id }}</code>
          </div>
          <div class="detail-item">
            <span class="detail-label">Auth URL</span>
            <span class="detail-value text-sm">{{ provider.auth_url }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Token URL</span>
            <span class="detail-value text-sm">{{ provider.token_url }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">User Info URL</span>
            <span class="detail-value text-sm">{{ provider.user_info_url }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Scopes</span>
            <span class="detail-value">
              <span v-for="s in provider.scopes" :key="s" class="badge badge-info" style="margin-right: 4px;">{{ s }}</span>
              <span v-if="!provider.scopes?.length" class="text-muted">None</span>
            </span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Field Mappings</span>
            <span class="detail-value text-sm text-mono">
              id={{ provider.user_id_field }}, email={{ provider.email_field }},
              name={{ provider.name_field }}, avatar={{ provider.avatar_field }}
            </span>
          </div>
        </div>
      </div>

      <!-- Edit / Create form -->
      <form v-if="editing || !provider" @submit.prevent="handleSave">
        <h3 v-if="!provider" style="margin-bottom: 16px;">Configure OAuth Provider</h3>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label">Provider Name</label>
            <input v-model="form.name" required class="form-input" placeholder="e.g. keycloak, gitea, azure" />
            <span class="form-hint">Identifier used in URLs</span>
          </div>
          <div class="form-group">
            <label class="form-label">Display Name</label>
            <input v-model="form.display_name" class="form-input" placeholder="e.g. Company SSO" />
            <span class="form-hint">Shown on the login button</span>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label">Client ID</label>
            <input v-model="form.client_id" required class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">Client Secret</label>
            <input v-model="form.client_secret" :required="!provider" type="password" class="form-input" :placeholder="provider ? '(unchanged)' : ''" />
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">Authorization URL</label>
          <input v-model="form.auth_url" required type="url" class="form-input" placeholder="https://provider.example.com/authorize" />
        </div>
        <div class="form-group">
          <label class="form-label">Token URL</label>
          <input v-model="form.token_url" required type="url" class="form-input" placeholder="https://provider.example.com/token" />
        </div>
        <div class="form-group">
          <label class="form-label">User Info URL</label>
          <input v-model="form.user_info_url" required type="url" class="form-input" placeholder="https://provider.example.com/userinfo" />
        </div>
        <div class="form-group">
          <label class="form-label">Scopes</label>
          <input v-model="form.scopes_raw" class="form-input" placeholder="openid email profile" />
          <span class="form-hint">Space-separated list of OAuth scopes</span>
        </div>

        <details class="field-mappings">
          <summary>Field Mappings (advanced)</summary>
          <div class="form-row" style="margin-top: 12px;">
            <div class="form-group">
              <label class="form-label">User ID Field</label>
              <input v-model="form.user_id_field" class="form-input" placeholder="sub" />
            </div>
            <div class="form-group">
              <label class="form-label">Email Field</label>
              <input v-model="form.email_field" class="form-input" placeholder="email" />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Name Field</label>
              <input v-model="form.name_field" class="form-input" placeholder="name" />
            </div>
            <div class="form-group">
              <label class="form-label">Avatar Field</label>
              <input v-model="form.avatar_field" class="form-input" placeholder="picture" />
            </div>
          </div>
        </details>

        <div class="form-group" style="margin-top: 16px;">
          <label class="checkbox-label">
            <input v-model="form.enabled" type="checkbox" />
            Enabled
          </label>
        </div>

        <div class="form-actions">
          <button v-if="editing" type="button" class="btn btn-secondary" @click="cancelEdit">Cancel</button>
          <button type="submit" class="btn btn-primary" :disabled="saving">
            {{ saving ? 'Saving...' : 'Save' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { oauthApi, type OAuthProviderDetail } from '@/api/oauth'
import { useConfirm } from '@/composables/useConfirm'

const { confirm } = useConfirm()

const loading = ref(true)
const saving = ref(false)
const editing = ref(false)
const provider = ref<OAuthProviderDetail | null>(null)

const form = reactive({
  name: '',
  display_name: '',
  client_id: '',
  client_secret: '',
  auth_url: '',
  token_url: '',
  user_info_url: '',
  scopes_raw: '',
  user_id_field: 'sub',
  email_field: 'email',
  name_field: 'name',
  avatar_field: 'picture',
  enabled: true,
})

function populateForm(p: OAuthProviderDetail) {
  form.name = p.name
  form.display_name = p.display_name
  form.client_id = p.client_id
  form.client_secret = ''
  form.auth_url = p.auth_url
  form.token_url = p.token_url
  form.user_info_url = p.user_info_url
  form.scopes_raw = (p.scopes ?? []).join(' ')
  form.user_id_field = p.user_id_field || 'sub'
  form.email_field = p.email_field || 'email'
  form.name_field = p.name_field || 'name'
  form.avatar_field = p.avatar_field || 'picture'
  form.enabled = p.enabled
}

function startEdit() {
  if (provider.value) populateForm(provider.value)
  editing.value = true
}

function cancelEdit() {
  editing.value = false
}

async function handleSave() {
  saving.value = true
  try {
    const scopes = form.scopes_raw.split(/\s+/).filter(Boolean)

    // When editing, if client_secret is empty, re-use existing
    // The backend should handle empty secret as "keep existing"
    const secret = form.client_secret || (provider.value ? '__unchanged__' : '')

    const res = await oauthApi.saveProvider({
      name: form.name,
      display_name: form.display_name,
      client_id: form.client_id,
      client_secret: secret,
      auth_url: form.auth_url,
      token_url: form.token_url,
      user_info_url: form.user_info_url,
      scopes,
      user_id_field: form.user_id_field,
      email_field: form.email_field,
      name_field: form.name_field,
      avatar_field: form.avatar_field,
      enabled: form.enabled,
    })
    provider.value = res.data
    editing.value = false
  } catch { /* handled */ } finally {
    saving.value = false
  }
}

async function confirmDelete() {
  const confirmed = await confirm({
    title: 'Delete OAuth Provider',
    message: 'Are you sure you want to remove the OAuth provider? Users who signed in via OAuth will need to reset their password to continue using local login.',
    confirmText: 'Delete',
    variant: 'danger',
  })
  if (!confirmed) return
  try {
    await oauthApi.deleteProvider()
    provider.value = null
    editing.value = false
  } catch { /* handled */ }
}

async function fetchProvider() {
  loading.value = true
  try {
    const res = await oauthApi.getProvider()
    provider.value = res.data ?? null
  } catch {
    provider.value = null
  } finally {
    loading.value = false
  }
}

onMounted(fetchProvider)
</script>

<style scoped>
.provider-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}
.provider-header > div:first-child {
  display: flex;
  align-items: center;
  gap: 12px;
}
.provider-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}
.provider-actions {
  display: flex;
  gap: 8px;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.detail-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}
.detail-value {
  font-size: 14px;
  color: var(--text-primary);
  word-break: break-all;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.field-mappings {
  margin-top: 16px;
  padding: 12px 16px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
}
.field-mappings summary {
  cursor: pointer;
  font-weight: 500;
  color: var(--text-secondary);
  font-size: 14px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
}

.action-delete {
  color: var(--danger-500);
}
.action-delete:hover {
  color: var(--danger-700);
}
</style>
