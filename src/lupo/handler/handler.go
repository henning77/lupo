package handler

import (
	"lupo/out"
	"lupo/event"
	"lupo/stream"
	"lupo/scanner"
	"io"
	"net"
)

type Handler struct {
	cid event.ConnId
	dst net.Conn
	src net.Conn
	send *stream.Stream
	rcv *stream.Stream
}

func NewHandler(dst net.Conn, src net.Conn, cid event.ConnId) *Handler {
	result := &Handler{
		cid:cid,
		dst:dst,
		src:src,
		send:stream.NewStream(cid, event.Send),
		rcv:stream.NewStream(cid, event.Receive)}

	// Set up scanners
	result.send.Listener = scanner.NewScanner(result.send)
	result.rcv.Listener = scanner.NewScanner(result.rcv)

	return result
}

func (h *Handler) Handle() {
	event.PostConnect(h.cid)

	// Copy & create events in both directions
	go copyWithBuffer(h.dst, h.src, h.cid, h.send)
	go copyWithBuffer(h.src, h.dst, h.cid, h.rcv)

	// TODO ensure both connections are closed in the end? (e.g. one closes, should cause the other one to close as well)
}

// Copy to dst and to event queue
func copyWithBuffer(dst io.Writer, src io.Reader, cid event.ConnId, s *stream.Stream) {
	multi := io.MultiWriter(dst, s)
	if _, err := io.Copy(multi, src); err != nil {
		out.Printf("Closed connection: %v", err)
	}

	event.PostDisconnect(cid)
}
