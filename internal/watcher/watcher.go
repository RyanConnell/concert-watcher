package watcher

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/RyanConnell/concert-watcher/pkg/discord"
	"github.com/RyanConnell/concert-watcher/pkg/set"
	"github.com/RyanConnell/concert-watcher/pkg/ticketmaster"
)

type Watcher struct {
	ticketmasterAPI        ticketmaster.API
	ticketmasterConfigFile string
	trackedArtists         *set.Set[string]
	discordWebhookURL      string
	diffFile               string
}

func NewWatcher(ticketmasterAPI ticketmaster.API, ticketmasterConfigFile string,
	trackedArtists []string, discordWebhookURL string, diffFile string) *Watcher {
	return &Watcher{
		ticketmasterAPI:        ticketmasterAPI,
		ticketmasterConfigFile: ticketmasterConfigFile,
		trackedArtists:         set.New(trackedArtists...),
		discordWebhookURL:      discordWebhookURL,
		diffFile:               diffFile,
	}
}

func (w *Watcher) FindEvents(diffMode, includePartialMatch bool) ([]*ticketmaster.Event,
	[]string, error) {

	searchCriteria, err := NewSearchCriteria(w.ticketmasterConfigFile)
	if err != nil {
		return nil, nil, err
	}
	events, err := w.ticketmasterAPI.GetEvents(searchCriteria)
	if err != nil {
		return nil, nil, err
	}

	seen := set.New[string]()
	if diffMode {
		previousIDs, err := w.previousEventIDs()
		if err != nil {
			return nil, nil, err
		}
		seen.Add(previousIDs...)
	}

	// Filter events based on our artist list.
	var matchingEvents []*ticketmaster.Event
	var matchingEventIDs []string
	var partialEventIDs []string
	for _, event := range events {
		if seen.Contains(event.ID) {
			// Keep track of which events are showing up again so that our 'seen' file will be
			// pruned automatically.
			matchingEventIDs = append(matchingEventIDs, event.ID)
			continue
		}
		for _, artist := range event.Embedded.Attractions {
			// Check if we can find an exact match for an artist.
			if w.trackedArtists.Contains(artist.Name) {
				matchingEvents = append(matchingEvents, event)
				matchingEventIDs = append(matchingEventIDs, event.ID)
				break
			}
			if !includePartialMatch {
				continue
			}

			var partialMatch bool
			for _, trackedArtist := range w.trackedArtists.Values() {
				if strings.Contains(artist.Name, trackedArtist) {
					matchingEvents = append(matchingEvents, event)
					partialEventIDs = append(partialEventIDs, event.ID)
					partialMatch = true
					break
				}
			}
			if partialMatch {
				break
			}
		}
	}

	if diffMode {
		if err := w.saveEventIDs(append(matchingEventIDs, partialEventIDs...)); err != nil {
			return nil, nil, err
		}
	}
	return matchingEvents, partialEventIDs, nil
}

func (w *Watcher) Notify(events []*ticketmaster.Event, partialEventIDs *set.Set[string]) error {
	if w.discordWebhookURL == "" {
		fmt.Println("No discord webhook URL was provided; Skipping POST to discord")
		return nil
	}
	if len(events) == 0 {
		return nil
	}

	webhookBody := discord.WebhookBody{
		Username: "Concert Watcher",
		Content:  "Take a look at your weekly dump of upcoming concerts!",
	}

	for _, event := range events {
		fieldName := fmt.Sprintf("%s", event.Embedded.Attractions[0].Name)

		embed := &discord.WebhookEmbed{
			Color: 24576,
			Footer: discord.WebhookEmbedFooter{
				Text: fmt.Sprintf("%s @ %s", event.Date(), event.Embedded.Venues[0]),
			},
		}
		if partialEventIDs.Contains(event.ID) {
			embed.Color = 12535447
			embed.Title = fmt.Sprintf("__Partial Match__ (on %q)",
				event.Embedded.Attractions[0].Name)
		}

		// Add information about band and supporting acts.
		embed.Fields = append(embed.Fields, discord.WebhookEmbedField{
			Name:  fieldName,
			Value: strings.Join(event.SupportingActs(), ",\n"),
		})

		// Add a purchasing link.
		embed.Fields = append(embed.Fields, discord.WebhookEmbedField{
			Value: fmt.Sprintf("[Click here for tickets](%s)", event.URL),
		})

		// Add a thumbnail image.
		if len(event.Images) != 0 {
			embed.Thumbnail = discord.URL{URL: event.Images[0].URL}
		}

		webhookBody.Embeds = append(webhookBody.Embeds, embed)
	}

	webhook := &discord.Webhook{URL: w.discordWebhookURL, Body: webhookBody}
	return webhook.Send()
}

// Query for a list of events that we've previously notiied on.
func (w *Watcher) previousEventIDs() ([]string, error) {
	file, err := os.Open(w.diffFile)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var ids []string
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}
	return ids, nil
}

// Save seen event IDs so that we don't send duplicate notifications each time we run.
func (w *Watcher) saveEventIDs(ids []string) error {
	file, err := os.Create(w.diffFile)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(strings.Join(ids, "\n") + "\n")
	return nil
}
