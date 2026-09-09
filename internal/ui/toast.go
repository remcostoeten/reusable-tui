package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

const ToastLifetime = 4 * time.Second

type ToastKind int

const (
	ToastInfo ToastKind = iota
	ToastSuccess
	ToastWarning
	ToastError
)

type ToastExpiredMsg struct {
	Seq int
}

type Toast struct {
	Text    string
	Kind    ToastKind
	Seq     int
	Visible bool
}

func (t Toast) token(th theme.Theme) theme.StatusToken {
	switch t.Kind {
	case ToastSuccess:
		return th.Status.Success
	case ToastWarning:
		return th.Status.Warning
	case ToastError:
		return th.Status.Danger
	default:
		return th.Status.Info
	}
}

func (t Toast) Render(th theme.Theme, width int) string {
	token := t.token(th)
	surface := lipgloss.NewStyle().Background(th.Base.Overlay)
	badge := lipgloss.NewStyle().Foreground(token.Fg).Background(th.Base.Overlay)
	text := lipgloss.NewStyle().Foreground(th.Text.Primary).Background(th.Base.Overlay)
	pad := surface.Render(Repeat(" ", th.Space.ToastPadX))
	head := badge.Render(token.Glyph + " " + token.Label)
	body := text.Render(" " + t.Text)
	used := th.Space.ToastPadX*2 + Width(token.Glyph+" "+token.Label) + Width(" "+t.Text)
	return pad + head + body + surface.Render(Repeat(" ", width-used)) + pad
}

func ExpireToast(seq int) tea.Cmd {
	return tea.Tick(ToastLifetime, toastExpiry(seq))
}

func toastExpiry(seq int) func(time.Time) tea.Msg {
	return func(time.Time) tea.Msg { return ToastExpiredMsg{Seq: seq} }
}
