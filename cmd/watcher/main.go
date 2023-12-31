package main

import (
	"log"

	"github.com/RyanConnell/concert-watcher/internal/commands"

	"github.com/alecthomas/kong"
)

// CLI contains all of our parsers.
var CLI struct {
	Scan commands.ScanCmd `cmd:"" help:"Scan Ticketmaster for matching events"`
}

type flags struct {
	apiKey                 string
	artistFile             string
	ticketmasterConfigFile string
	discordWebhookURL      string
	diffMode               bool
	diffFile               string
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "scan":
		CLI.Scan.Run()
	default:
		log.Fatalf("unknown command: %s", ctx.Command())
	}
}
