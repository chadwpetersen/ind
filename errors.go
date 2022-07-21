package ind

import (
	"errors"
)

var (
	ErrNoAvailableSlots  = errors.New("no valid slots available")
	ErrInvalidHTTPStatus = errors.New("invalid http status")
	ErrTooManyPeople     = errors.New("too many people")
)
