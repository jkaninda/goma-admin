package services

import (
	"net/http"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
)

// NormalizePageParams clamps page/size and computes offset.
func NormalizePageParams(page, size int) (int, int, int) {
	if size <= 0 || size > 100 {
		size = 20
	}
	if page < 0 {
		page = 0
	}
	offset := page * size
	return page, size, offset
}

// Paginated writes a paginated JSON response matching Posta's format.
func Paginated[T any](c *okapi.Context, items []T, total int64, page, size int) error {
	if items == nil {
		items = []T{}
	}
	totalPages := 0
	if size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}
	return c.JSON(http.StatusOK, dto.PageableResponse[T]{
		Success: true,
		Data:    items,
		Pageable: dto.Pageable{
			CurrentPage:   page,
			Size:          size,
			TotalPages:    totalPages,
			TotalElements: total,
			Empty:         len(items) == 0,
		},
	})
}
