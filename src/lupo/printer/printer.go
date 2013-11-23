package printer

import (
	"lupo/event"
	"lupo/out"
	"fmt"
)

const maxPayloadCharsToPrint = 32
const maxPayloadBytesToPrint = 16


func Accept() {
	for {
		select {
		case e:= <-event.Events:
			printEvent(e)
		}
	}
}

// Print the event.
//
// Generic event examples:
// 15:04:05.000  [1    Opened from localhost:23123
// 15:04:05.000 ->1    some text data
// 15:04:05.000 <-10   32 bytes [81 4f d3 c2 ...]
// 15:04:05.000  ]10   Closed
//
// HTTP event examples:
// 15:04:05.000 ->1    GET / HTTP/1.0
// 15:04:05.000 <-1    HTTP/1.0 OK
func printEvent(o interface{}) {
	ev := o.(event.Event)
	out.Stamp(ev.Stamp)
	printKind(ev.Kind)
	out.Out.WriteString(fmt.Sprintf("%-4d ", ev.Cid))

	switch t := o.(type) {
		case event.Event:
			printDesc(&t)
		case event.HttpEvent:
			printHttpDesc(&t)
	}
}

func printKind(k event.EventKind) {
	switch k {
		case event.Connect: out.Out.WriteString(" [")
		case event.Disconnect: out.Out.WriteString(" ]")
		case event.Send: out.Out.WriteString("->")
		case event.Receive: out.Out.WriteString("<-")
	}
}

func printDesc(e *event.Event) {
	switch e.Kind {
		case event.Connect:
			out.Out.WriteString("Opened from ")
			out.Out.Write(e.Payload)
		case event.Disconnect:
			out.Out.WriteString("Closed")
		case event.Send: fallthrough
		case event.Receive:
			printPayload(e.Payload)
	}	
}

func printHttpDesc(e *event.HttpEvent) {
	// Can only be Send or Receive
	out.Out.Write(e.Start)

	// TODO make configurable if headers are printed

	printPayload(e.Body)
}

func printPayload(d []byte) {
	textual := d[:min(len(d), maxPayloadCharsToPrint)]
	if isPrintable(textual) {
		out.WriteWithoutNewlines(textual)
		out.Out.WriteString("\n")
	} else {
		out.Out.WriteString(fmt.Sprintf("%d bytes [", len(d)))
		printBinary(d[:min(len(d), maxPayloadBytesToPrint)])
		out.Out.WriteString(fmt.Sprintf("]\n"))
	}
}

const hextable = "0123456789abcdef"

func printBinary(d []byte) {
	for i, b := range d {
		if i>0 && i%8 == 0 {
			out.Out.WriteString(" ")
		}
		// TODO ugly
		out.Out.WriteString(string(hextable[b>>4]))
		out.Out.WriteString(string(hextable[b&0x0f]))
	}
}

func isPrintable(d []byte) bool {
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
