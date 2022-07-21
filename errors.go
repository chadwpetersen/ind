package ind

import (
	"github.com/chadwpetersen/ind/errors"
)

var (
	ErrNoAvailableSlots = errors.New("no valid slots available")
	ErrTooManyPeople    = errors.New("too many people")
)
