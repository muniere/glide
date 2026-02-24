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

## Shell Completion

**zsh**: Add `share/zsh/site-functions` to your `fpath`:

```zsh
fpath=(/path/to/glide/share/zsh/site-functions $fpath)
autoload -U compinit && compinit
```

**bash**: Source `share/bash-completion/completions/glide`:

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
