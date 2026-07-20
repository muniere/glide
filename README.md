# glide

CLI to manage Git worktrees.

## Requirements

- Go 1.22+
- Git (with `git worktree`)

## Build

```sh
make build
```

This creates `./glide`.

## Run

```sh
./glide --help
```

Or without building:

```sh
make run ARGS="list"
```

## Commands

- `glide list` / `glide ls`
- `glide list --porcelain`
- `glide find <branch>`
- `glide add <branch>`
- `glide remove <branch> [branch ...]`
- `glide rm <branch> [branch ...]`

## Examples

```sh
glide list
glide find feature/login
glide add feature/login
glide remove feature/login
```

By default, `add` places worktrees under a `.worktrees/` directory inside the repository, preserving the branch name as a nested path (e.g. `.worktrees/feature/login`).

## Configuration

Place a config file at `$XDG_CONFIG_HOME/glide/config` (global) or `.glide/config` (per-project). The local config takes precedence.

### hierarchy (default)

Worktrees are placed as children of a root container directory, preserving the
branch name (including path separators) as a nested path.
`root` accepts an absolute path or a relative path from the repository root,
and defaults to `.worktrees` (relative to the repository root).

```toml
strategy = "hierarchy"

[hierarchy]
root = ".worktrees"           # default: ".worktrees" (relative to repo root)
# root = "/path/to/container" # absolute path
```

With a relative path (`root = ".worktrees"`, the default):

```
~/projects/myrepo/
├── .worktrees/
│   ├── main/
│   ├── feature/
│   │   ├── login/
│   │   └── signup/
│   └── bugfix/
│       └── issue-42/
└── ... (source files)
```

With an absolute path (`root = "/path/to/container"`):

```
/path/to/container/
├── main/
├── feature/
│   ├── login/
│   └── signup/
└── bugfix/
    └── issue-42/
```

### flat

Worktrees are placed as siblings of the repository directory,
with slashes in branch names replaced by `separator`.

```toml
strategy = "flat"

[flat]
separator = "-"  # default: "-"
```

```
~/projects/
├── myrepo/          # main worktree
├── myrepo-main/
├── myrepo-feature-login/
└── myrepo-feature-signup/
```

## Shell Completion

### zsh

Add `share/zsh/site-functions` to your `fpath`:

```zsh
fpath=(/path/to/glide/share/zsh/site-functions $fpath)
autoload -U compinit && compinit
```

### bash

Source `share/bash-completion/completions/glide`:

```bash
source /path/to/glide/share/bash-completion/completions/glide
```

## Development

```sh
make fmt
make test
make tidy
make clean
```
