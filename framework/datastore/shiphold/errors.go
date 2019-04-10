package shiphold

import (
	"errors"
)

var (
	errNoID              = errors.New("missing struct tag id or id field")
	errZeroID            = errors.New("id field must not be a zero value")
	errBadType           = errors.New("provided data must be a struct or a pointer to struct")
	errAlreadyExists     = errors.New("already exists")
	errNilParam          = errors.New("param must not be nil")
	errUnknownTag        = errors.New("unknown tag")
	errIdxNotFound       = errors.New("index not found")
	errSlicePtrNeeded    = errors.New("provided target must be a pointer to slice")
	errStructPtrNeeded   = errors.New("provided target must be a pointer to struct")
	errPtrNeeded         = errors.New("provided target must be a pointer to a valid variable")
	errNoName            = errors.New("provided target must have a name")
	errNotFound          = errors.New("not found")
	errNotInTransaction  = errors.New("not in transaction")
	errIncompatibleValue = errors.New("incompatible value")
	errDifferentCodec    = errors.New("the selected codec is incompatible with this bucket")
)
