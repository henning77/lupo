package server

import (
	"crypto/tls"
	"lupo/event"
	"lupo/handler"
	"lupo/out"
	"net"
)

// Host:port to listen from
var From string = ""

// Host:port to forward to
var To string = ""

// If true, use SSL/TLS connections
var Ssl bool = false

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
		out.Fatalf("Could not read server certificate (cert.pem, key.pem): %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func tlsClientConfig() *tls.Config {
	// Simply accept everything
	return &tls.Config{InsecureSkipVerify: true}
}

func handleConnection(src net.Conn) {
	cid := <-nextConnId

	var dst net.Conn
	var err error
	if Ssl {
		dst, err = tls.Dial("tcp", To, tlsClientConfig())
	} else {
		dst, err = net.Dial("tcp", To)
	}

	if err != nil {		
		src.Close()
		out.Fatalf("Error connecting to dest: %v", err)
	}

	handler := handler.NewHandler(dst, src, cid)
	handler.Handle()
}

func Listen() {
	event.PostGlobalf("Listening to [%v], forwarding to [%v]", From, To)

	var ln net.Listener
	var err error
	if Ssl {
		ln, err = tls.Listen("tcp", From, tlsServerConfig())
	} else {
		ln, err = net.Listen("tcp", From)
	}

	if err != nil {
		out.Fatalf("Could not open port: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}
