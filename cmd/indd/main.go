package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/chadwpetersen/ind"
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
	venue    = ind.VenueAmsterdam
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
		log.Println("Running with together strategy")
		run(customers...)
		return
	}

	log.Println("Running with individual strategy")
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
			log.Fatalln(err)
			return
		}
		if len(slots) == 0 {
			log.Println("No slots found")
			continue
		}

		dates := make([]string, len(slots))
		for _, slot := range slots {
			dates = append(dates, fmt.Sprintf("\n\t %s (%s - %s)", slot.Date, slot.Start, slot.End))
		}
		log.Printf("Found some slots: %v\n", dates)

		slot, err = ind.PickSlot(slots, before, after, strictly)
		if errors.Is(err, ind.ErrNoAvailableSlots) {
			log.Printf("No available slots to pick")
			continue
		} else if err != nil {
			log.Fatalln(err)
			return
		}
		log.Printf("Pick Slot: %v \n", slot)

		break
	}

	rdata, err := ind.ReserveAppointment(venue, *slot)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Reserve: %v \n", rdata)

	data, err := ind.CreateAppointment(venue, makeAppRequest(*slot, cl))
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Printf("Appointment: %s \n", string(data))

	err = makeOutput(*slot, data)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("Output successfully created")
}

func makeOutput(slot ind.Slot, raw []byte) error {
	t, err := template.ParseFiles("data/template.txt")
	// Capture any error
	if err != nil {
		return err
	}

	f, err := os.Create("data/output/say.txt")
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", slot.Date)
	if err != nil {
		return err
	}

	slot.Date = date.Format("January, 2 2006")

	err = t.Execute(f, slot)
	if err != nil {
		return err
	}

	rf, err := os.Create("data/output/raw.txt")
	if err != nil {
		return err
	}

	if _, err := rf.WriteString(string(raw)); err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		cmd := exec.Command("/usr/bin/say", "-f", "data/output/say.txt")
		err := cmd.Start()
		if err != nil {
			return err
		}

		err = cmd.Wait()
		if err != nil {
			return err
		}
	}

	return nil
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
