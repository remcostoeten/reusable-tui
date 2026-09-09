package theme

type Registry struct {
	order  []string
	byName map[string]Theme
}

func NewRegistry(themes ...Theme) *Registry {
	r := &Registry{byName: make(map[string]Theme, len(themes))}
	for _, t := range themes {
		r.Add(t)
	}
	return r
}

func Builtin() *Registry {
	return NewRegistry(VioletDark(), Monochrome(), HighContrast())
}

func (r *Registry) Add(t Theme) {
	if _, exists := r.byName[t.Name]; !exists {
		r.order = append(r.order, t.Name)
	}
	r.byName[t.Name] = t
}

func (r *Registry) Get(name string) (Theme, bool) {
	t, ok := r.byName[name]
	return t, ok
}

func (r *Registry) Names() []string {
	out := make([]string, len(r.order))
	copy(out, r.order)
	return out
}

func (r *Registry) Default() Theme {
	if len(r.order) == 0 {
		return VioletDark()
	}
	return r.byName[r.order[0]]
}

func (r *Registry) Resolve(name string, f Fidelity) Theme {
	t, ok := r.Get(name)
	if !ok {
		t = r.Default()
	}
	return Adapt(t, f)
}
