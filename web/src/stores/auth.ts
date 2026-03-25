import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi, type LoginRequest } from '@/api/auth'
import router from '@/router'

interface User {
  id: string
  email: string
  name: string
  role: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('access_token') || '')
  const user = ref<User | null>(JSON.parse(localStorage.getItem('user') || 'null'))
  const loading = ref(false)
  const error = ref('')

  const isAuthenticated = computed(() => !!token.value)

  async function login(data: LoginRequest) {
    loading.value = true
    error.value = ''
    try {
      const res = await authApi.login(data)
      console.log('[auth] login response:', JSON.stringify(res.data))
      token.value = res.data.access_token
      user.value = res.data.user
      localStorage.setItem('access_token', res.data.access_token)
      localStorage.setItem('user', JSON.stringify(res.data.user))
      console.log('[auth] token saved:', localStorage.getItem('access_token')?.substring(0, 20) + '...')
      router.push('/')
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { message?: string } } }
      error.value = axiosErr.response?.data?.message || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('user')
    router.push('/auth/login')
  }

  return { token, user, loading, error, isAuthenticated, login, logout }
})
