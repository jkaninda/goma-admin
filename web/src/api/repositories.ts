import api from './client'

export interface Repository {
  id: number
  name: string
  url: string
  branch: string
  authType: string
  hasAuth: boolean
  lastSyncedAt: string | null
  lastCommit: string
  status: string
  statusMessage: string
  createdAt: string
  updatedAt: string
}

export interface RepositoryCreateRequest {
  name: string
  url: string
  branch?: string
  authType?: string
  authValue?: string
}

export interface BrowseEntry {
  name: string
  isDir: boolean
  path: string
}

export interface SyncResult {
  status: string
  commit?: string
}

export const repositoriesApi = {
  list() {
    return api.get<Repository[]>('/repositories')
  },
  get(id: number) {
    return api.get<Repository>(`/repositories/${id}`)
  },
  create(data: RepositoryCreateRequest) {
    return api.post<Repository>('/repositories', data)
  },
  update(id: number, data: Partial<RepositoryCreateRequest>) {
    return api.put<Repository>(`/repositories/${id}`, data)
  },
  delete(id: number) {
    return api.delete(`/repositories/${id}`)
  },
  sync(id: number) {
    return api.post<SyncResult>(`/repositories/${id}/sync`)
  },
  browse(id: number, path?: string) {
    const params = path ? { path } : {}
    return api.get<BrowseEntry[]>(`/repositories/${id}/browse`, { params })
  },
}
