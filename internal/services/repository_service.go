package services

import (
	"context"

	"github.com/jkaninda/goma-admin/internal/crypto"
	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/models"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"github.com/jkaninda/okapi"
	"gorm.io/gorm"
)

// RepositoryService handles CRUD and sync for git repositories.
type RepositoryService struct {
	repo *repository.RepositoryRepository
	git  *GitService
}

// NewRepositoryService creates a new RepositoryService.
func NewRepositoryService(db *gorm.DB, git *GitService) *RepositoryService {
	return &RepositoryService{
		repo: repository.NewRepositoryRepository(db),
		git:  git,
	}
}

func (s *RepositoryService) List(c *okapi.Context) error {
	repos, err := s.repo.List(c.Request().Context())
	if err != nil {
		return c.AbortInternalServerError("Failed to list repositories", err)
	}
	return c.OK(repos)
}

func (s *RepositoryService) Get(c *okapi.Context, input *dto.RepositoryByIDRq) error {
	repo, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Repository not found", err)
	}
	return c.OK(repo)
}

func (s *RepositoryService) Create(c *okapi.Context, input *dto.CreateRepositoryRq) error {
	branch := input.Body.Branch
	if branch == "" {
		branch = "main"
	}

	repo := &models.Repository{
		Name:     input.Body.Name,
		URL:      input.Body.URL,
		Branch:   branch,
		AuthType: input.Body.AuthType,
		Status:   "pending",
	}

	if input.Body.AuthValue != "" {
		encrypted, err := crypto.Encrypt(input.Body.AuthValue)
		if err != nil {
			return c.AbortInternalServerError("Failed to encrypt credentials", err)
		}
		repo.AuthValue = encrypted
	}

	if err := s.repo.Create(c.Request().Context(), repo); err != nil {
		return c.AbortInternalServerError("Failed to create repository", err)
	}

	// Clone asynchronously — use background context since the HTTP request will be done
	go func() {
		bgCtx := context.Background()
		commit, err := s.git.Clone(repo.ID, repo.URL, repo.Branch, input.Body.AuthType, input.Body.AuthValue)
		if err != nil {
			logger.Error("Failed to clone repository", "repo", repo.Name, "error", err)
			_ = s.repo.UpdateSyncStatus(bgCtx, repo.ID, "error", err.Error(), "")
		} else {
			_ = s.repo.UpdateSyncStatus(bgCtx, repo.ID, "synced", "", commit)
		}
	}()

	return c.Created(repo)
}

func (s *RepositoryService) Update(c *okapi.Context, input *dto.UpdateRepositoryRq) error {
	ctx := c.Request().Context()

	existing, err := s.repo.GetByID(ctx, uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Repository not found", err)
	}

	existing.Name = input.Body.Name
	existing.URL = input.Body.URL
	if input.Body.Branch != "" {
		existing.Branch = input.Body.Branch
	}
	existing.AuthType = input.Body.AuthType

	if input.Body.AuthType == "" {
		existing.AuthValue = ""
	} else if input.Body.AuthValue != "" {
		encrypted, err := crypto.Encrypt(input.Body.AuthValue)
		if err != nil {
			return c.AbortInternalServerError("Failed to encrypt credentials", err)
		}
		existing.AuthValue = encrypted
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return c.AbortInternalServerError("Failed to update repository", err)
	}

	return c.OK(existing)
}

func (s *RepositoryService) Delete(c *okapi.Context, input *dto.RepositoryByIDRq) error {
	ctx := c.Request().Context()

	if err := s.repo.Delete(ctx, uint(input.ID)); err != nil {
		return c.AbortNotFound("Repository not found", err)
	}

	// Remove local clone
	if err := s.git.RemoveClone(uint(input.ID)); err != nil {
		logger.Error("Failed to remove repo clone", "repoID", input.ID, "error", err)
	}

	return c.NoContent()
}

func (s *RepositoryService) Sync(c *okapi.Context, input *dto.RepositoryByIDRq) error {
	ctx := c.Request().Context()

	repo, err := s.repo.GetByID(ctx, uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Repository not found", err)
	}

	var commit string
	if s.git.Exists(repo.ID) {
		commit, err = s.git.Pull(repo.ID, repo.Branch, repo.AuthType, repo.AuthValue)
	} else {
		commit, err = s.git.Clone(repo.ID, repo.URL, repo.Branch, repo.AuthType, repo.AuthValue)
	}

	if err != nil {
		_ = s.repo.UpdateSyncStatus(ctx, repo.ID, "error", err.Error(), "")
		return c.AbortInternalServerError("Sync failed", err)
	}

	_ = s.repo.UpdateSyncStatus(ctx, repo.ID, "synced", "", commit)

	return c.OK(dto.SyncResult{Status: "synced", Commit: commit})
}

func (s *RepositoryService) Browse(c *okapi.Context, input *dto.RepositoryByIDRq) error {
	repo, err := s.repo.GetByID(c.Request().Context(), uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Repository not found", err)
	}

	if !s.git.Exists(repo.ID) {
		return c.AbortBadRequest("Repository has not been synced yet", nil)
	}

	path := c.Query("path")
	entries, err := s.git.Browse(repo.ID, path)
	if err != nil {
		return c.AbortBadRequest("Failed to browse path", err)
	}

	return c.OK(entries)
}

// Webhook handles push notifications from Git providers (GitHub, GitLab).
func (s *RepositoryService) Webhook(c *okapi.Context, input *dto.RepositoryByIDRq) error {
	ctx := c.Request().Context()

	repo, err := s.repo.GetByID(ctx, uint(input.ID))
	if err != nil {
		return c.AbortNotFound("Repository not found", err)
	}

	var commit string
	if s.git.Exists(repo.ID) {
		commit, err = s.git.Pull(repo.ID, repo.Branch, repo.AuthType, repo.AuthValue)
	} else {
		commit, err = s.git.Clone(repo.ID, repo.URL, repo.Branch, repo.AuthType, repo.AuthValue)
	}

	if err != nil {
		_ = s.repo.UpdateSyncStatus(ctx, repo.ID, "error", err.Error(), "")
		return c.AbortInternalServerError("Sync failed", err)
	}

	_ = s.repo.UpdateSyncStatus(ctx, repo.ID, "synced", "", commit)

	return c.OK(dto.SyncResult{Status: "synced", Commit: commit})
}
