package services

import (
	"context"
	"time"

	"github.com/jkaninda/goma-admin/internal/dto"
	"github.com/jkaninda/goma-admin/internal/repository"
	"github.com/jkaninda/logger"
	"gorm.io/gorm"
)

// RepoSyncer is a background service that periodically pulls git repositories
// and syncs configs into instances that have auto_sync enabled.
type RepoSyncer struct {
	repoRepo     *repository.RepositoryRepository
	instanceRepo *repository.InstanceRepository
	configSvc    *InstanceConfigService
	git          *GitService
	interval     time.Duration
}

// NewRepoSyncer creates a new RepoSyncer.
func NewRepoSyncer(db *gorm.DB, git *GitService, configSvc *InstanceConfigService, interval time.Duration) *RepoSyncer {
	return &RepoSyncer{
		repoRepo:     repository.NewRepositoryRepository(db),
		instanceRepo: repository.NewInstanceRepository(db),
		configSvc:    configSvc,
		git:          git,
		interval:     interval,
	}
}

// Start begins the periodic sync loop. It blocks until ctx is cancelled.
func (rs *RepoSyncer) Start(ctx context.Context) error {
	logger.Info("Repo syncer started", "interval", rs.interval)

	ticker := time.NewTicker(rs.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Repo syncer stopped")
			return ctx.Err()
		case <-ticker.C:
			rs.syncAll(ctx)
		}
	}
}

// syncAll pulls all synced repos and imports configs for auto-sync instances.
func (rs *RepoSyncer) syncAll(ctx context.Context) {
	repos, err := rs.repoRepo.List(ctx)
	if err != nil {
		logger.Error("Repo syncer: failed to list repositories", "error", err)
		return
	}

	for _, repo := range repos {
		// Skip repos that are actively being cloned (no other states to skip)
		if repo.URL == "" {
			continue
		}

		// Pull latest
		var commit string
		var pullErr error
		if rs.git.Exists(repo.ID) {
			commit, pullErr = rs.git.Pull(repo.ID, repo.Branch, repo.AuthType, repo.AuthValue)
		} else {
			commit, pullErr = rs.git.Clone(repo.ID, repo.URL, repo.Branch, repo.AuthType, repo.AuthValue)
		}

		if pullErr != nil {
			_ = rs.repoRepo.UpdateSyncStatus(ctx, repo.ID, "error", pullErr.Error(), "")
			logger.Error("Repo syncer: pull failed", "repo", repo.Name, "error", pullErr)
			continue
		}

		// Check if commit changed
		if commit == repo.LastCommit {
			continue
		}

		_ = rs.repoRepo.UpdateSyncStatus(ctx, repo.ID, "synced", "", commit)
		logger.Info("Repo syncer: new commit detected", "repo", repo.Name, "commit", commit)

		// Find all instances linked to this repo with auto_sync enabled
		instances, err := rs.instanceRepo.List(ctx)
		if err != nil {
			logger.Error("Repo syncer: failed to list instances", "error", err)
			continue
		}

		for _, inst := range instances {
			if inst.RepositoryID == nil || *inst.RepositoryID != repo.ID || !inst.AutoSync {
				continue
			}

			// Read and import YAML files
			yamlFiles, err := rs.git.ReadYAMLFiles(repo.ID, inst.RepositoryPath)
			if err != nil {
				logger.Error("Repo syncer: failed to read config files", "instance", inst.Name, "error", err)
				continue
			}

			result := dto.ImportResult{}
			for _, data := range yamlFiles {
				partial := rs.configSvc.importYAMLBytes(ctx, inst.ID, data)
				result.Created += partial.Created
				result.Updated += partial.Updated
				result.Errors = append(result.Errors, partial.Errors...)
			}

			if result.Created > 0 || result.Updated > 0 {
				logger.Info("Repo syncer: config imported",
					"instance", inst.Name,
					"repo", repo.Name,
					"commit", commit,
					"created", result.Created,
					"updated", result.Updated,
				)
			}
		}
	}
}
