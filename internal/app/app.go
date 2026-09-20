// Package app is the composition root. It builds the registrar, lets every
// module write into it, freezes it and hands the result to the runtime.
package app

import (
	"github.com/remcostoeten/reusable-tui/internal/tui/core/config"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/registry"
	"github.com/remcostoeten/reusable-tui/internal/tui/core/runtime"
	"github.com/remcostoeten/reusable-tui/internal/tui/kernel"
	"github.com/remcostoeten/reusable-tui/internal/tui/shell"
	"github.com/remcostoeten/reusable-tui/internal/tui/theme"
	"github.com/remcostoeten/reusable-tui/internal/tui/theme/themes"
)

// Options configures a build.
type Options struct {
	Name    string
	Version string
	Config  config.Config
	Log     kernel.Logger
}

// Build assembles the application. Every failure it can detect — a duplicate
// route, a malformed binding, a keymap conflict — is reported here, before a
// single frame is drawn.
func Build(o Options) (runtime.Model, error) {
	if o.Log == nil {
		o.Log = kernel.NopLogger()
	}

	r := registry.NewRegistrar()
	registerGlobals(r)

	deps := Deps{Config: o.Config}
	for _, m := range Modules(deps) {
		m.Register(r)
	}

	reg, err := r.Freeze()
	if err != nil {
		return runtime.Model{}, err
	}

	palettes := theme.NewRegistry(themes.All()...)
	for _, t := range reg.Themes() {
		palettes.Add(t)
	}

	return runtime.New(runtime.Options{
		Registry: reg,
		Themes:   palettes,
		Chrome: shell.Chrome{
			Name:     o.Name,
			Version:  o.Version,
			Subtitle: o.Config.Path(),
		},
		Config: o.Config,
		Log:    o.Log,
	})
}
