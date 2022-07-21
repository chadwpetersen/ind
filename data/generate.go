package data

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/chadwpetersen/ind"
	"github.com/chadwpetersen/ind/alert"
)

// Generate creates the data output files when an appointment
// gets booked successfully.
func Generate(slot ind.Slot, raw []byte) error {
	t, err := template.ParseFiles("data/template.txt")
	// Capture any error
	if err != nil {
		return err
	}

	var (
		ts    = time.Now().Unix()
		sName = fmt.Sprintf("data/output/say-%d.txt", ts)
		rName = fmt.Sprintf("data/output/raw-%d.txt", ts)
	)

	f, err := os.Create(sName)
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

	rf, err := os.Create(rName)
	if err != nil {
		return err
	}

	if _, err := rf.WriteString(string(raw)); err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		if err := alert.Say("-f", sName); err != nil {
			return err
		}
	}

	return nil
}
