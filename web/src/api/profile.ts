import api from './client'

export interface Profile {
  id: string
  name: string
  email: string
  role: string
  created_at: string
  updated_at: string
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
}
