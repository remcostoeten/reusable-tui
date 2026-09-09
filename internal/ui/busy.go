package ui

import (
	"sort"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

type BusyMsg struct {
	ID    string
	Label string
}

type IdleMsg struct {
	ID string
}

func Busy(id, label string) tea.Cmd {
	return func() tea.Msg { return BusyMsg{ID: id, Label: label} }
}

func Idle(id string) tea.Cmd {
	return func() tea.Msg { return IdleMsg{ID: id} }
}

type BusySet struct {
	labels map[string]string
}

func NewBusySet() BusySet {
	return BusySet{labels: map[string]string{}}
}

func (b *BusySet) Start(id, label string) {
	b.labels[id] = label
}

func (b *BusySet) Stop(id string) {
	delete(b.labels, id)
}

func (b BusySet) Active() bool {
	return len(b.labels) > 0
}

func (b BusySet) Label() string {
	ids := make([]string, 0, len(b.labels))
	for id := range b.labels {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return ""
	}
	label := b.labels[ids[0]]
	if len(ids) > 1 {
		label += " +" + strconv.Itoa(len(ids)-1)
	}
	return label
}

func SpinnerFrame(t theme.Theme, tick int) string {
	frames := []rune(t.Marker.Spinner)
	if len(frames) == 0 {
		return ""
	}
	return string(frames[tick%len(frames)])
}
