import api from './client'

export interface Instance {
  id: number
  name: string
  environment: string
  description: string
  endpoint: string
  metricsEndpoint: string
  healthEndpoint: string
  version: string
  region: string
  tags: string[]
  lastSeen: string | null
  status: string
  enabled: boolean
  metadata: Record<string, unknown>
  createdAt: string
  updatedAt: string
  routes?: Route[]
}

export interface Route {
  id: number
  name: string
  path: string
  rewrite?: string
  priority: number
  enabled: boolean
  methods: string[]
  hosts: string[]
  target?: string
  disableMetrics: boolean
  backends?: Backend[]
  maintenance?: Maintenance
  healthCheck?: HealthCheck
  security?: Security
  middlewares?: string[]
}

export interface Backend {
  endpoint: string
  weight: number
  exclusive: boolean
}

export interface Maintenance {
  enabled: boolean
  statusCode: number
  message: string
}

export interface HealthCheck {
  path?: string
  interval?: string
  timeout?: string
  healthyStatuses?: number[]
}

export interface Security {
  forwardHostHeaders: boolean
  enableExploitProtection: boolean
  tls?: {
    insecureSkipVerify: boolean
    rootCAs?: string
    clientCert?: string
    clientKey?: string
  }
}

export interface InstanceCreateRequest {
  name: string
  environment: string
  description?: string
  endpoint: string
  metricsEndpoint?: string
  healthEndpoint?: string
  version?: string
  region?: string
  tags?: string[]
}

export interface ImportResult {
  created: number
  updated: number
  errors?: string[]
}

export const instancesApi = {
  list() {
    return api.get<Instance[]>('/instances')
  },
  get(id: number) {
    return api.get<Instance>(`/instances/${id}`)
  },
  create(data: InstanceCreateRequest) {
    return api.post<Instance>('/instances', data)
  },
  update(id: number, data: Partial<InstanceCreateRequest>) {
    return api.put<Instance>(`/instances/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/instances/${id}`)
  },
  getStats() {
    return api.get('/instances/stats')
  },
  getHealthy() {
    return api.get<Instance[]>('/instances/healthy')
  },
  getRoutes(id: number) {
    return api.get<Route[]>(`/instances/${id}/routes`)
  },
  exportConfig(id: number) {
    return api.get(`/instances/${id}/export`, {
      responseType: 'text',
      headers: { Accept: 'application/x-yaml' },
    })
  },
  importConfig(id: number, yamlContent: string) {
    return api.post<ImportResult>(`/instances/${id}/import`, yamlContent, {
      headers: { 'Content-Type': 'application/x-yaml' },
    })
  },
  copyTo(sourceId: number, targetId: number) {
    return api.post<ImportResult>(`/instances/${sourceId}/copy-to/${targetId}`)
  },
}
