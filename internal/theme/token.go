package theme

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	TokenBaseBackground = "base.background"
	TokenBaseSurface    = "base.surface"
	TokenBaseOverlay    = "base.overlay"

	TokenTextPrimary   = "text.primary"
	TokenTextSecondary = "text.secondary"
	TokenTextMuted     = "text.muted"
	TokenTextDisabled  = "text.disabled"
	TokenTextInverted  = "text.inverted"

	TokenAccentActive = "accent.active"
	TokenAccentMid    = "accent.mid"
	TokenAccentDim    = "accent.dim"

	TokenBorderFocused      = "border.focused"
	TokenBorderUnfocused    = "border.unfocused"
	TokenBorderSubtle       = "border.subtle"
	TokenBorderFocusedSet   = "border.focused.set"
	TokenBorderUnfocusedSet = "border.unfocused.set"

	TokenStatusSuccessFg    = "status.success.fg"
	TokenStatusSuccessLabel = "status.success.label"
	TokenStatusSuccessGlyph = "status.success.glyph"
	TokenStatusWarningFg    = "status.warning.fg"
	TokenStatusWarningLabel = "status.warning.label"
	TokenStatusWarningGlyph = "status.warning.glyph"
	TokenStatusDangerFg     = "status.danger.fg"
	TokenStatusDangerLabel  = "status.danger.label"
	TokenStatusDangerGlyph  = "status.danger.glyph"
	TokenStatusInfoFg       = "status.info.fg"
	TokenStatusInfoLabel    = "status.info.label"
	TokenStatusInfoGlyph    = "status.info.glyph"
)

const (
	BorderRounded = "rounded"
	BorderNormal  = "normal"
	BorderThick   = "thick"
	BorderDouble  = "double"
	BorderBlock   = "block"
	BorderHidden  = "hidden"
)

var tokenOrder = []string{
	TokenBaseBackground,
	TokenBaseSurface,
	TokenBaseOverlay,
	TokenTextPrimary,
	TokenTextSecondary,
	TokenTextMuted,
	TokenTextDisabled,
	TokenTextInverted,
	TokenAccentActive,
	TokenAccentMid,
	TokenAccentDim,
	TokenBorderFocused,
	TokenBorderUnfocused,
	TokenBorderSubtle,
	TokenBorderFocusedSet,
	TokenBorderUnfocusedSet,
	TokenStatusSuccessFg,
	TokenStatusSuccessLabel,
	TokenStatusSuccessGlyph,
	TokenStatusWarningFg,
	TokenStatusWarningLabel,
	TokenStatusWarningGlyph,
	TokenStatusDangerFg,
	TokenStatusDangerLabel,
	TokenStatusDangerGlyph,
	TokenStatusInfoFg,
	TokenStatusInfoLabel,
	TokenStatusInfoGlyph,
}

var (
	ErrUnknownToken = errors.New("theme: unknown token")
	ErrInvalidColor = errors.New("theme: invalid color")
	ErrInvalidValue = errors.New("theme: invalid value")
)

var hexPattern = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

func Tokens() []string {
	out := make([]string, len(tokenOrder))
	copy(out, tokenOrder)
	return out
}

func ParseColor(value string) (lipgloss.TerminalColor, error) {
	if hexPattern.MatchString(value) {
		return lipgloss.Color(strings.ToUpper(value)), nil
	}
	if index, err := strconv.Atoi(value); err == nil && index >= 0 && index <= 255 {
		return lipgloss.Color(value), nil
	}
	return nil, tokenError(ErrInvalidColor, value)
}

func FormatColor(color lipgloss.TerminalColor) string {
	if plain, ok := color.(lipgloss.Color); ok {
		return string(plain)
	}
	return ""
}

func ParseBorder(name string) (lipgloss.Border, error) {
	switch name {
	case BorderRounded:
		return lipgloss.RoundedBorder(), nil
	case BorderNormal:
		return lipgloss.NormalBorder(), nil
	case BorderThick:
		return lipgloss.ThickBorder(), nil
	case BorderDouble:
		return lipgloss.DoubleBorder(), nil
	case BorderBlock:
		return lipgloss.BlockBorder(), nil
	case BorderHidden:
		return lipgloss.HiddenBorder(), nil
	}
	return lipgloss.Border{}, tokenError(ErrInvalidValue, name)
}

func FormatBorder(border lipgloss.Border) string {
	switch border.TopLeft {
	case lipgloss.RoundedBorder().TopLeft:
		return BorderRounded
	case lipgloss.ThickBorder().TopLeft:
		return BorderThick
	case lipgloss.DoubleBorder().TopLeft:
		return BorderDouble
	case lipgloss.BlockBorder().TopLeft:
		return BorderBlock
	case lipgloss.HiddenBorder().TopLeft:
		return BorderHidden
	}
	return BorderNormal
}

func ReadToken(t Theme, path string) (string, error) {
	switch path {
	case TokenBaseBackground:
		return FormatColor(t.Base.Background), nil
	case TokenBaseSurface:
		return FormatColor(t.Base.Surface), nil
	case TokenBaseOverlay:
		return FormatColor(t.Base.Overlay), nil
	case TokenTextPrimary:
		return FormatColor(t.Text.Primary), nil
	case TokenTextSecondary:
		return FormatColor(t.Text.Secondary), nil
	case TokenTextMuted:
		return FormatColor(t.Text.Muted), nil
	case TokenTextDisabled:
		return FormatColor(t.Text.Disabled), nil
	case TokenTextInverted:
		return FormatColor(t.Text.Inverted), nil
	case TokenAccentActive:
		return FormatColor(t.Accent.Active), nil
	case TokenAccentMid:
		return FormatColor(t.Accent.Mid), nil
	case TokenAccentDim:
		return FormatColor(t.Accent.Dim), nil
	case TokenBorderFocused:
		return FormatColor(t.Border.Focused), nil
	case TokenBorderUnfocused:
		return FormatColor(t.Border.Unfocused), nil
	case TokenBorderSubtle:
		return FormatColor(t.Border.Subtle), nil
	case TokenBorderFocusedSet:
		return FormatBorder(t.Border.FocusedSet), nil
	case TokenBorderUnfocusedSet:
		return FormatBorder(t.Border.UnfocusedSet), nil
	case TokenStatusSuccessFg:
		return FormatColor(t.Status.Success.Fg), nil
	case TokenStatusSuccessLabel:
		return t.Status.Success.Label, nil
	case TokenStatusSuccessGlyph:
		return t.Status.Success.Glyph, nil
	case TokenStatusWarningFg:
		return FormatColor(t.Status.Warning.Fg), nil
	case TokenStatusWarningLabel:
		return t.Status.Warning.Label, nil
	case TokenStatusWarningGlyph:
		return t.Status.Warning.Glyph, nil
	case TokenStatusDangerFg:
		return FormatColor(t.Status.Danger.Fg), nil
	case TokenStatusDangerLabel:
		return t.Status.Danger.Label, nil
	case TokenStatusDangerGlyph:
		return t.Status.Danger.Glyph, nil
	case TokenStatusInfoFg:
		return FormatColor(t.Status.Info.Fg), nil
	case TokenStatusInfoLabel:
		return t.Status.Info.Label, nil
	case TokenStatusInfoGlyph:
		return t.Status.Info.Glyph, nil
	}
	return "", tokenError(ErrUnknownToken, path)
}

func WriteToken(t Theme, path, value string) (Theme, error) {
	if strings.HasSuffix(path, ".set") {
		return writeBorderSet(t, path, value)
	}
	if strings.HasSuffix(path, ".label") || strings.HasSuffix(path, ".glyph") {
		return writeStatusText(t, path, value)
	}
	return writeColorToken(t, path, value)
}

func writeColorToken(t Theme, path, value string) (Theme, error) {
	color, err := ParseColor(value)
	if err != nil {
		return t, fmt.Errorf("%s: %w", path, err)
	}
	switch path {
	case TokenBaseBackground:
		t.Base.Background = color
	case TokenBaseSurface:
		t.Base.Surface = color
	case TokenBaseOverlay:
		t.Base.Overlay = color
	case TokenTextPrimary:
		t.Text.Primary = color
	case TokenTextSecondary:
		t.Text.Secondary = color
	case TokenTextMuted:
		t.Text.Muted = color
	case TokenTextDisabled:
		t.Text.Disabled = color
	case TokenTextInverted:
		t.Text.Inverted = color
	case TokenAccentActive:
		t.Accent.Active = color
	case TokenAccentMid:
		t.Accent.Mid = color
	case TokenAccentDim:
		t.Accent.Dim = color
	case TokenBorderFocused:
		t.Border.Focused = color
	case TokenBorderUnfocused:
		t.Border.Unfocused = color
	case TokenBorderSubtle:
		t.Border.Subtle = color
	case TokenStatusSuccessFg:
		t.Status.Success.Fg = color
	case TokenStatusWarningFg:
		t.Status.Warning.Fg = color
	case TokenStatusDangerFg:
		t.Status.Danger.Fg = color
	case TokenStatusInfoFg:
		t.Status.Info.Fg = color
	default:
		return t, tokenError(ErrUnknownToken, path)
	}
	return t, nil
}

func writeBorderSet(t Theme, path, value string) (Theme, error) {
	border, err := ParseBorder(value)
	if err != nil {
		return t, fmt.Errorf("%s: %w", path, err)
	}
	switch path {
	case TokenBorderFocusedSet:
		t.Border.FocusedSet = border
	case TokenBorderUnfocusedSet:
		t.Border.UnfocusedSet = border
	default:
		return t, tokenError(ErrUnknownToken, path)
	}
	return t, nil
}

func writeStatusText(t Theme, path, value string) (Theme, error) {
	if value == "" {
		return t, tokenError(ErrInvalidValue, path)
	}
	switch path {
	case TokenStatusSuccessLabel:
		t.Status.Success.Label = value
	case TokenStatusSuccessGlyph:
		t.Status.Success.Glyph = value
	case TokenStatusWarningLabel:
		t.Status.Warning.Label = value
	case TokenStatusWarningGlyph:
		t.Status.Warning.Glyph = value
	case TokenStatusDangerLabel:
		t.Status.Danger.Label = value
	case TokenStatusDangerGlyph:
		t.Status.Danger.Glyph = value
	case TokenStatusInfoLabel:
		t.Status.Info.Label = value
	case TokenStatusInfoGlyph:
		t.Status.Info.Glyph = value
	default:
		return t, tokenError(ErrUnknownToken, path)
	}
	return t, nil
}

func tokenError(base error, detail string) error {
	return fmt.Errorf("%w: %s", base, detail)
}
