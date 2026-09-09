package ui

import "github.com/remcostoeten/reusable-tui/internal/theme"

func Row(t theme.Theme, blocks ...string) string {
	if len(blocks) == 0 {
		return ""
	}
	columns := make([][]string, len(blocks))
	height := 0
	for i, block := range blocks {
		columns[i] = Lines(block)
		if len(columns[i]) > height {
			height = len(columns[i])
		}
	}
	gutter := Repeat(" ", t.Space.Gutter)
	rows := make([]string, 0, height)
	for y := 0; y < height; y++ {
		line := ""
		for i, column := range columns {
			if i > 0 {
				line += gutter
			}
			if y < len(column) {
				line += column[y]
			}
		}
		rows = append(rows, line)
	}
	return Join(rows)
}

func SplitWidths(total, gutter int, weights ...int) []int {
	sum := 0
	for _, w := range weights {
		sum += w
	}
	available := total - gutter*(len(weights)-1)
	out := make([]int, len(weights))
	used := 0
	for i := 0; i < len(weights)-1; i++ {
		out[i] = available * weights[i] / sum
		used += out[i]
	}
	out[len(weights)-1] = available - used
	return out
}
