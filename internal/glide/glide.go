package glide

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/muniere/glide/internal/shell"
)

type GitRef struct {
	Name string
}

func (r GitRef) Short() string {
	if _, after, ok := strings.Cut(r.Name, "refs/heads/"); ok {
		return after
	}

	return r.Name
}

type GitWorktree struct {
	Path string
	Head string
	Ref  GitRef
}

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

func Prepare(branch string) (string, error) {
	res := shell.Capture("git", "rev-parse", "--show-toplevel")
	if !res.Success {
		if res.Stderr != "" {
			return "", fmt.Errorf("%s", res.Stderr)
		}
		return "", fmt.Errorf("Error: not in a git repository")
	}

	repoRoot := res.Stdout
	parent := filepath.Dir(repoRoot)
	repo := filepath.Base(repoRoot)
	normalized := strings.ReplaceAll(branch, "/", "-")
	return filepath.Join(parent, repo+"-"+normalized), nil
}

func Add(args []string) shell.CallResult {
	return shell.Call("git", append([]string{"worktree", "add"}, args...))
}

func Remove(path string) shell.CallResult {
	return shell.Call("git", []string{"worktree", "remove", path})
}
