package dto

type UpdateProfileRq struct {
	Body struct {
		Name  string `json:"name" required:"true" minLength:"2" description:"Display name"`
		Email string `json:"email" required:"true" format:"email" description:"Email address"`
	} `json:"body"`
}

type ChangePasswordRq struct {
	Body struct {
		CurrentPassword string `json:"current_password" required:"true" description:"Current password"`
		NewPassword     string `json:"new_password" required:"true" minLength:"6" description:"New password"`
	} `json:"body"`
}
