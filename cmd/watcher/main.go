package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/RyanConnell/concert-watcher/internal/watcher"
	"github.com/RyanConnell/concert-watcher/pkg/ticketmaster"
)

type flags struct {
	apiKey                 string
	artistFile             string
	ticketmasterConfigFile string
	discordWebhookURL      string
	diffMode               bool
	diffFile               string
}

func main() {
	f := parseFlags()

	// Ticketmaster will give us an authentication error if we don't have an API key.
	if f.apiKey == "" {
		log.Fatalf("Unable to run without a Ticketmaster API key. " +
			"Please pass one via the --apiKey option")
		return
	}

	// Read artists
	artists, err := readLines(f.artistFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return
	}

	reader := ticketmaster.NewReader(f.apiKey)
	watcher := watcher.NewWatcher(reader, f.ticketmasterConfigFile, artists,
		f.discordWebhookURL, f.diffFile)
	events, err := watcher.FindEvents(f.diffMode)
	if err != nil {
		log.Fatalf("Error retrieving events: %v", err)
		return
	}

	// Sort our events by date
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date() < events[j].Date()
	})
	fmt.Printf("Found %d matching events\n", len(events))
	for _, event := range events {
		fmt.Printf("- %s\n", event.String())
	}

	if err := watcher.Notify(events); err != nil {
		log.Fatalf("Error notifying discord: %v", err)
		return
	}
}

func parseFlags() *flags {
	f := &flags{}
	flag.StringVar(&f.apiKey, "apiKey", "", "Ticketmaster API Key")
	flag.StringVar(&f.artistFile, "artistFile", "artists", "Path to a file containing a list of "+
		"artists to search for")
	flag.StringVar(&f.discordWebhookURL, "discordWebhookURL", "", "Discord webhook URL")
	flag.StringVar(&f.ticketmasterConfigFile, "ticketmasterConfig", "ticketmaster.yaml",
		"Path to a file containing search criteria for ticketmaster event lookups")
	flag.BoolVar(&f.diffMode, "diff", false, "Run in diff mode which only notifies "+
		"for events we haven't already seen.")
	flag.StringVar(&f.diffFile, "diffFile", "previous-ids", "Path to a file that stores the list "+
		"of events that we have previously sent notifications for")
	flag.Parse()
	return f
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
