package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type StatusBar struct {
	Width int
	Left  []key.Binding
	Right []key.Binding
	Busy  string
}

func RenderStatusBar(t theme.Theme, s StatusBar) string {
	surface := lipgloss.NewStyle().Background(t.Base.Surface)
	keyStyle := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Surface)
	labelStyle := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Surface)
	busyStyle := lipgloss.NewStyle().Foreground(t.Status.Info.Fg).Background(t.Base.Surface)
	sep := surface.Render(Repeat(" ", t.Space.Gutter*2))

	budget := s.Width - 2*t.Space.StatusPadX
	busy := ""
	busyStyled := ""
	if s.Busy != "" {
		busy = Repeat(" ", t.Space.Gutter*2) + s.Busy
		busyStyled = busyStyle.Render(busy)
	}
	left := fitBindings(s.Left, budget-Width(busy), t)
	right := fitBindings(s.Right, budget-bindingsWidth(left, t)-Width(busy)-t.Space.Gutter*2, t)
	gap := budget - bindingsWidth(left, t) - Width(busy) - bindingsWidth(right, t)
	if gap < 0 {
		gap = 0
	}
	pad := surface.Render(Repeat(" ", t.Space.StatusPadX))
	return pad +
		renderBindings(left, keyStyle, labelStyle, sep) +
		busyStyled +
		surface.Render(Repeat(" ", gap)) +
		renderBindings(right, keyStyle, labelStyle, sep) +
		pad
}

func fitBindings(bindings []key.Binding, budget int, t theme.Theme) []key.Binding {
	enabled := make([]key.Binding, 0, len(bindings))
	for _, b := range bindings {
		if b.Enabled() {
			enabled = append(enabled, b)
		}
	}
	for len(enabled) > 0 && bindingsWidth(enabled, t) > budget {
		enabled = enabled[:len(enabled)-1]
	}
	return enabled
}

func bindingsWidth(bindings []key.Binding, t theme.Theme) int {
	total := 0
	for i, b := range bindings {
		if i > 0 {
			total += t.Space.Gutter * 2
		}
		total += Width(b.Help().Key) + 1 + Width(b.Help().Desc)
	}
	return total
}

func renderBindings(bindings []key.Binding, keyStyle, labelStyle lipgloss.Style, sep string) string {
	out := ""
	for i, b := range bindings {
		if i > 0 {
			out += sep
		}
		out += keyStyle.Render(b.Help().Key) + labelStyle.Render(" "+b.Help().Desc)
	}
	return out
}
