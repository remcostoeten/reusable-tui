package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const TickInterval = 180 * time.Millisecond

type TickMsg struct {
	Time time.Time
}

func newTick(at time.Time) tea.Msg {
	return TickMsg{Time: at}
}

func Tick() tea.Cmd {
	return tea.Tick(TickInterval, newTick)
}
