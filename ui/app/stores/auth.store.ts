import type { AuthCredentials, User } from "~/types"
export const useAuthStore = defineStore('auth', () => {
    const user = ref<User | null>(null)

    const token = ref<string | null>(null)
    const refreshToken = ref<string | null>(null)

    const { apiFetch } = useApiFetch()

    function setSession(data: any) {
        if (data.accessToken) {
            token.value = data.accessToken
            localStorage.setItem('token', data.accessToken)
        }

        if (data.refreshToken) {
            refreshToken.value = data.refreshToken
            localStorage.setItem('refreshToken', data.refreshToken)
        }
    }

    function clearSession() {
        token.value = null
        refreshToken.value = null
        user.value = null

        localStorage.removeItem('token')
        localStorage.removeItem('refreshToken')
    }


    async function fetchUser() {
        const res:any = await apiFetch('/auth/user', 'GET')
        user.value = res.data.value
    }


    async function login(credentials: AuthCredentials) {
        const res = await apiFetch('/auth/login', 'POST', {
            body: credentials
        })

        setSession(res)
        await fetchUser()
    }
    async function refreshSession() {

        const res = await apiFetch('/auth/refreshToken', 'POST', {
            body: { refreshToken: refreshToken.value }
        })

        setSession(res)
    }

    async function logout() {
        await apiFetch('/auth/logout','POST')
        clearSession()
        await navigateTo('/auth/login')
    }

    return {
        user,
        token,
        refreshToken,
        login,
        logout,
        fetchUser,
        refreshSession,
        clearSession
    }
})
