package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type Tab struct {
	ID    string
	Label string
}

type Header struct {
	App     string
	Version string
	Tabs    []Tab
	Active  int
	Width   int
}

type tabSpan struct {
	start int
	width int
}

func RenderHeader(t theme.Theme, h Header) string {
	appStyle := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background).Bold(true)
	versionStyle := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	activeStyle := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Background)
	inactiveStyle := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	gap := lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", t.Space.Gutter*2))

	lead := Repeat(" ", t.Space.HeaderPadX) + h.App + " " + h.Version
	row := lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", t.Space.HeaderPadX)) +
		appStyle.Render(h.App) + versionStyle.Render(" "+h.Version)

	cursor := Width(lead)
	spans := make([]tabSpan, 0, len(h.Tabs))
	for i, tab := range h.Tabs {
		row += gap
		cursor += Width(gap)
		label := tabLabel(t, tab.Label, i == h.Active)
		spans = append(spans, tabSpan{start: cursor, width: Width(label)})
		if i == h.Active {
			row += activeStyle.Render(label)
		} else {
			row += inactiveStyle.Render(label)
		}
		cursor += Width(label)
	}
	row += lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", h.Width-cursor))
	return Join([]string{row, headerRule(t, h, spans)})
}

func tabLabel(t theme.Theme, label string, active bool) string {
	if active {
		return t.Marker.FocusOpen + label + t.Marker.FocusClose
	}
	if t.Marker.FocusOpen == "" {
		return label
	}
	return Repeat(" ", Width(t.Marker.FocusOpen)) + label + Repeat(" ", Width(t.Marker.FocusClose))
}

func headerRule(t theme.Theme, h Header, spans []tabSpan) string {
	base := lipgloss.NewStyle().Foreground(t.Border.Subtle).Background(t.Base.Background)
	accent := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)
	if h.Active < 0 || h.Active >= len(spans) {
		return base.Render(Repeat(t.Marker.Rule, h.Width))
	}
	span := spans[h.Active]
	tail := h.Width - span.start - span.width
	if tail < 0 {
		tail = 0
	}
	return base.Render(Repeat(t.Marker.Rule, span.start)) +
		accent.Render(Repeat(t.Marker.Rule, span.width)) +
		base.Render(Repeat(t.Marker.Rule, tail))
}
