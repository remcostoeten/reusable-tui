package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func (m *Model) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if m.palette.Open() || m.jump.Active() {
		return nil
	}
	if delta := ui.WheelDelta(msg); delta != 0 {
		return m.handleWheel(msg, delta)
	}
	if !ui.IsClick(msg) {
		return nil
	}
	if msg.Y < m.theme.Space.HeaderRows {
		return m.handleHeaderClick(msg.X)
	}
	return m.handleBodyClick(msg)
}

func (m *Model) handleHeaderClick(x int) tea.Cmd {
	if id, ok := ui.TabAt(m.theme, m.header(), x); ok {
		m.selectTab(id)
	}
	return nil
}

func (m *Model) handleBodyClick(msg tea.MouseMsg) tea.Cmd {
	panelID, rect, ok := m.hits.At(msg.X, msg.Y-m.theme.Space.HeaderRows)
	if !ok {
		return nil
	}
	m.focusPanel(panelID)
	local := m.contentPoint(msg, rect)
	return m.ActiveScreen().Update(ui.ClickMsg{Panel: panelID, X: local.X, Y: local.Y})
}

func (m *Model) handleWheel(msg tea.MouseMsg, delta int) tea.Cmd {
	panelID, _, ok := m.hits.At(msg.X, msg.Y-m.theme.Space.HeaderRows)
	if !ok {
		return nil
	}
	return m.ActiveScreen().Update(ui.WheelMsg{Panel: panelID, Delta: delta})
}

type point struct {
	X int
	Y int
}

func (m *Model) contentPoint(msg tea.MouseMsg, rect ui.Rect) point {
	return point{
		X: msg.X - rect.X - 1 - m.theme.Space.PanelPadX,
		Y: msg.Y - m.theme.Space.HeaderRows - rect.Y - 1 - m.theme.Space.PanelPadY,
	}
}
