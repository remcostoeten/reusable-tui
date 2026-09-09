package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type Panel struct {
	Title   string
	Badge   string
	Hint    string
	Width   int
	Height  int
	Focused bool
	Body    string
}

func (p Panel) InnerWidth() int {
	return p.Width - 2
}

func (p Panel) ContentWidth(t theme.Theme) int {
	return p.InnerWidth() - 2*t.Space.PanelPadX
}

func (p Panel) ContentHeight(t theme.Theme) int {
	return p.Height - 2 - 2*t.Space.PanelPadY
}

func RenderPanel(t theme.Theme, p Panel) string {
	if p.Width < 4 || p.Height < 3 {
		return ""
	}
	set := t.BorderSet(p.Focused)
	edge := lipgloss.NewStyle().Foreground(t.BorderColor(p.Focused)).Background(t.Base.Background)
	rows := make([]string, 0, p.Height)
	rows = append(rows, panelTop(t, p, set, edge))
	for _, line := range panelBody(t, p) {
		rows = append(rows, edge.Render(set.Left)+line+edge.Render(set.Right))
	}
	rows = append(rows, edge.Render(set.BottomLeft+Repeat(set.Bottom, p.InnerWidth())+set.BottomRight))
	return Join(rows)
}

func panelTop(t theme.Theme, p Panel, set lipgloss.Border, edge lipgloss.Style) string {
	titleStyle := lipgloss.NewStyle().Foreground(t.TitleColor(p.Focused)).Background(t.Base.Background)
	badgeStyle := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	hintStyle := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)

	head := Repeat(set.Top, t.Space.TitleInset)
	plain := head
	styled := edge.Render(head)

	if p.Hint != "" {
		hint := t.Marker.HintOpen + p.Hint + t.Marker.HintClose
		plain += hint + Repeat(set.Top, t.Space.TitleInset)
		styled += hintStyle.Render(hint) + edge.Render(Repeat(set.Top, t.Space.TitleInset))
	}

	if p.Title != "" {
		title := " " + panelTitle(t, p) + " "
		plain += title
		styled += titleStyle.Render(title)
	}

	tail := ""
	tailStyled := ""
	if p.Badge != "" {
		badge := " " + p.Badge + " "
		tail = badge + Repeat(set.Top, t.Space.BadgeInset)
		tailStyled = badgeStyle.Render(badge) + edge.Render(Repeat(set.Top, t.Space.BadgeInset))
	}

	fill := p.InnerWidth() - Width(plain) - Width(tail)
	if fill < 0 {
		fill = 0
	}
	return edge.Render(set.TopLeft) + styled + edge.Render(Repeat(set.Top, fill)) + tailStyled + edge.Render(set.TopRight)
}

func panelTitle(t theme.Theme, p Panel) string {
	if p.Focused {
		return t.Marker.FocusOpen + p.Title + t.Marker.FocusClose
	}
	if t.Marker.FocusOpen == "" {
		return p.Title
	}
	return Repeat(" ", Width(t.Marker.FocusOpen)) + p.Title + Repeat(" ", Width(t.Marker.FocusClose))
}

func panelBody(t theme.Theme, p Panel) []string {
	pad := lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", t.Space.PanelPadX))
	content := Block(p.Body, p.ContentWidth(t), p.ContentHeight(t))
	rows := make([]string, 0, p.Height-2)
	blank := lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", p.InnerWidth()))
	for i := 0; i < t.Space.PanelPadY; i++ {
		rows = append(rows, blank)
	}
	for _, line := range content {
		rows = append(rows, pad+line+pad)
	}
	for len(rows) < p.Height-2 {
		rows = append(rows, blank)
	}
	return rows[:p.Height-2]
}

func StatusLabel(t theme.Theme, s theme.StatusToken) string {
	style := lipgloss.NewStyle().Foreground(s.Fg)
	return style.Render(s.Glyph + " " + s.Label)
}

func Rule(t theme.Theme, width int, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().Foreground(color).Background(t.Base.Background).Render(Repeat(t.Marker.Rule, width))
}
