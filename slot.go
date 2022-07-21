package ind

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

func FetchSlots(venue Venue, persons uint) ([]Slot, error) {
	if persons > 6 {
		return nil, errors.New("too many people")
	}

	resp, err := http.Get(slotURL(venue, persons))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid http status")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data = bytes.ReplaceAll(data, []byte(")]}',"), []byte(""))

	var slots = struct {
		Status string `json:"status"`
		Data   []Slot `json:"data"`
	}{}

	err = json.Unmarshal(data, &slots)
	if err != nil {
		return nil, err
	}

	return slots.Data, nil
}

func PickSlot(slots []Slot, before time.Time, after time.Time, strict bool) (*Slot, error) {
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

		if d.After(after) && d.Before(before) {
			return &slot, nil
		}
	}

	if !strict {
		for _, slot := range slots {
			d, err := time.Parse("2006-01-02", slot.Date)
			if err != nil {
				return nil, err
			}

			if d.After(after) {
				return &slot, nil
			}
		}
	}

	return nil, ErrNoAvailableSlots
}

func slotURL(venue Venue, persons uint) string {
	return fmt.Sprintf(
		"https://oap.ind.nl/oap/api/desks/%s/slots/?productKey=DOC&persons=%d",
		venue,
		persons,
	)
}
