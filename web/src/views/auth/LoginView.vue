<template>
  <div class="login-wrapper">
    <div class="card login-card">
      <div class="card-body">
        <div class="login-header">
          <img src="/logo.png" alt="Goma" class="login-logo" />
          <h1 class="login-title">Goma Admin</h1>
          <p class="login-subtitle">Sign in to your account</p>
        </div>

        <div v-if="authStore.error" class="login-error">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
          </svg>
          {{ authStore.error }}
        </div>

        <form @submit.prevent="handleLogin">
          <div class="form-group">
            <label class="form-label" for="email">Email</label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              class="form-input"
              placeholder="you@example.com"
              required
              autocomplete="email"
              :disabled="authStore.requires2FA"
            />
          </div>

          <div class="form-group">
            <label class="form-label" for="password">Password</label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              class="form-input"
              placeholder="Enter your password"
              required
              autocomplete="current-password"
              :disabled="authStore.requires2FA"
            />
          </div>

          <div v-if="authStore.requires2FA" class="form-group">
            <label class="form-label" for="tfa-code">Authentication Code</label>
            <input
              id="tfa-code"
              v-model="form.two_factor_code"
              type="text"
              class="form-input totp-input"
              placeholder="000000"
              maxlength="6"
              inputmode="numeric"
              autocomplete="one-time-code"
              required
            />
            <small class="form-hint">Enter the 6-digit code from your authenticator app</small>
          </div>

          <button type="submit" class="btn btn-primary login-submit" :disabled="authStore.loading">
            <span v-if="authStore.loading" class="spinner"></span>
            {{ authStore.loading ? 'Signing in...' : (authStore.requires2FA ? 'Verify & Sign in' : 'Sign in') }}
          </button>

          <button
            v-if="authStore.requires2FA"
            type="button"
            class="btn btn-secondary login-submit"
            style="margin-top: 8px"
            @click="reset2FA"
          >
            Back
          </button>
        </form>

        <!-- OAuth section below the form -->
        <template v-if="oauthProvider">
          <div class="login-divider">
            <span>OR</span>
          </div>

          <a :href="oauthProvider.auth_url" class="oauth-btn">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2" /><path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
            Sign in with {{ oauthProvider.display_name || oauthProvider.name }}
          </a>
        </template>

        <p class="login-footer">Contact your administrator for an account.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { oauthApi, type OAuthProviderInfo } from '@/api/oauth'

const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
  two_factor_code: '',
})

const oauthProvider = ref<OAuthProviderInfo | null>(null)

async function handleLogin() {
  try {
    await authStore.login({
      email: form.email,
      password: form.password,
      remember_me: false,
      two_factor_code: form.two_factor_code || undefined,
    })
  } catch {
    // error is handled in the store
  }
}

function reset2FA() {
  authStore.requires2FA = false
  authStore.error = ''
  form.two_factor_code = ''
}

onMounted(async () => {
  try {
    const res = await oauthApi.getProviderInfo()
    if (res.data && res.data.enabled) {
      oauthProvider.value = res.data
    }
  } catch {
    // OAuth not configured — no button shown
  }
})
</script>

<style scoped>
.login-wrapper {
  width: 100%;
  max-width: 400px;
}

.login-card {
  overflow: hidden;
}

.login-card .card-body {
  padding: 40px 36px;
}

.login-header {
  text-align: center;
  margin-bottom: 28px;
}

.login-logo {
  width: 52px;
  height: 52px;
  border-radius: var(--radius);
  margin-bottom: 12px;
}

.login-title {
  font-size: 22px;
  font-weight: 700;
  color: var(--text-primary);
  letter-spacing: -0.02em;
  margin-bottom: 4px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--text-muted);
}

.login-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  margin-bottom: 20px;
  border-radius: var(--radius);
  background: var(--danger-50);
  color: var(--danger-600);
  font-size: 13px;
  font-weight: 500;
}

.login-submit {
  width: 100%;
  margin-top: 8px;
}

.login-submit .spinner {
  width: 16px;
  height: 16px;
  border-width: 2px;
  border-color: rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
}

.login-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 24px 0;
  color: var(--text-muted);
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.login-divider::before,
.login-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-primary);
}

.oauth-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
  padding: 10px 16px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}

.oauth-btn:hover {
  background: var(--bg-secondary);
  border-color: var(--border-secondary);
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 13px;
  color: var(--text-muted);
}

.totp-input {
  font-size: 20px;
  text-align: center;
  letter-spacing: 8px;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
}

.form-hint {
  display: block;
  font-size: 12px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
