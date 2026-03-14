package vcs

import (
	"fmt"

	"github.com/muniere/glide/internal/shell"
)

// RepoRoot returns the absolute path of the root of the current repository.
func RepoRoot() (string, error) {
	res := shell.Capture("git", "rev-parse", "--show-toplevel")
	if !res.Success {
		if res.Stderr != "" {
			return "", fmt.Errorf("%s", res.Stderr)
		}
		return "", fmt.Errorf("Error: not in a git repository")
	}
	return res.Stdout, nil
}
