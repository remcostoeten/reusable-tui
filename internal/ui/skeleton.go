package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

var skeletonRatios = []int{100, 72, 88, 54, 94, 66}

func skeletonPulse(t theme.Theme, phase int) lipgloss.TerminalColor {
	shades := []lipgloss.TerminalColor{t.Base.Surface, t.Border.Subtle, t.Border.Unfocused, t.Border.Subtle}
	return shades[phase%len(shades)]
}

func Skeleton(t theme.Theme, width, height, phase int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	style := lipgloss.NewStyle().Foreground(skeletonPulse(t, phase)).Background(t.Base.Background)
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		if y%2 == 1 {
			rows = append(rows, Repeat(" ", width))
			continue
		}
		bar := width * skeletonRatios[(y/2)%len(skeletonRatios)] / 100
		rows = append(rows, style.Render(Repeat(t.Marker.Skeleton, bar))+Repeat(" ", width-bar))
	}
	return Join(rows)
}
