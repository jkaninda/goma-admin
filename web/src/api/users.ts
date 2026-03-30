import api from './client'
import type { PaginatedResponse } from './types'

export interface UserDetail {
  id: string
  email: string
  name: string
  avatar?: string
  role: string
  email_verified: boolean
  active: boolean
  two_factor_enabled: boolean
  oauth_provider?: string
  last_login_at?: string
  created_at: string
}

export interface CreateUserRequest {
  email: string
  name: string
  password: string
  role: string
}

export interface UpdateUserRequest {
  email?: string
  name?: string
  role?: string
  active?: boolean
}

export const usersApi = {
  list(page = 1, pageSize = 20, role?: string, search?: string) {
    return api.get<PaginatedResponse<UserDetail>>('/users', {
      params: { page, page_size: pageSize, role, search },
    })
  },

  get(id: string) {
    return api.get<UserDetail>(`/users/${id}`)
  },

  create(data: CreateUserRequest) {
    return api.post<UserDetail>('/users', data)
  },

  update(id: string, data: UpdateUserRequest) {
    return api.put<UserDetail>(`/users/${id}`, data)
  },

  delete(id: string) {
    return api.delete(`/users/${id}`)
  },

  disable2FA(id: string) {
    return api.delete(`/users/${id}/2fa`)
  },
}
