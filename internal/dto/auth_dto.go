package dto

type LoginRequest struct {
	Body struct {
		Email         string `json:"email" required:"true" format:"email" description:"User email address" example:"admin@example.com"`
		Password      string `json:"password" required:"true" minLength:"4" description:"User password"`
		RememberMe    bool   `json:"remember_me" description:"Extend session duration"`
		TwoFactorCode string `json:"two_factor_code" description:"TOTP code for 2FA verification"`
	} `json:"body"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token" description:"JWT access token"`
	ExpiresAt   int64        `json:"expires_at" description:"Token expiration timestamp"`
	TokenType   string       `json:"token_type" description:"Token type (Bearer)" example:"Bearer"`
	User        UserResponse `json:"user" description:"Authenticated user info"`
}

type UserResponse struct {
	ID               string `json:"id" description:"User UUID"`
	Email            string `json:"email" description:"User email"`
	Name             string `json:"name" description:"User display name"`
	Roles            string `json:"role" description:"User role"`
	TwoFactorEnabled bool   `json:"two_factor_enabled" description:"Whether 2FA is enabled"`
	OAuthProvider    string `json:"oauth_provider,omitempty" description:"OAuth provider name (empty for local accounts)"`
}

type Setup2FAResponse struct {
	Secret string `json:"secret" description:"TOTP secret key (Base32)"`
	URL    string `json:"url" description:"otpauth URL for QR code generation"`
}

type Verify2FARequest struct {
	Body struct {
		Code string `json:"code" required:"true" minLength:"6" maxLength:"6" description:"6-digit TOTP code"`
	} `json:"body"`
}

type Disable2FARequest struct {
	Body struct {
		Code string `json:"code" required:"true" minLength:"6" maxLength:"6" description:"6-digit TOTP code to confirm disable"`
	} `json:"body"`
}
