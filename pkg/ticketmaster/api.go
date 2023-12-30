package ticketmaster

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiURL       = "https://app.ticketmaster.com"
	getEventsURL = apiURL + "/discovery/v2/events?apikey=%s"
)

type API interface {
	GetEvents(map[string]string) ([]*Event, error)
}

// Ensure that our Reader implements the API interface.
var _ API = (*Reader)(nil)

type Link struct {
	URL string `json:"href"`
}

type getEventsResponse struct {
	Links    map[string]Link `json:"_links"`
	Embedded struct {
		Events []*Event `json:"events"`
	} `json:"_embedded"`
}

type Reader struct {
	apiKey string
}

func NewReader(apiKey string) *Reader {
	return &Reader{apiKey: apiKey}
}

func (r *Reader) GetEvents(params map[string]string) ([]*Event, error) {
	// Construct the URL from the params
	rootURL := fmt.Sprintf(getEventsURL, r.apiKey)
	for key, value := range params {
		rootURL += fmt.Sprintf("&%s=%s", key, value)
	}

	url := rootURL
	var events []*Event
	for {
		// The ticketmaster API only lets us paginate until we hit 1000 results. To circumvent this
		// we will make use of the startDate of the last event we found to start a new search.
		var eventsFromPagination int
		for {
			resp := &getEventsResponse{}
			err := query(url, resp)
			if err != nil {
				return nil, err
			}
			if len(resp.Embedded.Events) == 0 {
				break
			}
			events = append(events, resp.Embedded.Events...)
			eventsFromPagination += len(resp.Embedded.Events)

			next, ok := resp.Links["next"]
			if !ok {
				break
			}

			// If we continue to query after we've retrieved 1000 results we'll hit an error.
			if eventsFromPagination >= 1000 {
				break
			}
			url = fmt.Sprintf("%s%s&apikey=%s", apiURL, next.URL, r.apiKey)
		}

		// If we haven't found a multiple of 1000 then we weren't limited by ticketmaster and
		// instead simply ran out of events.
		if len(events)%1000 != 0 {
			break
		}

		url = fmt.Sprintf("%s&startDateTime=%s", rootURL, events[len(events)-1].Dates.Start.DateTime)
	}

	return events, nil
}

func query[T any](url string, target T) error {
	fmt.Printf("Visiting site: %q\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%q returned %q", url, resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}
