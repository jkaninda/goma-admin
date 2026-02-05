import type { HttpMethod } from "~/types/http"

export function useApiFetch() {
    const { $api } = useNuxtApp()

    async function apiFetch(
        url: string,
        method: HttpMethod,
        options: Record<string, any> = {}
    ) {
        return await $api(url, {
            method,
            ...options
        })
    }

    return { apiFetch }
}
