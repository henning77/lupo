package scanner

import (
	"lupo/stream"
	"lupo/event"
	"net/textproto"
	"time"
	"bufio"
	"strings"
	"bytes"
)

const (
	scanDelay = 20 * time.Millisecond
)

type Scanner struct {
	stream *stream.Stream
	timer *time.Timer
}

func NewScanner(s *stream.Stream) (*Scanner) {
	return &Scanner{stream:s}
}

// Notify scanner that the watched stream has changed & a rescan should take place.
func (s *Scanner) NotifyUpdate() {
	// If a trigger is/was already scheduled, reschedule it
	if s.timer != nil {
		s.timer.Reset(scanDelay)
	} else {
		// Trigger scan after delay, so adjacent chunks have time to come in.
		s.timer = time.AfterFunc(scanDelay, func() {
				s.scan()
			})
	}	

	// TODO Alternative strategy: don't reschedule, just let it fire at the scheduled time (i.e. prevent stream from clogging up)
}

func (s *Scanner) scan() {
	// Ensure incoming chunks doesn't modify the stream
	s.stream.Lock()

	for _, chunk := range s.stream.Chunks {
		// Copy the chunk for our own use
		d := make([]byte, len(chunk.Data))
		copy(d, chunk.Data)

		// HTTP chunk is either
		// Send stream:
		// 	   1) HTTP request (maybe with content)
		//     2) Some arbitrary content sent after the initial request
		// Receive stream:
		//     1) HTTP response (maybe with content)
		//     2) Some arbitrary content sent after the initial response
		r, h, b := tryHttp(chunk.Data)
		if r != nil {
			event.PostHttp(s.stream.Cid, s.stream.Direction, chunk.Stamp, d, r, h, b)
		} else {
			event.Post(s.stream.Cid, s.stream.Direction, chunk.Stamp, d)
		}
	}

	// Remove all chunks which have been processed
	s.stream.Chunks = s.stream.Chunks[0:0]
	s.stream.Data.Reset()

	// All modifications done, stream can proceed
	s.stream.Unlock()

	// TODO alternative to locking: Pass all chunks through channels. More difficult to merge adjacent chunks maybe?
}

func tryHttp(data []byte) (start []byte, headers textproto.MIMEHeader, body []byte) {
	var err error
	buf := bufio.NewReader(bytes.NewReader(data))
	tp := textproto.NewReader(buf)
	
	// Try to parse <Method> <URL> <HTTP/version>
	//           or <HTTP/version> <Code> <Status>
	start, err = tp.ReadLineBytes()
	f := strings.SplitN(string(start), " ", 3)
	if len(f) < 2 || (f[2] != "HTTP/1.0" && f[2] != "HTTP/1.1" && f[0] != "HTTP/1.0" && f[0] != "HTTP/1.1") {
		return nil, nil, nil
	}

	// Read headers
	headers, err = tp.ReadMIMEHeader()
	if err != nil {
		return nil, nil, nil
	}

	// The rest is content
	body = data[len(data) - buf.Buffered():]
	return
}
