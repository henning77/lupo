package server

import (
	"crypto/tls"
	"fmt"
	"lupo/event"
	"lupo/handler"
	"lupo/out"
	"net"
	"os"
)

// Host:port to listen from
var From string

// Host:port to forward to
var To string

// If true, use SSL/TLS connections
var Ssl bool

// Count the connections to make them easily identifiable
var nextConnId = make(chan event.ConnId)

func init() {
	go genConnectionIds()
}

func genConnectionIds() {
	var i event.ConnId = 1
	for {
		nextConnId <- i
		i++
	}
}

func tlsServerConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		out.Printf("Could not read server certificate (cert.pem, key.pem): %v", err)
		os.Exit(1)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func tlsClientConfig() *tls.Config {
	// Simply accept everything
	return &tls.Config{InsecureSkipVerify: true}
}

func handleConnection(src net.Conn) {
	cid := <-nextConnId
	out.Print(fmt.Sprintf("New connection: %v (from %v)", cid, src.RemoteAddr()))

	var dst net.Conn
	var err error
	if Ssl {
		dst, err = tls.Dial("tcp", To, tlsClientConfig())
	} else {
		dst, err = net.Dial("tcp", To)
	}

	if err != nil {
		out.Printf("Error connecting to dest: %v", err)
		panic(err)
	}

	handler := handler.NewHandler(dst, src, cid)
	handler.Handle()
}

func Listen() {
	out.Printf("Listening to [%v], forwarding to [%v]", From, To)

	var ln net.Listener
	var err error
	if Ssl {
		ln, err = tls.Listen("tcp", From, tlsServerConfig())
	} else {
		ln, err = net.Listen("tcp", From)
	}

	if err != nil {
		out.Printf("Could not open port: %v", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			out.Printf("Error accepting: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
