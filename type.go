package ind

const baseURL = "https://oap.ind.nl/oap/api/desks"

type Venue string

const (
	VenueAmsterdam Venue = "AM"
	VenueDenHaag   Venue = "DH"
	VenueZwolle    Venue = "ZW"
	VenueDenBosch  Venue = "DB"
)

func (v Venue) String() string {
	return string(v)
}

type Strategy int

const (
	StrategyIndividual Strategy = 1
	StrategyTogether   Strategy = 2
)

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
