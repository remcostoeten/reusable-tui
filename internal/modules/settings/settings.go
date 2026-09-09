// Package settings is a demo module showing a form built by hand — three
// fields do not justify a form framework.
package settings

import (
	tea "charm.land/bubbletea/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/command"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/config"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/render"
	"github.com/remcostoeten/reusable-tui/internal/tui/focus"
	"github.com/remcostoeten/reusable-tui/internal/tui/kernel"
	"github.com/remcostoeten/reusable-tui/internal/tui/layout"
	navigation "github.com/remcostoeten/reusable-tui/internal/tui/navigation"
	"github.com/remcostoeten/reusable-tui/internal/tui/ui"
)

const regionForm kernel.ID = "settings.form"

// Module is the settings screen.
type Module struct {
	cfg config.Config
}

// New builds the module around the loaded configuration.
func New(cfg config.Config) *Module {
	return &Module{cfg: cfg}
}

// ID names the module.
func (m *Module) ID() registry.ModuleID {
	return "settings"
}

// Register installs the route and the command that applies the form.
func (m *Module) Register(r *registry.Registrar) {
	r.Route(navigation.Route{ID: "settings", Title: "Settings", Order: 30},
		func(registry.Context) registry.View { return newView(m.cfg) })

	r.Command(command.Command{
		ID: "settings.apply", Title: "Apply Settings", Category: "Settings",
		When: command.OnRoute("settings"),
		Run: func(command.Scope) tea.Cmd {
			return registry.Notify(ui.Toast{Severity: ui.SeveritySuccess, Text: "Settings applied"})
		},
	})
}

// View is the settings form.
type View struct {
	theme      ui.Select
	mouse      ui.Checkbox
	follow     ui.Checkbox
	cursor     int
	configPath string
	rects      map[kernel.ID]kernel.Rect
}

func newView(cfg config.Config) *View {
	path := cfg.Path()
	if path == "" {
		path = "no config file — using defaults"
	}
	return &View{
		theme:      ui.NewSelect("Theme", "dark", "dim", "high-contrast", "ascii"),
		mouse:      ui.Checkbox{Label: "Mouse support", Checked: cfg.Mouse},
		follow:     ui.Checkbox{Label: "Follow new log output"},
		configPath: path,
		rects:      map[kernel.ID]kernel.Rect{},
	}
}

// Init has nothing to load.
func (v *View) Init() tea.Cmd {
	return nil
}

// FocusRegions declares one region: the form owns its own internal cursor.
func (v *View) FocusRegions() []focus.Region {
	return []focus.Region{{ID: regionForm, Order: 1, Label: "Settings", JumpKey: 's'}}
}

// RegionRects reports where the form was drawn.
func (v *View) RegionRects(kernel.Rect) map[kernel.ID]kernel.Rect {
	return v.rects
}

// Update moves between fields and edits the focused one. Switching the theme
// select applies it immediately, because a theme is a value.
func (v *View) Update(ctx registry.Context, msg tea.Msg) (registry.View, tea.Cmd) {
	next := *v
	if !ctx.Focus.Focused(regionForm) {
		return &next, nil
	}

	if key, ok := msg.(tea.KeyMsg); ok && !next.theme.IsOpen() {
		switch key.String() {
		case "up", "k":
			next.cursor = max(0, next.cursor-1)
			return next.focusField(), nil
		case "down", "j":
			next.cursor = min(2, next.cursor+1)
			return next.focusField(), nil
		case "space", "enter":
			switch next.cursor {
			case 1:
				next.mouse.Checked = !next.mouse.Checked
				return next.focusField(), nil
			case 2:
				next.follow.Checked = !next.follow.Checked
				return next.focusField(), nil
			}
		}
	}

	before := next.theme.Value()
	next = *next.focusField()
	next.theme, _ = next.theme.Update(msg)
	if next.theme.Value() != before {
		return &next, registry.SetTheme(next.theme.Value())
	}
	return &next, nil
}

// focusField points the widgets at whichever row the cursor is on.
func (v View) focusField() *View {
	v.theme.Focused = v.cursor == 0
	v.mouse.Focused = v.cursor == 1
	v.follow.Focused = v.cursor == 2
	return &v
}

// Render draws the three fields and the config path.
func (v *View) Render(rc render.Context) string {
	area := rc.Rect
	v.rects[regionForm] = area

	body := inner(area)
	rows := layout.Rows(body,
		layout.Fixed(v.theme.Height()),
		layout.Fixed(1),
		layout.Fixed(1),
		layout.Fixed(1),
		layout.Flex(1),
	)

	content := render.Compose(body,
		render.At(rows[0], v.theme.Render(rc.For(rows[0]))),
		render.At(rows[1], v.mouse.Render(rc.For(rows[1]))),
		render.At(rows[2], v.follow.Render(rc.For(rows[2]))),
		render.At(rows[4], rc.Styles().Subtle.Render(v.configPath)),
	)

	return ui.Panel{
		Title:   "Settings",
		Focused: rc.Focused(regionForm),
		Footer:  []ui.Hint{{Key: "space", Label: "Toggle"}, {Key: "enter", Label: "Open"}},
		Content: content,
	}.Render(rc)
}

func inner(r kernel.Rect) kernel.Rect {
	return r.Inset(1).InsetXY(1, 0)
}
