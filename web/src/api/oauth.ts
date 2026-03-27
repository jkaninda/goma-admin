import api from './client'

export interface OAuthProviderInfo {
  name: string
  display_name: string
  auth_url: string
  enabled: boolean
}

export interface OAuthProviderDetail {
  id: number
  name: string
  display_name: string
  client_id: string
  auth_url: string
  token_url: string
  user_info_url: string
  scopes: string[]
  user_id_field: string
  email_field: string
  name_field: string
  avatar_field: string
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface SaveOAuthProviderRequest {
  name: string
  display_name: string
  client_id: string
  client_secret: string
  auth_url: string
  token_url: string
  user_info_url: string
  scopes: string[]
  user_id_field: string
  email_field: string
  name_field: string
  avatar_field: string
  enabled: boolean
}

export interface OAuthLoginResponse {
  access_token: string
  expires_at: number
  token_type: string
  user: { id: string; email: string; name: string; role: string }
  new_user: boolean
}

export const oauthApi = {
  // Public: get provider info for login page
  getProviderInfo() {
    return api.get<OAuthProviderInfo>('/auth/oauth')
  },

  // Admin: get full provider config
  getProvider() {
    return api.get<OAuthProviderDetail>('/oauth-provider')
  },

  // Admin: create/update provider
  saveProvider(data: SaveOAuthProviderRequest) {
    return api.put<OAuthProviderDetail>('/oauth-provider', data)
  },

  // Admin: delete provider
  deleteProvider() {
    return api.delete('/oauth-provider')
  },
}
