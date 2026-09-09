package theme

import "github.com/charmbracelet/lipgloss"

const NameMonochrome = "monochrome"

func Monochrome() Theme {
	return Theme{
		Name:     NameMonochrome,
		Fidelity: FidelityTrueColor,
		Base: BaseTokens{
			Background: lipgloss.Color("#0A0A0A"),
			Surface:    lipgloss.Color("#131313"),
			Overlay:    lipgloss.Color("#1A1A1A"),
		},
		Text: TextTokens{
			Primary:   lipgloss.Color("#EDEDED"),
			Secondary: lipgloss.Color("#9E9E9E"),
			Muted:     lipgloss.Color("#5A5A5A"),
			Disabled:  lipgloss.Color("#363636"),
			Inverted:  lipgloss.Color("#0A0A0A"),
		},
		Accent: AccentTokens{
			Active: lipgloss.Color("#EDEDED"),
			Mid:    lipgloss.Color("#9E9E9E"),
			Dim:    lipgloss.Color("#5A5A5A"),
		},
		Border: BorderTokens{
			Focused:      lipgloss.Color("#EDEDED"),
			Unfocused:    lipgloss.Color("#363636"),
			Subtle:       lipgloss.Color("#1F1F1F"),
			FocusedSet:   lipgloss.ThickBorder(),
			UnfocusedSet: lipgloss.RoundedBorder(),
		},
		Status: StatusTokens{
			Success: StatusToken{Fg: lipgloss.Color("#EDEDED"), Label: "OK", Glyph: "+"},
			Warning: StatusToken{Fg: lipgloss.Color("#C4C4C4"), Label: "WARN", Glyph: "!"},
			Danger:  StatusToken{Fg: lipgloss.Color("#FFFFFF"), Label: "ERR", Glyph: "x"},
			Info:    StatusToken{Fg: lipgloss.Color("#9E9E9E"), Label: "INFO", Glyph: "i"},
		},
		Space:  defaultSpace(),
		Marker: defaultMarkers(),
	}
}
