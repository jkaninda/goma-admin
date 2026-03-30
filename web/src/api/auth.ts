import api from './client'

export interface LoginRequest {
  email: string
  password: string
  remember_me: boolean
  two_factor_code?: string
}

export interface AuthResponse {
  access_token: string
  expires_at: number
  token_type: string
  user: {
    id: string
    email: string
    name: string
    role: string
  }
}

export const authApi = {
  login(data: LoginRequest) {
    return api.post<AuthResponse>('/auth/login', data)
  },
  logout() {
    return api.post('/auth/logout')
  },
}
