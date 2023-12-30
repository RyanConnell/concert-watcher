package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/RyanConnell/concertwatch/internal/watcher"
	"github.com/RyanConnell/concertwatch/pkg/ticketmaster"
)

type flags struct {
	apiKey            string
	artistFile        string
	discordWebhookURL string
}

func main() {
	flags := parseFlags()

	// Read artists
	artists, err := readLines(flags.artistFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return
	}

	reader := ticketmaster.NewReader(flags.apiKey)
	watcher := watcher.NewWatcher(reader, artists, flags.discordWebhookURL)

	events, err := watcher.FindEvents()
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
