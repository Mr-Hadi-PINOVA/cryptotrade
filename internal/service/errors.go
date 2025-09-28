package service

import "errors"

// ErrValidation indicates the input payload failed validation.
var ErrValidation = errors.New("validation error")
