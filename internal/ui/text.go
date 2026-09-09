package ui

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func Width(s string) int {
	return ansi.StringWidth(s)
}

func Repeat(glyph string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(glyph, n)
}

func PadRight(s string, width int) string {
	return s + Repeat(" ", width-Width(s))
}

func PadLeft(s string, width int) string {
	return Repeat(" ", width-Width(s)) + s
}

func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if Width(s) <= width {
		return s
	}
	if width == 1 {
		return "."
	}
	return ansi.Truncate(s, width-1, "") + "."
}

func Fit(s string, width int) string {
	return PadRight(Truncate(s, width), width)
}

func Center(s string, width int) string {
	pad := width - Width(s)
	if pad <= 0 {
		return Truncate(s, width)
	}
	left := pad / 2
	return Repeat(" ", left) + s + Repeat(" ", pad-left)
}

func Lines(s string) []string {
	return strings.Split(s, "\n")
}

func Block(s string, width, height int) []string {
	raw := Lines(s)
	out := make([]string, 0, height)
	for i := 0; i < height; i++ {
		if i < len(raw) {
			out = append(out, Fit(raw[i], width))
			continue
		}
		out = append(out, Repeat(" ", width))
	}
	return out
}

func Join(lines []string) string {
	return strings.Join(lines, "\n")
}
