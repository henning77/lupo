package handler

import (
	"io"
	"lupo/event"
	"lupo/out"
	"lupo/scanner"
	"lupo/stream"
	"net"
)

type Handler struct {
	cid  event.ConnId
	dst  net.Conn
	src  net.Conn
	send *stream.Stream
	rcv  *stream.Stream
}

func NewHandler(dst net.Conn, src net.Conn, cid event.ConnId) *Handler {
	result := &Handler{
		cid:  cid,
		dst:  dst,
		src:  src,
		send: stream.NewStream(cid, event.Send),
		rcv:  stream.NewStream(cid, event.Receive)}

	// Set up scanners
	result.send.Listener = scanner.NewScanner(result.send)
	result.rcv.Listener = scanner.NewScanner(result.rcv)

	return result
}

func (h *Handler) Handle() {
	defer h.src.Close()
	defer h.dst.Close()

	event.PostConnect(h.cid, h.src.RemoteAddr().String())

	done := make(chan bool)

	// Copy & create events in both directions
	go copyWithBuffer(h.dst, h.src, h.cid, h.send, done)
	go copyWithBuffer(h.src, h.dst, h.cid, h.rcv, done)

	// Wait for connection to close
	<-done
	event.PostDisconnect(h.cid)
	// Wait for the other stream to close as well
	<-done
}

// Copy to dst and to event queue
func copyWithBuffer(dst io.Writer, src io.Reader, cid event.ConnId, s *stream.Stream, done chan<- bool) {
	multi := io.MultiWriter(dst, s)
	if _, err := io.Copy(multi, src); err != nil {
		out.Printf("Closed connection: %v", err)
	}

	done <- true
}
