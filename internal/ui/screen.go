package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type RenderContext struct {
	Theme   theme.Theme
	Keys    *keymap.Registry
	Width   int
	Height  int
	Tick    int
	Focused string
	Jump    Jump
}

type Screen interface {
	ID() string
	Title() string
	Panels() []string
	Init() tea.Cmd
	Focus(panelID string)
	Capturing() bool
	Update(msg tea.Msg) tea.Cmd
	View(ctx RenderContext) string
}
