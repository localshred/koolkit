# koolkit

koolkit is `kk`, a self-contained dev-utility plugin manager, in the style
of `krew` or `gh extension`. `kk` itself does one thing: it installs and
manages standalone plugin binaries in a managed directory and execs them
by name. All actual functionality (going to a directory, listing PRs,
driving Claude, etc.) lives in separate plugin binaries that `kk`
downloads and dispatches to.

## Install

```
brew install localshred/tap/koolkit
```

## The plugin model

A koolkit plugin is a separate repository that builds a binary named
`kk-<name>` (for example `kk-goto`) and publishes it as a GitHub release
asset via goreleaser. Plugins are not compiled into `kk`; they are
downloaded on demand and run as subprocesses.

`kk` finds out what a plugin does and needs through two flags every
plugin must support:

- `--kk-describe` — prints a one-line description, shown by `kk list`.
- `--kk-requires` — prints, one per line, the external tools the plugin
  needs at runtime, checked by `kk doctor`.

When you run `kk <name> [args...]`, and `<name>` is not a built-in `kk`
command, `kk` looks up `kk-<name>` in its managed plugin directory (falling
back to `PATH` for local development), execs it with your args passed
through untouched, and forwards its exit code.

Plugins are resolved through registries: TOML files listing plugin name,
GitHub repo, and release asset pattern. `kk` ships pointed at the public
koolkit registry by default; see `docs/usage.md` for how to add others.

## Quick start

```
brew install localshred/tap/koolkit
kk install goto
kk goto
```

## Public plugins

The default registry (`registry.toml` in this repo) currently lists:

- `goto` — `localshred/kk-goto`
- `prs` — `localshred/kk-prs`
- `claude` — `localshred/kk-claude`

Install any of them with `kk install <name>`.

## More

See [docs/usage.md](docs/usage.md) for the full command reference,
including `install`, `list`, `uninstall`, `registry add`/`list`/`remove`,
`bootstrap`, and `doctor`.
