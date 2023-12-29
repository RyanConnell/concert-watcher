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
	url := fmt.Sprintf(getEventsURL, r.apiKey)
	for key, value := range params {
		url += fmt.Sprintf("&%s=%s", key, value)
	}

	// TODO: Offset by startDate=events[-1].date to avoid the 1k results limit.
	var events []*Event
	for {
		resp := &getEventsResponse{}
		err := query(url, resp)
		if err != nil {
			return nil, err
		}
		events = append(events, resp.Embedded.Events...)

		next, ok := resp.Links["next"]
		if !ok {
			break
		}
		url = fmt.Sprintf("%s%s&apikey=%s", apiURL, next.URL, r.apiKey)
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

	if err := json.NewDecoder(resp.Body).Decode(&target); err != nil {
		return err
	}

	return nil
}
