package app

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/features/dashboard"
	"github.com/remcostoeten/reusable-tui/internal/features/example"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func click(x, y int) tea.MouseMsg {
	return tea.MouseMsg{X: x, Y: y, Button: tea.MouseButtonLeft, Action: tea.MouseActionPress}
}

func wheel(x, y int, down bool) tea.MouseMsg {
	button := tea.MouseButtonWheelUp
	if down {
		button = tea.MouseButtonWheelDown
	}
	return tea.MouseMsg{X: x, Y: y, Button: button, Action: tea.MouseActionPress}
}

func rendered(t *testing.T, m *Model) *Model {
	t.Helper()
	m.View()
	return m
}

func TestMouseClickFocusesPanel(t *testing.T) {
	db := newTestStore(t)
	m := rendered(t, settled(t, newTestModel(t, db)))
	m.Update(click(testWidth-5, 8))
	if got := m.focus["example"]; got != keymap.PanelExampleDetail {
		t.Fatalf("click did not focus the detail panel, focused is %q", got)
	}
}

func TestMouseClickSelectsRow(t *testing.T) {
	db := newTestStore(t)
	seed(t, db)
	m := rendered(t, settled(t, newTestModel(t, db)))
	header := m.theme.Space.HeaderRows
	m.Update(click(4, header+1+m.theme.Space.PanelPadY+2))
	if got := selectedTitle(t, m); got != "ship the palette" {
		t.Fatalf("third row not selected after click, selected is %q", got)
	}
}

func selectedTitle(t *testing.T, m *Model) string {
	t.Helper()
	screen, ok := m.ActiveScreen().(*example.Model)
	if !ok {
		t.Fatal("active screen is not the example screen")
	}
	item, ok := screen.Selected()
	if !ok {
		t.Fatal("nothing selected")
	}
	return item.Title
}

func TestMouseWheelMovesCursor(t *testing.T) {
	db := newTestStore(t)
	seed(t, db)
	m := rendered(t, settled(t, newTestModel(t, db)))
	m.Update(wheel(4, 8, true))
	if got := selectedTitle(t, m); got != "render the status bar" {
		t.Fatalf("wheel down did not move the cursor, selected is %q", got)
	}
}

func TestMouseClickSwitchesTab(t *testing.T) {
	db := newTestStore(t)
	m := rendered(t, settled(t, newTestModel(t, db)))
	label := "Dashboard"
	x := strings.Index(ui.Lines(m.View())[0], label)
	if x < 0 {
		t.Fatal("dashboard tab not in header")
	}
	m.Update(click(ui.Width(ui.Lines(m.View())[0][:x])+1, 0))
	if got := m.ActiveScreen().ID(); got != dashboard.ScreenID {
		t.Fatalf("header click did not switch tab, active is %q", got)
	}
}

func TestMouseIgnoredWhilePaletteOpen(t *testing.T) {
	db := newTestStore(t)
	m := rendered(t, settled(t, newTestModel(t, db)))
	m.Update(special(tea.KeyCtrlK))
	m.Update(click(testWidth-5, 8))
	if got := m.focus["example"]; got != keymap.PanelExampleList {
		t.Fatalf("click changed focus while palette open, focused is %q", got)
	}
}

func TestBusySpinnerInStatusBar(t *testing.T) {
	db := newTestStore(t)
	m := settled(t, newTestModel(t, db))
	m.Update(ui.BusyMsg{ID: "test.sync", Label: "syncing"})
	if !strings.Contains(m.View(), "syncing") {
		t.Fatal("busy label not rendered in status bar")
	}
	m.Update(ui.IdleMsg{ID: "test.sync"})
	if strings.Contains(m.View(), "syncing") {
		t.Fatal("busy label still rendered after idle")
	}
	rows := ui.Lines(m.View())
	if width := ui.Width(rows[len(rows)-1]); width != testWidth {
		t.Fatalf("status bar width %d after busy toggle", width)
	}
}
