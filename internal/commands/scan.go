package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/RyanConnell/concert-watcher/internal/watcher"
	"github.com/RyanConnell/concert-watcher/pkg/set"
	"github.com/RyanConnell/concert-watcher/pkg/ticketmaster"
)

// ScanCmd defines arguments for the "./concert-watcher scan" sub-command.
type ScanCmd struct {
	APIKey              string `help:"Ticketmaster API Key"`
	ArtistFile          string `help:"Path to file containing a list of artists to monitor"`
	DiscordWebhookURL   string `help:"Discord Webhook URL"`
	TicketmasterConfig  string `help:"Path to file containing ticketmaster search criteria"`
	Diff                bool   `help:"Only notify for events we haven't already seen"`
	DiffFile            string `default:"previous-ids" help:"Path to a file that stores the list of events we have previously sent notifications for"`
	IncludePartialMatch bool   `help:"Includes any events that partially match the list of artists"`
}

func (s *ScanCmd) Run() {
	// Ticketmaster will give us an authentication error if we don't have an API key.
	if s.APIKey == "" {
		log.Fatalf("Unable to run without a Ticketmaster API key. " +
			"Please pass one via the --apiKey option")
		return
	}

	// Read artists
	artists, err := readLines(s.ArtistFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return
	}

	reader := ticketmaster.NewReader(s.APIKey)
	watcher := watcher.NewWatcher(reader, s.TicketmasterConfig, artists,
		s.DiscordWebhookURL, s.DiffFile)
	events, partialEventIDs, err := watcher.FindEvents(s.Diff, s.IncludePartialMatch)
	if err != nil {
		log.Fatalf("Error retrieving events: %v", err)
		return
	}

	// Sort our events by date
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date() < events[j].Date()
	})

	partialEventSet := set.New(partialEventIDs...)
	fmt.Printf("\nFound %d matching events\n", len(events))
	for _, event := range events {
		var partialStr string
		if partialEventSet.Contains(event.ID) {
			partialStr = " [Partial Match]"
		}
		fmt.Printf("- %s%s\n", event.String(), partialStr)
	}

	if err := watcher.Notify(events, partialEventSet); err != nil {
		log.Fatalf("Error notifying discord: %v", err)
		return
	}
}

func readLines(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}
