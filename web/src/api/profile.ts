import api from './client'

export interface Profile {
  id: string
  name: string
  email: string
  role: string
  two_factor_enabled: boolean
  oauth_provider?: string
  created_at: string
  updated_at: string
}

export interface Setup2FAResponse {
  secret: string
  url: string
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
}

export const profileApi = {
  getProfile() {
    return api.get<Profile>('/profile')
  },

  updateProfile(data: UpdateProfileRequest) {
    return api.put<Profile>('/profile', data)
  },

  changePassword(data: ChangePasswordRequest) {
    return api.put('/profile/password', data)
  },

  setup2FA() {
    return api.post<Setup2FAResponse>('/profile/2fa/setup')
  },

  verify2FA(code: string) {
    return api.post('/profile/2fa/verify', { code })
  },

  disable2FA(code: string) {
    return api.post('/profile/2fa/disable', { code })
  },
}
