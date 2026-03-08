package dto

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	ExpiresAt   int64        `json:"expires_at"`
	TokenType   string       `json:"token_type"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Roles string `json:"role"`
}
