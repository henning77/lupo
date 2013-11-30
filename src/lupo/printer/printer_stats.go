package printer

import (
	"lupo/event"
	"lupo/out"
	"time"
	"fmt"
)

// Interval to collect (ms)
var Interval = 1000

func Stats() {
	printStatHeader()

	stats := stats{stamp: time.Now()}
	ticker := time.Tick(time.Duration(Interval) * time.Millisecond)

	for {
		select {
		case o := <-event.Events:
			switch ev := o.(type) {
			case *event.DataEvent:
				recordTransfer(&stats, ev)
			case *event.HttpEvent:
				recordTransfer(&stats, &ev.DataEvent)				
			case *event.ConnectEvent:
				stats.connCount++
				stats.connOpened++
			case *event.DisconnectEvent:
				stats.connCount--
				stats.connClosed++
			case *event.MessageEvent:
				// Do nothing
			default:
				panic("Unexpected event")
			}
		case <-ticker:
			printStats(&stats)
			resetStats(&stats)
		}
	}
}

type stats struct {
	stamp time.Time
	connCount uint32
	connOpened uint32
	connClosed uint32
	sent uint32
	received uint32
}

func recordTransfer(s *stats, e *event.DataEvent) {
	switch e.Kind {
		case event.Send:
			s.sent += uint32(len(e.Payload))
		case event.Receive:
			s.received += uint32(len(e.Payload))
		default:
			panic("Unexpected event kind")
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
	out.Out.WriteString("Date;ConnCount;ConnOpened;ConnClosed;Sent;Received;TotalTransferred\n")
}

func printStats(s *stats) {
	out.Out.WriteString(s.stamp.Format(out.Stampf))
	fmt.Fprintf(out.Out, ";%d;%d;%d;%d;%d;%d\n", s.connCount, s.connOpened, s.connClosed, s.sent, s.received, s.sent + s.received)
}