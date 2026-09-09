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

func Wrap(s string, width int) []string {
	if width <= 0 {
		return nil
	}
	out := make([]string, 0, 4)
	for _, paragraph := range Lines(s) {
		out = append(out, wrapLine(paragraph, width)...)
	}
	return out
}

func wrapLine(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	rows := make([]string, 0, 4)
	current := words[0]
	for _, word := range words[1:] {
		if Width(current)+1+Width(word) > width {
			rows = append(rows, current)
			current = word
			continue
		}
		current += " " + word
	}
	return append(rows, current)
}
