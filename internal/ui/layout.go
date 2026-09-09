package ui

import "github.com/remcostoeten/reusable-tui/internal/theme"

func Row(t theme.Theme, blocks ...string) string {
	return joinColumns(blocks, Repeat(" ", t.Space.Gutter))
}

func Column(blocks ...string) string {
	rows := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if block == "" {
			continue
		}
		rows = append(rows, Lines(block)...)
	}
	return Join(rows)
}

func joinColumns(blocks []string, gutter string) string {
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
	return split(total, gutter, weights)
}

func SplitHeights(total, gutter int, weights ...int) []int {
	return split(total, gutter, weights)
}

func split(total, gutter int, weights []int) []int {
	sum := 0
	for _, w := range weights {
		sum += w
	}
	available := total - gutter*(len(weights)-1)
	out := make([]int, len(weights))
	used := 0
	for i, weight := range weights {
		out[i] = available * weight / sum
		used += out[i]
	}
	for i := 0; used < available; i = (i + 1) % len(out) {
		out[i]++
		used++
	}
	return out
}
