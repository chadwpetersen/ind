package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/chadwpetersen/ind"
	"github.com/chadwpetersen/ind/data"
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
	if strategy == ind.StrategyTogether {
		log.Info("Running with together strategy")
		run(customers...)
		return
	}

	log.Info("Running with individual strategy")
	for _, customer := range customers {
		run(customer)
	}
}

func run(cl ...ind.Customer) {
	var (
		slot *ind.Slot
		err  error
	)
	for {
		time.Sleep(10 * time.Second)

		slots, err := ind.FetchSlots(venue, uint(len(cl)))
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
		for _, slot := range slots {
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
				"venue":     venue,
				"slot":      slot,
				"customers": cl,
			}))
		break
	}

	rslot, err := ind.ReserveAppointment(venue, *slot)
	if err != nil {
		log.Error("Failed to reserve appointment", err)
		return
	}
	log.Info("Reserved appointment", log.WithLabels(
		map[string]any{
			"venue":         venue,
			"reserved_slot": rslot,
			"customers":     cl,
		}))

	apt, err := ind.CreateAppointment(venue, makeAppRequest(*slot, cl))
	if err != nil {
		log.Error("Failed to create appointment", err)
		return
	}
	log.Pass("Created appointment", log.WithLabels(
		map[string]any{
			"venue":               venue,
			"appointment_details": string(apt),
			"customers":           cl,
		}))

	err = data.Generate(*slot, apt)
	if err != nil {
		log.Error("Failed to create output", err)
		return
	}

	log.Pass("Output successfully created")
}

func makeAppRequest(slot ind.Slot, cl []ind.Customer) ind.AppointmentReq {
	return ind.AppointmentReq{
		Slot: slot,
		Apt: ind.Appointment{
			ProductKey: "DOC",
			Date:       slot.Date,
			Start:      slot.Start,
			End:        slot.End,
			Email:      email,
			Phone:      phone,
			Language:   "en",
			Customers:  cl,
		},
	}
}
