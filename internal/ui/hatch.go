package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

const hatchPeriod = 3

func hatchRow(t theme.Theme, width, row int) []string {
	cells := make([]string, width)
	for x := 0; x < width; x++ {
		if (x+row)%hatchPeriod == 0 {
			cells[x] = t.Marker.Hatch
			continue
		}
		cells[x] = " "
	}
	return cells
}

func Hatch(t theme.Theme, width, height int) string {
	style := lipgloss.NewStyle().Foreground(t.Border.Subtle).Background(t.Base.Background)
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		rows = append(rows, style.Render(Join(hatchRow(t, width, y))))
	}
	return Join(rows)
}

func EmptyState(t theme.Theme, width, height int, label string) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	hatchStyle := lipgloss.NewStyle().Foreground(t.Border.Subtle).Background(t.Base.Background)
	labelStyle := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	text := " " + Truncate(label, width-2) + " "
	mid := height / 2
	start := (width - Width(text)) / 2
	if start < 0 {
		start = 0
	}
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		cells := hatchRow(t, width, y)
		if y != mid {
			rows = append(rows, hatchStyle.Render(joinCells(cells)))
			continue
		}
		left := hatchStyle.Render(joinCells(cells[:start]))
		right := hatchStyle.Render(joinCells(cells[start+Width(text):]))
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
