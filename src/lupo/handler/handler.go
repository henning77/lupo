package handler

import (
	"lupo/util"
	"lupo/event"
	"io"
	"net"
)

type sendOrRcv func(cid event.ConnId, payload []byte)

type EventGen struct {
	cid event.ConnId
	action sendOrRcv
}

func (g *EventGen) Write(p []byte) (n int, err error) {
	g.action(g.cid, p)
	return len(p), nil
}

func Handle(dst net.Conn, src net.Conn, cid event.ConnId) {
	event.Connected(cid)

	// Copy & create events in both directions
	go copyWithEvents(dst, src, cid, event.Sent)
	go copyWithEvents(src, dst, cid, event.Received)

	// TODO ensure both connections are closed in the end? (e.g. one closes, should cause the other one to close as well)
}

// Copy to dst and to event queue
func copyWithEvents(dst io.Writer, src io.Reader, cid event.ConnId, action sendOrRcv) {
	gen := &EventGen{cid:cid, action:action}
	multi := io.MultiWriter(dst, gen)
	if _, err := io.Copy(multi, src); err != nil {
		util.Printf("Closed connection: %v", err)
	}

	event.Disconnected(cid)
}
