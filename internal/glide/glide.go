package glide

import (
	"fmt"
	"strings"

	"github.com/muniere/glide/internal/shell"
)

// GitRef represents a git reference such as a branch or tag.
type GitRef struct {
	Name string
}

// Short returns the short name of the ref, stripping the "refs/heads/" prefix if present.
func (r GitRef) Short() string {
	if _, after, ok := strings.Cut(r.Name, "refs/heads/"); ok {
		return after
	}

	return r.Name
}

// GitWorktree represents a single git worktree entry.
type GitWorktree struct {
	Path string
	Head string
	Ref  GitRef
}

// List returns all worktrees in the current repository.
func List() ([]GitWorktree, error) {
	res := shell.Capture("git", "worktree", "list", "--porcelain")
	if !res.Success {
		if res.Stderr != "" {
			return nil, fmt.Errorf("%s", res.Stderr)
		}
		return nil, fmt.Errorf("Error: failed to list worktrees")
	}

	lines := make([]string, 0)
	for _, line := range strings.Split(res.Stdout, "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}

	chunks := make([][]string, 0, len(lines)/3)
	for i := 0; i+2 < len(lines); i += 3 {
		chunks = append(chunks, lines[i:i+3])
	}

	rows := make([]GitWorktree, 0, len(chunks))
	for _, chunk := range chunks {
		path, ok := strings.CutPrefix(chunk[0], "worktree ")
		if !ok {
			continue
		}
		head, ok := strings.CutPrefix(chunk[1], "HEAD ")
		if !ok {
			continue
		}
		branch, ok := strings.CutPrefix(chunk[2], "branch ")
		if !ok {
			continue
		}

		rows = append(rows, GitWorktree{Path: path, Head: head, Ref: GitRef{Name: branch}})
	}

	return rows, nil
}

// Find returns the path of the worktree for the given branch, or nil if not found.
func Find(branch string) (*string, error) {
	entries, err := List()
	if err != nil {
		return nil, err
	}
	for _, it := range entries {
		if it.Ref.Short() == branch {
			return &it.Path, nil
		}
	}
	return nil, nil
}

// Add creates a new worktree with the given arguments.
func Add(args []string) shell.CallResult {
	return shell.Call("git", append([]string{"worktree", "add"}, args...))
}

// Remove deletes the worktree at the given path.
func Remove(path string) shell.CallResult {
	return shell.Call("git", []string{"worktree", "remove", path})
}
