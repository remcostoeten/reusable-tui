package app_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/exp/golden"

	"github.com/remcostoeten/reusable-tui/internal/app"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/config"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/render"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/runtime"
)

func build(t *testing.T, themeName string) runtime.Model {
	t.Helper()

	cfg := config.Defaults()
	cfg.Theme = themeName
	cfg.Mouse = false

	m, err := app.Build(app.Options{Name: "demo", Version: "0.1.0", Config: cfg})
	if err != nil {
		t.Fatalf("Build(): %v", err)
	}
	return m
}

// resize drives the model the way the terminal would, then applies a key
// sequence, returning the model that results.
func resize(t *testing.T, m runtime.Model, w, h int, keys ...string) runtime.Model {
	t.Helper()

	model, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	for _, k := range keys {
		next, cmd := model.Update(keyPress(k))
		model = drain(next, cmd)
	}
	return model.(runtime.Model)
}

// drain runs the commands a message produced and feeds the results back, the
// way the event loop does, so a test sees the state the user would.
func drain(model tea.Model, cmd tea.Cmd) tea.Model {
	for range 8 {
		if cmd == nil {
			return model
		}
		msg := cmd()
		if msg == nil {
			return model
		}
		if batch, ok := msg.(tea.BatchMsg); ok {
			for _, c := range batch {
				model = drain(model, c)
			}
			return model
		}
		model, cmd = model.Update(msg)
	}
	return model
}

func TestBuildSucceedsWithNoKeymapConflicts(t *testing.T) {
	build(t, "dark")
}

func TestBuildRejectsAnUnknownTheme(t *testing.T) {
	cfg := config.Defaults()
	cfg.Theme = "does-not-exist"

	m, err := app.Build(app.Options{Name: "demo", Config: cfg})
	if err != nil {
		t.Fatalf("an unknown theme must fall back, not fail: %v", err)
	}
	if got := m.Theme().Name; got != "dark" {
		t.Errorf("theme = %q, want the default \"dark\"", got)
	}
}

func TestFrameFillsTheTerminal(t *testing.T) {
	sizes := []struct{ w, h int }{
		{60, 20}, {80, 24}, {120, 40}, {200, 60},
	}

	for _, size := range sizes {
		m := resize(t, build(t, "dark"), size.w, size.h)
		lines := render.Lines(m.Render())

		if len(lines) != size.h {
			t.Fatalf("%dx%d: frame is %d rows, want %d", size.w, size.h, len(lines), size.h)
		}
		for i, l := range lines {
			if w := render.Width(l); w > size.w {
				t.Fatalf("%dx%d: row %d is %d cells, want at most %d", size.w, size.h, i, w, size.w)
			}
		}
	}
}

func TestTerminalTooSmall(t *testing.T) {
	m := resize(t, build(t, "ascii"), 30, 8)
	if !strings.Contains(m.Render(), "too small") {
		t.Errorf("a terminal below the minimum must say so, got:\n%s", m.Render())
	}
}

func TestNavigationChangesTheRoute(t *testing.T) {
	tests := []struct {
		name string
		keys []string
		want string
	}{
		{name: "starts on home", want: "home"},
		{name: "next section", keys: []string{"]"}, want: "manager"},
		{name: "twice", keys: []string{"]", "]"}, want: "settings"},
		{name: "wraps", keys: []string{"]", "]", "]"}, want: "home"},
		{name: "previous wraps backwards", keys: []string{"["}, want: "settings"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := resize(t, build(t, "dark"), 120, 40, tt.keys...)
			if got := string(m.Route().Route); got != tt.want {
				t.Errorf("route = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTabCyclesFocusWithinTheRoute(t *testing.T) {
	m := resize(t, build(t, "dark"), 120, 40)
	first := m.Focus().Current()

	m = resize(t, build(t, "dark"), 120, 40, "tab")
	second := m.Focus().Current()

	if first == second {
		t.Fatalf("tab did not move focus; it stayed on %q", first)
	}
	if first.IsZero() || second.IsZero() {
		t.Fatalf("focus landed nowhere: %q then %q", first, second)
	}
}

func TestPaletteOpensAndCloses(t *testing.T) {
	m := resize(t, build(t, "dark"), 120, 40, "ctrl+k")
	if got := m.Overlay(); got != "palette" {
		t.Fatalf("overlay = %q, want \"palette\"", got)
	}
	if !strings.Contains(m.Render(), "Commands") {
		t.Error("the palette did not draw over the frame")
	}
	if got := m.Focus().Depth(); got != 2 {
		t.Errorf("focus depth = %d, want 2 — the overlay must trap focus", got)
	}

	m = resize(t, build(t, "dark"), 120, 40, "ctrl+k", "esc")
	if got := m.Overlay(); got != "" {
		t.Errorf("overlay = %q, want none after esc", got)
	}
	if got := m.Focus().Depth(); got != 1 {
		t.Errorf("focus depth = %d, want 1 restored", got)
	}
}

func TestHelpIsGeneratedFromTheRegistry(t *testing.T) {
	m := resize(t, build(t, "ascii"), 120, 40, "?")
	frame := m.Render()

	for _, want := range []string{"Command Palette", "Quit", "APPLICATION", "ctrl+k"} {
		if !strings.Contains(frame, want) {
			t.Errorf("the help overlay is missing %q", want)
		}
	}
}

func TestThemeCycling(t *testing.T) {
	tests := []struct {
		presses int
		want    string
	}{
		{presses: 0, want: "dark"},
		{presses: 1, want: "dim"},
		{presses: 2, want: "high-contrast"},
		{presses: 3, want: "ascii"},
		{presses: 4, want: "dark"},
	}

	for _, tt := range tests {
		keys := make([]string, tt.presses)
		for i := range keys {
			keys[i] = "ctrl+t"
		}
		m := resize(t, build(t, "dark"), 120, 40, keys...)
		if got := m.Theme().Name; got != tt.want {
			t.Errorf("after %d presses theme = %q, want %q", tt.presses, got, tt.want)
		}
	}
}

func TestDrillInAndBack(t *testing.T) {
	m := resize(t, build(t, "dark"), 120, 40, "]", "enter")
	if got := string(m.Route().Route); got != "manager.detail" {
		t.Fatalf("route = %q, want \"manager.detail\"", got)
	}

	m = resize(t, build(t, "dark"), 120, 40, "]", "enter", "esc")
	if got := string(m.Route().Route); got != "manager" {
		t.Errorf("route = %q, want \"manager\" after esc", got)
	}
}

func TestEscAtTheRootDoesNotNavigate(t *testing.T) {
	m := resize(t, build(t, "dark"), 120, 40, "esc")
	if got := string(m.Route().Route); got != "home" {
		t.Errorf("route = %q, want \"home\" — esc has nowhere to go at the root", got)
	}
}

func TestDangerousCommandAsksFirst(t *testing.T) {
	m := resize(t, build(t, "dark"), 120, 40, "]", "d")
	if got := m.Overlay(); got != "confirm" {
		t.Fatalf("overlay = %q, want \"confirm\" — a dangerous command must ask", got)
	}
	if !strings.Contains(m.Render(), "cannot be undone") {
		t.Error("the confirmation did not draw")
	}
}

func TestJumpModeBadgesEveryRegion(t *testing.T) {
	m := resize(t, build(t, "ascii"), 120, 40, "v")
	if got := m.Overlay(); got != "jump" {
		t.Fatalf("overlay = %q, want \"jump\"", got)
	}

	frame := m.Render()
	if !strings.Contains(frame, "Sources") {
		t.Error("jump mode painted over the frame instead of badging it")
	}
}

func TestEveryCommandIsReachable(t *testing.T) {
	m := resize(t, build(t, "dark"), 120, 40)
	ctx := m.Context()

	commands := ctx.Commands.All()
	if len(commands) == 0 {
		t.Fatal("no commands were registered")
	}

	// The palette is always bound, and it lists everything available, so a
	// command is reachable as long as it carries a title to find it by.
	for _, c := range commands {
		if c.Title == "" {
			t.Errorf("command %q has no title, so the palette cannot offer it", c.ID)
		}
		if c.Run == nil {
			t.Errorf("command %q has no implementation", c.ID)
		}
	}
}

func TestFrameGolden(t *testing.T) {
	var b strings.Builder
	for _, size := range []struct{ w, h int }{{80, 24}, {120, 40}} {
		for _, name := range []string{"dark", "ascii"} {
			m := resize(t, build(t, name), size.w, size.h)
			b.WriteString("--- " + name + " ---\n")
			b.WriteString(m.Render())
			b.WriteString("\n")
		}
	}
	golden.RequireEqual(t, []byte(b.String()))
}
