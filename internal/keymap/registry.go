package keymap

import "github.com/charmbracelet/bubbles/key"

type Set struct {
	PanelID  string
	Label    string
	Bindings []key.Binding
}

type Registry struct {
	global []key.Binding
	sets   map[string]Set
	order  []string
}

func NewRegistry() *Registry {
	return &Registry{sets: make(map[string]Set)}
}

func (r *Registry) RegisterGlobal(bindings ...key.Binding) {
	r.global = append(r.global, bindings...)
}

func (r *Registry) Register(panelID, label string, bindings ...key.Binding) {
	if _, exists := r.sets[panelID]; !exists {
		r.order = append(r.order, panelID)
	}
	r.sets[panelID] = Set{PanelID: panelID, Label: label, Bindings: bindings}
}

func (r *Registry) Global() []key.Binding {
	return r.global
}

func (r *Registry) For(panelID string) []key.Binding {
	return r.sets[panelID].Bindings
}

func (r *Registry) Set(panelID string) (Set, bool) {
	s, ok := r.sets[panelID]
	return s, ok
}

func (r *Registry) Sets() []Set {
	out := make([]Set, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.sets[id])
	}
	return out
}
