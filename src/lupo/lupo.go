package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

// Host:port to listen from
var from string

// Host:port to forward to
var to string

// Count the connections to make them easily identifiable
var nextConnId = make(chan int)

// Timestamp format for logging
const stamp = "15:04:05.000"

// Special Writer which writes head + nicely formatted binary / textual chunks + tail
type transferLog struct {
	head string
	tail string
}

// Write nicely formatted binary / textual chunk
func (l *transferLog) Write(p []byte) (n int, err error) {
	logPrint(l.head)

	binChunk := false
	chunkStart := 0
	printableCount := 0

	for i, b := range p {
		if isPrintable(b) {
			printableCount++

			// char chunks have a minimum length
			if binChunk && printableCount >= 5 {
				// Write the previous binary chunk
				chunkEnd := i - printableCount + 1
				writeHexChunk(p[chunkStart:chunkEnd])

				binChunk = false
				chunkStart = chunkEnd
			}
		} else {
			printableCount = 0

			if !binChunk {
				// Write the previous char chunk
				os.Stdout.Write(p[chunkStart:i])

				binChunk = true
				chunkStart = i
			}
		}
	}

	// Final chunk
	if binChunk {
		writeHexChunk(p[chunkStart:])
	} else {
		os.Stdout.Write(p[chunkStart:])
	}

	os.Stdout.WriteString(l.tail)
	return len(p), nil
}

func writeHexChunk(p []byte) {
	dumper := hex.Dumper(os.Stdout)
	dumper.Write(p)
	dumper.Close()
}

func isPrintable(b byte) bool {
	return b == 0x0d || b == 0x0a || (b >= 0x20 && b <= 0x7e)
}

// Log with optimized timestamp
func logPrintf(s string, a ...interface{}) {
	os.Stdout.WriteString(time.Now().Format(stamp))
	os.Stdout.WriteString(" ")
	os.Stdout.WriteString(fmt.Sprintf(s, a...))
	os.Stdout.WriteString("\n")
}

// Log with optimized timestamp
func logPrint(s string) {
	os.Stdout.WriteString(time.Now().Format(stamp))
	os.Stdout.WriteString(" ")
	os.Stdout.WriteString(s)
	os.Stdout.WriteString("\n")
}

// Copy to dst and to logger
func copyWithLog(dst io.Writer, src io.Reader, head string, tail string) {
	logger := &transferLog{head: head, tail: tail}
	multi := io.MultiWriter(dst, logger)
	if _, err := io.Copy(multi, src); err != nil {
		logPrintf("Error copying: %v", err)
	}
}

func handleConnection(src net.Conn) {
	connId := <-nextConnId
	logPrint(fmt.Sprintf("New connection: %v", connId))

	dst, err := net.Dial("tcp", to)
	if err != nil {
		logPrintf("Error connecting to dest: %v", err)
		panic(err)
	}

	// Copy & log in both directions
	go copyWithLog(dst, src, fmt.Sprintf("->%v", connId), "\n")
	go copyWithLog(src, dst, fmt.Sprintf("<-%v", connId), "\n")
}

func genConnectionIds() {
	i := 1
	for {
		nextConnId <- i
		i++
	}
}

func init() {
	flag.StringVar(&from, "from", ":8081", "Source host/port to listen to")
	flag.StringVar(&to, "to", "localhost:8080", "Destination host/port to forward to")
}

func main() {
	flag.Parse()

	logPrintf("Listening to [%v], forwarding to [%v]", from, to)

	ln, err := net.Listen("tcp", from)
	if err != nil {
		logPrintf("Could not open port: %v", err)
		os.Exit(1)
	}

	go genConnectionIds()

	for {
		conn, err := ln.Accept()
		if err != nil {
			logPrintf("Error accepting: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}
