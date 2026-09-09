package app

import (
	"fmt"

	"github.com/remcostoeten/reusable-tui/internal/theme"
)

func errUnknownTheme(name string) error {
	return fmt.Errorf("%w: %s", theme.ErrUnknownTheme, name)
}
