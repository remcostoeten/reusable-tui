package app

import "github.com/remcostoeten/reusable-tui/internal/config"

func (m *Model) applyFocus() {
	screen := m.ActiveScreen()
	screen.Focus(m.focus[screen.ID()])
}

func (m *Model) cycleTab(delta int) {
	m.active = (m.active + delta + len(m.screens)) % len(m.screens)
	m.jump.Hide()
	m.applyFocus()
}

func (m *Model) selectTab(screenID string) {
	for i, screen := range m.screens {
		if screen.ID() == screenID {
			m.active = i
			m.jump.Hide()
			m.applyFocus()
			return
		}
	}
}

func (m *Model) cycleFocus(delta int) {
	screen := m.ActiveScreen()
	panels := screen.Panels()
	current := indexOf(panels, m.focus[screen.ID()])
	if current < 0 {
		current = 0
	}
	next := (current + delta + len(panels)) % len(panels)
	m.focus[screen.ID()] = panels[next]
	m.applyFocus()
}

func (m *Model) focusPanel(panelID string) {
	screen := m.ActiveScreen()
	if indexOf(screen.Panels(), panelID) < 0 {
		return
	}
	m.focus[screen.ID()] = panelID
	m.applyFocus()
}

func indexOf(values []string, target string) int {
	for i, value := range values {
		if value == target {
			return i
		}
	}
	return -1
}

func (m *Model) restoreSession(session config.Session) {
	for screenID, panelID := range session.Panels {
		for _, screen := range m.screens {
			if screen.ID() == screenID && indexOf(screen.Panels(), panelID) >= 0 {
				m.focus[screenID] = panelID
			}
		}
	}
	if session.Screen != "" {
		m.selectTab(session.Screen)
	}
}

func (m *Model) session() config.Session {
	panels := make(map[string]string, len(m.focus))
	for screenID, panelID := range m.focus {
		panels[screenID] = panelID
	}
	return config.Session{Screen: m.ActiveScreen().ID(), Panels: panels}
}
