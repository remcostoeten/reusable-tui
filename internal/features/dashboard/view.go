package dashboard

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const (
	calendarMinWidth = 21
	statRows         = 3
)

func (m *Model) View(ctx ui.RenderContext) string {
	return ui.RenderDashboard(ctx.Theme, ui.Dashboard{
		Width:  ctx.Width,
		Height: ctx.Height,
		Stacks: []ui.Stack{
			{Weight: 5, Regions: []ui.Region{m.accountsRegion(ctx), m.insightsRegion(ctx)}},
			{Weight: 6, Regions: []ui.Region{m.modeRegion(ctx), m.periodRegion(ctx)}},
			{Weight: 9, Regions: []ui.Region{m.overviewRegion(ctx)}},
		},
	})
}

func (m *Model) accountsRegion(ctx ui.RenderContext) ui.Region {
	return ui.Region{
		Weight: 2,
		Panel: ui.Panel{
			Title:   "Accounts",
			Badge:   money(m.balance()),
			Hint:    ctx.Jump.Hint(keymap.PanelDashboardAccounts),
			Focused: ctx.Focused == keymap.PanelDashboardAccounts,
		},
		Body: m.accountsBody,
	}
}

func (m *Model) insightsRegion(ctx ui.RenderContext) ui.Region {
	return ui.Region{
		Weight: 3,
		Panel: ui.Panel{
			Title:   "Insights",
			Badge:   m.periodLabel(),
			Hint:    ctx.Jump.Hint(keymap.PanelDashboardInsights),
			Focused: ctx.Focused == keymap.PanelDashboardInsights,
		},
		Body: m.insightsBody,
	}
}

func (m *Model) modeRegion(ctx ui.RenderContext) ui.Region {
	focused := ctx.Focused == keymap.PanelDashboardMode
	return ui.Region{
		Weight: 2,
		Panel: ui.Panel{
			Title:   "View and add",
			Hint:    ctx.Jump.Hint(keymap.PanelDashboardMode),
			Focused: focused,
		},
		Body: m.modeBody(focused),
	}
}

func (m *Model) periodRegion(ctx ui.RenderContext) ui.Region {
	focused := ctx.Focused == keymap.PanelDashboardPeriod
	return ui.Region{
		Weight: 3,
		Panel: ui.Panel{
			Title:   "Period",
			Hint:    ctx.Jump.Hint(keymap.PanelDashboardPeriod),
			Focused: focused,
		},
		Body: m.periodBody(focused),
	}
}

func (m *Model) overviewRegion(ctx ui.RenderContext) ui.Region {
	return ui.Region{
		Panel: ui.Panel{
			Title:   "Overview",
			Hint:    ctx.Jump.Hint(keymap.PanelDashboardOverview),
			Focused: ctx.Focused == keymap.PanelDashboardOverview,
		},
		Body: m.overviewBody,
	}
}

func (m *Model) accountsBody(t theme.Theme, width, height int) string {
	if len(m.accounts) == 0 {
		return ui.EmptyState(t, width, height, "no accounts yet")
	}
	rows := make([]string, 0, height)
	for i, account := range m.accounts {
		if i >= height {
			break
		}
		rows = append(rows, accountRow(t, account, i == m.cursor, width))
	}
	return ui.Join(rows)
}

func accountRow(t theme.Theme, a Account, selected bool, width int) string {
	marker := t.Marker.UnselectedPad
	style := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(t.Base.Background)
	if selected {
		marker = t.Marker.SelectedLeft
		style = lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Accent.Dim)
	}
	value := money(a.Balance)
	name := ui.Fit(a.Name, width-ui.Width(marker)-1-ui.Width(value))
	return style.Render(marker + " " + name + value)
}

func (m *Model) insightsBody(t theme.Theme, width, height int) string {
	if height < statRows {
		return ui.EmptyState(t, width, height, "no room for insights")
	}
	stats := ui.RenderStats(t, width,
		ui.Stat{Label: m.mode.String(), Value: money(m.total())},
		ui.Stat{Label: "per day", Value: money(m.perDay())},
	)
	rows := ui.Lines(stats)
	rows = append(rows, ui.Rule(t, width, t.Border.Subtle))
	chart := ui.RenderBarChart(t, ui.BarChart{
		Bars:   m.series(),
		Width:  width,
		Height: height - len(rows),
		Labels: true,
	})
	if chart == "" {
		return ui.Join(rows)
	}
	return ui.Join(append(rows, ui.Lines(chart)...))
}

func (m *Model) modeBody(focused bool) ui.BodyFunc {
	return func(t theme.Theme, width, height int) string {
		rows := []string{
			ui.RenderSegment(t, ui.Segment{
				Options: []string{ModeExpense.String(), ModeIncome.String()},
				Active:  int(m.mode),
				Width:   width,
				Focused: focused,
			}),
		}
		if height > 1 {
			hint := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
			rows = append(rows, hint.Render(ui.Fit("h/l switches the view", width)))
		}
		return ui.Join(rows)
	}
}

func (m *Model) periodBody(focused bool) ui.BodyFunc {
	return func(t theme.Theme, width, height int) string {
		rows := []string{ui.RenderPeriod(t, width, m.periodLabel(), focused)}
		if width < calendarMinWidth || height < 3 {
			return ui.Join(rows)
		}
		calendar := ui.RenderCalendar(t, ui.Calendar{
			Month:     m.selected,
			Selected:  m.selected,
			Today:     m.clock(),
			Width:     width,
			Focused:   focused,
			WeekBand:  true,
			MarkToday: true,
		})
		return ui.Join(append(rows, ui.Lines(calendar)...))
	}
}

func (m *Model) overviewBody(t theme.Theme, width, height int) string {
	m.overview.Width = width
	m.overview.Height = height
	m.overview.SetContent(ui.RenderDoc(t, width, m.doc()...))
	return m.overview.View()
}

func (m *Model) doc() []ui.DocBlock {
	return []ui.DocBlock{
		ui.Heading("The dashboard layout"),
		ui.Spacer(),
		ui.Text("Three stacked columns share one body area. Every panel keeps its own focus state, jump hint and keymap set."),
		ui.Spacer(),
		ui.Divider(),
		ui.Spacer(),
		ui.Heading("Layout"),
		ui.Bullet("Stacks split the width by weight, regions split each column by height."),
		ui.Bullet("Region bodies receive their measured content size, so widgets never guess."),
		ui.Bullet("Swap the stack list to get a sidebar, a split view or a single pane."),
		ui.Spacer(),
		ui.Heading("Widgets"),
		ui.Bullet("Segment control, calendar, stat row, bar chart and this document renderer."),
		ui.Bullet("All of them read the active theme, so a theme switch repaints everything."),
		ui.Spacer(),
		ui.Callout("Hint: press f for jump mode, then the letter shown on a panel border."),
	}
}

func money(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
