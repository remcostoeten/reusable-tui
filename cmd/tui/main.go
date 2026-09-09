package main

import (
	"errors"
	"flag"
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
	f, err := parseFlags(os.Args[1:], os.Stderr)
	switch {
	case errors.Is(err, flag.ErrHelp):
		return
	case errors.Is(err, errShowVersion):
		fmt.Println(app.Name, app.Version)
		return
	case err != nil:
		os.Exit(2)
	}
	if err := run(f); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(f flags) error {
	configPath, err := resolveConfigPath(f)
	if err != nil {
		return err
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	if f.themeName != "" {
		cfg.Theme = f.themeName
	}
	dbPath, err := resolveDBPath(f)
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

	program := tea.NewProgram(model, programOptions(f)...)
	_, err = program.Run()
	return err
}

func programOptions(f flags) []tea.ProgramOption {
	opts := []tea.ProgramOption{tea.WithAltScreen()}
	if !f.noMouse {
		opts = append(opts, tea.WithMouseCellMotion())
	}
	return opts
}

func resolveConfigPath(f flags) (string, error) {
	if f.configPath != "" {
		return f.configPath, nil
	}
	return config.Path(app.Name)
}

func resolveDBPath(f flags) (string, error) {
	if f.dbPath != "" {
		return f.dbPath, nil
	}
	return store.DefaultPath(app.Name)
}
