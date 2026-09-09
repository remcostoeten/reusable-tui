package dashboard

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
)

const (
	ScreenID    = "dashboard"
	ScreenTitle = "Dashboard"
)

type Mode int

const (
	ModeExpense Mode = iota
	ModeIncome
)

type Account struct {
	Name    string
	Balance float64
}

type Model struct {
	keys     keymap.Keys
	clock    func() time.Time
	accounts []Account
	cursor   int
	mode     Mode
	selected time.Time
	overview viewport.Model
	focused  string
}

func New(keys keymap.Keys, clock func() time.Time) *Model {
	if clock == nil {
		clock = time.Now
	}
	return &Model{
		keys:     keys,
		clock:    clock,
		accounts: defaultAccounts(),
		mode:     ModeExpense,
		selected: clock(),
		overview: viewport.New(0, 0),
		focused:  keymap.PanelDashboardAccounts,
	}
}

func defaultAccounts() []Account {
	return []Account{
		{Name: "Checking", Balance: 2480.15},
		{Name: "Savings", Balance: 9120.00},
		{Name: "Cash", Balance: 140.50},
	}
}

func (m *Model) ID() string {
	return ScreenID
}

func (m *Model) Title() string {
	return ScreenTitle
}

func (m *Model) Panels() []string {
	return []string{
		keymap.PanelDashboardAccounts,
		keymap.PanelDashboardInsights,
		keymap.PanelDashboardMode,
		keymap.PanelDashboardPeriod,
		keymap.PanelDashboardOverview,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Focus(panelID string) {
	m.focused = panelID
}

func (m *Model) Capturing() bool {
	return false
}

func (m *Model) Selected() (Account, bool) {
	if m.cursor < 0 || m.cursor >= len(m.accounts) {
		return Account{}, false
	}
	return m.accounts[m.cursor], true
}

func (m *Model) moveCursor(delta int) {
	if len(m.accounts) == 0 {
		return
	}
	m.cursor = (m.cursor + delta + len(m.accounts)) % len(m.accounts)
}

func (m *Model) shiftPeriod(delta int) {
	m.selected = m.selected.AddDate(0, 0, delta*daysPerWeek)
}

func (m *Model) shiftMode(delta int) {
	m.mode = Mode((int(m.mode) + delta + modeCount) % modeCount)
}

func (m *Model) today() {
	m.selected = m.clock()
}
