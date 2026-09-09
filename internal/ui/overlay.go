package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func Overlay(background, box string, x, y int, dim lipgloss.Style) string {
	bg := Lines(ansi.Strip(background))
	fg := Lines(box)
	boxWidth := 0
	for _, line := range fg {
		if w := Width(line); w > boxWidth {
			boxWidth = w
		}
	}
	rows := make([]string, 0, len(bg))
	for i, line := range bg {
		if i < y || i >= y+len(fg) {
			rows = append(rows, dim.Render(line))
			continue
		}
		rows = append(rows, overlayRow(line, Fit(fg[i-y], boxWidth), x, boxWidth, dim))
	}
	return Join(rows)
}

func overlayRow(bg, fg string, x, boxWidth int, dim lipgloss.Style) string {
	runes := []rune(bg)
	left := ""
	if x <= len(runes) {
		left = string(runes[:x])
	} else {
		left = bg + Repeat(" ", x-len(runes))
	}
	right := ""
	if end := x + boxWidth; end < len(runes) {
		right = string(runes[end:])
	}
	return dim.Render(left) + fg + dim.Render(right)
}
