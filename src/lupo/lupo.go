package main

import (
	"flag"
	"os"
	"lupo/server"
	"lupo/printer"
)

func init() {
	flag.StringVar(&server.From, "from", "", "Source host/port to listen to. Example: ':8081'")
	flag.StringVar(&server.To, "to", "", "Destination host/port to forward to. Example: 'localhost:8080'")
	flag.BoolVar(&server.Ssl, "ssl", false, "If true, expect and provide SSL/TLS connections. Needs cert.pem + key.pem in the same directory")
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
