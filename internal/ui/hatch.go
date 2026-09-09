package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

const (
	hatchPeriod = 2
	hatchSlope  = 1
	dotStepX    = 3
	dotStepY    = 2
)

func slashCell(t theme.Theme, x, y int) string {
	if (x+y*hatchSlope)%hatchPeriod == 0 {
		return t.Marker.Hatch
	}
	return " "
}

func dotCell(t theme.Theme, x, y int) string {
	if x%dotStepX == 0 && y%dotStepY == 0 {
		return t.Marker.Dot
	}
	return " "
}

func patternCell(t theme.Theme, x, y int) string {
	switch t.Marker.Empty {
	case theme.EmptyPatternDots:
		return dotCell(t, x, y)
	case theme.EmptyPatternNone:
		return " "
	default:
		return slashCell(t, x, y)
	}
}

func patternRow(t theme.Theme, width, row int) []string {
	cells := make([]string, width)
	for x := 0; x < width; x++ {
		cells[x] = patternCell(t, x, row)
	}
	return cells
}

func Hatch(t theme.Theme, width, height int) string {
	style := lipgloss.NewStyle().Foreground(t.Border.Subtle).Background(t.Base.Background)
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		rows = append(rows, style.Render(joinCells(patternRow(t, width, y))))
	}
	return Join(rows)
}

func EmptyState(t theme.Theme, width, height int, label string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	fillStyle := lipgloss.NewStyle().Foreground(t.Border.Subtle).Background(t.Base.Background)
	labelStyle := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	text := " " + Truncate(label, width-2) + " "
	mid := height / 2
	start := (width - Width(text)) / 2
	if start < 0 {
		start = 0
	}
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		cells := patternRow(t, width, y)
		if y != mid {
			rows = append(rows, fillStyle.Render(joinCells(cells)))
			continue
		}
		left := fillStyle.Render(joinCells(cells[:start]))
		right := fillStyle.Render(joinCells(cells[start+Width(text):]))
		rows = append(rows, left+labelStyle.Render(text)+right)
	}
	return Join(rows)
}

func joinCells(cells []string) string {
	out := ""
	for _, c := range cells {
		out += c
	}
	return out
}
