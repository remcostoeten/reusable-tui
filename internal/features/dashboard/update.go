package dashboard

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
)

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch m.focused {
	case keymap.PanelDashboardAccounts:
		return m.handleAccountsKey(keyMsg)
	case keymap.PanelDashboardMode:
		return m.handleModeKey(keyMsg)
	case keymap.PanelDashboardPeriod, keymap.PanelDashboardInsights:
		return m.handlePeriodKey(keyMsg)
	case keymap.PanelDashboardOverview:
		return m.handleOverviewKey(keyMsg)
	}
	return nil
}

func (m *Model) handleAccountsKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Dashboard.Up):
		m.moveCursor(-1)
	case key.Matches(msg, m.keys.Dashboard.Down):
		m.moveCursor(1)
	}
	return nil
}

func (m *Model) handleModeKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Dashboard.Prev):
		m.shiftMode(-1)
	case key.Matches(msg, m.keys.Dashboard.Next):
		m.shiftMode(1)
	}
	return nil
}

func (m *Model) handlePeriodKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Dashboard.Prev):
		m.shiftPeriod(-1)
	case key.Matches(msg, m.keys.Dashboard.Next):
		m.shiftPeriod(1)
	case key.Matches(msg, m.keys.Dashboard.Today):
		m.today()
	}
	return nil
}

func (m *Model) handleOverviewKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Dashboard.Up):
		m.overview.ScrollUp(1)
	case key.Matches(msg, m.keys.Dashboard.Down):
		m.overview.ScrollDown(1)
	}
	return nil
}
