package example

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const (
	ScreenID    = "example"
	ScreenTitle = "Example"
	busyLoad    = "example.load"
)

type Model struct {
	db      *sql.DB
	keys    keymap.Keys
	input   textinput.Model
	items   []Item
	cursor  int
	focused string
	loading bool
	adding  bool
}

func New(db *sql.DB, keys keymap.Keys) *Model {
	in := textinput.New()
	in.Prompt = "+ "
	in.Placeholder = "item title"
	in.CharLimit = 80
	return &Model{
		db:      db,
		keys:    keys,
		input:   in,
		loading: true,
		focused: keymap.PanelExampleList,
	}
}

func (m *Model) ID() string {
	return ScreenID
}

func (m *Model) Title() string {
	return ScreenTitle
}

func (m *Model) Panels() []string {
	return []string{keymap.PanelExampleList, keymap.PanelExampleDetail}
}

func (m *Model) Init() tea.Cmd {
	return m.reload()
}

func (m *Model) reload() tea.Cmd {
	m.loading = true
	return tea.Batch(ui.Busy(busyLoad, "loading items"), loadItems(m.db))
}

func (m *Model) selectRow(row int) {
	if m.adding {
		row--
	}
	if row < 0 || row >= len(m.items) {
		return
	}
	m.cursor = row
}

func (m *Model) Focus(panelID string) {
	m.focused = panelID
}

func (m *Model) Selected() (Item, bool) {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		return Item{}, false
	}
	return m.items[m.cursor], true
}

func (m *Model) Adding() bool {
	return m.adding
}

func (m *Model) startAdding() tea.Cmd {
	m.adding = true
	m.input.SetValue("")
	return m.input.Focus()
}

func (m *Model) stopAdding() {
	m.adding = false
	m.input.Blur()
	m.input.SetValue("")
}

func (m *Model) clampCursor() {
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

func (m *Model) moveCursor(delta int) {
	if len(m.items) == 0 {
		return
	}
	m.cursor = (m.cursor + delta + len(m.items)) % len(m.items)
}

func (m *Model) notifySelected() tea.Cmd {
	item, ok := m.Selected()
	if !ok {
		return ui.Info("nothing selected")
	}
	return ui.Notify("Example item", item.Title)
}

func (m *Model) Capturing() bool {
	return m.adding
}
