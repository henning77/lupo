package main

import (
	"flag"
	"lupo/printer"
	"lupo/server"
	"os"
	"fmt"
)

const version = "v0.2.0"

var style string

func appHeader() {
	fmt.Printf("lupo %v\n\n", version)
}

func usageExamples() {
	fmt.Println("Examples:")
	fmt.Println("  lupo -from :8080 -to myserver:80")
	fmt.Println("  lupo -from :8443 -to myserver:443 -ssl")
	fmt.Println("  lupo -style full -from :8080 -to myserver:80")	
	fmt.Println("  lupo -style stats -from :8080 -to myserver:80")
}

func init() {
	flag.StringVar(&style, "style", "short", "Style of output. 'short' prints one line per transfer. 'full' prints everything. 'stats' prints a csv with stats")
	flag.StringVar(&server.From, "from", server.From, "Source host/port to listen to")
	flag.StringVar(&server.To, "to", server.To, "Destination host/port to forward to")
	flag.BoolVar(&server.Ssl, "ssl", server.Ssl, "If true, expect and provide SSL/TLS connections. Needs cert.pem + key.pem in the same directory")
	//flag.IntVar(&printer.MaxPayloadCharsToPrint, "trunc", printer.MaxPayloadCharsToPrint, "Truncate content (number of chars to print)")
	//flag.BoolVar(&printer.Headers, "headers", printer.Headers, "Print HTTP headers")
	// TODO out.Out
}

func main() {
	flag.Parse()
	if len(server.From) == 0 || len(server.To) == 0 {
		appHeader()
		flag.Usage()
		fmt.Println()
		usageExamples()
		os.Exit(1)
	}

	if style == "short" {
		go printer.Short()
	} else if style == "full" {
		go printer.Full()
	} else if style == "stats" {
		go printer.Stats()
	} else {
		panic("Unknown style")
	}
	
	server.Listen()
}
