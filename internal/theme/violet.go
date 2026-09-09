package theme

import "github.com/charmbracelet/lipgloss"

const NameVioletDark = "violet-dark"

func VioletDark() Theme {
	return Theme{
		Name:     NameVioletDark,
		Fidelity: FidelityTrueColor,
		Base: BaseTokens{
			Background: lipgloss.Color("#0B0B0E"),
			Surface:    lipgloss.Color("#121216"),
			Overlay:    lipgloss.Color("#17171D"),
		},
		Text: TextTokens{
			Primary:   lipgloss.Color("#E6E6EA"),
			Secondary: lipgloss.Color("#9A9AA5"),
			Muted:     lipgloss.Color("#56565F"),
			Disabled:  lipgloss.Color("#34343B"),
			Inverted:  lipgloss.Color("#0B0B0E"),
		},
		Accent: AccentTokens{
			Active: lipgloss.Color("#8B5CF6"),
			Mid:    lipgloss.Color("#6D4AAF"),
			Dim:    lipgloss.Color("#3C2E5C"),
		},
		Border: BorderTokens{
			Focused:      lipgloss.Color("#8B5CF6"),
			Unfocused:    lipgloss.Color("#34343B"),
			Subtle:       lipgloss.Color("#1E1E24"),
			FocusedSet:   lipgloss.RoundedBorder(),
			UnfocusedSet: lipgloss.RoundedBorder(),
		},
		Status: StatusTokens{
			Success: StatusToken{Fg: lipgloss.Color("#3FB950"), Label: "OK", Glyph: "+"},
			Warning: StatusToken{Fg: lipgloss.Color("#D29922"), Label: "WARN", Glyph: "!"},
			Danger:  StatusToken{Fg: lipgloss.Color("#F85149"), Label: "ERR", Glyph: "x"},
			Info:    StatusToken{Fg: lipgloss.Color("#58A6FF"), Label: "INFO", Glyph: "i"},
		},
		Space:  defaultSpace(),
		Marker: defaultMarkers(),
	}
}
