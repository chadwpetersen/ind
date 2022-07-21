package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chadwpetersen/ind"
	"github.com/chadwpetersen/ind/data"
	"github.com/chadwpetersen/ind/errors"
	"github.com/chadwpetersen/ind/log"
)

var (
	// This helps filters the range of dates to select.
	after  = time.Date(2022, time.September, 1, 0, 0, 0, 0, time.UTC)
	before = time.Date(2022, time.September, 15, 0, 0, 0, 0, time.UTC)

	// Setting strictly to false allows for other dates outside
	// the date range to be selected as well. With strictly set
	// to false this will pick the first next available slot that
	// is detected.
	strictly = true

	// Select a venue and a strategy to find an appointment.
	venue    = ind.VenueDenHaag
	strategy = ind.StrategyIndividual

	// Your personal details along with the customer information
	// should go here.
	email     = "abc@example.com"
	phone     = "+27912345678"
	customers = []ind.Customer{
		{
			VNumber: "",
			Name:    "",
			Surname: "",
		},
		{
			VNumber: "",
			Name:    "",
			Surname: "",
		},
		{
			VNumber: "",
			Name:    "",
			Surname: "",
		},
	}
)

func main() {
	ctx := context.Background()

	if strategy == ind.StrategyTogether {
		log.Info("Running with together strategy")
		run(ctx, customers...)
		return
	}

	log.Info("Running with individual strategy")
	for _, customer := range customers {
		run(ctx, customer)
	}
}

func run(ctx context.Context, cl ...ind.Customer) {
	var (
		slot *ind.Slot
		err  error
	)

	client := ind.NewClient()
	for {
		time.Sleep(10 * time.Second)

		slots, err := client.Find(ctx, venue, len(cl))
		if err != nil {
			log.Error("Failed to fetch slots", err, log.WithLabels(
				map[string]any{
					"venue":     venue,
					"customers": cl,
				}))
			return
		}
		if len(slots) == 0 {
			log.Warn("No slots found", log.WithLabels(
				map[string]any{
					"venue":     venue,
					"customers": cl,
				}))
			continue
		}

		dates := make([]string, 0)
		for i, slot := range slots {
			slots[i].Venue = venue
			dates = append(dates, fmt.Sprintf("%s (%s - %s)", slot.Date, slot.Start, slot.End))
		}

		log.Info("Found some slots", log.WithLabels(
			map[string]any{
				"venue":     venue,
				"customers": cl,
				"dates":     dates,
			}))

		slot, err = ind.PickSlot(slots, before, after, strictly)
		if errors.Is(err, ind.ErrNoAvailableSlots) {
			log.Warn("No available slots to pick")
			continue
		} else if err != nil {
			log.Error("Failed to pick a slot", err)
			return
		}

		log.Info("Picked a slot", log.WithAlert(), log.WithLabels(
			map[string]any{
				"slot":      slot,
				"customers": cl,
			}))
		break
	}

	err = client.Reserve(ctx, slot)
	if err != nil {
		log.Error("Failed to reserve appointment", err)
		return
	}
	log.Info("Reserved appointment", log.WithLabels(
		map[string]any{
			"reserved_slot": slot,
			"customers":     cl,
		}))

	raw, err := client.Book(ctx, email, phone, slot, cl)
	if err != nil {
		log.Error("Failed to create appointment", err)
		return
	}

	err = data.Generate(*slot, raw)
	if err != nil {
		log.Error("Failed to create output", err)
		return
	}

	log.Pass("Output successfully created")
}
