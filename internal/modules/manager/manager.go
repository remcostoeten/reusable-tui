// Package manager is a demo module showing a master-detail screen with a
// destructive, confirmed command and a drill-in route.
package manager

import (
	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/command"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/render"
	"github.com/remcostoeten/reusable-tui/internal/tui/focus"
	"github.com/remcostoeten/reusable-tui/internal/tui/input"
	"github.com/remcostoeten/reusable-tui/internal/tui/kernel"
	"github.com/remcostoeten/reusable-tui/internal/tui/layout"
	navigation "github.com/remcostoeten/reusable-tui/internal/tui/navigation"
	"github.com/remcostoeten/reusable-tui/internal/tui/shell/overlay"
	"github.com/remcostoeten/reusable-tui/internal/tui/ui"
)

const (
	regionGroups kernel.ID = "manager.groups"
	regionItems  kernel.ID = "manager.items"
)

// Store is the module's data source. A real module would hold a client here;
// the point is that the shell never sees it.
type Store struct {
	groups []string
	items  map[string][][]string
}

// NewStore builds the demo data.
func NewStore() *Store {
	return &Store{
		groups: []string{"Running", "Pending", "Failed"},
		items: map[string][][]string{
			"Running": {
				{"api-server-7d9f", "Running", "4d"},
				{"worker-queue-2b1a", "Running", "17m"},
				{"cache-warmer-91c", "Running", "2h"},
			},
			"Pending": {
				{"migration-4f2", "Pending", "31s"},
			},
			"Failed": {
				{"scheduler-0", "CrashLoopBackOff", "2h"},
				{"exporter-3a", "Error", "9m"},
			},
		},
	}
}

// Groups lists the categories.
func (s *Store) Groups() []string {
	return s.groups
}

// Items lists the rows in a category.
func (s *Store) Items(group string) [][]string {
	return s.items[group]
}

// Module is the manager screen.
type Module struct {
	store *Store
}

// New builds the module around its store, which is a constructor argument
// rather than something resolved from a container.
func New(store *Store) *Module {
	return &Module{store: store}
}

// ID names the module.
func (m *Module) ID() registry.ModuleID {
	return "manager"
}

// Register installs the route, a drill-in route, a dangerous command and the
// key that reaches it.
func (m *Module) Register(r *registry.Registrar) {
	r.Route(navigation.Route{ID: "manager", Title: "Manager", Order: 20},
		func(registry.Context) registry.View { return newView(m.store) })

	r.Route(navigation.Route{ID: "manager.detail", Title: "Detail", Order: 21, Hidden: true, Ephemeral: true},
		func(ctx registry.Context) registry.View { return newDetail(ctx.Route.Params.Get("item")) })

	r.Command(command.Command{
		ID: "manager.open", Title: "Open Selected Item", Category: "Manager",
		Description: "Drill into the highlighted row",
		When:        command.OnRoute("manager"),
		Run: func(s command.Scope) tea.Cmd {
			return navigation.Push("manager.detail", navigation.Params{"item": s.Region})
		},
	})

	r.Command(command.Command{
		ID: "manager.delete", Title: "Delete Selected Item", Category: "Manager",
		Description: "Remove the highlighted row",
		When:        command.OnRoute("manager"),
		Dangerous:   true,
		Run: func(command.Scope) tea.Cmd {
			return overlay.Ask("Delete item", "This cannot be undone.", "manager.delete.confirmed", true)
		},
	})

	r.Command(command.Command{
		ID: "manager.delete.confirmed", Title: "Delete (confirmed)", Category: "Manager",
		When: command.OnRoute("manager"),
		Run: func(command.Scope) tea.Cmd {
			return registry.Notify(ui.Toast{Severity: ui.SeveritySuccess, Text: "Item deleted"})
		},
	})

	r.Bind(input.LayerScreen, input.Binding{
		Keys: []string{"enter"}, Command: "manager.open",
		When: command.OnRoute("manager"), Hint: "Open", Priority: 60,
	})
	r.Bind(input.LayerScreen, input.Binding{
		Keys: []string{"d"}, Command: "manager.delete",
		When: command.OnRoute("manager"), Hint: "Delete", Priority: 55,
	})

	r.StatusItem(registry.StatusItem{
		Slot: registry.SlotRight, Order: 10,
		Render: func(rc render.Context) string {
			return rc.Styles().Success.Render(rc.Glyphs.Success) + " " + rc.Styles().StatusLabel.Render("ready")
		},
	})
}

// View is the master-detail screen.
type View struct {
	store  *Store
	groups ui.List
	items  ui.Table
	rects  map[kernel.ID]kernel.Rect
}

func newView(store *Store) *View {
	items := make([]ui.ListItem, len(store.Groups()))
	for i, g := range store.Groups() {
		items[i] = ui.ListItem{ID: kernel.ID(g), Title: g}
	}

	columns := []ui.Column{
		{Title: "NAME", Width: layout.Flex(2)},
		{Title: "STATUS", Width: layout.Fixed(18)},
		{Title: "AGE", Width: layout.Fixed(5), Align: render.Right},
	}

	v := &View{
		store:  store,
		groups: ui.NewList(items...),
		items:  ui.NewTable(columns, nil),
		rects:  map[kernel.ID]kernel.Rect{},
	}
	return v.sync()
}

// sync re-reads the store for whichever group is selected.
func (v *View) sync() *View {
	group, ok := v.groups.Selected()
	if !ok {
		return v
	}
	v.items = v.items.SetRows(v.store.Items(group.Title))
	return v
}

// Init has nothing to load: the store is already in memory.
func (v *View) Init() tea.Cmd {
	return nil
}

// FocusRegions declares the two panels.
func (v *View) FocusRegions() []focus.Region {
	return []focus.Region{
		{ID: regionGroups, Order: 1, Label: "Groups", JumpKey: 'g'},
		{ID: regionItems, Order: 2, Label: "Items", JumpKey: 'i', Skip: len(v.items.Rows) == 0},
	}
}

// RegionRects reports where the panels were drawn.
func (v *View) RegionRects(kernel.Rect) map[kernel.ID]kernel.Rect {
	return v.rects
}

// Update forwards to the focused panel and re-reads the store when the group
// selection moves.
func (v *View) Update(ctx registry.Context, msg tea.Msg) (registry.View, tea.Cmd) {
	next := *v
	next.groups.Focused = ctx.Focus.Focused(regionGroups)
	next.items.Focused = ctx.Focus.Focused(regionItems)

	before := next.groups.Cursor()
	next.groups, _ = next.groups.Update(msg)
	next.items, _ = next.items.Update(msg)

	if next.groups.Cursor() != before {
		return next.sync(), nil
	}
	return &next, nil
}

// Render lays out the two panels.
func (v *View) Render(rc render.Context) string {
	area := rc.Rect
	cols := layout.Cols(area, layout.Fixed(20), layout.Flex(1))

	v.rects[regionGroups] = cols[0]
	v.rects[regionItems] = cols[1]

	groups := v.groups.SetHeight(inner(cols[0]).Height)
	items := v.items.SetHeight(inner(cols[1]).Height)

	return render.Compose(area,
		render.At(cols[0], ui.Panel{
			Title:   "Groups",
			Focused: rc.Focused(regionGroups),
			Scroll:  groups.ScrollPos(),
			Content: groups.Render(rc.For(inner(cols[0]))),
		}.Render(rc.For(cols[0]))),
		render.At(cols[1], ui.Panel{
			Title:    "Items",
			Subtitle: count(len(v.items.Rows)),
			Focused:  rc.Focused(regionItems),
			Footer:   []ui.Hint{{Key: "enter", Label: "Open"}, {Key: "d", Label: "Delete"}},
			Scroll:   items.ScrollPos(),
			Content:  items.Render(rc.For(inner(cols[1]))),
		}.Render(rc.For(cols[1]))),
	)
}

func inner(r kernel.Rect) kernel.Rect {
	return r.Inset(1).InsetXY(1, 0)
}
