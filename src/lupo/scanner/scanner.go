package scanner

import (
	"bufio"
	"bytes"
	"lupo/event"
	"lupo/stream"
	"net/textproto"
	"time"
)

const (
	scanDelay = 20 * time.Millisecond
)

type Scanner struct {
	stream *stream.Stream
	timer  *time.Timer
}

func NewScanner(s *stream.Stream) *Scanner {
	return &Scanner{stream: s}
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
			event.PostData(s.stream.Cid, s.stream.Direction, chunk.Stamp, d)
		}
	}

	// Remove all chunks which have been processed
	s.stream.Chunks = s.stream.Chunks[0:0]
	s.stream.Data.Reset()

	// All modifications done, stream can proceed
	s.stream.Unlock()

	// TODO alternative to locking: Pass all chunks through channels. More difficult to merge adjacent chunks maybe?
}

// All possible starting bytes of a HTTP request or response
// These are byte arrays to make the comparison easier
var httpStart = [][]byte{[]byte("HTTP"), []byte("GET"), []byte("POST"), []byte("HEAD"), []byte("PUT"), []byte("OPTIONS"), 
				 		 []byte("DELETE"), []byte("TRACE"), []byte("CONNECT"), []byte("MOVE")}
const httpMin = "HTTP/1.1 200 OK"

// Look for: HTTP, GET, POST, HEAD, PUT, OPTIONS, DELETE, TRACE, CONNECT, MOVE
func checkHttpStart(d []byte) bool {
	if len(d) < len(httpMin) {
		return false
	}
	for _, v := range httpStart {
		if bytes.Equal(v, d[:len(v)]) {
			return true
		}		
	}
	return false
}

func tryHttp(data []byte) (start []byte, headers textproto.MIMEHeader, body []byte) {
	// Fail early by checking the first bytes
	if !checkHttpStart(data) {
		return
	}

	var err error
	buf := bufio.NewReader(bytes.NewReader(data))
	tp := textproto.NewReader(buf)

	// Try to parse <Method> <URL> <HTTP/version>
	//           or <HTTP/version> <Code> <Status>
	start, err = tp.ReadLineBytes()

	// Read headers
	headers, err = tp.ReadMIMEHeader()
	if err != nil {
		return nil, nil, nil
	}

	// The rest is content
	body = data[len(data)-buf.Buffered():]
	return
}
