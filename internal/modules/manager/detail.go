package manager

import (
	"strconv"

	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/render"
	"github.com/remcostoeten/reusable-tui/internal/tui/ui"
)

// Detail is the drill-in screen. Pushing it is the same mechanism as switching
// sections, so esc unwinds it and focus restoration works at any depth.
type Detail struct {
	item string
}

func newDetail(item string) *Detail {
	return &Detail{item: item}
}

// Init has nothing to load.
func (d *Detail) Init() tea.Cmd {
	return nil
}

// Update ignores everything: esc is a global binding gated on drill-in depth.
func (d *Detail) Update(registry.Context, tea.Msg) (registry.View, tea.Cmd) {
	return d, nil
}

// Render draws a single panel.
func (d *Detail) Render(rc render.Context) string {
	name := d.item
	if name == "" {
		name = "the selected item"
	}
	return ui.Panel{
		Title:   "Detail",
		Focused: true,
		Footer:  []ui.Hint{{Key: "esc", Label: "Back"}},
		Content: ui.EmptyState{
			Title: "Drilled into " + name,
			Hint:  "esc unwinds one level",
		}.Render(rc.For(inner(rc.Rect))),
	}.Render(rc)
}

func count(n int) string {
	return strconv.Itoa(n)
}
