package watcher

import (
	"github.com/RyanConnell/concertwatch/pkg/set"
	"github.com/RyanConnell/concertwatch/pkg/ticketmaster"
)

type Watcher struct {
	ticketmasterAPI ticketmaster.API
	trackedArtists  *set.Set[string]
}

func NewWatcher(ticketmasterAPI ticketmaster.API, trackedArtists []string) *Watcher {
	return &Watcher{
		ticketmasterAPI: ticketmasterAPI,
		trackedArtists:  set.New(trackedArtists...),
	}
}

func (w *Watcher) FindEvents() ([]*ticketmaster.Event, error) {
	searchCriteria := map[string]string{
		"size":               "200",   // Maximum size allowed by ticketmaster API.
		"countryCode":        "IE",    // Filter results to Ireland.
		"classificationName": "music", // Filter results to concerts.
	}
	events, err := w.ticketmasterAPI.GetEvents(searchCriteria)
	if err != nil {
		return nil, err
	}

	// Filter events based on our artist list.
	var matchingEvents []*ticketmaster.Event
	for _, event := range events {
		for _, artist := range event.Embedded.Attractions {
			if w.trackedArtists.Contains(artist.Name) {
				matchingEvents = append(matchingEvents, event)
			}
		}
	}

	return matchingEvents, nil
}
