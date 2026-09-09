package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type Stat struct {
	Label string
	Value string
	Token *theme.StatusToken
}

func RenderStats(t theme.Theme, width int, stats ...Stat) string {
	if len(stats) == 0 || width <= 0 {
		return ""
	}
	widths := SplitWidths(width, 0, segmentWeights(len(stats))...)
	labels := ""
	values := ""
	for i, stat := range stats {
		labels += statLabel(t, stat, widths[i], i == len(stats)-1)
		values += statValue(t, stat, widths[i], i == len(stats)-1)
	}
	return Join([]string{labels, values})
}

func statLabel(t theme.Theme, s Stat, width int, right bool) string {
	style := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	return style.Render(align(s.Label, width, right))
}

func statValue(t theme.Theme, s Stat, width int, right bool) string {
	color := t.Text.Primary
	if s.Token != nil {
		color = s.Token.Fg
	}
	style := lipgloss.NewStyle().Foreground(color).Background(t.Base.Background).Bold(true)
	return style.Render(align(s.Value, width, right))
}

func align(s string, width int, right bool) string {
	if right {
		return PadLeft(Truncate(s, width), width)
	}
	return Fit(s, width)
}
