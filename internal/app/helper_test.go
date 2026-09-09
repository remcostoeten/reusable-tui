package app

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/muesli/termenv"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/store"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const (
	testWidth  = 100
	testHeight = 28
)

func TestMain(m *testing.M) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	m.Run()
}

func newTestStore(t *testing.T) *store.Store {
	t.Helper()
	db, err := store.Open(t.TempDir() + "/state.db")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(closer(db))
	return db
}

func closer(db *store.Store) func() {
	return func() { db.Close() }
}

func seed(t *testing.T, db *store.Store) {
	t.Helper()
	rows := []struct {
		title string
		done  int
	}{
		{"wire up the root model", 1},
		{"render the status bar", 0},
		{"ship the palette", 0},
	}
	for _, row := range rows {
		_, err := db.DB().Exec(
			`INSERT INTO example_items (title, done, created_at) VALUES (?, ?, '2024-01-01 00:00:00')`,
			row.title, row.done,
		)
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
	}
}

func newTestModel(t *testing.T, db *store.Store) *Model {
	t.Helper()
	return New(Options{
		Clock:    fixedClock,
		Store:    db,
		Config:   config.Defaults(),
		Notifier: notify.Discard{},
		Themes:   theme.Builtin(),
		Fidelity: theme.FidelityTrueColor,
	})
}

func fixedClock() time.Time {
	return time.Date(2024, time.January, 17, 9, 0, 0, 0, time.UTC)
}

func settled(t *testing.T, m *Model) *Model {
	t.Helper()
	m.Update(tea.WindowSizeMsg{Width: testWidth, Height: testHeight})
	for _, screen := range m.screens {
		drain(m, screen.Init())
	}
	return m
}

func drain(m *Model, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, part := range batch {
			drain(m, part)
		}
		return
	}
	switch msg.(type) {
	case ui.TickMsg, ui.ToastExpiredMsg, nil:
		return
	}
	_, next := m.Update(msg)
	drain(m, next)
}

func runProgram(t *testing.T, m *Model, keys ...tea.Msg) *Model {
	t.Helper()
	return runProgramUntil(t, m, "", keys...)
}

func runProgramUntil(t *testing.T, m *Model, want string, keys ...tea.Msg) *Model {
	t.Helper()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(testWidth, testHeight))
	for _, msg := range keys {
		tm.Send(msg)
	}
	if want != "" {
		teatest.WaitFor(t, tm.Output(), contains(want), teatest.WithDuration(5*time.Second))
	}
	if err := tm.Quit(); err != nil {
		t.Fatalf("quit: %v", err)
	}
	final, ok := tm.FinalModel(t, teatest.WithFinalTimeout(5*time.Second)).(*Model)
	if !ok {
		t.Fatal("final model type mismatch")
	}
	return final
}

func contains(want string) func([]byte) bool {
	return func(out []byte) bool { return bytes.Contains(out, []byte(want)) }
}

func runes(value string) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func special(kind tea.KeyType) tea.Msg {
	return tea.KeyMsg{Type: kind}
}
