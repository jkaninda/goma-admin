import api from './client'

export interface Middleware {
  id: number
  instanceId: number
  name: string
  type: string
  config: Record<string, unknown>
  createdAt: string
  updatedAt: string
}

export interface MiddlewareCreateRequest {
  name: string
  type: string
  config: Record<string, unknown>
}

export interface ImportResult {
  created: number
  updated: number
  errors?: string[]
}

export const middlewaresApi = {
  list(page = 0, size = 20) {
    return api.get('/middlewares', { params: { page, size } })
  },
  get(id: number) {
    return api.get<Middleware>(`/middlewares/${id}`)
  },
  create(data: MiddlewareCreateRequest) {
    return api.post<Middleware>('/middlewares', data)
  },
  update(id: number, data: MiddlewareCreateRequest) {
    return api.put<Middleware>(`/middlewares/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/middlewares/${id}`)
  },
  search(query: string) {
    return api.get('/middlewares/search', { params: { q: query } })
  },
  importMiddlewares(yamlContent: string) {
    return api.post<ImportResult>('/import/middlewares', yamlContent, {
      headers: { 'Content-Type': 'application/x-yaml' },
    })
  },
}
