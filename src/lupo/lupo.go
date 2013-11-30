package main

import (
	"flag"
	"lupo/printer"
	"lupo/server"
	"os"
)


// Usage examples
// lupo -from :8080 -to myserver:80
// lupo -from :8443 -to myserver:443 -ssl
// lupo -style full -from :8080 -to myserver:80

var style string

func init() {
	flag.StringVar(&style, "style", "short", "Style of output. 'short' prints one line per transfer. 'full' prints everything. 'stats' prints a csv with stats.")
	flag.StringVar(&server.From, "from", server.From, "Source host/port to listen to. Example: ':8081'")
	flag.StringVar(&server.To, "to", server.To, "Destination host/port to forward to. Example: 'localhost:8080'")
	flag.BoolVar(&server.Ssl, "ssl", server.Ssl, "If true, expect and provide SSL/TLS connections. Needs cert.pem + key.pem in the same directory")
	//flag.IntVar(&printer.MaxPayloadCharsToPrint, "trunc", printer.MaxPayloadCharsToPrint, "Truncate content (number of chars to print)")
	//flag.BoolVar(&printer.Headers, "headers", printer.Headers, "Print HTTP headers")
	// TODO out.Out
}

func main() {
	flag.Parse()
	if len(server.From) == 0 || len(server.To) == 0 {
		flag.Usage()
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
