package services

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitService handles git clone, pull, and browse operations on local clones.
type GitService struct {
	baseDir string // base directory for all repo clones
}

// NewGitService creates a new GitService.
// baseDir is the directory where repos are cloned (e.g. /var/lib/goma-admin/repos).
func NewGitService(baseDir string) *GitService {
	return &GitService{baseDir: baseDir}
}

// repoDir returns the local clone directory for a given repository ID.
func (g *GitService) repoDir(repoID uint) string {
	return filepath.Join(g.baseDir, fmt.Sprintf("%d", repoID))
}

// authMethod builds go-git transport auth from type and value.
func authMethod(authType, authValue string) transport.AuthMethod {
	switch authType {
	case "token":
		return &http.BasicAuth{Username: "x-token-auth", Password: authValue}
	case "basic":
		parts := strings.SplitN(authValue, ":", 2)
		if len(parts) == 2 {
			return &http.BasicAuth{Username: parts[0], Password: parts[1]}
		}
		return &http.BasicAuth{Username: authValue}
	}
	return nil
}

// Clone clones a repository. If it already exists locally, it removes and re-clones.
func (g *GitService) Clone(repoID uint, url, branch, authType, authValue string) (string, error) {
	dir := g.repoDir(repoID)

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dir), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Remove existing clone if present
	if err := os.RemoveAll(dir); err != nil {
		return "", fmt.Errorf("failed to clean clone directory: %w", err)
	}

	opts := &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Depth:         1,
	}
	if auth := authMethod(authType, authValue); auth != nil {
		opts.Auth = auth
	}

	repo, err := git.PlainClone(dir, false, opts)
	if err != nil {
		_ = os.RemoveAll(dir)
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return headCommit(repo)
}

// Pull pulls the latest changes for a cloned repository.
func (g *GitService) Pull(repoID uint, branch, authType, authValue string) (string, error) {
	dir := g.repoDir(repoID)

	repo, err := git.PlainOpen(dir)
	if err != nil {
		return "", fmt.Errorf("failed to open repo: %w", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	opts := &git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Force:         true,
	}
	if auth := authMethod(authType, authValue); auth != nil {
		opts.Auth = auth
	}

	err = wt.Pull(opts)
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return "", fmt.Errorf("git pull failed: %w", err)
	}

	return headCommit(repo)
}

// Browse lists files and directories at a given path within the cloned repo.
func (g *GitService) Browse(repoID uint, subpath string) ([]BrowseEntry, error) {
	dir := g.repoDir(repoID)
	target := filepath.Join(dir, filepath.Clean(subpath))

	// Security: ensure we don't escape the repo dir
	if !strings.HasPrefix(target, dir) {
		return nil, fmt.Errorf("invalid path")
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var result []BrowseEntry
	for _, e := range entries {
		name := e.Name()
		// Skip hidden files/dirs (.git, etc.)
		if strings.HasPrefix(name, ".") {
			continue
		}
		entryPath := subpath
		if entryPath != "" {
			entryPath += "/"
		}
		entryPath += name
		result = append(result, BrowseEntry{
			Name:  name,
			IsDir: e.IsDir(),
			Path:  entryPath,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		// Dirs first, then alphabetical
		if result[i].IsDir != result[j].IsDir {
			return result[i].IsDir
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// BrowseEntry represents a file or directory in the repo.
type BrowseEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Path  string `json:"path"`
}

// ReadFile reads a file from the cloned repo at the given path.
func (g *GitService) ReadFile(repoID uint, filePath string) ([]byte, error) {
	dir := g.repoDir(repoID)
	target := filepath.Join(dir, filepath.Clean(filePath))

	if !strings.HasPrefix(target, dir) {
		return nil, fmt.Errorf("invalid path")
	}

	return os.ReadFile(target)
}

// ReadYAMLFiles reads all .yaml/.yml files in a directory within the cloned repo.
func (g *GitService) ReadYAMLFiles(repoID uint, subpath string) ([][]byte, error) {
	dir := g.repoDir(repoID)
	target := filepath.Join(dir, filepath.Clean(subpath))

	if !strings.HasPrefix(target, dir) {
		return nil, fmt.Errorf("invalid path")
	}

	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files [][]byte
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
			data, err := os.ReadFile(filepath.Join(target, name))
			if err != nil {
				return nil, fmt.Errorf("failed to read %s: %w", name, err)
			}
			files = append(files, data)
		}
	}

	return files, nil
}

// RemoveClone deletes the local clone for a repository.
func (g *GitService) RemoveClone(repoID uint) error {
	return os.RemoveAll(g.repoDir(repoID))
}

// Exists checks if a local clone exists for the repository.
func (g *GitService) Exists(repoID uint) bool {
	_, err := os.Stat(g.repoDir(repoID))
	return err == nil
}

func headCommit(repo *git.Repository) (string, error) {
	ref, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}
	return ref.Hash().String()[:7], nil
}
