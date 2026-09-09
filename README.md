# reusable-tui

A production-grade base for terminal UI applications in Go. It ships a full
application shell (tabs, focus, jump mode, command palette, keymap registry,
theming, mouse support, session restore, busy indicator, SQLite persistence,
migrations, notifications, toasts) and exactly one
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
cmd/tui/flags.go
cmd/tui/flags_test.go
cmd/tui/main.go
internal/app/command.go
internal/app/error.go
internal/app/focus.go
internal/app/format.go
internal/app/golden_test.go
internal/app/help_screen.go
internal/app/helper_test.go
internal/app/mouse.go
internal/app/mouse_test.go
internal/app/root.go
internal/app/session_test.go
internal/app/shell_test.go
internal/app/testdata/TestGoldenDashboardScreen/high-contrast.golden
internal/app/testdata/TestGoldenDashboardScreen/monochrome.golden
internal/app/testdata/TestGoldenDashboardScreen/violet-dark.golden
internal/app/testdata/TestGoldenExampleScreen/high-contrast.golden
internal/app/testdata/TestGoldenExampleScreen/monochrome.golden
internal/app/testdata/TestGoldenExampleScreen/violet-dark.golden
internal/app/theme_test.go
internal/app/update.go
internal/app/view.go
internal/config/config.go
internal/features/dashboard/model.go
internal/features/dashboard/query.go
internal/features/dashboard/update.go
internal/features/dashboard/view.go
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
internal/theme/file.go
internal/theme/high_contrast.go
internal/theme/monochrome.go
internal/theme/registry.go
internal/theme/theme.go
internal/theme/token.go
internal/theme/token_test.go
internal/theme/violet.go
internal/ui/busy.go
internal/ui/calendar.go
internal/ui/chart.go
internal/ui/command.go
internal/ui/dashboard.go
internal/ui/doc.go
internal/ui/hatch.go
internal/ui/hatch_test.go
internal/ui/header.go
internal/ui/hit.go
internal/ui/jump.go
internal/ui/layout.go
internal/ui/message.go
internal/ui/modal.go
internal/ui/overlay.go
internal/ui/palette.go
internal/ui/panel.go
internal/ui/screen.go
internal/ui/scroll.go
internal/ui/segment.go
internal/ui/skeleton.go
internal/ui/stat.go
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

On the dashboard screen, `k` / `j` move within a list, `h` / `l` step the period
or switch the view, and `t` jumps back to today.

### Mouse

Mouse support is on by default (`--no-mouse` turns it off). Clicking a header
tab switches screens, clicking a panel focuses it, and the wheel scrolls
whatever is under the pointer. Screens never see raw `tea.MouseMsg`: the root
model translates clicks and wheel events into `ui.ClickMsg` and `ui.WheelMsg`,
each carrying the panel id and content-local coordinates.

Hit testing works from a `ui.HitMap` that is rebuilt on every render. A screen
registers each panel's rectangle in body coordinates from its `View`:

```go
ctx.Hits.Add(keymap.PanelReportsList, ui.Rect{X: 0, Y: 0, Width: widths[0], Height: ctx.Height})
```

`ui.RenderDashboard` does this for you when a region's `Panel.ID` is set. A
screen that wants row selection handles `ui.ClickMsg` and reads `Y` as the row
index; one that wants wheel scrolling handles `ui.WheelMsg` and applies `Delta`
(`ui.Scroll` does that for a `viewport.Model`).

### Session restore

The active screen and the focused panel of every screen are written to the
`session` block of the config file on quit, and restored on the next launch.
Both `q` and the palette's Quit command go through the same path, so nothing
else needs to know. Unknown screen or panel ids in the file are ignored.

### Busy indicator

Long running commands announce themselves with `ui.Busy(id, label)` and clear
with `ui.Idle(id)`. While any id is busy the status bar shows a spinner and the
label, frames taken from `theme.Marker.Spinner` (ASCII in the no-color
fidelity). The example screen wraps its load:

```go
func (m *Model) reload() tea.Cmd {
    m.loading = true
    return tea.Batch(ui.Busy(busyLoad, "loading items"), loadItems(m.db))
}
```

and returns `ui.Idle(busyLoad)` from its result handler on both the success and
the error path, so a failed query never leaves the spinner running.

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

### Flags

| Flag | Effect |
| --- | --- |
| `--config path` | use this config file instead of the XDG default |
| `--db path` | use this SQLite file instead of the XDG default |
| `--theme name` | start with this theme for the run, without persisting it |
| `--no-mouse` | leave mouse reporting off |
| `--version` | print the version and exit |

Point `--config` and `--db` at a scratch directory to run a second instance or
a fixture next to your real state.

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

## The dashboard layout

`internal/features/dashboard` is an optional second example screen: a three
column dashboard with stacked panels, the shape most tracker style TUIs use.
It is data free and store free, so it stays a layout reference.

```
+-------------+---------------+-----------------------+
| Accounts    | View and add  | Overview              |
+-------------+---------------+                       |
| Insights    | Period        |                       |
+-------------+---------------+-----------------------+
```

`ui.RenderDashboard` takes stacks and regions instead of manual widths:

```go
ui.RenderDashboard(ctx.Theme, ui.Dashboard{
    Width:  ctx.Width,
    Height: ctx.Height,
    Stacks: []ui.Stack{
        {Weight: 5, Regions: []ui.Region{accounts, insights}},
        {Weight: 6, Regions: []ui.Region{mode, period}},
        {Weight: 9, Regions: []ui.Region{overview}},
    },
})
```

Stack weights split the width, region weights split each column's height, and a
region's `Body` is called with the measured content size, so a widget never
guesses its own box:

```go
ui.Region{
    Weight: 3,
    Panel:  ui.Panel{Title: "Insights", Focused: ctx.Focused == keymap.PanelDashboardInsights},
    Body:   func(t theme.Theme, width, height int) string { return chart(t, width, height) },
}
```

One stack with one region is a full screen pane, two stacks are a sidebar plus a
main area. Nothing else in the shell changes: focus cycling, jump mode, the
status bar and the help screen read the same registries as any other screen.

### Widgets

| Component | Purpose |
| --- | --- |
| `ui.RenderSegment` | inline tab strip inside a panel |
| `ui.RenderPeriod` | `<<< label >>>` period stepper |
| `ui.RenderCalendar` | month grid with a selected day and a week band |
| `ui.RenderStats` | label and value pairs across a row |
| `ui.RenderBarChart` | block glyph bars with optional labels |
| `ui.RenderDoc` | headings, wrapped text, bullets, rules and callouts |

They take a `theme.Theme` and explicit sizes, so they work in any panel, not
only in this screen.

### Removing it

Delete `internal/features/dashboard/`, drop the `dashboard.New(...)` entry from
the `screens` slice in `internal/app/root.go`, and remove the `PanelDashboard*`
constants and `DashboardKeys` from `internal/keymap/bindings.go`. The layout and
widget primitives in `internal/ui` stay usable on their own.

## Adding a keybinding

Touch exactly one file: `internal/keymap/bindings.go`.

1. Add the field to `GlobalKeys`, `ExampleKeys` or your feature's key struct.
2. Give it a value in `Default()`.
3. List it in the relevant `r.Register(...)` call inside `Registry()`.

The status bar and the help screen are both generated from that registry, so
neither needs editing. Match it in an update function with
`key.Matches(msg, m.keys.Reports.Archive)`.

## Theming

Colors are fully customizable at runtime. There are three routes, in
increasing order of effort.

### 1. Override tokens in the config file

`~/.config/reusable-tui/config.json` can recolor any registered theme without
replacing it. Overrides are keyed by theme name and survive theme switching.

```json
{
  "theme": "violet-dark",
  "overrides": {
    "violet-dark": {
      "accent.active": "#22D3EE",
      "border.focused": "#22D3EE"
    }
  }
}
```

### 2. Drop a theme file in the themes directory

Any `*.json` file in `~/.config/reusable-tui/themes/` becomes a registered
theme, listed in the palette like a builtin. `extends` inherits every token
from an existing theme so you only state what differs.

```json
{
  "name": "cyan",
  "extends": "violet-dark",
  "tokens": {
    "accent.active": "#22D3EE",
    "accent.mid": "#1A9EB4",
    "accent.dim": "#123F49",
    "border.focused": "#22D3EE"
  }
}
```

Omit `extends` to inherit from the default theme. Run **Export active theme to
a file** from the palette to write the current theme out in full as a starting
point, and **Reload themes from disk** to pick up edits without restarting.

A file that fails to parse or names an unknown token never blocks startup: the
theme is skipped and the problem is reported as a warning toast.

### 3. Add a theme in Go

1. Create `internal/theme/<name>.go` returning a `Theme` from a named function.
   Fill every token; reuse `defaultSpace()` and `defaultMarkers()` unless the
   theme genuinely changes spacing or glyphs.
2. Register it in `Builtin()` in `internal/theme/registry.go`.

It immediately appears in the palette, is persisted to config when selected,
and is picked up by the golden test, which renders the example screen once per
registered theme.

### Token reference

`theme.Tokens()` returns the authoritative list. Every one of them is settable
from a config override or a theme file.

| Group | Tokens |
| --- | --- |
| `base` | `background`, `surface`, `overlay` |
| `text` | `primary`, `secondary`, `muted`, `disabled`, `inverted` |
| `accent` | `active`, `mid`, `dim` |
| `border` | `focused`, `unfocused`, `subtle`, `focused.set`, `unfocused.set` |
| `status.<state>` | `fg`, `label`, `glyph` for `success`, `warning`, `danger`, `info` |
| `marker` | `empty`, `hatch`, `dot` |

Color values are `#RGB`, `#RRGGBB`, or an ANSI index `0`-`255`. Border set
values are `rounded`, `normal`, `thick`, `double`, `block` or `hidden`.

`marker.empty` picks the fill behind an empty panel body: `slash` (the default
diagonal hatch), `dots` (a sparse dot grid) or `none` (blank, label only).
`marker.hatch` and `marker.dot` swap the glyphs those two patterns draw with.

```json
{
  "theme": "violet-dark",
  "overrides": { "violet-dark": { "marker.empty": "dots" } }
}
```

Every semantic state carries a color, a text label and an ASCII glyph, so no
state is signalled by color alone; recoloring a status token never makes it
unreadable. Spacing stays in Go, in `theme.SpaceTokens`.

### Adding a token

Adding a token to the customizable set touches one file,
`internal/theme/token.go`: add the constant, list it in `tokenOrder`, and add
its case to `ReadToken` and to the matching writer. `TestEveryTokenRoundTrips`
fails if a listed token is not readable and writable, so the list cannot drift
from the switches.

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
palette search, theme switching, window resize, the minimum-size guard and the
dashboard layout's size integrity, plus token round-tripping, user theme files,
config overrides, export and reload. Golden-file tests render the example screen
and the dashboard screen once per builtin theme, so an
unintended aesthetic change fails CI. Regenerate the goldens with `make golden`
only when the change is intentional.

## Removing the example

Delete `internal/features/example/`, drop the `example.New(...)` entry from the
`screens` slice in `internal/app/root.go`, remove the `PanelExample*` constants
and `ExampleKeys` from `internal/keymap/bindings.go`, and delete
`internal/store/migrations/00001_example_items.sql`. Nothing else refers to it.
