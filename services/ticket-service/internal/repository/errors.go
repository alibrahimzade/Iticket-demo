package repository

import "errors"

var ErrNotFound = errors.New("event not found")
var ErrInvalidState = errors.New("invalid state")
