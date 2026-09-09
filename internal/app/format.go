package app

import "strconv"

func itoa(value int) string {
	return strconv.Itoa(value)
}

func joinWarnings(warnings []error) string {
	if len(warnings) == 0 {
		return ""
	}
	if len(warnings) == 1 {
		return warnings[0].Error()
	}
	return warnings[0].Error() + " (+" + itoa(len(warnings)-1) + " more)"
}
