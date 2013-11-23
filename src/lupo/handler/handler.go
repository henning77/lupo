package handler

import (
	"lupo/util"
	"lupo/event"
	"io"
	"net"
	"time"
	"bytes"
	"sync"
)

type Chunk struct {
	Stamp time.Time
	Data []byte
}

type Stream struct {
	sync.Mutex

	Cid event.ConnId

	// Either Send or Receive
	Direction event.EventKind

	Data bytes.Buffer
	Chunks []Chunk

	// Every stream needs a scanner which scans it & generates events from it
	Scanner scanner.Scanner
}

func NewStream(cid event.ConnId, dir event.EventKind) (*Stream) {
	result := &Stream{
		Cid:cid,
		Direction:dir,
		// We should never have more than 10 chunks in a stream (scanner should remove them before the limit is reached)
		Chunks:make([]Chunk, 0, 10)
	}
	result.Scanner = scanner.NewScanner(result)
	return result
}

func (s *Stream) Write(p []byte) (n int, err error) {
	// Aquire lock (if scanner might be active)
	locked := false
	if len(s.Chunks) > 0 {
		s.Lock()
		locked = true
	}

	data.Write(p)
	d := data.Bytes()

	// Check if it is worthwile to create another chunk (i.e. time diff is sufficient)
	now := time.Now()
	lastChunk := ...
	if now.Sub(lastChunk.stamp) > 10 * time.Milliseconds {
		// Create new chunk
		chunk := Chunk{Stamp:time.Now(), Data:d[len(d) - len(p):]}
		s.Chunks = s.Chunks[:len(s.Chunks + 1)]
	} else {
		// Extend the last chunk
		lastChunk.data = d[len(d) - len(p) - len(lastChunk.data):]
	}

	// Release stream for scanner
	if locked {
		s.Unlock()
	}

	// Notify scanner
	s.Scanner.NotifyUpdate()

	return len(p), nil
}


type Handler struct {
	cid event.ConnId
	dst net.Conn
	src net.Conn
	send stream
	rcv stream
}

func NewHandler(dst net.Conn, src net.Conn, cid event.ConnId) *Handler {
	return &Handler{
		cid:cid,
		dst:dst,
		src:src,
		send:NewStream(cid, event.Send),
		rcv:NewStream(cid, event.Receive)
	}
}

func (h *Handler) Handle() {
	event.Connected(h.cid)

	// Copy & create events in both directions
	go copyWithBuffer(h.dst, h.src, h.send)
	go copyWithBuffer(h.src, h.dst, h.rcv)

	// TODO ensure both connections are closed in the end? (e.g. one closes, should cause the other one to close as well)
}

// Copy to dst and to event queue
func copyWithBuffer(dst io.Writer, src io.Reader, s *stream) {
	multi := io.MultiWriter(dst, s)
	if _, err := io.Copy(multi, src); err != nil {
		util.Printf("Closed connection: %v", err)
	}

	event.Disconnected(cid)
}
