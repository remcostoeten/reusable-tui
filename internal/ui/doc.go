package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type BlockKind int

const (
	BlockHeading BlockKind = iota
	BlockText
	BlockBullet
	BlockCallout
	BlockRule
	BlockSpacer
)

type DocBlock struct {
	Kind BlockKind
	Text string
}

func Heading(text string) DocBlock {
	return DocBlock{Kind: BlockHeading, Text: text}
}

func Text(text string) DocBlock {
	return DocBlock{Kind: BlockText, Text: text}
}

func Bullet(text string) DocBlock {
	return DocBlock{Kind: BlockBullet, Text: text}
}

func Callout(text string) DocBlock {
	return DocBlock{Kind: BlockCallout, Text: text}
}

func Divider() DocBlock {
	return DocBlock{Kind: BlockRule}
}

func Spacer() DocBlock {
	return DocBlock{Kind: BlockSpacer}
}

func RenderDoc(t theme.Theme, width int, blocks ...DocBlock) string {
	if width <= 0 {
		return ""
	}
	rows := make([]string, 0, len(blocks)*2)
	for _, block := range blocks {
		rows = append(rows, docBlock(t, width, block)...)
	}
	return Join(rows)
}

func docBlock(t theme.Theme, width int, b DocBlock) []string {
	switch b.Kind {
	case BlockHeading:
		return []string{headingRow(t, width, b.Text)}
	case BlockBullet:
		return bulletRows(t, width, b.Text)
	case BlockCallout:
		return calloutRows(t, width, b.Text)
	case BlockRule:
		return []string{Rule(t, width, t.Border.Unfocused)}
	case BlockSpacer:
		return []string{blankRow(t, width)}
	}
	return textRows(t, width, b.Text)
}

func headingRow(t theme.Theme, width int, text string) string {
	style := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background).Underline(true)
	return style.Render(Fit(text, width))
}

func textRows(t theme.Theme, width int, text string) []string {
	style := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	lines := Wrap(text, width)
	rows := make([]string, 0, len(lines))
	for _, line := range lines {
		rows = append(rows, style.Render(Fit(line, width)))
	}
	return rows
}

func bulletRows(t theme.Theme, width int, text string) []string {
	marker := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)
	style := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	indent := Width(t.Marker.Bullet) + 1
	lines := Wrap(text, width-indent)
	rows := make([]string, 0, len(lines))
	for i, line := range lines {
		lead := marker.Render(t.Marker.Bullet + " ")
		if i > 0 {
			lead = style.Render(Repeat(" ", indent))
		}
		rows = append(rows, lead+style.Render(Fit(line, width-indent)))
	}
	return rows
}

func calloutRows(t theme.Theme, width int, text string) []string {
	bar := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Surface)
	style := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Surface)
	indent := 2
	lines := Wrap(text, width-indent-1)
	rows := make([]string, 0, len(lines))
	for _, line := range lines {
		rows = append(rows, bar.Render("│ ")+style.Render(Fit(line, width-indent)))
	}
	return rows
}

func blankRow(t theme.Theme, width int) string {
	return lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", width))
}
