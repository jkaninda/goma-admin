import api from './client'

export interface ApiKey {
  id: number
  user_id: string
  instance_id?: string
  name: string
  key_prefix: string
  created_at: string
  expires_at: string | null
  last_used_at: string | null
  revoked: boolean
  allowed_ips: string[] | null
}

export interface ApiKeyCreateResponse {
  key: string
  id: number
  name: string
  prefix: string
  expires_at: string | null
  message: string
}

export const apiKeysApi = {
  list(page = 0, size = 20) {
    return api.get<{ data: ApiKey[]; meta: { page: number; size: number; total: number } }>(
      '/api-keys',
      { params: { page, size } }
    )
  },
  create(name: string, allowedIPs?: string[], expiresInDays?: number) {
    const body: Record<string, unknown> = { name }
    if (allowedIPs && allowedIPs.length > 0) body.allowed_ips = allowedIPs
    if (expiresInDays !== undefined) body.expires_in_days = expiresInDays
    return api.post<ApiKeyCreateResponse>('/api-keys', body)
  },
  revoke(id: number) {
    return api.put(`/api-keys/${id}/revoke`)
  },
  delete(id: number) {
    return api.delete(`/api-keys/${id}`)
  },
}
