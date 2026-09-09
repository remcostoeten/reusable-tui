package ui

import (
	"strings"
	"testing"

	"github.com/remcostoeten/reusable-tui/internal/theme"
)

const emptyLabel = "nothing here"

func patternTheme(t *testing.T, pattern string) theme.Theme {
	t.Helper()
	built, err := theme.Apply(theme.VioletDark(), map[string]string{theme.TokenMarkerEmpty: pattern})
	if err != nil {
		t.Fatalf("apply %q: %v", pattern, err)
	}
	return built
}

func plainRows(rendered string) []string {
	rows := Lines(rendered)
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, stripStyles(row))
	}
	return out
}

func stripStyles(row string) string {
	plain := ""
	inEscape := false
	for _, r := range row {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		plain += string(r)
	}
	return plain
}

func TestEmptyStateHonoursThePattern(t *testing.T) {
	cases := map[string]string{
		string(theme.EmptyPatternSlash): "/",
		string(theme.EmptyPatternDots):  "·",
	}
	for pattern, glyph := range cases {
		t.Run(pattern, patternCase(pattern, glyph))
	}
}

func patternCase(pattern, glyph string) func(*testing.T) {
	return func(t *testing.T) {
		rows := plainRows(EmptyState(patternTheme(t, pattern), 40, 7, emptyLabel))
		body := strings.Join(rows, "\n")
		if !strings.Contains(body, glyph) {
			t.Fatalf("pattern %q did not render %q", pattern, glyph)
		}
		if strings.Contains(body, otherGlyph(glyph)) {
			t.Fatalf("pattern %q leaked %q", pattern, otherGlyph(glyph))
		}
	}
}

func otherGlyph(glyph string) string {
	if glyph == "/" {
		return "·"
	}
	return "/"
}

func TestEmptyStateHiddenPatternShowsOnlyTheLabel(t *testing.T) {
	rows := plainRows(EmptyState(patternTheme(t, string(theme.EmptyPatternNone)), 40, 7, emptyLabel))
	for i, row := range rows {
		trimmed := strings.TrimSpace(row)
		if i == len(rows)/2 {
			if trimmed != emptyLabel {
				t.Fatalf("label row is %q", trimmed)
			}
			continue
		}
		if trimmed != "" {
			t.Fatalf("row %d is not blank: %q", i, trimmed)
		}
	}
}

func TestEmptyStateKeepsItsBoxWhateverThePattern(t *testing.T) {
	for _, pattern := range []string{"slash", "dots", "none"} {
		rows := plainRows(EmptyState(patternTheme(t, pattern), 40, 7, emptyLabel))
		if len(rows) != 7 {
			t.Fatalf("pattern %q rendered %d rows", pattern, len(rows))
		}
		for i, row := range rows {
			if width := Width(row); width != 40 {
				t.Fatalf("pattern %q row %d is %d cells wide", pattern, i, width)
			}
		}
	}
}
