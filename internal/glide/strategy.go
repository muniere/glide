package glide

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ResolveOptions holds optional parameters for resolving a worktree path.
type ResolveOptions struct {
	RepoRoot string
}

// Strategy determines where a new worktree is placed.
type Strategy interface {
	Resolve(branch string, opts ResolveOptions) (string, error)
}

// flatStrategy places worktrees as siblings of the repository directory,
// using the repo name and branch name joined by a separator.
type flatStrategy struct {
	separator string
}

// Resolve returns the path for a new worktree for the given branch.
func (s flatStrategy) Resolve(branch string, opts ResolveOptions) (string, error) {
	parent := filepath.Dir(opts.RepoRoot)
	repo := filepath.Base(opts.RepoRoot)
	normalized := strings.ReplaceAll(branch, "/", s.separator)
	return filepath.Join(parent, repo+s.separator+normalized), nil
}

// hierarchyStrategy places worktrees as children of a root container directory,
// preserving the branch name as-is (including path separators).
type hierarchyStrategy struct {
	root string
}

// Resolve returns the path for a new worktree for the given branch.
func (s hierarchyStrategy) Resolve(branch string, opts ResolveOptions) (string, error) {
	if s.root == "" {
		return "", fmt.Errorf("Error: strategy.root is required for hierarchy strategy")
	}
	root := s.root
	if !filepath.IsAbs(root) {
		root = filepath.Join(opts.RepoRoot, root)
	}
	return filepath.Join(root, branch), nil
}
