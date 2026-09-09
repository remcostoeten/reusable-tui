package app

import (
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/theme"
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
		ID:    "theme.reload",
		Label: "Reload themes from disk",
		Group: "theme",
		Run:   reloadThemes(m.opts.ThemeDir, m.opts.Config.Overrides),
	})
	m.commands.Add(ui.Command{
		ID:    "theme.export",
		Label: "Export active theme to a file",
		Group: "theme",
		Run:   exportActive(),
	})
	m.commands.Add(ui.Command{
		ID:    "app.quit",
		Label: "Quit",
		Group: "app",
		Run:   requestQuit,
	})
}

type screenMsg struct {
	ID string
}

type quitMsg struct{}

func requestQuit() tea.Msg {
	return quitMsg{}
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

type themesReloadedMsg struct {
	registry *theme.Registry
	warnings []error
}

func reloadThemes(dir string, overrides map[string]map[string]string) tea.Cmd {
	return func() tea.Msg { return reloadResult(dir, overrides) }
}

func reloadResult(dir string, overrides map[string]map[string]string) tea.Msg {
	registry, warnings := theme.Compose(dir, overrides)
	return themesReloadedMsg{registry: registry, warnings: warnings}
}

type exportThemeMsg struct{}

func exportActive() tea.Cmd {
	return func() tea.Msg { return exportThemeMsg{} }
}

func exportTheme(dir string, active theme.Theme) tea.Cmd {
	return func() tea.Msg { return exportResult(dir, active) }
}

func exportResult(dir string, active theme.Theme) tea.Msg {
	path := filepath.Join(dir, active.Name+".json")
	if err := theme.WriteFile(path, theme.Export(active)); err != nil {
		return ui.ErrorMsg{Err: err}
	}
	return ui.SuccessMsg{Text: "wrote " + path}
}
