package printer

import (
	"lupo/event"
	"lupo/out"
	"encoding/hex"
)

func Full() {
	for {
		select {
		case o := <-event.Events:
			switch ev := o.(type) {
			case *event.DataEvent:
				printDataEventFull(ev)
			case *event.HttpEvent:
				printDataEventFull(&ev.DataEvent)
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

func printDataEventFull(ev *event.DataEvent) {
	out.EntryBegin(ev.Stamp, ev.Kind, ev.Cid, len(ev.Payload))
	printPayloadFull(ev.Payload)
	out.EntryEnd()
}

func printPayloadFull(d []byte) {
	binChunk := false
	chunkStart := 0
	printableCount := 0

	for i, b := range d {
		if isPrintable(b) {
			printableCount++

			// char chunks have a minimum length
			if binChunk && printableCount >= 5 {
				// Write the previous binary chunk
				chunkEnd := i - printableCount + 1
				writeHexChunk(d[chunkStart:chunkEnd])

				binChunk = false
				chunkStart = chunkEnd
			}
		} else {
			printableCount = 0

			if !binChunk {
				// Write the previous char chunk
				out.Out.Write(d[chunkStart:i])

				binChunk = true
				chunkStart = i
			}
		}
	}

	// Final chunk
	if binChunk {
		writeHexChunk(d[chunkStart:])
	} else {
		out.Out.Write(d[chunkStart:])
	}
}

func writeHexChunk(p []byte) {
	dumper := hex.Dumper(out.Out)
	dumper.Write(p)
	dumper.Close()
}

func isPrintable(b byte) bool {
	return b == 0x0d || b == 0x0a || (b >= 0x20 && b <= 0x7e)
}
