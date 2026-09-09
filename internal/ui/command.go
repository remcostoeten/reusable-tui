package ui

import tea "github.com/charmbracelet/bubbletea"

type Command struct {
	ID    string
	Label string
	Group string
	Run   tea.Cmd
}

type CommandRegistry struct {
	order []string
	byID  map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{byID: make(map[string]Command)}
}

func (r *CommandRegistry) Add(c Command) {
	if _, exists := r.byID[c.ID]; !exists {
		r.order = append(r.order, c.ID)
	}
	r.byID[c.ID] = c
}

func (r *CommandRegistry) Get(id string) (Command, bool) {
	c, ok := r.byID[id]
	return c, ok
}

func (r *CommandRegistry) All() []Command {
	out := make([]Command, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.byID[id])
	}
	return out
}

type commandSource []Command

func (s commandSource) String(i int) string {
	return s[i].Group + " " + s[i].Label
}

func (s commandSource) Len() int {
	return len(s)
}
