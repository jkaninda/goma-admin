export default defineNuxtPlugin(() => {
    const config = useRuntimeConfig()

    const api = $fetch.create({
        baseURL: config.public.apiUrl,
        credentials: "include",

        async onRequest({ options }) {
            const auth = useAuthStore()

            const headers = new Headers(options.headers || {})

            headers.set("Accept", "application/json")
            headers.set("Referer", config.public.appUrl)

            if (auth.token) {
                headers.set("Authorization", `Bearer ${auth.token}`)
            }

            if (import.meta.server) {
                const cookie = useRequestHeaders(["cookie"])
                if (cookie.cookie) headers.set("cookie", cookie.cookie)
            }

            options.headers = headers
        },

        async onResponseError({ response }): Promise<void> {
            const auth = useAuthStore()

            if (response.status === 401 && auth.refreshToken) {
                try {
                    await auth.refreshSession()
                } catch {
                    auth.clearSession()
                }
            }
        }
    })

    return { provide: { api } }
})
