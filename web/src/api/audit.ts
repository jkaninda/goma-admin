import api from './client'

export interface ConfigSnapshot {
  id: number
  instanceId: number
  userId?: string
  action: string
  resource: string
  resourceId: number
  name: string
  before?: Record<string, unknown>
  after?: Record<string, unknown>
  createdAt: string
}

export interface PageableResponse<T> {
  success: boolean
  data: T[]
  pageable: {
    current_page: number
    size: number
    total_pages: number
    total_elements: number
    empty: boolean
  }
}

export const auditApi = {
  list(page = 0, size = 20) {
    return api.get<PageableResponse<ConfigSnapshot>>('/audit/snapshots', {
      params: { page, size },
    })
  },

  get(id: number) {
    return api.get<ConfigSnapshot>(`/audit/snapshots/${id}`)
  },

  rollback(id: number) {
    return api.post(`/audit/snapshots/${id}/rollback`)
  },
}
