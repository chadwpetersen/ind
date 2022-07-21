package ind

import (
	"context"
)

type API interface {
	Find(ctx context.Context, venue Venue, amount int) ([]*Slot, error)

	Reserve(ctx context.Context, slot *Slot) error

	Book(ctx context.Context, email string, phone string, slot *Slot, customers []Customer) ([]byte, error)
}
