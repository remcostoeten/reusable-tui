package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type Modal struct {
	Title  string
	Body   string
	Width  int
	Height int
}

func RenderModal(t theme.Theme, m Modal) string {
	set := t.Border.FocusedSet
	edge := lipgloss.NewStyle().Foreground(t.Border.Focused).Background(t.Base.Overlay)
	title := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Overlay)
	surface := lipgloss.NewStyle().Background(t.Base.Overlay)
	inner := m.Width - 2

	head := Repeat(set.Top, t.Space.TitleInset) + " " + m.Title + " "
	fill := inner - Width(head)
	if fill < 0 {
		fill = 0
	}
	top := edge.Render(set.TopLeft+Repeat(set.Top, t.Space.TitleInset)) +
		title.Render(" "+m.Title+" ") +
		edge.Render(Repeat(set.Top, fill)+set.TopRight)

	rows := []string{top}
	content := Block(m.Body, inner-2*t.Space.ModalPadX, m.Height-2-2*t.Space.ModalPadY)
	blank := surface.Render(Repeat(" ", inner))
	pad := surface.Render(Repeat(" ", t.Space.ModalPadX))
	for i := 0; i < t.Space.ModalPadY; i++ {
		rows = append(rows, edge.Render(set.Left)+blank+edge.Render(set.Right))
	}
	for _, line := range content {
		rows = append(rows, edge.Render(set.Left)+pad+line+pad+edge.Render(set.Right))
	}
	for len(rows) < m.Height-1 {
		rows = append(rows, edge.Render(set.Left)+blank+edge.Render(set.Right))
	}
	rows = append(rows, edge.Render(set.BottomLeft+Repeat(set.Bottom, inner)+set.BottomRight))
	return Join(rows[:m.Height])
}

func PresentModal(t theme.Theme, background string, m Modal, width, height int) string {
	dim := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	box := RenderModal(t, m)
	x := (width - m.Width) / 2
	y := (height - m.Height) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	return Overlay(background, box, x, y, dim)
}
