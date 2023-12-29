package ticketmaster

import "fmt"

type Attraction struct {
	Name string `json:"name"`
}

type Venue struct {
	Name string `json:"name"`
}

type Date struct {
	LocalDate string `json:"localDate"`
	LocalTime string `json:"localTime"`
}

type Event struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
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

func (e *Event) String() string {
	artists := e.Embedded.Attractions
	str := fmt.Sprintf("%v: %q are playing", e.Date(), artists[0].Name)
	if len(artists) > 1 {
		str += fmt.Sprintf(" with %q", artists[1].Name)
		for i, artist := range artists[2:] {
			if i == len(artists[2:])-1 {
				str += fmt.Sprintf(", and %q", artist.Name)
			} else {
				str += fmt.Sprintf(", %q", artist.Name)
			}
		}
	}
	str += fmt.Sprintf(" in %q", e.Embedded.Venues[0].Name)
	return str
}
