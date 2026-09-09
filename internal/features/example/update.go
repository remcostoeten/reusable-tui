package example

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch typed := msg.(type) {
	case itemsLoadedMsg:
		m.items = typed.items
		m.loading = false
		m.clampCursor()
		return nil
	case itemChangedMsg:
		return tea.Batch(loadItems(m.db), ui.Success(typed.action+" "+typed.title))
	case tea.KeyMsg:
		return m.handleKey(typed)
	}
	return nil
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	if m.adding {
		return m.handleAddKey(msg)
	}
	switch m.focused {
	case keymap.PanelExampleList:
		return m.handleListKey(msg)
	case keymap.PanelExampleDetail:
		return m.handleDetailKey(msg)
	}
	return nil
}

func (m *Model) handleAddKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Global.Cancel):
		m.stopAdding()
		return nil
	case key.Matches(msg, m.keys.Global.Confirm):
		title := m.input.Value()
		m.stopAdding()
		return addItem(m.db, title)
	}
	in, cmd := m.input.Update(msg)
	m.input = in
	return cmd
}

func (m *Model) handleListKey(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keys.Example.Up):
		m.moveCursor(-1)
	case key.Matches(msg, m.keys.Example.Down):
		m.moveCursor(1)
	case key.Matches(msg, m.keys.Example.Add):
		return m.startAdding()
	case key.Matches(msg, m.keys.Example.Toggle):
		return m.mutateSelected(flipSelected)
	case key.Matches(msg, m.keys.Example.Delete):
		return m.mutateSelected(removeSelected)
	case key.Matches(msg, m.keys.Example.Refresh):
		m.loading = true
		return loadItems(m.db)
	}
	return nil
}

func (m *Model) handleDetailKey(msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, m.keys.Example.Notify) {
		return m.notifySelected()
	}
	return nil
}

func (m *Model) mutateSelected(run func(*Model, Item) tea.Cmd) tea.Cmd {
	item, ok := m.Selected()
	if !ok {
		return ui.Info("nothing selected")
	}
	return run(m, item)
}

func flipSelected(m *Model, item Item) tea.Cmd {
	return flipItem(m.db, item.ID)
}

func removeSelected(m *Model, item Item) tea.Cmd {
	return removeItem(m.db, item.ID, item.Title)
}
