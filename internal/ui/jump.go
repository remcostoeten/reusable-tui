package ui

const jumpLetters = "asdfghjkl"

type Jump struct {
	active bool
	hints  map[string]string
	target map[string]string
}

func NewJump() Jump {
	return Jump{hints: map[string]string{}, target: map[string]string{}}
}

func (j Jump) Active() bool {
	return j.active
}

func (j *Jump) Show(panelIDs []string) {
	j.active = true
	j.hints = make(map[string]string, len(panelIDs))
	j.target = make(map[string]string, len(panelIDs))
	for i, id := range panelIDs {
		if i >= len(jumpLetters) {
			return
		}
		letter := string(jumpLetters[i])
		j.hints[id] = letter
		j.target[letter] = id
	}
}

func (j *Jump) Hide() {
	j.active = false
	j.hints = map[string]string{}
	j.target = map[string]string{}
}

func (j Jump) Hint(panelID string) string {
	if !j.active {
		return ""
	}
	return j.hints[panelID]
}

func (j Jump) Target(letter string) (string, bool) {
	id, ok := j.target[letter]
	return id, ok
}
