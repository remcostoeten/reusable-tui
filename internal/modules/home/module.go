// Package home is a demo module. Delete it and the shell still runs.
package home

import (
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	navigation "github.com/remcostoeten/reusable-tui/internal/tui/navigation"
)

// Module is the home screen.
type Module struct{}

// New builds the module.
func New() *Module {
	return &Module{}
}

// ID names the module.
func (m *Module) ID() registry.ModuleID {
	return "home"
}

// Register installs the route. Home has no verbs of its own: everything it
// demonstrates is a shell feature.
func (m *Module) Register(r *registry.Registrar) {
	r.Route(navigation.Route{
		ID:    "home",
		Title: "Home",
		Order: 10,
	}, func(registry.Context) registry.View {
		return newView()
	})
}
