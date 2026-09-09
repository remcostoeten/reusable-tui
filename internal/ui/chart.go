package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

var barGlyphs = []string{" ", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

type Bar struct {
	Label string
	Value float64
}

type BarChart struct {
	Bars   []Bar
	Width  int
	Height int
	Labels bool
}

func RenderBarChart(t theme.Theme, c BarChart) string {
	if len(c.Bars) == 0 || c.Width <= 0 || c.Height <= 1 {
		return ""
	}
	plot := c.Height
	if c.Labels {
		plot--
	}
	widths := SplitWidths(c.Width, 0, segmentWeights(len(c.Bars))...)
	levels := barLevels(c.Bars, plot)
	rows := make([]string, 0, c.Height)
	for row := 0; row < plot; row++ {
		rows = append(rows, barRow(t, levels, widths, plot-row))
	}
	if c.Labels {
		rows = append(rows, barLabels(t, c.Bars, widths))
	}
	return Join(rows)
}

func barRow(t theme.Theme, levels []int, widths []int, cell int) string {
	style := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)
	row := ""
	for i, level := range levels {
		row += style.Render(Repeat(barGlyph(level, cell), widths[i]-1)) +
			lipgloss.NewStyle().Background(t.Base.Background).Render(" ")
	}
	return row
}

func barGlyph(level, cell int) string {
	steps := len(barGlyphs) - 1
	filled := level - (cell-1)*steps
	if filled <= 0 {
		return barGlyphs[0]
	}
	if filled >= steps {
		return barGlyphs[steps]
	}
	return barGlyphs[filled]
}

func barLevels(bars []Bar, plot int) []int {
	peak := 0.0
	for _, bar := range bars {
		if bar.Value > peak {
			peak = bar.Value
		}
	}
	steps := plot * (len(barGlyphs) - 1)
	out := make([]int, len(bars))
	for i, bar := range bars {
		if peak <= 0 {
			continue
		}
		out[i] = int(bar.Value / peak * float64(steps))
	}
	return out
}

func barLabels(t theme.Theme, bars []Bar, widths []int) string {
	style := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	row := ""
	for i, bar := range bars {
		row += style.Render(Center(Truncate(bar.Label, widths[i]), widths[i]))
	}
	return row
}
