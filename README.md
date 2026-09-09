# reusable-tui

A production-grade base for terminal UI applications in Go. It ships a full
application shell (tabs, focus, jump mode, command palette, keymap registry,
theming, SQLite persistence, migrations, notifications, toasts) and exactly one
throwaway example feature that exercises every pattern end to end. Delete the
example, keep the shell.

## Stack

| Concern | Package |
| --- | --- |
| Runtime | `charmbracelet/bubbletea` v1 |
| Styling | `charmbracelet/lipgloss` v1 |
| Widgets | `charmbracelet/bubbles` v1 |
| Fuzzy search | `sahilm/fuzzy` |
| Database | `modernc.org/sqlite` (cgo-free) |
| Migrations | `pressly/goose/v3`, embedded with `go:embed` |
| Paths | `adrg/xdg` |
| Secrets | `zalando/go-keyring` |
| Notifications | `gen2brain/beeep` |
| Tests | `charmbracelet/x/exp/teatest`, `charmbracelet/x/exp/golden` |
| Tooling | `golangci-lint`, `goreleaser` |

## File tree

```
.github/workflows/ci.yml
.golangci.yml
.goreleaser.yaml
Makefile
README.md
cmd/tui/main.go
internal/app/command.go
internal/app/error.go
internal/app/focus.go
internal/app/format.go
internal/app/golden_test.go
internal/app/help_screen.go
internal/app/helper_test.go
internal/app/root.go
internal/app/shell_test.go
internal/app/testdata/TestGoldenExampleScreen/high-contrast.golden
internal/app/testdata/TestGoldenExampleScreen/monochrome.golden
internal/app/testdata/TestGoldenExampleScreen/violet-dark.golden
internal/app/update.go
internal/app/view.go
internal/config/config.go
internal/features/example/message.go
internal/features/example/model.go
internal/features/example/mutation.go
internal/features/example/query.go
internal/features/example/update.go
internal/features/example/view.go
internal/keymap/bindings.go
internal/keymap/registry.go
internal/notify/beeep.go
internal/notify/notify.go
internal/secret/secret.go
internal/store/errors.go
internal/store/migrate.go
internal/store/migrations/00001_example_items.sql
internal/store/store.go
internal/theme/high_contrast.go
internal/theme/monochrome.go
internal/theme/registry.go
internal/theme/theme.go
internal/theme/violet.go
internal/ui/command.go
internal/ui/hatch.go
internal/ui/header.go
internal/ui/jump.go
internal/ui/layout.go
internal/ui/message.go
internal/ui/modal.go
internal/ui/overlay.go
internal/ui/palette.go
internal/ui/panel.go
internal/ui/screen.go
internal/ui/skeleton.go
internal/ui/statusbar.go
internal/ui/text.go
internal/ui/tick.go
```

## Architecture

The root model in `internal/app` owns focus, tab selection, the active theme,
the toast, the palette and jump mode. It composes screens and routes messages.

```
internal/app  ->  internal/features/*  ->  internal/{ui,theme,keymap,store,config,notify,secret}
```

Dependencies flow one way. Features import internal packages; internal packages
never import features. Screens implement `ui.Screen`, which lives in `internal/ui`
so features never need to import `internal/app`.

Features never reach into each other. They emit `tea.Msg` values (`ui.ErrorMsg`,
`ui.SuccessMsg`, `ui.NotifyMsg`, `ui.ThemeMsg`, `ui.FocusMsg`) that the root
model dispatches.

Every async operation is a `tea.Cmd` returning a typed message. No goroutine ever
writes to a model field.

### Keys

| Key | Action |
| --- | --- |
| `tab` / `shift+tab` | next / previous screen |
| `ctrl+n` / `ctrl+p` | next / previous panel |
| `f` | jump mode, then press a panel's hint letter |
| `ctrl+k` | command palette |
| `?` | help screen |
| `q` / `ctrl+c` | quit |

### Data and errors

SQLite lives under `xdg.DataHome`. Migrations are embedded and run on startup.
Repository functions return `(T, error)` with the sentinel errors in
`internal/store/errors.go`. Nothing panics into the UI: errors become
`ui.ErrorMsg`, and the root model renders them as a transient toast above the
status bar.

### Rendering

View code reads semantic tokens only. There is no hex string and no magic
spacing number outside `internal/theme`. Spacing is expressed in
`theme.SpaceTokens` units: one cell horizontally, one row vertically.

Color profile is detected with `lipgloss.ColorProfile()` and mapped to a
`theme.Fidelity`. In the no-color fidelity, focus stops depending on hue: the
focused border switches to a thick box-drawing set and the title is wrapped in
brackets.

## Running

```
make run
make test
make lint
make golden      # rewrite the golden files after an intentional visual change
make snapshot    # goreleaser static binaries for darwin/linux, amd64/arm64
```

## Adding a screen

1. Create `internal/features/<name>/` with `model.go`, `update.go`, `view.go`,
   `query.go`, `mutation.go`, `message.go`.
2. Implement `ui.Screen` on a pointer model:

```go
func (m *Model) ID() string          { return ScreenID }
func (m *Model) Title() string       { return ScreenTitle }
func (m *Model) Panels() []string    { return []string{keymap.PanelReportsList} }
func (m *Model) Init() tea.Cmd       { return loadReports(m.db) }
func (m *Model) Focus(id string)     { m.focused = id }
func (m *Model) Capturing() bool     { return m.editing }
func (m *Model) Update(tea.Msg) tea.Cmd
func (m *Model) View(ui.RenderContext) string
```

3. Append it to the `screens` slice in `internal/app/root.go`. Tab cycling, focus
   state, the status bar, the help view and a `navigate` palette entry all follow
   automatically.

`Capturing()` tells the root model that the screen owns raw keystrokes, so global
single-letter bindings are suppressed while a text input is focused.

## Adding a panel

Panels are values, not types. Build one in a screen's `View`:

```go
panel := ui.Panel{
    Title:   "Reports",
    Badge:   strconv.Itoa(len(m.reports)),
    Hint:    ctx.Jump.Hint(keymap.PanelReportsList),
    Width:   widths[0],
    Height:  ctx.Height,
    Focused: ctx.Focused == keymap.PanelReportsList,
}
panel.Body = m.listBody(ctx.Theme, panel.ContentWidth(ctx.Theme), panel.ContentHeight(ctx.Theme))
return ui.RenderPanel(ctx.Theme, panel)
```

Then add the panel id to the screen's `Panels()` so focus cycling and jump mode
reach it. Use `ui.SplitWidths` and `ui.Row` for multi-panel layouts, and
`ui.EmptyState` / `ui.Skeleton` for the empty and loading bodies. Never hardcode
a width: everything derives from `ctx.Width` and `ctx.Height`, which come from
`tea.WindowSizeMsg`.

## Adding a keybinding

Touch exactly one file: `internal/keymap/bindings.go`.

1. Add the field to `GlobalKeys`, `ExampleKeys` or your feature's key struct.
2. Give it a value in `Default()`.
3. List it in the relevant `r.Register(...)` call inside `Registry()`.

The status bar and the help screen are both generated from that registry, so
neither needs editing. Match it in an update function with
`key.Matches(msg, m.keys.Reports.Archive)`.

## Adding a theme

1. Create `internal/theme/<name>.go` returning a `Theme` from a named function.
   Fill every token; reuse `defaultSpace()` and `defaultMarkers()` unless the
   theme genuinely changes spacing or glyphs.
2. Register it in `Builtin()` in `internal/theme/registry.go`.

It immediately appears in the palette as `Theme: <name>`, is persisted to config
when selected, and is picked up by the golden test, which renders the example
screen once per registered theme.

Every semantic state carries a color, a text label and an ASCII glyph, so no
state is signalled by color alone.

## Adding a command to the palette

Commands are `{id, label, group, tea.Cmd}`. Register them in
`registerCommands` in `internal/app/command.go`:

```go
m.commands.Add(ui.Command{
    ID:    "reports.export",
    Label: "Export reports",
    Group: "reports",
    Run:   exportReports(),
})
```

The palette fuzzy-matches against `group + " " + label`. Theme switching and
screen navigation are registered the same way, which is what proves the registry
is the only path a command needs.

## Tests

`go test ./...` covers, through `teatest`, tab cycling, focus cycling, jump mode,
palette search, theme switching, window resize and the minimum-size guard. A
golden-file test renders the example screen once per builtin theme, so an
unintended aesthetic change fails CI. Regenerate the goldens with `make golden`
only when the change is intentional.

## Removing the example

Delete `internal/features/example/`, drop the `example.New(...)` entry from the
`screens` slice in `internal/app/root.go`, remove the `PanelExample*` constants
and `ExampleKeys` from `internal/keymap/bindings.go`, and delete
`internal/store/migrations/00001_example_items.sql`. Nothing else refers to it.
