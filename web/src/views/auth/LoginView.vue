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
              placeholder="admin@example.com"
              required
              autocomplete="email"
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
            />
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="form.remember_me" type="checkbox" />
              Remember me
            </label>
          </div>

          <button type="submit" class="btn btn-primary login-submit" :disabled="authStore.loading">
            <span v-if="authStore.loading" class="spinner"></span>
            {{ authStore.loading ? 'Signing in...' : 'Sign in' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: '',
  remember_me: false,
})

async function handleLogin() {
  try {
    await authStore.login(form)
  } catch {
    // error is handled in the store
  }
}
</script>

<style scoped>
.login-wrapper {
  width: 100%;
  max-width: 380px;
}

.login-card {
  overflow: hidden;
}

.login-card .card-body {
  padding: 36px 32px;
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
  margin-top: 4px;
}

.login-submit .spinner {
  width: 16px;
  height: 16px;
  border-width: 2px;
  border-color: rgba(255, 255, 255, 0.3);
  border-top-color: #ffffff;
}
</style>
