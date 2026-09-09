package app

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const paletteWidthRatio = 3

func (m *Model) View() string {
	if m.width < m.theme.Space.MinWidth || m.height < m.theme.Space.MinHeight {
		return m.tooSmall()
	}
	frame := ui.Join(m.frameRows())
	if !m.palette.Open() {
		return frame
	}
	return ui.PresentModal(m.theme, frame, m.palette.Modal(m.theme, m.paletteWidth()), m.width, m.height)
}

func (m *Model) frameRows() []string {
	rows := ui.Lines(ui.RenderHeader(m.theme, m.header()))
	m.hits.Reset()
	rows = append(rows, ui.Lines(m.ActiveScreen().View(m.context()))...)
	if m.toast.Visible {
		rows = append(rows, m.toast.Render(m.theme, m.width))
	}
	rows = append(rows, ui.RenderStatusBar(m.theme, ui.StatusBar{
		Width: m.width,
		Left:  m.bindings.For(m.focus[m.ActiveScreen().ID()]),
		Right: m.bindings.Global(),
		Busy:  m.busyLabel(),
	}))
	return rows
}

func (m *Model) busyLabel() string {
	if !m.busy.Active() {
		return ""
	}
	return ui.SpinnerFrame(m.theme, m.tick) + " " + m.busy.Label()
}

func (m *Model) header() ui.Header {
	return ui.Header{
		App:     Name,
		Version: Version,
		Tabs:    m.tabs(),
		Active:  m.active,
		Width:   m.width,
	}
}

func (m *Model) paletteWidth() int {
	width := m.width * 2 / paletteWidthRatio
	if width < m.theme.Space.MinWidth/2 {
		width = m.theme.Space.MinWidth / 2
	}
	return width
}

func (m *Model) tooSmall() string {
	style := lipgloss.NewStyle().Foreground(m.theme.Text.Secondary).Background(m.theme.Base.Background)
	accent := lipgloss.NewStyle().Foreground(m.theme.Accent.Active).Background(m.theme.Base.Background)
	rows := []string{
		accent.Render("terminal too small"),
		style.Render("minimum " + itoa(m.theme.Space.MinWidth) + "x" + itoa(m.theme.Space.MinHeight)),
		style.Render("current " + itoa(m.width) + "x" + itoa(m.height)),
	}
	return ui.Join(rows)
}
