package render

import (
	lipgloss "charm.land/lipgloss/v2"

	"github.com/remcostoeten/reusable-tui/internal/tui/kernel"
)

// Block is rendered content together with the rectangle it occupies. Widgets
// return plain strings; a Block is how a caller states where one goes.
type Block struct {
	Rect    kernel.Rect
	Content string
}

// At pairs content with its rectangle.
func At(r kernel.Rect, content string) Block {
	return Block{Rect: r, Content: content}
}

// Compose draws blocks onto a blank area in the order given, later blocks
// covering earlier ones. Each block is clipped to its own rectangle first, so
// a widget that overruns cannot corrupt its neighbours.
//
// Trailing whitespace is trimmed from every output row. That keeps golden
// files diffable and costs nothing: a terminal row ends where its last glyph is.
func Compose(area kernel.Rect, blocks ...Block) string {
	if area.IsEmpty() {
		return ""
	}

	layers := make([]*lipgloss.Layer, 0, len(blocks)+1)
	layers = append(layers, lipgloss.NewLayer(Blank(area.Size())))
	for i, b := range blocks {
		if b.Rect.IsEmpty() {
			continue
		}
		layers = append(layers, lipgloss.NewLayer(Clip(b.Content, b.Rect.Size())).
			X(b.Rect.X-area.X).
			Y(b.Rect.Y-area.Y).
			Z(i+1))
	}
	return lipgloss.NewCompositor(layers...).Render()
}

// Overlay composites over on top of base at the given position. Cell accuracy
// across escape sequences and wide runes is lipgloss's problem, not ours.
func Overlay(base, over string, at kernel.Rect) string {
	if over == "" || at.IsEmpty() {
		return base
	}
	return lipgloss.NewCompositor(
		lipgloss.NewLayer(base),
		lipgloss.NewLayer(Clip(over, at.Size())).X(at.X).Y(at.Y).Z(1),
	).Render()
}
