import api from './client'

export interface User {
  id: string
  name: string
  email: string
  role: string
  created_at: string
  updated_at: string
}

export interface UserListResponse {
  users: User[]
  total: number
  page: number
  limit: number
}

export interface UpdateUserRequest {
  name?: string
  email?: string
  role?: string
}

export interface UpdatePasswordRequest {
  password: string
}

export const usersApi = {
  listUsers(page = 1, limit = 20) {
    return api.get<UserListResponse>('/users', { params: { page, limit } })
  },

  getUser(id: string) {
    return api.get<User>(`/users/${id}`)
  },

  updateUser(id: string, data: UpdateUserRequest) {
    return api.put<User>(`/users/${id}`, data)
  },

  deleteUser(id: string) {
    return api.delete(`/users/${id}`)
  },

  updatePassword(id: string, data: UpdatePasswordRequest) {
    return api.put(`/users/${id}/password`, data)
  },
}
