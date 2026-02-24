package cli

import (
	"fmt"
	"strings"

	"github.com/muniere/glide/internal/glide"
	"github.com/muniere/glide/internal/shell"
	"github.com/spf13/cobra"
)

var Version = "1.0.0"

func Execute(args []string) error {
	cmd := &cobra.Command{
		Use:           "glide",
		Short:         "Manage git worktrees",
		Version:       Version,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("Error: command required")
		},
	}

	cmd.AddCommand(func() *cobra.Command {
		ctx := listContext{}
		cmd := &cobra.Command{
			Use:     "list",
			Aliases: []string{"ls"},
			Short:   "List all worktrees",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx.args = args
				return list.execute(ctx)
			},
		}

		cmd.Flags().BoolVar(&ctx.porcelain, "porcelain", false, "machine readable output")
		return cmd
	}())

	cmd.AddCommand(func() *cobra.Command {
		ctx := findContext{}
		return &cobra.Command{
			Use:   "find <branch>",
			Short: "Resolve existing worktree path",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx.args = args
				return find.execute(ctx)
			},
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return find.predict(args, toComplete)
			},
		}
	}())

	cmd.AddCommand(func() *cobra.Command {
		ctx := addContext{}
		return &cobra.Command{
			Use:   "add <branch>",
			Short: "Create a worktree for an existing branch",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx.args = args
				return add.execute(ctx)
			},
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return add.predict(args, toComplete)
			},
		}
	}())

	cmd.AddCommand(func() *cobra.Command {
		ctx := removeContext{}
		return &cobra.Command{
			Use:     "remove <branch> [branch ...]",
			Aliases: []string{"rm"},
			Short:   "Remove worktrees",
			RunE: func(cmd *cobra.Command, args []string) error {
				ctx.args = args
				return remove.execute(ctx)
			},
			ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return remove.predict(args, toComplete)
			},
		}
	}())

	cmd.SetArgs(args)
	return cmd.Execute()
}

type listCommand struct {}

var list listCommand

type listContext struct {
	args      []string
	porcelain bool
}

func (_ listCommand) execute(ctx listContext) error {
	entries, err := glide.List()
	if err != nil {
		return err
	}

	if ctx.porcelain {
		for _, it := range entries {
			fmt.Printf("worktree %s\n", it.Path)
			fmt.Printf("HEAD %s\n", it.Head)
			fmt.Printf("branch %s\n", it.Ref.Short())
		}
		return nil
	}

	pad := 0
	for _, it := range entries {
		if len(it.Path) > pad {
			pad = len(it.Path)
		}
	}
	for _, it := range entries {
		fmt.Printf("%-*s %s %s\n", pad, it.Path, it.Head, it.Ref.Short())
	}
	return nil
}

type findCommand struct{}

var find findCommand

type findContext struct {
	args []string
}

func (_ findCommand) execute(ctx findContext) error {
	if len(ctx.args) == 0 {
		return fmt.Errorf("Error: branch name required")
	}
	if len(ctx.args) > 1 {
		return fmt.Errorf("Error: unexpected argument '%s'", ctx.args[1])
	}

	branch := ctx.args[0]
	path, err := glide.Find(branch)
	if err != nil {
		return err
	}
	if path == nil {
		return fmt.Errorf("Error: no worktree found for branch '%s'", branch)
	}

	fmt.Println(*path)
	return nil
}

func (_ findCommand) predict(args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	entries, err := glide.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	branches := make([]string, 0, len(entries))
	for _, it := range entries {
		branches = append(branches, it.Ref.Short())
	}
	return branches, cobra.ShellCompDirectiveNoFileComp
}

type addCommand struct{}

var add addCommand

type addContext struct {
	args []string
}

func (_ addCommand) execute(ctx addContext) error {
	if len(ctx.args) == 0 {
		return fmt.Errorf("Error: branch name required")
	}
	if len(ctx.args) > 1 {
		return fmt.Errorf("Error: unexpected argument '%s'", ctx.args[1])
	}

	branch := ctx.args[0]
	path, err := glide.Find(branch)
	if err != nil {
		return err
	}
	if path != nil {
		return fmt.Errorf("Error: worktree already exists for branch '%s' at '%s'", branch, *path)
	}

	target, err := glide.Prepare(branch)
	if err != nil {
		return err
	}
	status := glide.Add([]string{target, branch})
	if !status.Success {
		return fmt.Errorf("Error: failed to create worktree")
	}
	return nil
}

func (_ addCommand) predict(args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	branches, err := completion.branches()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	entries, err := glide.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	inWorktree := make(map[string]bool, len(entries))
	for _, it := range entries {
		inWorktree[it.Ref.Short()] = true
	}
	result := make([]string, 0, len(branches))
	for _, b := range branches {
		if !inWorktree[b] {
			result = append(result, b)
		}
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}

type removeCommand struct{}

var remove removeCommand

type removeContext struct {
	args []string
}

func (_ removeCommand) execute(ctx removeContext) error {
	if len(ctx.args) == 0 {
		return fmt.Errorf("Error: branch name required")
	}

	for _, branch := range ctx.args {
		path, err := glide.Find(branch)
		if err != nil {
			return err
		}
		if path == nil {
			return fmt.Errorf("Error: no worktree found for branch '%s'", branch)
		}

		status := glide.Remove(*path)
		if !status.Success {
			return fmt.Errorf("Error: failed to remove worktree")
		}
	}
	return nil
}

func (_ removeCommand) predict(args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	entries, err := glide.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	branches := make([]string, 0, len(entries))
	for _, it := range entries {
		branches = append(branches, it.Ref.Short())
	}
	return branches, cobra.ShellCompDirectiveNoFileComp
}

type completionCommand struct{}

var completion completionCommand

func (_ completionCommand) branches() ([]string, error) {
	res := shell.Capture("git", "branch", "--format=%(refname:short)")
	if !res.Success {
		if res.Stderr != "" {
			return nil, fmt.Errorf("%s", res.Stderr)
		}
		return nil, fmt.Errorf("Error: failed to list branches")
	}
	branches := make([]string, 0)
	for _, line := range strings.Split(res.Stdout, "\n") {
		if strings.TrimSpace(line) != "" {
			branches = append(branches, strings.TrimSpace(line))
		}
	}
	return branches, nil
}
