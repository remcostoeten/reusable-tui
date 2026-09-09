package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func TestGoldenExampleScreen(t *testing.T) {
	for _, name := range theme.Builtin().Names() {
		t.Run(name, goldenCase(name))
	}
}

func goldenCase(name string) func(*testing.T) {
	return func(t *testing.T) { renderGolden(t, name) }
}

func renderGolden(t *testing.T, name string) {
	t.Helper()
	db := newTestStore(t)
	seed(t, db)
	m := settled(t, newTestModel(t, db))
	m.Update(ui.ThemeMsg{Name: name})
	golden.RequireEqual(t, []byte(m.View()))
}

func TestGoldenDashboardScreen(t *testing.T) {
	for _, name := range theme.Builtin().Names() {
		t.Run(name, dashboardGoldenCase(name))
	}
}

func dashboardGoldenCase(name string) func(*testing.T) {
	return func(t *testing.T) { renderDashboardGolden(t, name) }
}

func renderDashboardGolden(t *testing.T, name string) {
	t.Helper()
	db := newTestStore(t)
	m := settled(t, newTestModel(t, db))
	m.Update(ui.ThemeMsg{Name: name})
	m.Update(tea.KeyMsg{Type: tea.KeyTab})
	golden.RequireEqual(t, []byte(m.View()))
}
