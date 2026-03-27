import api from './client'

export interface Route {
  id: number
  instanceId: number
  name: string
  config: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface RouteCreateRequest {
  name: string
  config: Record<string, unknown>
}

export interface ImportResult {
  created: number
  updated: number
  errors?: string[]
}

export const routesApi = {
  list(page = 0, size = 20, search = '') {
    const params: Record<string, unknown> = { page, size }
    if (search) params.q = search
    return api.get('/routes', { params })
  },
  get(id: number) {
    return api.get<Route>(`/routes/${id}`)
  },
  create(data: RouteCreateRequest) {
    return api.post<Route>('/routes', data)
  },
  update(id: number, data: RouteCreateRequest) {
    return api.put<Route>(`/routes/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/routes/${id}`)
  },
  findByPath(path: string) {
    return api.get('/routes/find', { params: { path } })
  },
  importRoutes(yamlContent: string) {
    return api.post<ImportResult>('/import/routes', yamlContent, {
      headers: { 'Content-Type': 'application/x-yaml' },
    })
  },
}
