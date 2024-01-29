# go-hooks: Go linter selection as managed pre-commit hooks

This repository provides a small collection of [pre-commit](https://pre-commit.com/) hooks to build,
test, and lint Go code.

These hooks use the native Go environment that `pre-commit` provides. As a consequence, `pre-commit`
automatically installs and updates them into an isolated environment. This plays to `pre-commit`'s
strengths as being a package manager for deterministic, portable hooks. Most other existing hooks
for Go do not leverage `pre-commit`'s builtin Go support and depend on a system-wide or `$GOBIN`
installation of the tools instead.

## Supported hooks

The following hooks are currently supported:

- `build` to make sure commits always compile, runs `go build ./...` by default.
- `test` to ensure the tests pass, runs `go test ./...` by default.
- `vet` for running Go's integrated linter, runs `go vet ./...` by default.
- `mod-tidy` to run `go mod tidy`.
- [`goimports`](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) for formatting and managing
  package `import`s. By default, this hook mutates the content of files.
- [`staticcheck`](https://staticcheck.dev), a popular, configurable linter with many checks.
- [`errcheck`](https://github.com/kisielk/errcheck), which enforces error checking.
- [`godot`](https://github.com/tetafro/godot), which ensures API comments end with a dot.
- [`revive`](https://revive.run), an established `go lint` successor.
- [`go-critic`](https://github.com/go-critic/go-critic) for additional checks that are not provided
  by other common tools.
- [`errname`](https://github.com/Antonboom/errname) to enforce naming conventions around error
  variables and types.
- [`ineffassign`](https://github.com/gordonklaus/ineffassign) to complain about ineffectual
  assignments.

If you favourite linter is missing, please open an issue or file a PR!

## Setup
Make sure `pre-commit` is installed on your system according to the [installation
instructions](https://pre-commit.com/#install). Then, pull in the hooks by editing
`.pre-commit-config.yaml` in your repository root. For example:
```yaml
repos:
  - repo: https://github.com/lubgr/go-pre-commit-hooks
    rev: v0.1.0
    hooks:
      - id: goimports
      - id: build
      - id: test
        args: [-cover, -race, ./...]
      - id: vet
      - id: mod-tidy
      - id: staticcheck
      - id: errcheck
        args: [-ignoretests, -blank, ./...]
      - id: godot
      - id: go-critic
        args: [check, '-enable=#diagnostic,#performance', ./...]
      - id: errname
      - id: ineffassign
```
Run `pre-commit install` to pull in the hooks into their own isolated and cached environment. You
probably want to run `pre-commit run --all-files` next, to get an idea what this configuration
complains about.

The hooks are configured to run the tools with their default behaviour. You can override this by
specifying the arguments `args: [...]`. For example, to pass a configuration file to `revive`:
```yaml
      - id: revive
        args: [-config, .revive-config.toml, ./...]
```
Or, to have `goimports` list non-compliant files instead of overwriting them:
```yaml
      - id: goimports
        args: [-l]
```

## Conceptual Limitations

Go tooling works with packages, and individual files are mostly irrelevant. This contrasts
`pre-commit`'s approach, where a hook is expected to operate on one or more files. In fact,
`pre-commit` can end up passing files from different packages to a hook - most Go tools will refuse
to process these. Considering this mismatch, the default of the hooks in this repository is to
ignore specific files (`pass_filenames: false`) and run the tool on the entire repository, most
often using the `./...` argument. This is simple enough, but has downsides (they can partly be
addressed, but not in a very satisfying way):

- It is not suitable for large repositories. One way to address this is to duplicate hooks in
  `.pre-commit-config.yaml` with different per-package filters based on the `files` regex, and
  passing the relevant package through the `args` array rather than the default `./...`. When the
  number of packages grow, the hook entries for `.pre-commit-config.yaml` file should be generated.
- Untracked files impact the tools, e.g. issues are found in untracked files meant for later
  commits, or warnings about unused symbols are omitted because they are used in untracked files. In
  order to reproduce the exact outcome of a `pre-commit` run e.g. in CI, one can automate a shallow
  clone of the repository into a temporary directory, port staged changes but not untracked files,
  and run `pre-commit` there.

## Supported Hooks

The number of hooks provided here is not exhaustive, and is not meant to be exhaustive (issues/PRs
are welcome, though!).

As a reference, [`golangci-lint`](https://golangci-lint.run) provides many more linters along with
pleasant usability and a unified configuration approach. However, the number of integrated tools
seem to complicate builds, and using `go install` is [discouraged
from](https://golangci-lint.run/usage/install/#install-from-source) (the executable can crash
straight away). Hence, `golangci-lint` requires system-wide or otherwise external installation and
cannot be integrated as a native Go `pre-commit` hook.
