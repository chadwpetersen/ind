package flag

import (
	"errors"
	"strings"
	"time"

	"github.com/chadwpetersen/ind"
)

type venue ind.Venue

func (v *venue) Set(val string) error {
	if !ind.Venue(val).IsValid() {
		return errors.New("supported values are [\"AM\", \"DB\", \"DH\", \"ZW\"]")
	}
	*v = venue(val)
	return nil
}

func (v *venue) Get() any {
	return v.Value().Code()
}

func (v *venue) Value() ind.Venue {
	return ind.Venue(*v)
}

func (v *venue) String() string {
	return ""
}

type strategy ind.Strategy

func (s *strategy) Set(val string) error {
	if val == ind.StrategyIndividual.String() {
		*s = strategy(1)
		return nil
	}
	if val == ind.StrategyTogether.String() {
		*s = strategy(2)
		return nil
	}
	return errors.New("supported values are [\"individual\", \"together\"]")
}

func (s *strategy) Get() any {
	return s.Value().String()
}

func (s *strategy) Value() ind.Strategy {
	return ind.Strategy(*s)
}

func (s *strategy) String() string {
	*s = strategy(ind.StrategyTogether)
	return string(*s)
}

type email string

func (e *email) Set(val string) error {
	*e = email(val)
	return nil
}

func (e *email) Get() any {
	return e.Value()
}

func (e *email) Value() string {
	return string(*e)
}

func (e *email) String() string {
	return ""
}

type phone string

func (p *phone) Set(val string) error {
	*p = phone(val)
	return nil
}

func (p *phone) Get() any {
	return p.Value()
}

func (p *phone) Value() string {
	return string(*p)
}

func (p *phone) String() string {
	return ""
}

type before time.Time

func (b *before) Set(val string) error {
	t, err := time.Parse("02/Jan/2006", val)
	if err != nil {
		return errors.New("date format should be dd/mmm/yyyy")
	}
	*b = before(t)
	return nil
}

func (b *before) Get() any {
	return b.Value()
}

func (b *before) Value() *time.Time {
	t := time.Time(*b)
	if t.IsZero() {
		return nil
	}
	return &t
}

func (b *before) String() string {
	return ""
}

type after time.Time

func (a *after) Set(val string) error {
	t, err := time.Parse("02/Jan/2006", val)
	if err != nil {
		err = errors.New("date format should be dd/mmm/yyyy")
	}
	*a = after(t)
	return err
}

func (a *after) Get() any {
	return a.Value()
}

func (a *after) Value() *time.Time {
	t := time.Time(*a)
	if t.IsZero() {
		return nil
	}
	return &t
}

func (a *after) String() string {
	y, m, d := time.Now().Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	*a = after(date)
	return date.Format("02/Jan/2006")
}

type date time.Time

func (d *date) Set(val string) error {
	t, err := time.Parse("02/Jan/2006", val)
	if err != nil {
		err = errors.New("date format should be dd/mmm/yyyy")
	}
	*d = date(t)
	return err
}

func (d *date) Get() any {
	return d.Value()
}

func (d *date) Value() *time.Time {
	t := time.Time(*d)
	if t.IsZero() {
		return nil
	}
	return &t
}

func (d *date) String() string {
	y, m, dt := time.Now().Date()
	t := time.Date(y, m, dt, 0, 0, 0, 0, time.UTC)
	*d = date(t)
	return t.Format("02/Jan/2006")
}

type stringList []string

func (sl *stringList) Set(val string) error {
	*sl = strings.Split(val, ",")
	return nil
}

func (sl *stringList) Get() any {
	return sl.Value()
}

func (sl *stringList) Value() []string {
	return *sl
}

func (sl *stringList) String() string {
	return ""
}
