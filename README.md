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

By default, `add` creates a sibling directory named like `<repo>-<branch>` (slashes in branch names are replaced with `-`).

## Configuration

Place a config file at `$XDG_CONFIG_HOME/glide/config` (global) or `.glide/config` (per-project). The local config takes precedence.

### flat (default)

Worktrees are placed as siblings of the repository directory.

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

### hierarchy

Worktrees are placed as children of a root container directory.
`root` accepts an absolute path or a relative path from the repository root.

```toml
strategy = "hierarchy"

[hierarchy]
root = "/path/to/container"   # absolute path
# root = ".worktrees"         # relative to repo root
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

With a relative path (`root = ".worktrees"`):

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
