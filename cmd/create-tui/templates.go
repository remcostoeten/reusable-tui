package main

import (
	"fmt"
	"strings"
)

// readme writes the new project's own README, replacing this repository's —
// which documents the scaffolder and would survive rewriting as nonsense.
//
// It runs after the rewrite, so the link back to the template survives.
func readme(o options) string {
	screens := "internal/modules/hello/"
	if o.Demo {
		screens = "internal/modules/*/"
	}

	return fmt.Sprintf(`# %[1]s

A terminal application, scaffolded from
[reusable-tui](https://github.com/%[4]s).

## Run

    make run

## Layout

%[2]s

Adding a screen is a directory under %[3]sinternal/modules%[3]s and a line in
%[3]sregister.go%[3]s. A module registers its routes, commands, key bindings and
status items into a %[3]sRegistrar%[3]s; the command palette, the help overlay and
the status bar are generated from that registry, so a verb you add shows up in
all three without being listed anywhere else.

%[3]sapp.Build%[3]s freezes the registry and reports duplicate routes, malformed
bindings and keymap conflicts before the first frame is drawn.

## Keys

%[3]stab%[3]s cycles focus · %[3]sv%[3]s badges every region and jumps to the one you press ·
%[3]salt+hjkl%[3]s moves geometrically · %[3]s[%[3]s and %[3]s]%[3]s change section · %[3]sctrl+k%[3]s opens the
palette · %[3]s?%[3]s lists every command · %[3]sctrl+t%[3]s cycles %[3]sdark%[3]s, %[3]sdim%[3]s, %[3]shigh-contrast%[3]s
and %[3]sascii%[3]s · %[3]sq%[3]s quits.

## Development

    make test     # go test ./...
    make lint     # golangci-lint run ./...
    make build    # bin/%[1]s

%[3]sinternal/tui/boundaries_test.go%[3]s enforces the layering: the shell never
imports a module, the runtime never imports the shell.
`, o.Name, tree(o.Name, screens), "`", sourceRepo)
}

// tree renders the layout block as an aligned two-column list, so that a long
// application name does not shear the comments out of line.
func tree(name, screens string) string {
	rows := [][2]string{
		{"cmd/" + name + "/main.go", "entrypoint: config, logging, runtime"},
		{"internal/app/app.go", "composition root"},
		{"internal/app/register.go", "the module list — one line per screen"},
		{screens, "your screens"},
		{"internal/tui/", "the shell; you rarely edit this"},
	}

	width := 0
	for _, row := range rows {
		if n := len([]rune(row[0])); n > width {
			width = n
		}
	}

	var b strings.Builder
	for _, row := range rows {
		pad := strings.Repeat(" ", width-len([]rune(row[0]))+2)
		b.WriteString("    " + row[0] + pad + row[1] + "\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

// The files below replace the demo when scaffolding a pruned project. They
// still carry the source module path because prune runs before rewrite, which
// is what fixes every import in one pass.
//
// None of them may contain a backtick: they are raw literals here and ordinary
// Go source in the generated tree.

const registerGo = `package app

import (
	"github.com/remcostoeten/reusable-tui/internal/modules/hello"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/config"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
)

// Deps are the services a module needs, wired here rather than resolved from a
// container: constructor arguments, checked by the compiler.
type Deps struct {
	Config config.Config
}

// Modules is the whole module list. Adding a screen is a line here plus a
// directory under internal/modules.
func Modules(deps Deps) []registry.Module {
	return []registry.Module{
		hello.New(),
	}
}
`

const helloModuleGo = `// Package hello is the starting screen. Copy it to add another.
package hello

import (
	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/command"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/input"
	navigation "github.com/remcostoeten/reusable-tui/internal/tui/navigation"
	"github.com/remcostoeten/reusable-tui/internal/tui/ui"
)

// Module is the hello screen.
type Module struct{}

// New builds the module.
func New() *Module {
	return &Module{}
}

// ID names the module. Command and region identifiers are namespaced under it.
func (m *Module) ID() registry.ModuleID {
	return "hello"
}

// Register installs everything the module contributes: a route, a verb, and
// the key that reaches the verb. The palette, the help overlay and the status
// bar are all generated from this, so nothing has to be listed twice.
func (m *Module) Register(r *registry.Registrar) {
	r.Route(navigation.Route{
		ID:    "hello",
		Title: "Hello",
		Order: 10,
	}, func(registry.Context) registry.View {
		return newView()
	})

	r.Command(command.Command{
		ID:          "hello.greet",
		Title:       "Say Hello",
		Category:    "Hello",
		Description: "Show a notification",
		When:        command.OnRoute("hello"),
		Run: func(command.Scope) tea.Cmd {
			return registry.Notify(ui.Toast{Severity: ui.SeveritySuccess, Text: "Hello"})
		},
	})

	r.Bind(input.LayerScreen, input.Binding{
		Keys:     []string{"g"},
		Command:  "hello.greet",
		When:     command.OnRoute("hello"),
		Hint:     "Greet",
		Priority: 60,
	})
}
`

const helloViewGo = `package hello

import (
	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/render"
	"github.com/remcostoeten/reusable-tui/internal/tui/focus"
	"github.com/remcostoeten/reusable-tui/internal/tui/kernel"
	"github.com/remcostoeten/reusable-tui/internal/tui/layout"
	"github.com/remcostoeten/reusable-tui/internal/tui/ui"
)

const (
	regionGreeting kernel.ID = "hello.greeting"
	regionNotes    kernel.ID = "hello.notes"
)

// View is the hello screen: two panels, so that tab has somewhere to go.
type View struct {
	notes ui.ScrollArea
	rects map[kernel.ID]kernel.Rect
}

func newView() *View {
	return &View{
		notes: ui.NewScrollArea(notesText),
		rects: map[kernel.ID]kernel.Rect{},
	}
}

// Init has nothing to load. Return a command here to fetch on mount.
func (v *View) Init() tea.Cmd {
	return nil
}

// FocusRegions declares the ring that tab cycles and v jumps into.
func (v *View) FocusRegions() []focus.Region {
	return []focus.Region{
		{ID: regionGreeting, Order: 1, Label: "Greeting", JumpKey: 'g'},
		{ID: regionNotes, Order: 2, Label: "Notes", JumpKey: 'n'},
	}
}

// RegionRects reports where the panels were drawn, which is what makes the
// regions clickable and reachable with alt+hjkl.
func (v *View) RegionRects(kernel.Rect) map[kernel.ID]kernel.Rect {
	return v.rects
}

// Update forwards to whichever panel holds focus.
func (v *View) Update(ctx registry.Context, msg tea.Msg) (registry.View, tea.Cmd) {
	next := *v
	next.notes.Focused = ctx.Focus.Focused(regionNotes)

	var cmd tea.Cmd
	next.notes, cmd = next.notes.Update(msg)
	return &next, cmd
}

// Render lays the screen out and records where each region landed.
func (v *View) Render(rc render.Context) string {
	area := rc.Rect
	rows := layout.Rows(area, layout.Fixed(5), layout.Flex(1))

	v.rects[regionGreeting] = rows[0]
	v.rects[regionNotes] = rows[1]

	// The prose is wrapped to whatever width the panel ended up with, rather
	// than being stored pre-formatted for one terminal size.
	inner := rows[1].Inset(1).InsetXY(1, 0)
	notes := v.notes.SetContent(render.Wrap(notesText, inner.Width)).SetHeight(inner.Height)

	greeting := ui.Panel{
		Title:   "Greeting",
		Focused: rc.Focused(regionGreeting),
		Footer:  []ui.Hint{{Key: "g", Label: "Greet"}},
		Content: rc.Styles().Base.Render("Hello."),
	}

	return render.Compose(area,
		render.At(rows[0], greeting.Render(rc.For(rows[0]))),
		render.At(rows[1], ui.Panel{
			Title:   "Notes",
			Focused: rc.Focused(regionNotes),
			Scroll:  notes.ScrollPos(),
			Content: notes.Render(rc.For(inner)),
		}.Render(rc.For(rows[1]))),
	)
}

const notesText = "This is your screen. It lives in internal/modules/hello.\n" +
	"\n" +
	"Adding another is a directory beside this one and a line in\n" +
	"internal/app/register.go. Nothing else in the shell has to change.\n" +
	"\n" +
	"tab cycles focus, v badges every region and jumps to the one you press,\n" +
	"alt+hjkl moves geometrically. ctrl+k opens the command palette, ? lists\n" +
	"every command with the key that reaches it, and ctrl+t cycles the theme.\n" +
	"All of it is generated from what modules registered, so a verb you add\n" +
	"shows up in the palette, the help screen and the status bar at once.\n" +
	"\n" +
	"Narrow the terminal: the sidebar collapses at 120 columns and disappears\n" +
	"below 80. Under 40x10 the frame is replaced by one legible message."
`

const appTestGo = `package app_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/app"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/config"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/runtime"
)

// build assembles the real application. It is the strongest cheap test there
// is: Build reports every duplicate route, malformed binding and keymap
// conflict, so a module that registers something contradictory fails here
// rather than at the first keypress.
func build(t *testing.T, themeName string) runtime.Model {
	t.Helper()

	cfg := config.Defaults()
	cfg.Theme = themeName
	cfg.Mouse = false

	m, err := app.Build(app.Options{Name: "app", Version: "0.1.0", Config: cfg})
	if err != nil {
		t.Fatalf("Build(): %v", err)
	}
	return m
}

// resize drives the model the way the terminal would.
func resize(t *testing.T, m runtime.Model, w, h int) runtime.Model {
	t.Helper()

	model, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return model.(runtime.Model)
}

func TestBuildSucceeds(t *testing.T) {
	build(t, "dark")
}

func TestFrameRendersAtEverySize(t *testing.T) {
	for _, size := range []struct{ w, h int }{{80, 24}, {120, 40}} {
		m := resize(t, build(t, "ascii"), size.w, size.h)
		if frame := m.Render(); !strings.Contains(frame, "Hello") {
			t.Errorf("%dx%d frame does not show the hello screen:\n%s", size.w, size.h, frame)
		}
	}
}
`
