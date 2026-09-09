package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = typed.Width
		m.height = typed.Height
		return m, m.broadcast(msg)
	case ui.TickMsg:
		m.tick++
		return m, tea.Batch(ui.Tick(), m.broadcast(msg))
	case ui.ToastExpiredMsg:
		m.expireToast(typed.Seq)
		return m, nil
	case ui.ErrorMsg:
		return m, m.showToast(typed.Err.Error(), ui.ToastError)
	case ui.InfoMsg:
		return m, m.showToast(typed.Text, ui.ToastInfo)
	case ui.SuccessMsg:
		return m, m.showToast(typed.Text, ui.ToastSuccess)
	case ui.WarnMsg:
		return m, m.showToast(typed.Text, ui.ToastWarning)
	case themesReloadedMsg:
		return m, m.adoptThemes(typed)
	case exportThemeMsg:
		return m, exportTheme(m.opts.ThemeDir, m.theme)
	case ui.NotifyMsg:
		return m, notifyCmd(m.opts.Notifier, typed.Title, typed.Body)
	case ui.ThemeMsg:
		return m, m.switchTheme(typed.Name)
	case ui.FocusMsg:
		m.focusPanel(typed.PanelID)
		return m, nil
	case screenMsg:
		m.selectTab(typed.ID)
		return m, nil
	case tea.KeyMsg:
		return m, m.handleKey(typed)
	}
	return m, m.broadcast(msg)
}

func (m *Model) broadcast(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.screens))
	for _, screen := range m.screens {
		cmds = append(cmds, screen.Update(msg))
	}
	return tea.Batch(cmds...)
}

func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	if m.palette.Open() {
		return m.palette.Update(msg)
	}
	if m.jump.Active() {
		return m.handleJumpKey(msg)
	}
	if m.ActiveScreen().Capturing() {
		return m.ActiveScreen().Update(msg)
	}
	if cmd, handled := m.handleGlobalKey(msg); handled {
		return cmd
	}
	return m.ActiveScreen().Update(msg)
}

func (m *Model) handleJumpKey(msg tea.KeyMsg) tea.Cmd {
	if key.Matches(msg, m.keys.Global.Cancel) {
		m.jump.Hide()
		return nil
	}
	if target, ok := m.jump.Target(msg.String()); ok {
		m.jump.Hide()
		m.focusPanel(target)
	}
	return nil
}

func (m *Model) handleGlobalKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch {
	case key.Matches(msg, m.keys.Global.Quit):
		return tea.Quit, true
	case key.Matches(msg, m.keys.Global.NextTab):
		m.cycleTab(1)
		return nil, true
	case key.Matches(msg, m.keys.Global.PrevTab):
		m.cycleTab(-1)
		return nil, true
	case key.Matches(msg, m.keys.Global.NextFocus):
		m.cycleFocus(1)
		return nil, true
	case key.Matches(msg, m.keys.Global.PrevFocus):
		m.cycleFocus(-1)
		return nil, true
	case key.Matches(msg, m.keys.Global.Jump):
		m.jump.Show(m.ActiveScreen().Panels())
		return nil, true
	case key.Matches(msg, m.keys.Global.Palette):
		m.palette.Show()
		return nil, true
	case key.Matches(msg, m.keys.Global.Help):
		m.selectTab(helpScreenID)
		return nil, true
	}
	return nil, false
}

func (m *Model) showToast(text string, kind ui.ToastKind) tea.Cmd {
	m.toastSeq++
	m.toast = ui.Toast{Text: text, Kind: kind, Seq: m.toastSeq, Visible: true}
	return ui.ExpireToast(m.toastSeq)
}

func (m *Model) expireToast(seq int) {
	if m.toast.Seq == seq {
		m.toast.Visible = false
	}
}

func (m *Model) switchTheme(name string) tea.Cmd {
	if _, ok := m.opts.Themes.Get(name); !ok {
		return ui.Fail(errUnknownTheme(name))
	}
	m.theme = m.opts.Themes.Resolve(name, m.opts.Fidelity)
	m.opts.Config.Theme = name
	return persistConfig(m.opts.ConfigPath, m.opts.Config)
}

func (m *Model) adoptThemes(msg themesReloadedMsg) tea.Cmd {
	m.opts.Themes = msg.registry
	m.theme = msg.registry.Resolve(m.opts.Config.Theme, m.opts.Fidelity)
	m.commands = ui.NewCommandRegistry()
	m.registerCommands()
	m.palette = ui.NewPalette(m.commands)
	if warning := joinWarnings(msg.warnings); warning != "" {
		return ui.Warn(warning)
	}
	return ui.Success("themes reloaded")
}
