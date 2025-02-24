package models

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("internal repository error")
	ErrValidation = errors.New("validation error")
)
