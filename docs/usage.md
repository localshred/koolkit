# kk usage

`kk` has a small set of built-in commands for managing plugins, plus
plugin dispatch: any other bare name is looked up as an installed plugin
and exec'd. This doc covers the built-ins.

## install

Install one or more plugins by name (resolved against configured
registries), or by a direct `github.com/owner/repo` path.

```
kk install goto
kk install goto prs claude
kk install github.com/localshred/kk-goto
```

The direct-path form is for plugins not listed in any registry: `kk`
derives the plugin name from the repo (stripping a `kk-` prefix) and
downloads `<repo>_<os>_<arch>.tar.gz` from that repo's latest GitHub
release, the same asset layout goreleaser produces for registry plugins.

Installing verifies the downloaded archive against the release's
`checksums.txt`, extracts the `kk-<name>` binary, and places it into the
managed plugin directory (see below).

## uninstall

Remove one or more installed plugins.

```
kk uninstall goto
kk uninstall goto prs
```

This deletes `kk-<name>` from the managed plugin directory and drops it
from the install manifest.

## list

Show installed plugins (with their `--kk-describe` text) and plugins
available in configured registries but not yet installed.

```
kk list
```

## registry

Manage the set of registry sources `kk install`, `kk list`, and
`kk bootstrap` resolve plugins against. A source is either an HTTP(S) URL
to a registry TOML file or a local file path.

```
kk registry list
kk registry add https://raw.githubusercontent.com/someorg/kk-plugins/main/registry.toml
kk registry remove https://raw.githubusercontent.com/someorg/kk-plugins/main/registry.toml
```

Sources are persisted in the user config file (see below). By default
`kk` is configured with one source, the public koolkit registry:

```
https://raw.githubusercontent.com/localshred/koolkit/main/registry.toml
```

### How registries resolve

Each registry is a TOML file listing plugins:

```
[[plugin]]
name = "goto"
repo = "localshred/kk-goto"
asset = "kk-goto_{{.Os}}_{{.Arch}}.tar.gz"
```

`repo` is `owner/repo` on the release host, and `asset` is a Go template
for the goreleaser archive name, expanded with `.Os` (title-cased, e.g.
`Darwin`, `Linux`) and `.Arch` (lowercase, e.g. `amd64`, `arm64`).

Multiple sources merge like brew taps: `kk` loads every configured source
and flattens them into one name-to-plugin map, with earlier sources
taking precedence on name collisions. A source that fails to load (bad
URL, malformed TOML) is skipped with a warning, not fatal.

Each plugin entry may also set `host`, the release host to fetch
from, defaulting to `github.com` when unset. This is what would let a
registry point at plugins hosted on another GitHub Enterprise instance
(e.g. an internal GHE) instead of the public github.com. No such internal
registry or internal plugins are wired up yet; that is future work.

## bootstrap

Fuzzy-select from all registry plugins not already installed, and
install whichever ones you pick.

```
kk bootstrap
```

If everything available is already installed, it reports that and exits
without prompting.

## doctor

Report on the "subject tools" installed plugins declare via
`--kk-requires` (e.g. a plugin that wraps `git` declaring it needs `git`
on `PATH`), showing whether each is found and which plugins need it.

```
kk doctor
```

If no installed plugin declares any requirements, it reports that and
exits.

## Managed plugin directory and config paths

`kk` follows the XDG Base Directory spec:

- Plugin binaries (`kk-<name>`) and the install manifest
  (`manifest.toml`) live in `$XDG_DATA_HOME/koolkit/bin`, defaulting to
  `~/.local/share/koolkit/bin` when `XDG_DATA_HOME` is unset.
- User config (`config.toml`, holding registry sources) lives in
  `$XDG_CONFIG_HOME/koolkit`, defaulting to `~/.config/koolkit` when
  `XDG_CONFIG_HOME` is unset.

## Plugin dispatch

Any invocation `kk <name> [args...]` where `<name>` is not a built-in
command or alias is treated as a plugin call. `kk` looks up `kk-<name>`
in the managed plugin directory, falling back to `PATH` (for local plugin
development), execs it with the remaining args passed through untouched,
and exits with the plugin's own exit code. Plugin flags are never parsed
by `kk`; only a bare name in the first argument position triggers
dispatch, so a plugin can use any flags it wants, including ones that
collide with `kk`'s own.

```
kk goto some-project
kk prs --mine
```

If no matching plugin binary is found, `kk` prints a hint to install it
or run `kk list`, and exits non-zero.
