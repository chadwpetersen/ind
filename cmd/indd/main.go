package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chadwpetersen/ind"
	"github.com/chadwpetersen/ind/data"
	"github.com/chadwpetersen/ind/errors"
	"github.com/chadwpetersen/ind/flag"
	"github.com/chadwpetersen/ind/log"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	if flag.Strategy.Value() == ind.StrategyTogether {
		log.Info("Running with together strategy")
		run(ctx, flag.Customers...)
		return
	}

	log.Info("Running with individual strategy")
	for _, customer := range flag.Customers {
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

		slots, err := client.Find(ctx, flag.Venue.Value(), len(cl))
		if err != nil {
			log.Error("Failed to fetch slots", err, log.WithLabels(
				map[string]any{
					"venue":     flag.Venue,
					"customers": cl,
				}))
			return
		}
		if len(slots) == 0 {
			log.Warn("No slots found", log.WithLabels(
				map[string]any{
					"venue":     flag.Venue,
					"customers": cl,
				}))
			continue
		}

		dates := make([]string, 0)
		for i, slot := range slots {
			slots[i].Venue = flag.Venue.Value()
			dates = append(dates, fmt.Sprintf("%s (%s - %s)", slot.Date, slot.Start, slot.End))
		}

		log.Info("Found some slots", log.WithLabels(
			map[string]any{
				"venue":     flag.Venue,
				"customers": cl,
				"dates":     dates,
			}))

		slot, err = client.Pick(ctx, slots, flag.Before.Value(), flag.After.Value())
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

	raw, err := client.Book(ctx, flag.Email.Value(), flag.Phone.Value(), slot, cl)
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
