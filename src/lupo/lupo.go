package main

import (
	"flag"
	"lupo/printer"
	"lupo/server"
	"os"
)

func init() {
	flag.StringVar(&server.From, "from", server.From, "Source host/port to listen to. Example: ':8081'")
	flag.StringVar(&server.To, "to", server.To, "Destination host/port to forward to. Example: 'localhost:8080'")
	flag.BoolVar(&server.Ssl, "ssl", server.Ssl, "If true, expect and provide SSL/TLS connections. Needs cert.pem + key.pem in the same directory")
	flag.IntVar(&printer.MaxPayloadCharsToPrint, "trunc", printer.MaxPayloadCharsToPrint, "Truncate content (number of chars to print)")
	flag.BoolVar(&printer.Headers, "headers", printer.Headers, "Print HTTP headers")
	// TODO out.Out
}

func main() {
	flag.Parse()
	if len(server.From) == 0 || len(server.To) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	go printer.Accept()
	server.Listen()
}
