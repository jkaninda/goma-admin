package dto

// CreateUserRequest is used by admins to create a new user.
type CreateUserRequest struct {
	Body struct {
		Email    string `json:"email" required:"true" format:"email" description:"User email address"`
		Name     string `json:"name" required:"true" minLength:"2" description:"Display name"`
		Password string `json:"password" required:"true" minLength:"6" description:"User password"`
		Role     string `json:"role" description:"User role (viewer, user, admin)" example:"user"`
	} `json:"body"`
}

// UpdateUserRequest is used by admins to update a user.
type UpdateUserRequest struct {
	ID   string `param:"id" required:"true" description:"User UUID"`
	Body struct {
		Email  string `json:"email" format:"email" description:"User email address"`
		Name   string `json:"name" description:"Display name"`
		Role   string `json:"role" description:"User role"`
		Active *bool  `json:"active" description:"Whether the user is active"`
	} `json:"body"`
}

// UserByIDRequest identifies a user by ID path param.
type UserByIDRequest struct {
	ID string `param:"id" required:"true" description:"User UUID"`
}

// ListUsersRequest is used to list/paginate users.
type ListUsersRequest struct {
	Page     int    `query:"page" default:"1" description:"Page number" example:"1"`
	PageSize int    `query:"page_size" default:"20" description:"Items per page" example:"20"`
	Role     string `query:"role" description:"Filter by role"`
	Search   string `query:"search" description:"Search by name or email"`
}

// UserDetailResponse includes full user details for admin views.
type UserDetailResponse struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	Avatar        string  `json:"avatar,omitempty"`
	Role          string  `json:"role"`
	EmailVerified bool    `json:"email_verified"`
	Active        bool    `json:"active"`
	OAuthProvider string  `json:"oauth_provider,omitempty"`
	LastLoginAt   *string `json:"last_login_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
}
