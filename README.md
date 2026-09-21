# reusable-tui

A base for terminal applications in Go. It ships an application shell — routing
and a section bar, a focus ring with jump mode, a command palette, a generated
help screen, a keymap registry that refuses to build on a conflict, four
themes, a responsive layout solver and a widget set — plus a demo that
exercises all of it. You scaffold a project from it; you do not depend on it.

Built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 and
[Lip Gloss](https://github.com/charmbracelet/lipgloss) v2.

## Start a project

```sh
go run github.com/remcostoeten/reusable-tui/cmd/create-tui@latest myapp
```

It asks for a module path and an application name, downloads the template,
rewrites every import to your module, replaces the demo with one starter
screen, runs `go mod tidy` and commits. Then:

```sh
cd myapp && make run
```

Non-interactively:

```sh
create-tui -module github.com/you/myapp -name myapp -y myapp
create-tui -demo -ref v0.2.0 playground   # keep the demo, pin a tag
```

| Flag | Meaning |
| --- | --- |
| `-module` | Module path for the new project. Defaults to a guess from your git config. |
| `-name` | Application and binary name. Defaults to the directory name. |
| `-demo` | Keep `home`, `manager` and `settings` instead of pruning to one screen. |
| `-ref` | Branch or tag to scaffold from. Defaults to `master`. |
| `-local` | Scaffold from a checkout on disk instead of downloading. |
| `-y` | Take every default without prompting. |

The scaffolder imports nothing outside the standard library, so `go run` on it
resolves in a second rather than pulling the shell's dependency graph.

## Migrate an existing CLI

An agent skill, `reusable-tui-migrate`, ports an existing Go CLI onto this
template. It scaffolds a project, moves the domain logic over, turns
subcommands into modules and keeps headless flags working. Install it in the
old CLI's repository:

```sh
npx skills add remcostoeten/reusable-tui
```

Then ask your agent to migrate the CLI. Scaffolded projects already ship the
skill in `.claude/skills/`.

## What a scaffolded project looks like

```
cmd/myapp/main.go                 entrypoint: config, logging, runtime
internal/app/app.go               composition root
internal/app/register.go          the module list — one line per screen
internal/modules/hello/           your first screen
internal/tui/                     the shell; you rarely edit this
```

Adding a screen is a directory under `internal/modules` and a line in
`register.go`. A module registers its routes, commands, key bindings, status
items and themes into a `Registrar`; the palette, the help overlay and the
status bar are generated from that registry, so a verb you add shows up in all
three without being listed anywhere else. `app.Build` freezes the registry and
reports duplicate routes, malformed bindings and keymap conflicts before the
first frame is drawn — which is why `TestBuildSucceeds` is worth as much as it
is.

## Keys

`tab` cycles focus · `v` badges every region and jumps to the one you press ·
`alt+hjkl` moves geometrically · `[` and `]` change section · `ctrl+k` opens the
palette · `?` lists every command with the key that reaches it · `ctrl+t` cycles
`dark`, `dim`, `high-contrast` and `ascii` · `q` quits.

The `ascii` theme has no colour at all. If the app is usable there, nothing is
being said with colour alone.

## Development

```sh
make run      # go run ./cmd/tui
make test     # go test ./...
make golden   # re-record the golden frame
make lint     # golangci-lint run ./...
```

`internal/tui/boundaries_test.go` enforces the layering — the shell never
imports a module, the runtime never imports the shell — so the architecture is
checked rather than described.

## Releasing

`go run ...@latest` resolves the newest semver tag, so cutting one is what makes
the command above work for other people:

```sh
git tag v0.1.0 && git push --tags
```

## License

MIT. See [LICENSE](LICENSE).
