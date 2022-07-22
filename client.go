package ind

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/chadwpetersen/ind/http"
	"github.com/chadwpetersen/ind/log"
)

type client struct {
	httpClient *http.Client
}

func NewClient() client {
	return client{
		httpClient: http.NewClient(baseURL),
	}
}

func (c client) Find(ctx context.Context, venue Venue,
	amount int) ([]*Slot, error) {

	if amount > 6 {
		return nil, ErrTooManyPeople
	}

	url := fmt.Sprintf("%s/slots/?productKey=DOC&persons=%d", venue.Code(), amount)

	data, _, err := c.httpClient.GET(ctx, url)
	if err != nil {
		return nil, err
	}

	var slots = struct {
		Status string  `json:"status"`
		Data   []*Slot `json:"data"`
	}{}

	err = json.Unmarshal(data, &slots)
	if err != nil {
		return nil, err
	}

	return slots.Data, nil
}

func (c client) Pick(_ context.Context, slots []*Slot,
	date, before, after *time.Time) (*Slot, error) {

	if len(slots) == 0 {
		return nil, ErrNoAvailableSlots
	}

	sort.Slice(slots, func(i, j int) bool {
		aDate, err := time.Parse("2006-01-02", slots[i].Date)
		if err != nil {
			return false
		}

		bDate, err := time.Parse("2006-01-02", slots[j].Date)
		if err != nil {
			return false
		}

		return aDate.Before(bDate)
	})

	for _, slot := range slots {
		d, err := time.Parse("2006-01-02", slot.Date)
		if err != nil {
			return nil, err
		}

		if date != nil {
			if d.Equal(*date) {
				return slot, nil
			}
		} else if before != nil && after != nil {
			if d.Before(*before) && d.After(*after) {
				return slot, nil
			}
		} else if before != nil {
			if d.Before(*before) {
				return slot, nil
			}
		} else if after != nil {
			if d.After(*after) {
				return slot, nil
			}
		}
	}

	return nil, ErrNoAvailableSlots
}

func (c client) Reserve(ctx context.Context, slot *Slot) error {
	url := fmt.Sprintf("%s/slots/%s", slot.Venue.Code(), slot.Key)

	data, _, err := c.httpClient.POST(ctx, url, slot)
	if err != nil {
		return err
	}

	var reserve = struct {
		Status string `json:"status"`
		Data   Slot   `json:"data"`
	}{}

	err = json.Unmarshal(data, &reserve)
	if err != nil {
		return err
	}

	return nil
}

func (c client) Book(ctx context.Context, email, phone string,
	slot *Slot, customers []Customer) ([]byte, error) {

	var (
		url  = fmt.Sprintf("%s/appointments", slot.Venue.Code())
		aReq = struct {
			Slot *Slot       `json:"bookableSlot"`
			Apt  Appointment `json:"appointment"`
		}{
			Slot: slot,
			Apt: Appointment{
				ProductKey: "DOC",
				Date:       slot.Date,
				Start:      slot.Start,
				End:        slot.End,
				Email:      email,
				Phone:      phone,
				Language:   "en",
				Customers:  customers,
			},
		}
	)

	data, _, err := c.httpClient.POST(ctx, url, aReq)
	if err != nil {
		return nil, err
	}

	log.Pass("Booked appointment", log.WithLabels(
		map[string]any{
			"request":  aReq,
			"response": data,
		},
	))

	return data, nil
}
