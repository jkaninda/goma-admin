package dto

type LoginRequest struct {
	Body struct {
		Email      string `json:"email" required:"true" format:"email" description:"User email address" example:"admin@example.com"`
		Password   string `json:"password" required:"true" minLength:"4" description:"User password"`
		RememberMe bool   `json:"remember_me" description:"Extend session duration"`
	} `json:"body"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token" description:"JWT access token"`
	ExpiresAt   int64        `json:"expires_at" description:"Token expiration timestamp"`
	TokenType   string       `json:"token_type" description:"Token type (Bearer)" example:"Bearer"`
	User        UserResponse `json:"user" description:"Authenticated user info"`
}

type UserResponse struct {
	ID    string `json:"id" description:"User UUID"`
	Email string `json:"email" description:"User email"`
	Name  string `json:"name" description:"User display name"`
	Roles string `json:"role" description:"User role"`
}
