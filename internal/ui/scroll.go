package ui

import "github.com/charmbracelet/bubbles/viewport"

func Scroll(v *viewport.Model, delta int) {
	if delta > 0 {
		v.ScrollDown(delta)
		return
	}
	v.ScrollUp(-delta)
}
