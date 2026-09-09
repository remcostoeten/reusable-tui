package theme

import "github.com/charmbracelet/lipgloss"

const NameHighContrast = "high-contrast"

func HighContrast() Theme {
	return Theme{
		Name:     NameHighContrast,
		Fidelity: FidelityTrueColor,
		Base: BaseTokens{
			Background: lipgloss.Color("#000000"),
			Surface:    lipgloss.Color("#0A0A0A"),
			Overlay:    lipgloss.Color("#101010"),
		},
		Text: TextTokens{
			Primary:   lipgloss.Color("#FFFFFF"),
			Secondary: lipgloss.Color("#D4D4D8"),
			Muted:     lipgloss.Color("#A1A1AA"),
			Disabled:  lipgloss.Color("#71717A"),
			Inverted:  lipgloss.Color("#000000"),
		},
		Accent: AccentTokens{
			Active: lipgloss.Color("#C4B5FD"),
			Mid:    lipgloss.Color("#A78BFA"),
			Dim:    lipgloss.Color("#7C3AED"),
		},
		Border: BorderTokens{
			Focused:      lipgloss.Color("#C4B5FD"),
			Unfocused:    lipgloss.Color("#A1A1AA"),
			Subtle:       lipgloss.Color("#71717A"),
			FocusedSet:   lipgloss.ThickBorder(),
			UnfocusedSet: lipgloss.RoundedBorder(),
		},
		Status: StatusTokens{
			Success: StatusToken{Fg: lipgloss.Color("#4ADE80"), Label: "OK", Glyph: "+"},
			Warning: StatusToken{Fg: lipgloss.Color("#FACC15"), Label: "WARN", Glyph: "!"},
			Danger:  StatusToken{Fg: lipgloss.Color("#FB7185"), Label: "ERR", Glyph: "x"},
			Info:    StatusToken{Fg: lipgloss.Color("#7DD3FC"), Label: "INFO", Glyph: "i"},
		},
		Space:  defaultSpace(),
		Marker: defaultMarkers(),
	}
}
