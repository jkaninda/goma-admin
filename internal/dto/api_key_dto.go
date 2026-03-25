package dto

type CreateAPIKeyRq struct {
	Body struct {
		Name          string   `json:"name" required:"true" minLength:"2" description:"API key name" example:"Production Key"`
		AllowedIPs    []string `json:"allowed_ips,omitempty" description:"Allowed IP addresses or CIDR blocks"`
		ExpiresInDays *int     `json:"expires_in_days,omitempty" description:"Expiration in days (0 or omit = never)"`
	} `json:"body"`
}

type APIKeyByIDRq struct {
	ID int `param:"id" required:"true" description:"API key ID"`
}

type APIKeyCreateResponse struct {
	Key       string  `json:"key" description:"Full API key (shown only once)"`
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Prefix    string  `json:"prefix"`
	ExpiresAt *string `json:"expires_at,omitempty"`
	Message   string  `json:"message"`
}
