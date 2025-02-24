package errors

import "errors"

var (
	ErrInvalidNewsId   = errors.New("invalid news id")
	ErrInvalidCategory = errors.New("invalid category id")
	ErrNotFoundNews    = errors.New("news not found")
)
