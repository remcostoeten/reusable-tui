package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/sahilm/fuzzy"
)

const paletteRows = 8

type Palette struct {
	commands *CommandRegistry
	input    textinput.Model
	matches  []Command
	cursor   int
	open     bool
}

func NewPalette(commands *CommandRegistry) Palette {
	in := textinput.New()
	in.Prompt = "> "
	in.Placeholder = "search commands"
	in.CharLimit = 64
	p := Palette{commands: commands, input: in}
	p.refresh()
	return p
}

func (p Palette) Open() bool {
	return p.open
}

func (p *Palette) Show() {
	p.open = true
	p.cursor = 0
	p.input.SetValue("")
	p.input.Focus()
	p.refresh()
}

func (p *Palette) Hide() {
	p.open = false
	p.input.Blur()
}

func (p *Palette) refresh() {
	all := p.commands.All()
	query := p.input.Value()
	if query == "" {
		p.matches = all
		return
	}
	found := fuzzy.FindFrom(query, commandSource(all))
	matches := make([]Command, 0, len(found))
	for _, m := range found {
		matches = append(matches, all[m.Index])
	}
	p.matches = matches
	if p.cursor >= len(p.matches) {
		p.cursor = 0
	}
}

func (p *Palette) Update(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}
	switch keyMsg.String() {
	case "esc":
		p.Hide()
		return nil
	case "up", "ctrl+p":
		p.moveCursor(-1)
		return nil
	case "down", "ctrl+n":
		p.moveCursor(1)
		return nil
	case "enter":
		return p.run()
	}
	in, cmd := p.input.Update(msg)
	p.input = in
	p.refresh()
	return cmd
}

func (p *Palette) moveCursor(delta int) {
	if len(p.matches) == 0 {
		return
	}
	p.cursor = (p.cursor + delta + len(p.matches)) % len(p.matches)
}

func (p *Palette) run() tea.Cmd {
	if p.cursor >= len(p.matches) {
		p.Hide()
		return nil
	}
	selected := p.matches[p.cursor]
	p.Hide()
	return selected.Run
}

func (p Palette) Matches() []Command {
	return p.matches
}

func (p Palette) Query() string {
	return p.input.Value()
}

func (p Palette) Modal(t theme.Theme, width int) Modal {
	return Modal{
		Title:  "Commands",
		Body:   p.body(t, width-2-2*t.Space.ModalPadX),
		Width:  width,
		Height: paletteRows + 4 + 2*t.Space.ModalPadY,
	}
}

func (p Palette) body(t theme.Theme, width int) string {
	prompt := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Overlay)
	p.input.Width = width - Width(p.input.Prompt) - 1
	rows := []string{prompt.Render(Fit(p.input.View(), width))}
	rows = append(rows, Rule(t, width, t.Border.Subtle))
	rows = append(rows, p.rows(t, width)...)
	return Join(rows)
}

func (p Palette) rows(t theme.Theme, width int) []string {
	active := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Accent.Dim)
	idle := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Overlay)
	group := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Overlay)
	if len(p.matches) == 0 {
		return []string{idle.Render(Fit("  no matching commands", width))}
	}
	rows := make([]string, 0, paletteRows)
	start := paletteWindow(p.cursor)
	for offset := 0; offset < paletteRows; offset++ {
		i := start + offset
		if i >= len(p.matches) {
			rows = append(rows, idle.Render(Repeat(" ", width)))
			continue
		}
		item := p.matches[i]
		label := " " + t.Marker.Bullet + " " + item.Label
		if i == p.cursor {
			label = " " + t.Marker.SelectedLeft + " " + item.Label
			rows = append(rows, active.Render(Fit(label, width)))
			continue
		}
		rows = append(rows, idle.Render(Fit(label, width-Width(item.Group)-1))+group.Render(" "+item.Group))
	}
	return rows
}

func paletteWindow(cursor int) int {
	if cursor < paletteRows {
		return 0
	}
	return cursor - paletteRows + 1
}
