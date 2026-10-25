# Gitframe

Gitframe is a Go command-line tool for Git contributor statistics. It reports lines, commits, and files per author for a selected repository snapshot.

## Installation

Requires **Go 1.27.2+** and **Git** on your `PATH`.

```sh
go install github.com/punkertron/gitframe/cmd/gitframe@latest
```

Or, from a local checkout:

```sh
go install ./cmd/gitframe
```

Make sure your Go binary directory (`GOBIN`, or `GOPATH/bin` by default) is on your `PATH`.

## Usage

```sh
# Analyze Go and Markdown files in the current repository.
gitframe --extensions='.go,.md'

# Export Go contributor statistics from another repository.
gitframe --repository=../my-project --languages=Go --format=json --progress=false
```

Example output:

```text
Name  Lines Commits Files
Alice 1200  18      12
Bob   450   7       5
```

## Flags

| Flag | Default | Description |
| --- | --- | --- |
| `--repository` | `.` | Path to the Git repository. |
| `--revision` | `HEAD` | Commit, branch, or tag to analyze. |
| `--order-by` | `lines` | Sort by `lines`, `commits`, or `files`. |
| `--use-committer` | `false` | Group contributions by committer instead of author. |
| `--format` | `tabular` | Output as `tabular`, `csv`, `json`, or `json-lines`. |
| `--extensions` | None | Comma-separated extensions, including the dot: `.go,.md`. |
| `--languages` | None | Comma-separated language names, such as `Go,Markdown`. |
| `--exclude` | None | Comma-separated glob patterns for paths to skip. |
| `--restrict-to` | None | Comma-separated glob patterns for paths to include. |
| `--progress` | `true` | Show progress; use `--progress=false` to disable it. |

Use `--help` for help or `--version` for the version.

## Behavior

- Counts use committed files at the selected revision. Uncommitted code changes are excluded, and the repository is left unchanged.
- Lines, including blank lines, belong to their last modifying author. Commits and files are counted uniquely per person; a file can count toward several authors. Empty files contribute zero lines and use their last changing commit.
- Default sorting is by lines, commits, then files, all descending; final ties use names ascending. `--order-by` moves the chosen field first.
- Filters combine; exclusions take priority. Globs match repository-relative paths, and `*` does not cross directory separators. Quote patterns to prevent shell expansion.
- Language names are case-insensitive and use the [embedded extension mapping](configs/language_extensions.json). Unknown languages and languages without extensions are warned about and ignored.
- Processing uses bounded concurrency and cancels remaining work on the first error. Results go to stdout; progress, warnings, and errors go to stderr. Progress is cleared on completion. Failures return a nonzero exit code.
