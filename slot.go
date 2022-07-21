package ind

import (
	"sort"
	"time"
)

func PickSlot(slots []*Slot, before time.Time, after time.Time, strict bool) (*Slot, error) {
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

	if !strict {
		return slots[0], nil
	}

	for _, slot := range slots {
		d, err := time.Parse("2006-01-02", slot.Date)
		if err != nil {
			return nil, err
		}

		if d.After(after) && d.Before(before) {
			return slot, nil
		}
	}

	for _, slot := range slots {
		d, err := time.Parse("2006-01-02", slot.Date)
		if err != nil {
			return nil, err
		}

		if d.After(after) {
			return slot, nil
		}
	}

	return nil, ErrNoAvailableSlots
}
