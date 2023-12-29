package ticketmaster

type Attraction struct {
	Name string `json:"name"`
}

type Venue struct {
	Name string `json:"name"`
	//Address map[string]string `json:"address"`
}

type Date struct {
	LocalDate string `json:"localDate"`
	LocalTime string `json:"localTime"`
}

type Event struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
	// 'Classifications' would be good
	// 'PriceRanges' would be good too
	Embedded struct {
		Venues      []Venue      `json:"venues"`
		Attractions []Attraction `json:"attractions"`
	} `json:"_embedded"`
	Dates struct {
		Start *Date `json:"start"`
	} `json:"dates"`
}

func (e *Event) Date() string {
	return e.Dates.Start.LocalDate
}
