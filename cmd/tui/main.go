package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/app"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/store"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	configPath, err := config.Path(app.Name)
	if err != nil {
		return err
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	dbPath, err := store.DefaultPath(app.Name)
	if err != nil {
		return err
	}
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	themeDir := theme.UserDir(app.Name)
	themes, warnings := theme.Compose(themeDir, cfg.Overrides)

	model := app.New(app.Options{
		Store:      db,
		Config:     cfg,
		ConfigPath: configPath,
		Notifier:   notify.NewDesktop(app.Name),
		Themes:     themes,
		ThemeDir:   themeDir,
		Fidelity:   theme.FidelityFor(lipgloss.ColorProfile()),
		Warnings:   warnings,
	})

	program := tea.NewProgram(model, tea.WithAltScreen())
	_, err = program.Run()
	return err
}
