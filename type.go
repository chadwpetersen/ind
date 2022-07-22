package ind

const baseURL = "https://oap.ind.nl/oap/api/desks"

type Venue string

const (
	VenueAmsterdam Venue = "AM"
	VenueDenHaag   Venue = "DH"
	VenueZwolle    Venue = "ZW"
	VenueDenBosch  Venue = "DB"
)

var venues = map[Venue]string{
	VenueAmsterdam: "Amsterdam",
	VenueDenHaag:   "Den Haag",
	VenueZwolle:    "Zwolle",
	VenueDenBosch:  "Den Bosch",
}

func (v Venue) String() string {
	name, ok := venues[v]
	if !ok {
		return ""
	}

	return name
}

func (v Venue) Code() string {
	if !v.IsValid() {
		return ""
	}

	return string(v)
}

func (v Venue) IsValid() bool {
	_, ok := venues[v]
	return ok
}

type Strategy int

const (
	StrategyIndividual Strategy = 1
	StrategyTogether   Strategy = 2
)

var strategies = map[Strategy]string{
	StrategyIndividual: "individual",
	StrategyTogether:   "together",
}

func (s Strategy) String() string {
	name, ok := strategies[s]
	if !ok {
		return ""
	}

	return name
}

func (s Strategy) IsValid() bool {
	_, ok := strategies[s]
	return ok
}

type Slot struct {
	Key    string `json:"key"`
	Venue  Venue  `json:"-"`
	Date   string `json:"date"`
	Start  string `json:"startTime"`
	End    string `json:"endTime"`
	Parts  uint   `json:"parts"`
	Booked bool   `json:"booked,omitempty"`
}

type Appointment struct {
	ProductKey string     `json:"productKey"`
	Date       string     `json:"date"`
	Start      string     `json:"startTime"`
	End        string     `json:"endTime"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Language   string     `json:"language"`
	Customers  []Customer `json:"customers"`
}

type Customer struct {
	VNumber string `json:"vNumber"`
	Name    string `json:"firstName"`
	Surname string `json:"lastName"`
}
