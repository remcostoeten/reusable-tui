package store

import "errors"

var (
	ErrNotFound = errors.New("store: record not found")
	ErrConflict = errors.New("store: conflicting record")
	ErrInvalid  = errors.New("store: invalid input")
)
