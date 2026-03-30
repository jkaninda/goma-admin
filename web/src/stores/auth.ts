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
  const requires2FA = ref(false)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin' || user.value?.role === 'superadmin')

  async function login(data: LoginRequest) {
    loading.value = true
    error.value = ''
    try {
      const res = await authApi.login(data)
      setSession(res.data.access_token, res.data.user)
      requires2FA.value = false
      router.push('/')
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { message?: string; requires_2fa?: boolean } } }
      if (axiosErr.response?.data?.requires_2fa) {
        requires2FA.value = true
        error.value = ''
      } else {
        error.value = axiosErr.response?.data?.message || 'Login failed'
      }
      throw err
    } finally {
      loading.value = false
    }
  }

  function loginWithOAuth(accessToken: string, userData: User) {
    setSession(accessToken, userData)
    router.push('/')
  }

  function setSession(accessToken: string, userData: User) {
    token.value = accessToken
    user.value = userData
    localStorage.setItem('access_token', accessToken)
    localStorage.setItem('user', JSON.stringify(userData))
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('user')
    router.push('/auth/login')
  }

  return { token, user, loading, error, requires2FA, isAuthenticated, isAdmin, login, loginWithOAuth, logout }
})
