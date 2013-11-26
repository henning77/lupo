package main

import (
	"flag"
	"log"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"crypto/tls"
	"os/signal"
)

// Host:port to listen from
var from string

// Host:port to forward to
var to string

// If true, use SSL/TLS connections
var ssl bool

// Interval to collect (ms)
var interval int

// Count the connections to make them easily identifiable
var nextConnId = make(chan int)

// Timestamp format for logging
const stampf = "02.01.2006 15:04:05.000"

var events = make(chan event, 1000)

type eventKind byte
const (
	_                 = iota
	connect eventKind = iota
	disconnect
	send
	receive
)
type event struct {
	kind eventKind
	data uint32
}

type stats struct {
	stamp time.Time
	connCount uint32
	connOpened uint32
	connClosed uint32
	sent uint32
	received uint32
}

func statsCollector() {
	stats := stats{stamp: time.Now()}
	ticker := time.Tick(time.Duration(interval) * time.Millisecond)
	for {
		select {
			case e := <-events:
				switch e.kind {
					case connect:
						stats.connCount++
						stats.connOpened++
					case disconnect:
						stats.connCount--
						stats.connClosed++
					case send:
						stats.sent += e.data
					case receive:
						stats.received += e.data					
				}
			case <-ticker:
				printStats(&stats)
				resetStats(&stats)
		}
	}
}

func resetStats(s *stats) {
	// Note: connCount is not reset, as we want to track the absolute value
	s.stamp = time.Now()
	s.connOpened = 0
	s.connClosed = 0
	s.sent = 0
	s.received = 0
}

func printStatHeader() {
	os.Stdout.WriteString("Date;ConnCount;ConnOpened;ConnClosed;Sent;Received;TotalTransferred\n")
}

func printStats(s *stats) {
	os.Stdout.WriteString(s.stamp.Format(stampf))
	fmt.Printf(";%d;%d;%d;%d;%d;%d\n", s.connCount, s.connOpened, s.connClosed, s.sent, s.received, s.sent + s.received)
}

type statsWriter struct {
	direction eventKind
}

func (s *statsWriter) Write(p []byte) (n int, err error) {
	events <- event{kind:s.direction, data:uint32(len(p))}
	return len(p), nil
}

// Copy to dst and to logger
func copyWithLog(dst io.Writer, src io.Reader, direction eventKind, done chan<- bool) {
	logger := &statsWriter{direction:direction}
	multi := io.MultiWriter(dst, logger)
	if _, err := io.Copy(multi, src); err != nil {
		//logPrintf("Closed connection: %v", err)
	}
	done <- true	
}

func handleConnection(src net.Conn, closeConn <-chan bool, connClosed chan<- bool) {
	//defer os.Stdout.WriteString("Closed conn")
	defer func() { connClosed <- true } ()
	defer src.Close()

	events <- event{kind:connect}

	var dst net.Conn
	var err error
	if ssl {
		dst, err = tls.Dial("tcp", to, tlsClientConfig())
	} else {
		dst, err = net.Dial("tcp", to)
	}
	
	if err != nil {
		log.Panicf("Error connecting to dest: %v", err)
	}
	defer dst.Close()

	done := make(chan bool)

	// Copy & log in both directions
	go copyWithLog(dst, src, send, done)
	go copyWithLog(src, dst, receive, done)

	openStreams := 2

HandlerLoop:
	for {
		select {
		case <-done:
			if openStreams == 2 {
				// First stream has been closed
				events <- event{kind:disconnect}
				openStreams--
			} else {
				// Second stream has been closed -> Connection is closed, can exit
				break HandlerLoop
			}

		case <-closeConn:
			// Server signals we should close
			break HandlerLoop
		}
	}
}

func tlsServerConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Panicf("Could not read server certificate (cert.pem, key.pem): %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func tlsClientConfig() *tls.Config {
	// Simply accept everything
	return &tls.Config{InsecureSkipVerify: true}
}

func listen(c chan<- net.Conn) {
	var ln net.Listener
	var err error
	if ssl {
		ln, err = tls.Listen("tcp", from, tlsServerConfig())
	} else {
		ln, err = net.Listen("tcp", from)
	}
	
	if err != nil {
		log.Panicf("Could not open port: %v", err)
	}

	for  {
		conn, err := ln.Accept()
		if err != nil {
			log.Panicf("Error accepting: %v", err)
		}
		c<-conn
	}
}

func server() {
	// Listen for SIGTERM (kill) and SIGINT (Ctrl+C)
	sigExit := make(chan os.Signal, 2)
	signal.Notify(sigExit, os.Interrupt, os.Kill)

	connCount := 0
	newConn := make(chan net.Conn)
	closeConn := make(chan bool)
	connClosed := make(chan bool)

	go listen(newConn)

ServerLoop:
	for {
		select {
		case c := <-newConn:
			connCount++
			go handleConnection(c, closeConn, connClosed)
		case <-connClosed:
			connCount--
		case <-sigExit:
			break ServerLoop
		}
	}

	// Signal close to all connection gorountines
	for i:=0; i<connCount; i++ {
		closeConn<-true
	}

	// Wait until all are closed (goroutines signal back)
	for i:=0; i<connCount; i++ {
		<-connClosed
	}
}

func init() {
	flag.StringVar(&from, "from", "", "Source host/port to listen to. Example: ':8081'")
	flag.StringVar(&to, "to", "", "Destination host/port to forward to. Example: 'localhost:8080'")
	flag.BoolVar(&ssl, "ssl", false, "If true, expect and provide SSL/TLS connections. Needs cert.pem + key.pem in the same directory")
	flag.IntVar(&interval, "interval", 1000, "Interval in which to collect (ms)")
}

func main() {	
	flag.Parse()
	if len(from) == 0 || len(to) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	printStatHeader()

	go statsCollector()
	server()

	os.Exit(0)
}
