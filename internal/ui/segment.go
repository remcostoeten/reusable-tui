package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type Segment struct {
	Options []string
	Active  int
	Width   int
	Focused bool
}

func RenderSegment(t theme.Theme, s Segment) string {
	if len(s.Options) == 0 || s.Width <= 0 {
		return ""
	}
	widths := SplitWidths(s.Width, 0, segmentWeights(len(s.Options))...)
	row := ""
	for i, option := range s.Options {
		row += segmentCell(t, option, widths[i], i == s.Active, s.Focused)
	}
	return row
}

func segmentCell(t theme.Theme, label string, width int, active, focused bool) string {
	style := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	if active {
		style = lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Surface)
	}
	if active && focused {
		style = lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Accent.Dim)
	}
	return style.Render(Center(Truncate(label, width), width))
}

func segmentWeights(count int) []int {
	out := make([]int, count)
	for i := range out {
		out[i] = 1
	}
	return out
}
