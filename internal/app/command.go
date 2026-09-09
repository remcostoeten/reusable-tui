package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func (m *Model) registerCommands() {
	for _, name := range m.opts.Themes.Names() {
		m.commands.Add(ui.Command{
			ID:    "theme." + name,
			Label: "Theme: " + name,
			Group: "theme",
			Run:   ui.SetTheme(name),
		})
	}
	for _, screen := range m.screens {
		m.commands.Add(ui.Command{
			ID:    "screen." + screen.ID(),
			Label: "Go to " + screen.Title(),
			Group: "navigate",
			Run:   goToScreen(screen.ID()),
		})
	}
	m.commands.Add(ui.Command{
		ID:    "app.quit",
		Label: "Quit",
		Group: "app",
		Run:   tea.Quit,
	})
}

type screenMsg struct {
	ID string
}

func goToScreen(id string) tea.Cmd {
	return func() tea.Msg { return screenMsg{ID: id} }
}

func notifyCmd(notifier notify.Notifier, title, body string) tea.Cmd {
	return func() tea.Msg { return notifyResult(notifier, title, body) }
}

func notifyResult(notifier notify.Notifier, title, body string) tea.Msg {
	if err := notifier.Notify(title, body); err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return nil
}

func persistConfig(path string, cfg config.Config) tea.Cmd {
	return func() tea.Msg { return persistResult(path, cfg) }
}

func persistResult(path string, cfg config.Config) tea.Msg {
	if path == "" {
		return nil
	}
	if err := config.Save(path, cfg); err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return nil
}
