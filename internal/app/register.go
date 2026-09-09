package app

import (
	"github.com/remcostoeten/reusable-tui/internal/modules/home"
	"github.com/remcostoeten/reusable-tui/internal/modules/manager"
	"github.com/remcostoeten/reusable-tui/internal/modules/settings"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/config"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
)

// Deps are the module services, wired here rather than resolved from a
// container: constructor arguments, checked by the compiler.
type Deps struct {
	Config config.Config
	Store  *manager.Store
}

// Modules is the whole module list. Adding one is a line here plus a
// directory — the honest version of "no giant central file".
func Modules(deps Deps) []registry.Module {
	return []registry.Module{
		home.New(),
		manager.New(deps.Store),
		settings.New(deps.Config),
	}
}
