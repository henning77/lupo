package stream

import (
	"bytes"
	"lupo/event"
	"sync"
	"time"
)

const (
	chunkMinSilence = 10 * time.Millisecond
)

type Listener interface {
	NotifyUpdate()
}

type Chunk struct {
	Stamp time.Time
	Data  []byte
}

type Stream struct {
	sync.Mutex

	Cid event.ConnId

	// Either Send or Receive
	Direction event.EventKind

	Data   bytes.Buffer
	Chunks []Chunk

	// Will be notified when the stream updates
	Listener
}

func NewStream(cid event.ConnId, dir event.EventKind) *Stream {
	return &Stream{
		Cid:       cid,
		Direction: dir,
		// We should never have more than 10 chunks in a stream (scanner should remove them before the limit is reached)
		Chunks: make([]Chunk, 0, 10)}
}

func (s *Stream) Write(p []byte) (n int, err error) {
	// Aquire lock (if scanner might be active)
	locked := false
	if len(s.Chunks) > 0 {
		s.Lock()
		locked = true
	}

	s.Data.Write(p)
	d := s.Data.Bytes()

	adjacentChunk := false
	if len(s.Chunks) > 0 {
		// Check if it is worthwile to create another chunk (i.e. time diff is sufficient)
		now := time.Now()
		lastChunk := s.Chunks[len(s.Chunks)-1]
		if now.Sub(lastChunk.Stamp) < chunkMinSilence {
			// Extend the last chunk
			lastChunk.Data = d[len(d)-len(p)-len(lastChunk.Data):]
			adjacentChunk = true
		}
	}
	if !adjacentChunk {
		// Create new chunk
		chunk := Chunk{Stamp: time.Now(), Data: d[len(d)-len(p):]}
		s.Chunks = append(s.Chunks, chunk)
	}

	// Release stream for scanner
	if locked {
		s.Unlock()
	}

	// Notify listener
	s.NotifyUpdate()

	return len(p), nil
}
