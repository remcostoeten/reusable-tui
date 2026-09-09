package app

import "errors"

func errUnknownTheme(name string) error {
	return errors.New("theme not registered: " + name)
}
