package keymap

import "github.com/charmbracelet/bubbles/key"

const (
	PanelExampleList   = "example.list"
	PanelExampleDetail = "example.detail"
	PanelHelp          = "help.body"

	PanelDashboardAccounts = "dashboard.accounts"
	PanelDashboardInsights = "dashboard.insights"
	PanelDashboardMode     = "dashboard.mode"
	PanelDashboardPeriod   = "dashboard.period"
	PanelDashboardOverview = "dashboard.overview"
)

type GlobalKeys struct {
	NextTab   key.Binding
	PrevTab   key.Binding
	NextFocus key.Binding
	PrevFocus key.Binding
	Jump      key.Binding
	Palette   key.Binding
	Help      key.Binding
	Cancel    key.Binding
	Confirm   key.Binding
	Quit      key.Binding
}

type ExampleKeys struct {
	Up      key.Binding
	Down    key.Binding
	Add     key.Binding
	Toggle  key.Binding
	Delete  key.Binding
	Refresh key.Binding
	Notify  key.Binding
}

type DashboardKeys struct {
	Up    key.Binding
	Down  key.Binding
	Prev  key.Binding
	Next  key.Binding
	Today key.Binding
}

type HelpKeys struct {
	ScrollUp   key.Binding
	ScrollDown key.Binding
}

type Keys struct {
	Global    GlobalKeys
	Example   ExampleKeys
	Dashboard DashboardKeys
	Help      HelpKeys
}

func Default() Keys {
	return Keys{
		Global: GlobalKeys{
			NextTab:   key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
			PrevTab:   key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("S-tab", "prev tab")),
			NextFocus: key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("C-n", "next panel")),
			PrevFocus: key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("C-p", "prev panel")),
			Jump:      key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "jump")),
			Palette:   key.NewBinding(key.WithKeys("ctrl+k"), key.WithHelp("C-k", "palette")),
			Help:      key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
			Cancel:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel")),
			Confirm:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm")),
			Quit:      key.NewBinding(key.WithKeys("ctrl+c", "q"), key.WithHelp("q", "quit")),
		},
		Example: ExampleKeys{
			Up:      key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
			Down:    key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
			Add:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
			Toggle:  key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
			Delete:  key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "delete")),
			Refresh: key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "refresh")),
			Notify:  key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "notify")),
		},
		Dashboard: DashboardKeys{
			Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "up")),
			Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "down")),
			Prev:  key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("h", "prev")),
			Next:  key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("l", "next")),
			Today: key.NewBinding(key.WithKeys("t"), key.WithHelp("t", "today")),
		},
		Help: HelpKeys{
			ScrollUp:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("k", "scroll up")),
			ScrollDown: key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("j", "scroll down")),
		},
	}
}

func (k Keys) Registry() *Registry {
	r := NewRegistry()
	r.RegisterGlobal(
		k.Global.Palette,
		k.Global.Jump,
		k.Global.NextTab,
		k.Global.PrevTab,
		k.Global.NextFocus,
		k.Global.PrevFocus,
		k.Global.Help,
		k.Global.Quit,
	)
	r.Register(PanelExampleList, "Items",
		k.Example.Up,
		k.Example.Down,
		k.Example.Add,
		k.Example.Toggle,
		k.Example.Delete,
		k.Example.Refresh,
	)
	r.Register(PanelExampleDetail, "Detail",
		k.Example.Notify,
	)
	r.Register(PanelDashboardAccounts, "Accounts",
		k.Dashboard.Up,
		k.Dashboard.Down,
	)
	r.Register(PanelDashboardInsights, "Insights",
		k.Dashboard.Prev,
		k.Dashboard.Next,
	)
	r.Register(PanelDashboardMode, "View and add",
		k.Dashboard.Prev,
		k.Dashboard.Next,
	)
	r.Register(PanelDashboardPeriod, "Period",
		k.Dashboard.Prev,
		k.Dashboard.Next,
		k.Dashboard.Today,
	)
	r.Register(PanelDashboardOverview, "Overview",
		k.Dashboard.Up,
		k.Dashboard.Down,
	)
	r.Register(PanelHelp, "Help",
		k.Help.ScrollUp,
		k.Help.ScrollDown,
	)
	return r
}
