package home

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/render"
	"github.com/remcostoeten/reusable-tui/internal/tui/focus"
	"github.com/remcostoeten/reusable-tui/internal/tui/kernel"
	"github.com/remcostoeten/reusable-tui/internal/tui/layout"
	"github.com/remcostoeten/reusable-tui/internal/tui/ui"
)

const (
	regionSources  kernel.ID = "home.sources"
	regionMetrics  kernel.ID = "home.metrics"
	regionOverview kernel.ID = "home.overview"
)

// View is the home screen: three panels exercising focus, tabs, an empty
// state, a sparkline and a scrolling body.
type View struct {
	sources  ui.List
	overview ui.ScrollArea
	tabs     int
	rects    map[kernel.ID]kernel.Rect
}

func newView() *View {
	return &View{
		sources:  ui.NewList(),
		overview: ui.NewScrollArea(overviewText),
		rects:    map[kernel.ID]kernel.Rect{},
	}
}

// Init has nothing to load.
func (v *View) Init() tea.Cmd {
	return nil
}

// FocusRegions declares the ring. Sources is skipped while it is empty, so tab
// does not stop on a panel with nothing to select.
func (v *View) FocusRegions() []focus.Region {
	return []focus.Region{
		{ID: regionSources, Order: 1, Label: "Sources", JumpKey: 's', Skip: len(v.sources.Items) == 0},
		{ID: regionMetrics, Order: 2, Label: "Metrics", JumpKey: 'm'},
		{ID: regionOverview, Order: 3, Label: "Overview", JumpKey: 'o'},
	}
}

// RegionRects reports where the panels were drawn, which is what makes the
// regions clickable and reachable with alt+hjkl.
func (v *View) RegionRects(kernel.Rect) map[kernel.ID]kernel.Rect {
	return v.rects
}

// Update forwards to whichever panel holds focus.
func (v *View) Update(ctx registry.Context, msg tea.Msg) (registry.View, tea.Cmd) {
	next := *v

	if key, ok := msg.(tea.KeyMsg); ok && key.String() == "c" {
		next.tabs = (next.tabs + 1) % len(tabTitles)
		return &next, nil
	}

	next.sources.Focused = ctx.Focus.Focused(regionSources)
	next.overview.Focused = ctx.Focus.Focused(regionOverview)

	var cmd tea.Cmd
	next.sources, _ = next.sources.Update(msg)
	next.overview, cmd = next.overview.Update(msg)
	return &next, cmd
}

// Render lays the screen out and records where each region landed.
func (v *View) Render(rc render.Context) string {
	area := rc.Rect
	cols := layout.Cols(area, layout.Fixed(24), layout.Flex(1))
	left := layout.Rows(cols[0], layout.Percent(45), layout.Flex(1))

	v.rects[regionSources] = left[0]
	v.rects[regionMetrics] = left[1]
	v.rects[regionOverview] = cols[1]

	sources := ui.Panel{
		Title:   "Sources",
		Focused: rc.Focused(regionSources),
		Content: ui.EmptyState{Title: "No sources yet", Hint: "nothing to load"}.Render(rc.For(inner(left[0]))),
	}

	metrics := ui.Panel{
		Title:    "Metrics",
		Subtitle: "0.0/s",
		Focused:  rc.Focused(regionMetrics),
		Content:  v.metrics(rc, inner(left[1])),
	}

	// The prose is wrapped to whatever width the panel ended up with, rather
	// than being stored pre-formatted for one terminal size.
	overview := v.overview.SetContent(render.Wrap(overviewText, inner(cols[1]).Width)).
		SetHeight(inner(cols[1]).Height - 2)
	body := render.Join(
		ui.Tabs{Items: tabTitles, Active: v.tabs}.Render(rc.For(inner(cols[1]))),
		"",
		overview.Render(rc.For(inner(cols[1]).InsetXY(0, 1))),
	)

	return render.Compose(area,
		render.At(left[0], sources.Render(rc.For(left[0]))),
		render.At(left[1], metrics.Render(rc.For(left[1]))),
		render.At(cols[1], ui.Panel{
			Title:   "Overview",
			Focused: rc.Focused(regionOverview),
			Footer:  []ui.Hint{{Key: "c", Label: "Cycle tabs"}},
			Scroll:  overview.ScrollPos(),
			Content: body,
		}.Render(rc.For(cols[1]))),
	)
}

// metrics draws a label over a sparkline.
func (v *View) metrics(rc render.Context, area kernel.Rect) string {
	rows := layout.Rows(area, layout.Fixed(1), layout.Fixed(1), layout.Flex(1))
	spark := ui.Sparkline{Values: sample}.Render(rc.For(rows[2]))

	return render.Compose(area,
		render.At(rows[0], rc.Styles().Muted.Render("Throughput")),
		render.At(rows[1], rc.Styles().Base.Render("0")),
		render.At(rows[2], strings.Repeat(spark+"\n", max(1, min(3, rows[2].Height)))),
	)
}

// inner is a panel's content area: inside the border and the density padding.
func inner(r kernel.Rect) kernel.Rect {
	return r.Inset(1).InsetXY(1, 0)
}

var tabTitles = []string{"Active", "Archived", "All"}

var sample = []float64{1, 4, 2, 8, 5, 3, 9, 6, 2, 7, 4, 1, 5, 8, 3, 6}

const overviewText = `Welcome to the shell.

An application shell that lives in your terminal.

Focus
  The focused region has a bright, heavier border and a bold title, so
  the state survives a monochrome terminal.
  tab cycles the ring, v badges every region and jumps to the one you
  press, alt+hjkl moves geometrically.

Commands
  ctrl+k opens the palette. Everything reachable by a key is reachable
  there too, which is what makes the keyboard guarantee testable.
  ? lists every command with the key that reaches it, generated from
  the registry rather than written by hand.

Themes
  ctrl+t cycles dark, dim, high-contrast and ascii. The last one has no
  colour at all: if the app is usable there, nothing is being said with
  colour alone.

Layout
  Narrow the terminal. The sidebar collapses to initials at 120 columns
  and disappears below 80; under 40x10 the frame is replaced by a single
  legible message rather than a garbled one.`
