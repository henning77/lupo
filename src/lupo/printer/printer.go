package filter

import (
	"lupo/event"
	"lupo/out"
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
func printEvent(e *event.Event) {
	out.Stamp(e.Stamp)
	printKind(e.Kind)
	out.Out.WriteString(fmt.Sprintf("%-4d ", e.Cid))
	e.printDesc()
}

func printKind(k EventKind) {
	switch k {
		case Connect: out.Out.WriteString(" [")
		case Disconnect: out.Out.WriteString(" ]")
		case Send: out.Out.WriteString("->")
		case Receive: out.Out.WriteString("<-")
	}
}

func (e *event.HttpEvent) printDesc() {
	// Can only be Send or Receive
	out.Out.Write(e.Start)

	// TODO make configurable if headers are printed

	printPayload(e.Body)
}

func (e *event.Event) printDesc() {
	switch k {
		case Connect:
			out.Out.WriteString("Opened from ")
			out.Out.Write(e.Payload)
		case Disconnect:
			out.Out.WriteString("Closed")
		case Send: fallthrough
		case Receive:
			printPayload(e.Payload)
	}	
}

func printPayload(d []byte) []byte {
	/*
	switch t := e.(type) {
		case event.HttpEvent:

		case event.Event:

	}*/
	textual := d[:math.Min(len(d), maxPayloadCharsToPrint)]
	if isPrintable(textual) {
		out.Out.WriteWithoutNewlines()
		out.Out.WriteString("\n")
	} else {
		out.Out.WriteString(fmt.Sprintf("%d bytes [", len(d)))
		printBinary(d[:math.Min(len(d), maxPayloadBytesToPrint)])
		out.Out.WriteString(fmt.Sprintf("]\n"))
	}
}

const hextable = "0123456789abcdef"

func printBinary(d []byte) {
	for i, b := range d {
		if i>0 && i%8 == 0 {
			out.Out.WriteString(" ")
		}
		out.Out.Write(hextable[b>>4])
		out.Out.Write(hextable[b&0x0f])
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
