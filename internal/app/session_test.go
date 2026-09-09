package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/features/dashboard"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/store"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

func newSessionModel(t *testing.T, db *store.Store, cfg config.Config, path string) *Model {
	t.Helper()
	return New(Options{
		Clock:      fixedClock,
		Store:      db,
		Config:     cfg,
		ConfigPath: path,
		Notifier:   notify.Discard{},
		Themes:     theme.Builtin(),
		Fidelity:   theme.FidelityTrueColor,
	})
}

func TestSessionRestore(t *testing.T) {
	db := newTestStore(t)
	cfg := config.Defaults()
	cfg.Session = config.Session{
		Screen: dashboard.ScreenID,
		Panels: map[string]string{
			dashboard.ScreenID: keymap.PanelDashboardPeriod,
			"example":          "example.bogus",
			"ghost":            "nowhere",
		},
	}
	m := newSessionModel(t, db, cfg, "")
	if got := m.ActiveScreen().ID(); got != dashboard.ScreenID {
		t.Fatalf("session screen not restored, active is %q", got)
	}
	if got := m.focus[dashboard.ScreenID]; got != keymap.PanelDashboardPeriod {
		t.Fatalf("session panel not restored, focused is %q", got)
	}
	if got := m.focus["example"]; got != keymap.PanelExampleList {
		t.Fatalf("invalid panel id was accepted, focused is %q", got)
	}
}

func TestSessionUnknownScreenIgnored(t *testing.T) {
	db := newTestStore(t)
	cfg := config.Defaults()
	cfg.Session.Screen = "ghost"
	m := newSessionModel(t, db, cfg, "")
	if got := m.ActiveScreen().ID(); got != "example" {
		t.Fatalf("unknown screen changed the active tab to %q", got)
	}
}

func TestQuitPersistsSession(t *testing.T) {
	db := newTestStore(t)
	path := t.TempDir() + "/config.json"
	m := newSessionModel(t, db, config.Defaults(), path)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(testWidth, testHeight))
	tm.Send(special(tea.KeyTab))
	tm.Send(special(tea.KeyCtrlN))
	tm.Send(runes("q"))
	tm.WaitFinished(t)
	saved, err := config.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if saved.Session.Screen != dashboard.ScreenID {
		t.Fatalf("session screen not persisted, got %q", saved.Session.Screen)
	}
	if got := saved.Session.Panels[dashboard.ScreenID]; got != keymap.PanelDashboardInsights {
		t.Fatalf("session panel not persisted, got %q", got)
	}
}
