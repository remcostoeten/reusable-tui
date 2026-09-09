package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/config"
	"github.com/remcostoeten/reusable-tui/internal/features/dashboard"
	"github.com/remcostoeten/reusable-tui/internal/features/example"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/notify"
	"github.com/remcostoeten/reusable-tui/internal/store"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const Name = "reusable-tui"

var Version = "v0.1.0"

type Options struct {
	Store      *store.Store
	Config     config.Config
	ConfigPath string
	Notifier   notify.Notifier
	Themes     *theme.Registry
	ThemeDir   string
	Fidelity   theme.Fidelity
	Warnings   []error
	Clock      func() time.Time
}

type Model struct {
	opts     Options
	keys     keymap.Keys
	bindings *keymap.Registry
	commands *ui.CommandRegistry
	theme    theme.Theme
	screens  []ui.Screen
	active   int
	focus    map[string]string
	palette  ui.Palette
	jump     ui.Jump
	toast    ui.Toast
	toastSeq int
	busy     ui.BusySet
	hits     *ui.HitMap
	width    int
	height   int
	tick     int
}

func New(opts Options) *Model {
	keys := keymap.Default()
	if opts.Clock == nil {
		opts.Clock = time.Now
	}
	screens := []ui.Screen{
		example.New(opts.Store.DB(), keys),
		dashboard.New(keys, opts.Clock),
		newHelpScreen(keys),
	}
	m := &Model{
		opts:     opts,
		keys:     keys,
		bindings: keys.Registry(),
		commands: ui.NewCommandRegistry(),
		theme:    opts.Themes.Resolve(opts.Config.Theme, opts.Fidelity),
		screens:  screens,
		focus:    map[string]string{},
		jump:     ui.NewJump(),
		busy:     ui.NewBusySet(),
		hits:     ui.NewHitMap(),
	}
	m.registerCommands()
	m.palette = ui.NewPalette(m.commands)
	for _, screen := range screens {
		m.focus[screen.ID()] = screen.Panels()[0]
	}
	m.restoreSession(opts.Config.Session)
	m.applyFocus()
	return m
}

func (m *Model) Init() tea.Cmd {
	cmds := []tea.Cmd{ui.Tick()}
	for _, screen := range m.screens {
		cmds = append(cmds, screen.Init())
	}
	if warning := joinWarnings(m.opts.Warnings); warning != "" {
		cmds = append(cmds, ui.Warn(warning))
	}
	return tea.Batch(cmds...)
}

func (m *Model) Theme() theme.Theme {
	return m.theme
}

func (m *Model) ActiveScreen() ui.Screen {
	return m.screens[m.active]
}

func (m *Model) tabs() []ui.Tab {
	out := make([]ui.Tab, 0, len(m.screens))
	for _, screen := range m.screens {
		out = append(out, ui.Tab{ID: screen.ID(), Label: screen.Title()})
	}
	return out
}

func (m *Model) context() ui.RenderContext {
	return ui.RenderContext{
		Theme:   m.theme,
		Keys:    m.bindings,
		Width:   m.width,
		Height:  m.bodyHeight(),
		Tick:    m.tick,
		Focused: m.focus[m.ActiveScreen().ID()],
		Jump:    m.jump,
		Hits:    m.hits,
	}
}

func (m *Model) bodyHeight() int {
	return m.height - m.theme.Space.HeaderRows - m.theme.Space.StatusRows - m.toastRows()
}

func (m *Model) toastRows() int {
	if m.toast.Visible {
		return 1
	}
	return 0
}
