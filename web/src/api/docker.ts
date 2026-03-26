import api from './client'

export interface DockerStatus {
  enabled: boolean
  connected: boolean
  swarmMode: boolean
  lastSync: string
  routeCount: number
}

export interface DockerEvent {
  type: string
  message: string
  routeCount?: number
  timestamp: string
}

export const dockerApi = {
  status() {
    return api.get<DockerStatus>('/docker/status')
  },
  sync() {
    return api.post<{ message: string }>('/docker/sync')
  },
  events(onEvent: (event: DockerEvent) => void, onStatus?: (status: DockerStatus) => void): EventSource {
    const token = localStorage.getItem('access_token')
    const url = `/api/v1/docker/events${token ? `?token=${encodeURIComponent(token)}` : ''}`
    const es = new EventSource(url)

    if (onStatus) {
      es.addEventListener('status', (e) => {
        onStatus(JSON.parse(e.data))
      })
    }

    for (const type of ['sync_started', 'sync_completed', 'routes_changed', 'sync_error', 'connected']) {
      es.addEventListener(type, (e) => {
        onEvent(JSON.parse(e.data))
      })
    }

    return es
  },
}
