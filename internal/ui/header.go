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

	row := lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", t.Space.HeaderPadX)) +
		appStyle.Render(h.App) + versionStyle.Render(" "+h.Version)

	spans := tabSpans(t, h)
	cursor := headerLeadWidth(t, h)
	for i, tab := range h.Tabs {
		row += gap
		label := tabLabel(t, tab.Label, i == h.Active)
		if i == h.Active {
			row += activeStyle.Render(label)
		} else {
			row += inactiveStyle.Render(label)
		}
		cursor = spans[i].start + spans[i].width
	}
	row += lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", h.Width-cursor))
	return Join([]string{row, headerRule(t, h, spans)})
}

func headerLeadWidth(t theme.Theme, h Header) int {
	return Width(Repeat(" ", t.Space.HeaderPadX) + h.App + " " + h.Version)
}

func tabSpans(t theme.Theme, h Header) []tabSpan {
	cursor := headerLeadWidth(t, h)
	gap := t.Space.Gutter * 2
	spans := make([]tabSpan, 0, len(h.Tabs))
	for i, tab := range h.Tabs {
		cursor += gap
		width := Width(tabLabel(t, tab.Label, i == h.Active))
		spans = append(spans, tabSpan{start: cursor, width: width})
		cursor += width
	}
	return spans
}

func TabAt(t theme.Theme, h Header, x int) (string, bool) {
	for i, span := range tabSpans(t, h) {
		if x >= span.start && x < span.start+span.width {
			return h.Tabs[i].ID, true
		}
	}
	return "", false
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
