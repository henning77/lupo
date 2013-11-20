package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"crypto/tls"
	"lupo"
)

func init() {
	flag.StringVar(&from, "from", "", "Source host/port to listen to. Example: ':8081'")
	flag.StringVar(&to, "to", "", "Destination host/port to forward to. Example: 'localhost:8080'")
	flag.BoolVar(&ssl, "ssl", false, "If true, expect and provide SSL/TLS connections. Needs cert.pem + key.pem in the same directory")
}

func main() {
	flag.Parse()
	if len(from) == 0 || len(to) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	lupo.Listen()
}
