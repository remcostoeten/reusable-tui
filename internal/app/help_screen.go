package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const (
	helpScreenID    = "help"
	helpScreenTitle = "Help"
)

type helpScreen struct {
	keys     keymap.Keys
	bindings *keymap.Registry
	view     viewport.Model
	focused  string
}

func newHelpScreen(keys keymap.Keys) *helpScreen {
	return &helpScreen{
		keys:     keys,
		bindings: keys.Registry(),
		view:     viewport.New(0, 0),
		focused:  keymap.PanelHelp,
	}
}

func (s *helpScreen) ID() string {
	return helpScreenID
}

func (s *helpScreen) Title() string {
	return helpScreenTitle
}

func (s *helpScreen) Panels() []string {
	return []string{keymap.PanelHelp}
}

func (s *helpScreen) Init() tea.Cmd {
	return nil
}

func (s *helpScreen) Focus(panelID string) {
	s.focused = panelID
}

func (s *helpScreen) Capturing() bool {
	return false
}

func (s *helpScreen) Update(msg tea.Msg) tea.Cmd {
	if wheel, ok := msg.(ui.WheelMsg); ok {
		ui.Scroll(&s.view, wheel.Delta)
		return nil
	}
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch {
	case key.Matches(keyMsg, s.keys.Help.ScrollUp):
		s.view.ScrollUp(1)
	case key.Matches(keyMsg, s.keys.Help.ScrollDown):
		s.view.ScrollDown(1)
	}
	return nil
}

func (s *helpScreen) View(ctx ui.RenderContext) string {
	ctx.Hits.Add(keymap.PanelHelp, ui.Rect{X: 0, Y: 0, Width: ctx.Width, Height: ctx.Height})
	panel := ui.Panel{
		ID:      keymap.PanelHelp,
		Title:   "Keybindings",
		Hint:    ctx.Jump.Hint(keymap.PanelHelp),
		Width:   ctx.Width,
		Height:  ctx.Height,
		Focused: ctx.Focused == keymap.PanelHelp,
	}
	s.view.Width = panel.ContentWidth(ctx.Theme)
	s.view.Height = panel.ContentHeight(ctx.Theme)
	s.view.SetContent(s.content(ctx.Theme, s.view.Width))
	panel.Body = s.view.View()
	return ui.RenderPanel(ctx.Theme, panel)
}

func (s *helpScreen) content(t theme.Theme, width int) string {
	heading := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Background)
	keyStyle := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)
	descStyle := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	rows := []string{heading.Render("global"), ui.Rule(t, width, t.Border.Subtle)}
	rows = append(rows, bindingRows(s.bindings.Global(), keyStyle, descStyle)...)
	for _, set := range s.bindings.Sets() {
		rows = append(rows, "", heading.Render(set.Label+" "+t.Marker.Bullet+" "+set.PanelID), ui.Rule(t, width, t.Border.Subtle))
		rows = append(rows, bindingRows(set.Bindings, keyStyle, descStyle)...)
	}
	return ui.Join(rows)
}

func bindingRows(bindings []key.Binding, keyStyle, descStyle lipgloss.Style) []string {
	rows := make([]string, 0, len(bindings))
	for _, b := range bindings {
		rows = append(rows, keyStyle.Render(ui.Fit(b.Help().Key, 10))+descStyle.Render(b.Help().Desc))
	}
	return rows
}
