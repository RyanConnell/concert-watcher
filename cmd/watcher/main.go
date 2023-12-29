package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/RyanConnell/concertwatch/pkg/ticketmaster"
)

type flags struct {
	apiKey     string
	artistFile string
}

func main() {
	flags := parseFlags()

	// Read artists
	artists, err := readLines(flags.artistFile)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return
	}
	artistSet := make(map[string]struct{})
	for _, artist := range artists {
		artistSet[strings.ToLower(artist)] = struct{}{}
	}

	// Search for events
	searchCriteria := map[string]string{
		"size":               "200", // Limit response size to 200 events. (max supported by API)
		"countryCode":        "IE",  // Filter results to Ireland
		"classificationName": "music",
	}
	reader := ticketmaster.NewReader(flags.apiKey)
	events, err := reader.GetEvents(searchCriteria)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	// Filter events based on our artist list.
	fmt.Printf("Found %d events\n", len(events))
	for _, event := range events {
		for _, attr := range event.Embedded.Attractions {
			if _, ok := artistSet[strings.ToLower(attr.Name)]; !ok {
				continue
			}

			fmt.Printf("%v: %s are playing ", event.Date(), attr.Name)
			if len(event.Embedded.Attractions) != 1 {
				fmt.Printf("with %v ", event.Embedded.Attractions)
			}
			fmt.Printf("at %v\n", event.Embedded.Venues)
		}
	}
}

func parseFlags() *flags {
	f := &flags{}
	flag.StringVar(&f.apiKey, "apiKey", "", "Ticketmaster API Key")
	flag.StringVar(&f.artistFile, "artistFile", "artists", "Path to a file containing a list of "+
		"artists to search for")
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
