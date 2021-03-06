package handler

import (
	"io"
	"lupo/event"
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
	event.PostConnect(h.cid, h.src.RemoteAddr().String())

	// Notify closed streams through this channel
	done := make(chan *stream.Stream)

	// Copy & create events in both directions
	go copyWithBuffer(h.dst, h.src, h.cid, h.send, done)
	go copyWithBuffer(h.src, h.dst, h.cid, h.rcv, done)

	// Wait for the first stream to close
	closedStream := <-done

	// Close both connections
	h.src.Close()
	h.dst.Close()
	event.PostDisconnect(h.cid, closedStream.Direction)

	// Wait for the second stream to quit
	<-done
				
	close(done)
}

// Copy to dst and to event queue
func copyWithBuffer(dst io.Writer, src io.Reader, cid event.ConnId, s *stream.Stream, done chan<- *stream.Stream) {
	multi := io.MultiWriter(dst, s)
	if _, err := io.Copy(multi, src); err != nil {
		//out.Printf("Closed connection: %v", err)
	}

	done <- s
}
