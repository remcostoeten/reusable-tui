package example

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func (m *Model) View(ctx ui.RenderContext) string {
	t := ctx.Theme
	widths := ui.SplitWidths(ctx.Width, t.Space.Gutter, 3, 2)
	list := ui.Panel{
		Title:   "Items",
		Badge:   strconv.Itoa(len(m.items)),
		Hint:    ctx.Jump.Hint(keymap.PanelExampleList),
		Width:   widths[0],
		Height:  ctx.Height,
		Focused: ctx.Focused == keymap.PanelExampleList,
	}
	list.Body = m.listBody(t, list.ContentWidth(t), list.ContentHeight(t), ctx.Tick)

	detail := ui.Panel{
		Title:   "Detail",
		Hint:    ctx.Jump.Hint(keymap.PanelExampleDetail),
		Width:   widths[1],
		Height:  ctx.Height,
		Focused: ctx.Focused == keymap.PanelExampleDetail,
	}
	detail.Body = m.detailBody(t, detail.ContentWidth(t), detail.ContentHeight(t))

	return ui.Row(t, ui.RenderPanel(t, list), ui.RenderPanel(t, detail))
}

func (m *Model) listBody(t theme.Theme, width, height, tick int) string {
	if m.loading {
		return ui.Skeleton(t, width, height, tick)
	}
	rows := make([]string, 0, height)
	if m.adding {
		prompt := lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)
		m.input.Width = width - ui.Width(m.input.Prompt) - 1
		rows = append(rows, prompt.Render(ui.Fit(m.input.View(), width)))
		height--
	}
	if len(m.items) == 0 {
		rows = append(rows, ui.Lines(ui.EmptyState(t, width, height, "no items yet, press a to add"))...)
		return ui.Join(rows)
	}
	for i, item := range m.items {
		if i >= height {
			break
		}
		rows = append(rows, m.itemRow(t, item, i == m.cursor, width))
	}
	return ui.Join(rows)
}

func (m *Model) itemRow(t theme.Theme, item Item, selected bool, width int) string {
	token := t.Status.Info
	if item.Done {
		token = t.Status.Success
	}
	marker := t.Marker.UnselectedPad
	style := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	if selected {
		marker = t.Marker.SelectedLeft
		style = lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Accent.Dim)
	}
	badge := lipgloss.NewStyle().Foreground(token.Fg).Background(style.GetBackground()).Render(token.Glyph)
	label := ui.Fit(item.Title, width-ui.Width(marker+" "+token.Glyph+" "))
	return style.Render(marker+" ") + badge + style.Render(" "+label)
}

func (m *Model) detailBody(t theme.Theme, width, height int) string {
	item, ok := m.Selected()
	if !ok {
		return ui.EmptyState(t, width, height, "select an item")
	}
	token := t.Status.Info
	state := "pending"
	if item.Done {
		token = t.Status.Success
		state = "done"
	}
	label := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	value := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Background)
	status := lipgloss.NewStyle().Foreground(token.Fg).Background(t.Base.Background)
	rows := []string{
		value.Render(ui.Fit(item.Title, width)),
		ui.Rule(t, width, t.Border.Subtle),
		label.Render("id      ") + value.Render(strconv.FormatInt(item.ID, 10)),
		label.Render("state   ") + status.Render(token.Glyph+" "+state),
		label.Render("created ") + value.Render(item.CreatedAt),
	}
	return ui.Join(rows)
}
