import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { instancesApi, type Instance } from '@/api/instances'
import api from '@/api/client'

export const useInstanceStore = defineStore('instance', () => {
  const instances = ref<Instance[]>([])

  const stored = localStorage.getItem('goma_instance_id')
  const currentInstanceId = ref<number | null>(
    stored ? Number(stored) : null
  )

  const currentInstance = computed(() =>
    instances.value.find((i) => i.id === currentInstanceId.value) ?? null
  )

  const contextLabel = computed(() =>
    currentInstance.value?.name ?? 'All Instances'
  )

  const isGlobal = computed(() => currentInstanceId.value === null)

  function setInstance(id: number | null) {
    currentInstanceId.value = id
    if (id !== null) {
      localStorage.setItem('goma_instance_id', String(id))
    } else {
      localStorage.removeItem('goma_instance_id')
    }
  }

  async function fetchInstances() {
    try {
      const res = await instancesApi.list()
      // list() returns a flat array (not paginated)
      instances.value = Array.isArray(res.data) ? res.data : (res.data as any).data || []
      // Validate stored instance still exists
      if (
        currentInstanceId.value !== null &&
        !instances.value.find((i) => i.id === currentInstanceId.value)
      ) {
        setInstance(null)
      }
    } catch {
      // silently fail
    }
  }

  return {
    instances,
    currentInstanceId,
    currentInstance,
    contextLabel,
    isGlobal,
    setInstance,
    fetchInstances,
  }
})

// Axios interceptor: inject X-Goma-Instance-Id header
api.interceptors.request.use((config) => {
  const stored = localStorage.getItem('goma_instance_id')
  if (stored) {
    config.headers['X-Goma-Instance-Id'] = stored
  }
  return config
})
