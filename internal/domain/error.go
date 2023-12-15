package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrRequired        = errors.New("required field")
	ErrInvalidVariable = errors.New("invalid variable")
	ErrInvalidValue    = errors.New("invalid value")
)
