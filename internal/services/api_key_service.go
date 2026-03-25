package services

import (
	"time"

	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(db *gorm.DB) *APIKeyService {
	return &APIKeyService{repo: repository.NewAPIKeyRepository(db)}
}

func (s *APIKeyService) List(c *okapi.Context, input *dto.ListRequest) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	instanceID := OptionalInstanceID(c)
	page, size, offset := NormalizePageParams(input.Page, input.Size)

	keys, total, err := s.repo.FindByUserAndInstance(c.Context(), userID, instanceID, size, offset)
	if err != nil {
		return c.AbortInternalServerError("Failed to list API keys", err)
	}

	return Paginated(c, keys, total, page, size)
}

func (s *APIKeyService) Create(c *okapi.Context, input *dto.CreateAPIKeyRq) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	instanceID := OptionalInstanceID(c)

	rawKey, keyHash, keyPrefix, err := models.GenerateAPIKey()
	if err != nil {
		return c.AbortInternalServerError("Failed to generate API key", err)
	}

	apiKey := &models.APIKey{
		UserID:     userID,
		InstanceID: instanceID,
		Name:       input.Body.Name,
		KeyHash:    keyHash,
		KeyPrefix:  keyPrefix,
		AllowedIPs: input.Body.AllowedIPs,
	}

	if input.Body.ExpiresInDays != nil && *input.Body.ExpiresInDays > 0 {
		exp := time.Now().AddDate(0, 0, *input.Body.ExpiresInDays)
		apiKey.ExpiresAt = &exp
	}

	if err := s.repo.Create(c.Context(), apiKey); err != nil {
		return c.AbortInternalServerError("Failed to create API key", err)
	}

	var expiresAt *string
	if apiKey.ExpiresAt != nil {
		s := apiKey.ExpiresAt.Format(time.RFC3339)
		expiresAt = &s
	}

	return c.Created(dto.APIKeyCreateResponse{
		Key:       rawKey,
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Prefix:    apiKey.KeyPrefix,
		ExpiresAt: expiresAt,
		Message:   "Save this API key securely. It will not be shown again.",
	})
}

func (s *APIKeyService) Revoke(c *okapi.Context, input *dto.APIKeyByIDRq) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	key, err := s.repo.FindByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("API key not found")
	}

	if key.UserID != userID {
		return c.AbortForbidden("Not authorized to revoke this key")
	}

	if err := s.repo.Revoke(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Failed to revoke API key", err)
	}

	return c.OK(okapi.M{"message": "API key revoked"})
}

func (s *APIKeyService) Delete(c *okapi.Context, input *dto.APIKeyByIDRq) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.AbortUnauthorized("Unauthorized")
	}

	key, err := s.repo.FindByID(c.Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("API key not found")
	}

	if key.UserID != userID {
		return c.AbortForbidden("Not authorized to delete this key")
	}

	if !key.Revoked && !key.IsExpired() {
		return c.AbortBadRequest("Can only delete revoked or expired keys", nil)
	}

	if err := s.repo.Delete(c.Context(), uint(input.ID)); err != nil {
		return c.AbortInternalServerError("Failed to delete API key", err)
	}

	return c.NoContent()
}
