package watcher

import (
	"fmt"
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
}

func NewWatcher(ticketmasterAPI ticketmaster.API, ticketmasterConfigFile string,
	trackedArtists []string, discordWebhookURL string) *Watcher {
	return &Watcher{
		ticketmasterAPI:        ticketmasterAPI,
		ticketmasterConfigFile: ticketmasterConfigFile,
		trackedArtists:         set.New(trackedArtists...),
		discordWebhookURL:      discordWebhookURL,
	}
}

func (w *Watcher) FindEvents() ([]*ticketmaster.Event, error) {
	searchCriteria, err := NewSearchCriteria(w.ticketmasterConfigFile)
	if err != nil {
		return nil, err
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

func (w *Watcher) Notify(events []*ticketmaster.Event) error {
	if w.discordWebhookURL == "" {
		fmt.Println("No discord webhook URL was provided; Skipping POST to discord")
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
