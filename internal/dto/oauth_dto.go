package dto

// OAuthCallbackRequest represents the OAuth callback with authorization code.
type OAuthCallbackRequest struct {
	Code  string `query:"code" required:"true" description:"Authorization code from OAuth provider"`
	State string `query:"state" required:"true" description:"CSRF state parameter"`
}

// OAuthLoginResponse is returned after successful OAuth login.
type OAuthLoginResponse struct {
	AccessToken string       `json:"access_token" description:"JWT access token"`
	ExpiresAt   int64        `json:"expires_at" description:"Token expiration timestamp"`
	TokenType   string       `json:"token_type" description:"Token type (Bearer)" example:"Bearer"`
	User        UserResponse `json:"user" description:"Authenticated user info"`
	NewUser     bool         `json:"new_user" description:"True if the user was just created"`
}

// OAuthProviderInfo describes the configured OAuth provider (public).
type OAuthProviderInfo struct {
	Name        string `json:"name" description:"Provider name"`
	DisplayName string `json:"display_name" description:"Human-friendly label"`
	AuthURL     string `json:"auth_url" description:"URL to redirect the user to for authorization"`
	Enabled     bool   `json:"enabled" description:"Whether the provider is enabled"`
}

// SaveOAuthProviderRequest is used by admins to configure the OAuth provider.
type SaveOAuthProviderRequest struct {
	Body struct {
		Name         string   `json:"name" required:"true" minLength:"2" description:"Provider identifier (e.g. keycloak, gitea)" example:"keycloak"`
		DisplayName  string   `json:"display_name" description:"Human-friendly label" example:"Company SSO"`
		ClientID     string   `json:"client_id" required:"true" description:"OAuth2 client ID"`
		ClientSecret string   `json:"client_secret" required:"true" description:"OAuth2 client secret"`
		AuthURL      string   `json:"auth_url" required:"true" description:"Authorization endpoint URL"`
		TokenURL     string   `json:"token_url" required:"true" description:"Token endpoint URL"`
		UserInfoURL  string   `json:"user_info_url" required:"true" description:"Userinfo endpoint URL"`
		Scopes       []string `json:"scopes" description:"OAuth scopes" example:"[\"openid\",\"email\",\"profile\"]"`
		UserIDField  string   `json:"user_id_field" description:"JSON field for provider user ID (default: sub)"`
		EmailField   string   `json:"email_field" description:"JSON field for email (default: email)"`
		NameField    string   `json:"name_field" description:"JSON field for display name (default: name)"`
		AvatarField  string   `json:"avatar_field" description:"JSON field for avatar URL (default: picture)"`
		Enabled      bool     `json:"enabled" description:"Whether the provider is enabled"`
	} `json:"body"`
}

// OAuthProviderDetailResponse is the admin view of the OAuth provider config.
type OAuthProviderDetailResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	ClientID    string   `json:"client_id"`
	AuthURL     string   `json:"auth_url"`
	TokenURL    string   `json:"token_url"`
	UserInfoURL string   `json:"user_info_url"`
	Scopes      []string `json:"scopes"`
	UserIDField string   `json:"user_id_field"`
	EmailField  string   `json:"email_field"`
	NameField   string   `json:"name_field"`
	AvatarField string   `json:"avatar_field"`
	Enabled     bool     `json:"enabled"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
