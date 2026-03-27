<template>
  <div class="login-wrapper">
    <div class="card login-card">
      <div class="card-body" style="text-align: center; padding: 48px 32px;">
        <div v-if="error" class="login-error">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" />
          </svg>
          {{ error }}
        </div>
        <div v-else>
          <div class="spinner" style="margin: 0 auto 16px;"></div>
          <p style="color: var(--text-muted); font-size: 14px;">Completing sign in...</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const route = useRoute()
const authStore = useAuthStore()
const error = ref('')

onMounted(async () => {
  // The backend OAuth callback returns JSON with the JWT.
  // We get here after the backend redirects, or the frontend handles
  // the callback by reading query params and exchanging them.
  // Since our backend returns JSON at /api/v1/auth/oauth/callback,
  // we need to call it from here with the code and state params.

  const code = route.query.code as string
  const state = route.query.state as string

  if (!code || !state) {
    error.value = 'Missing OAuth callback parameters'
    return
  }

  try {
    const { default: api } = await import('@/api/client')
    const res = await api.get('/auth/oauth/callback', { params: { code, state } })
    const data = res.data as {
      access_token: string
      user: { id: string; email: string; name: string; role: string }
    }
    authStore.loginWithOAuth(data.access_token, data.user)
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { message?: string } } }
    error.value = axiosErr.response?.data?.message || 'OAuth login failed'
    setTimeout(() => router.push('/auth/login'), 3000)
  }
})
</script>

<style scoped>
.login-wrapper {
  width: 100%;
  max-width: 380px;
}

.login-card {
  overflow: hidden;
}

.login-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: var(--radius);
  background: var(--danger-50);
  color: var(--danger-600);
  font-size: 13px;
  font-weight: 500;
}
</style>
