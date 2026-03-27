package dto

// Response is the standard API response envelope.
type Response[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

// PageableResponse is the paginated API response envelope.
type PageableResponse[T any] struct {
	Success  bool     `json:"success"`
	Data     []T      `json:"data"`
	Pageable Pageable `json:"pageable"`
}

// Pageable holds pagination metadata.
type Pageable struct {
	CurrentPage   int   `json:"current_page"`
	Size          int   `json:"size"`
	TotalPages    int   `json:"total_pages"`
	TotalElements int64 `json:"total_elements"`
	Empty         bool  `json:"empty"`
}

// ListRequest is the standard pagination request used by all list endpoints.
type ListRequest struct {
	Page   int    `query:"page" default:"0" description:"Page number (0-based)"`
	Size   int    `query:"size" default:"20" description:"Page size (max 100)"`
	Search string `query:"q" description:"Search filter"`
}
