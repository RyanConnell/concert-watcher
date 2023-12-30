package ticketmaster

import (
	"fmt"
	"strings"
)

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
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	URL      string       `json:"url"`
	Images   []EventImage `json:"images"`
	Embedded struct {
		Venues      []Venue      `json:"venues"`
		Attractions []Attraction `json:"attractions"`
	} `json:"_embedded"`
	Dates struct {
		Start *Date `json:"start"`
	} `json:"dates"`
}

type EventImage struct {
	URL string `json:"url"`
}

func (e *Event) Date() string {
	return e.Dates.Start.LocalDate
}

func (e *Event) SupportingActs() []string {
	support := make([]string, len(e.Embedded.Attractions[1:]))
	for i, attr := range e.Embedded.Attractions[1:] {
		support[i] = attr.Name
	}
	return support
}

func (e *Event) String() string {
	str := fmt.Sprintf("%v: %q are playing", e.Date(), e.Embedded.Attractions[0].Name)
	if len(e.Embedded.Attractions) > 1 {
		str += fmt.Sprintf(" with %s", strings.Join(e.SupportingActs(), ", "))
	}
	str += fmt.Sprintf(" in %q", e.Embedded.Venues[0].Name)
	return str
}
