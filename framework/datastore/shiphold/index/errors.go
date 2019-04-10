package index

import (
	"errors"
)

var (
	errNotFound      = errors.New("not found")
	errAlreadyExists = errors.New("already exists")
	errNilParam      = errors.New("param must not be nil")
)
