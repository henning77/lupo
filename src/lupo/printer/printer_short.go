package printer

import (
	"fmt"
	"lupo/event"
	"lupo/out"
)

var MaxPayloadCharsToPrint = 80
var MaxPayloadBytesToPrint = MaxPayloadCharsToPrint / 2
var Headers = false

func Short() {
	for {
		select {
		case o := <-event.Events:
			switch ev := o.(type) {
			case *event.DataEvent:
				printDataEventShort(ev)
			case *event.HttpEvent:
				printHttpEventShort(ev)
			case *event.ConnectEvent:
				printConnectEventShort(ev)
			case *event.DisconnectEvent:
				printDisconnectEventShort(ev)
			case *event.MessageEvent:
				printMessageEventShort(ev)
			default:
				panic("Unexpected event")
			}			
		}
	}
}

func printConnectEventShort(ev *event.ConnectEvent) {
	out.EntryBegin(ev.Stamp, ev.Kind, ev.Cid, 0)
	out.Out.WriteString("New connection from ")
	out.Out.WriteString(ev.From)
	out.EntryEnd()
}

func printDisconnectEventShort(ev *event.DisconnectEvent) {
	out.EntryBegin(ev.Stamp, ev.Kind, ev.Cid, 0)
	if ev.Initiator == event.Send {
		// Send stream was closed -> Client
		out.Out.WriteString("Client closed connection")
	} else {
		out.Out.WriteString("Server closed connection")
	}
	out.EntryEnd()
}

func printMessageEventShort(ev *event.MessageEvent) {
	out.EntryBegin(ev.Stamp, ev.Kind, ev.Cid, 0)
	out.Out.WriteString(ev.Message)
	out.EntryEnd()
}

// Print the event.
//
// Generic event examples:
// 15:04:05.000  [1    Opened from localhost:23123
// 15:04:05.000 ->1    some text data
// 15:04:05.000 <-10   32 bytes [81 4f d3 c2 ...]
// 15:04:05.000  ]10   Closed
func printDataEventShort(ev *event.DataEvent) {
	out.EntryBegin(ev.Stamp, ev.Kind, ev.Cid, len(ev.Payload))
	printPayloadShort(ev.Payload)
	out.EntryEnd()
}

// HTTP event examples:
// 15:04:05.000 ->1    GET / HTTP/1.0
// 15:04:05.000 <-1    HTTP/1.0 OK
func printHttpEventShort(ev *event.HttpEvent) {
	out.EntryBegin(ev.Stamp, ev.Kind, ev.Cid, len(ev.Payload))
	printHttpDescShort(ev)
	out.EntryEnd()
}

func printHttpDescShort(e *event.HttpEvent) {
	// Can only be Send or Receive
	out.Out.Write(e.Start)
	out.Out.WriteString(" ")

	if (Headers) {
		out.Out.WriteString(fmt.Sprintf("%v ", e.Headers))
	}

	printPayloadShort(e.Body)
}

func printPayloadShort(d []byte) {
	textual := d[:min(len(d), MaxPayloadCharsToPrint)]
	if isAllPrintable(textual) {
		out.WriteWithoutNewlines(textual)
		if len(d) > MaxPayloadCharsToPrint {
			out.Out.WriteString("(...)")
		}
	} else {
		printBinary(d[:min(len(d), MaxPayloadBytesToPrint)])		
		if len(d) > MaxPayloadBytesToPrint {
			out.Out.WriteString("(...)")
		}
	}
}

const hextable = "0123456789abcdef"

func printBinary(d []byte) {
	for i, b := range d {
		if i > 0 && i%8 == 0 {
			out.Out.WriteString(" ")
		}
		// TODO ugly
		out.Out.WriteString(string(hextable[b>>4]))
		out.Out.WriteString(string(hextable[b&0x0f]))
	}
}

func isAllPrintable(d []byte) bool {
	for _, b := range d {
		if !(b == 0x0d || b == 0x0a || (b >= 0x20 && b <= 0x7e)) {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

/*
// Special Writer which writes head + nicely formatted binary / textual chunks + tail
type transferLog struct {
	head string
	tail string
}


// Write nicely formatted binary / textual chunk
func (l *transferLog) Write(p []byte) (n int, err error) {
	util.Print(l.head)

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
*/
