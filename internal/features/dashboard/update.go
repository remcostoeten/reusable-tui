package dashboard

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch typed := msg.(type) {
	case ui.ClickMsg:
		return m.handleClick(typed)
	case ui.WheelMsg:
		return m.handleWheel(typed)
	case tea.KeyMsg:
		return m.handleKey(typed)
	}
	return nil
}

func (m *Model) handleClick(msg ui.ClickMsg) tea.Cmd {
	if msg.Panel == keymap.PanelDashboardAccounts && msg.Y >= 0 && msg.Y < len(m.accounts) {
		m.cursor = msg.Y
	}
	return nil
}

func (m *Model) handleWheel(msg ui.WheelMsg) tea.Cmd {
	switch msg.Panel {
	case keymap.PanelDashboardAccounts:
		m.moveCursor(msg.Delta)
	case keymap.PanelDashboardOverview:
		ui.Scroll(&m.overview, msg.Delta)
	case keymap.PanelDashboardPeriod, keymap.PanelDashboardInsights:
		m.shiftPeriod(msg.Delta)
	}
	return nil
}

func (m *Model) handleKey(keyMsg tea.KeyMsg) tea.Cmd {
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
