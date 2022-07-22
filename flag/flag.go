package flag

import (
	"flag"
	"fmt"
	"os"

	"github.com/chadwpetersen/ind"
	"github.com/chadwpetersen/ind/log"
)

var (
	Venue    venue
	Strategy strategy

	Phone phone
	Email email

	Before before
	After  after
	Date   date

	Customers []ind.Customer
	vNumbers  stringList
	names     stringList
	surnames  stringList
)

func init() {
	flag.Var(&Venue, "venue", "(Required). Venue to book the appointment. (supported values [\"AM\", \"DB\", \"DH\", \"ZW\"]).")
	flag.Var(&Strategy, "strategy", "Strategy to use when booking the appointment. (supported values [\"individual\", \"together\"])")
	flag.Var(&Phone, "phone", "Phone is the phone number used for the appointment.")
	flag.Var(&Email, "email", "Email is the email address used for the appointment.")
	flag.Var(&Before, "before", "Before is an optional value to filter the appointments that occur before this provided date.")
	flag.Var(&After, "after", "After is an optional value to filter the appointments that occur after this provided date.")
	flag.Var(&Date, "date", "Date is an optional value to filter the appointments that occur on this provided date.")
	flag.Var(&vNumbers, "vnumbers", "V Numbers to use when booking the appointment.")
	flag.Var(&names, "names", "Names to use when booking the appointment.")
	flag.Var(&surnames, "surnames", "Surnames to use when booking the appointment.")
}

// Parse is a wrapper around the flag package.
func Parse() {
	flag.Parse()
	parseArgs()

	log.Debug("Provided the following arguments", log.WithLabels(map[string]any{
		"venue":    Venue.Value().String(),
		"before":   Before.Value(),
		"after":    After.Value(),
		"date":     Date.Value(),
		"strategy": Strategy.Value().String(),
		"phone":    Phone.Value(),
		"email":    Email.Value(),
		"vnumbers": vNumbers.Value(),
		"names":    names.Value(),
		"surnames": surnames.Value(),
	}))
}

func parseArgs() {
	var (
		required = []string{"venue", "phone", "email", "vnumbers", "names", "surnames"}
		seen     = make(map[string]struct{})
	)

	flag.Visit(func(f *flag.Flag) {
		seen[f.Name] = struct{}{}
	})

	for _, name := range required {
		if _, ok := seen[name]; ok {
			continue
		}

		usage(2, "missing value for required flag -%s", name)
	}

	// Ignore Before and After values if date is provided.
	if Date.Value() != nil {
		Before = before{}
		After = after{}
	}

	if len(vNumbers) != len(names) {
		usage(2, "unequal amount of vnumbers and names")
	}

	for i, vn := range vNumbers {
		surname := surnames[0]
		if len(surnames) == len(vNumbers) {
			surname = surnames[i]
		}

		Customers = append(Customers, ind.Customer{
			VNumber: vn,
			Name:    names[i],
			Surname: surname,
		})
	}
}

func usage(code int, format string, a ...any) {
	fmt.Fprintf(flag.CommandLine.Output(), format+"\n", a...)
	flag.Usage()
	os.Exit(code)
}
