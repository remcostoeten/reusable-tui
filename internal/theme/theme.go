package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type Fidelity int

const (
	FidelityTrueColor Fidelity = iota
	FidelityBasic
	FidelityMono
)

type StatusToken struct {
	Fg    lipgloss.TerminalColor
	Label string
	Glyph string
}

type BaseTokens struct {
	Background lipgloss.TerminalColor
	Surface    lipgloss.TerminalColor
	Overlay    lipgloss.TerminalColor
}

type TextTokens struct {
	Primary   lipgloss.TerminalColor
	Secondary lipgloss.TerminalColor
	Muted     lipgloss.TerminalColor
	Disabled  lipgloss.TerminalColor
	Inverted  lipgloss.TerminalColor
}

type AccentTokens struct {
	Active lipgloss.TerminalColor
	Mid    lipgloss.TerminalColor
	Dim    lipgloss.TerminalColor
}

type BorderTokens struct {
	Focused      lipgloss.TerminalColor
	Unfocused    lipgloss.TerminalColor
	Subtle       lipgloss.TerminalColor
	FocusedSet   lipgloss.Border
	UnfocusedSet lipgloss.Border
}

type StatusTokens struct {
	Success StatusToken
	Warning StatusToken
	Danger  StatusToken
	Info    StatusToken
}

type SpaceTokens struct {
	Unit       int
	Gutter     int
	PanelPadX  int
	PanelPadY  int
	TitleInset int
	BadgeInset int
	StatusPadX int
	HeaderPadX int
	ModalPadX  int
	ModalPadY  int
	ToastPadX  int
	HeaderRows int
	StatusRows int
	MinWidth   int
	MinHeight  int
}

type EmptyPattern string

const (
	EmptyPatternSlash EmptyPattern = "slash"
	EmptyPatternDots  EmptyPattern = "dots"
	EmptyPatternNone  EmptyPattern = "none"
)

type MarkerTokens struct {
	FocusOpen     string
	FocusClose    string
	HintOpen      string
	HintClose     string
	Hatch         string
	Dot           string
	Empty         EmptyPattern
	Skeleton      string
	Rule          string
	Bullet        string
	SelectedLeft  string
	UnselectedPad string
	Spinner       string
}

type Theme struct {
	Name     string
	Fidelity Fidelity
	Base     BaseTokens
	Text     TextTokens
	Accent   AccentTokens
	Border   BorderTokens
	Status   StatusTokens
	Space    SpaceTokens
	Marker   MarkerTokens
}

func (t Theme) BorderColor(focused bool) lipgloss.TerminalColor {
	if focused {
		return t.Border.Focused
	}
	return t.Border.Unfocused
}

func (t Theme) BorderSet(focused bool) lipgloss.Border {
	if focused {
		return t.Border.FocusedSet
	}
	return t.Border.UnfocusedSet
}

func (t Theme) TitleColor(focused bool) lipgloss.TerminalColor {
	if focused {
		return t.Text.Primary
	}
	return t.Text.Secondary
}

func defaultSpace() SpaceTokens {
	return SpaceTokens{
		Unit:       1,
		Gutter:     1,
		PanelPadX:  1,
		PanelPadY:  0,
		TitleInset: 1,
		BadgeInset: 1,
		StatusPadX: 1,
		HeaderPadX: 1,
		ModalPadX:  2,
		ModalPadY:  1,
		ToastPadX:  1,
		HeaderRows: 2,
		StatusRows: 1,
		MinWidth:   64,
		MinHeight:  16,
	}
}

func defaultMarkers() MarkerTokens {
	return MarkerTokens{
		FocusOpen:     "",
		FocusClose:    "",
		HintOpen:      "(",
		HintClose:     ")",
		Hatch:         "/",
		Dot:           "·",
		Empty:         EmptyPatternSlash,
		Skeleton:      "█",
		Rule:          "─",
		Bullet:        "·",
		SelectedLeft:  ">",
		UnselectedPad: " ",
		Spinner:       "⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏",
	}
}

func FidelityFor(p termenv.Profile) Fidelity {
	switch p {
	case termenv.Ascii:
		return FidelityMono
	case termenv.ANSI:
		return FidelityBasic
	default:
		return FidelityTrueColor
	}
}

func Adapt(t Theme, f Fidelity) Theme {
	t.Fidelity = f
	if f == FidelityMono {
		t.Border.FocusedSet = lipgloss.ThickBorder()
		t.Border.UnfocusedSet = lipgloss.RoundedBorder()
		t.Marker.FocusOpen = "["
		t.Marker.FocusClose = "]"
		t.Marker.Spinner = "|/-\\"
	}
	return t
}
