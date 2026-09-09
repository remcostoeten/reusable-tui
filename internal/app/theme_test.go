package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func TestUserThemeIsSelectableFromThePalette(t *testing.T) {
	dir := t.TempDir()
	body := `{"name":"cyan","extends":"violet-dark","tokens":{"accent.active":"#22d3ee","border.focused":"#22d3ee"}}`
	if err := os.WriteFile(filepath.Join(dir, "cyan.json"), []byte(body), 0o600); err != nil {
		t.Fatalf("write theme: %v", err)
	}
	registry, failures := theme.Compose(dir, nil)
	if len(failures) != 0 {
		t.Fatalf("compose: %v", failures)
	}

	db := newTestStore(t)
	m := settled(t, New(Options{
		Store:    db,
		Config:   config.Defaults(),
		Notifier: notify.Discard{},
		Themes:   registry,
		ThemeDir: dir,
		Fidelity: theme.FidelityTrueColor,
	}))

	if !hasCommand(m, "theme.cyan") {
		t.Fatal("user theme did not register a palette command")
	}
	m.Update(ui.ThemeMsg{Name: "cyan"})
	if got := theme.FormatColor(m.Theme().Accent.Active); got != "#22D3EE" {
		t.Fatalf("accent is %q", got)
	}
	if !strings.Contains(m.View(), "34;211;238") {
		t.Fatal("custom accent is not present in the rendered frame")
	}
}

func TestConfigOverridesRecolorABuiltinTheme(t *testing.T) {
	registry, failures := theme.Compose(t.TempDir(), map[string]map[string]string{
		theme.NameVioletDark: {theme.TokenAccentActive: "#FF0080"},
	})
	if len(failures) != 0 {
		t.Fatalf("compose: %v", failures)
	}
	resolved := registry.Resolve(theme.NameVioletDark, theme.FidelityTrueColor)
	if got := theme.FormatColor(resolved.Accent.Active); got != "#FF0080" {
		t.Fatalf("override did not survive resolve, accent is %q", got)
	}
}

func TestStartupWarningsSurfaceAsAToast(t *testing.T) {
	db := newTestStore(t)
	m := New(Options{
		Store:    db,
		Config:   config.Defaults(),
		Notifier: notify.Discard{},
		Themes:   theme.Builtin(),
		Fidelity: theme.FidelityTrueColor,
		Warnings: []error{errors.New("cyan.json: bad hex"), errors.New("second problem")},
	})
	m.Update(tea.WindowSizeMsg{Width: testWidth, Height: testHeight})
	runCmd(t, m, m.Init())
	if !m.toast.Visible {
		t.Fatal("startup warning did not raise a toast")
	}
	if m.toast.Kind != ui.ToastWarning {
		t.Fatalf("toast kind is %v", m.toast.Kind)
	}
	if !strings.Contains(m.toast.Text, "cyan.json") || !strings.Contains(m.toast.Text, "+1 more") {
		t.Fatalf("toast text is %q", m.toast.Text)
	}
}

func TestExportWritesTheActiveTheme(t *testing.T) {
	dir := t.TempDir()
	db := newTestStore(t)
	m := settled(t, New(Options{
		Store:    db,
		Config:   config.Defaults(),
		Notifier: notify.Discard{},
		Themes:   theme.Builtin(),
		ThemeDir: dir,
		Fidelity: theme.FidelityTrueColor,
	}))
	m.Update(ui.ThemeMsg{Name: theme.NameMonochrome})
	runCmd(t, m, exportTheme(m.opts.ThemeDir, m.theme))

	written, err := theme.ReadFile(filepath.Join(dir, theme.NameMonochrome+".json"))
	if err != nil {
		t.Fatalf("read exported theme: %v", err)
	}
	if written.Name != theme.NameMonochrome {
		t.Fatalf("exported the wrong theme: %q", written.Name)
	}
	if len(written.Tokens) != len(theme.Tokens()) {
		t.Fatalf("exported %d tokens, expected %d", len(written.Tokens), len(theme.Tokens()))
	}
}

func TestReloadPicksUpANewThemeFile(t *testing.T) {
	dir := t.TempDir()
	db := newTestStore(t)
	registry, _ := theme.Compose(dir, nil)
	m := settled(t, New(Options{
		Store:    db,
		Config:   config.Defaults(),
		Notifier: notify.Discard{},
		Themes:   registry,
		ThemeDir: dir,
		Fidelity: theme.FidelityTrueColor,
	}))
	if hasCommand(m, "theme.amber") {
		t.Fatal("theme exists before it was written")
	}

	body := `{"name":"amber","extends":"violet-dark","tokens":{"accent.active":"#F59E0B"}}`
	if err := os.WriteFile(filepath.Join(dir, "amber.json"), []byte(body), 0o600); err != nil {
		t.Fatalf("write theme: %v", err)
	}
	runCmd(t, m, reloadThemes(dir, nil))

	if !hasCommand(m, "theme.amber") {
		t.Fatal("reload did not register the new theme")
	}
	m.Update(ui.ThemeMsg{Name: "amber"})
	if got := theme.FormatColor(m.Theme().Accent.Active); got != "#F59E0B" {
		t.Fatalf("accent is %q", got)
	}
}

func hasCommand(m *Model, id string) bool {
	_, ok := m.commands.Get(id)
	return ok
}

func runCmd(t *testing.T, m *Model, cmd tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("no command to run")
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, inner := range batch {
			runCmd(t, m, inner)
		}
		return
	}
	if msg == nil {
		return
	}
	m.Update(msg)
}
