package event

import (
	"time"
)

// All events from all connections go here
var Events = make(chan *Event, 1000)

// An event on a specific connection.
type Event struct {
	Cid     ConnId
	Kind    EventKind
	Stamp time.Time
	Payload []byte
}

type HttpEvent struct {
	Event
	Method []byte
	Headers []byte
	Content []byte
}

type ConnId int
type EventKind byte

const (
	_                   = iota
	Connect EventKind = iota
	Disconnect
	Send 
	Receive 
)

func Connected(cid ConnId) {
	Events <- &Event{Cid:cid, Kind:Connect, Stamp:time.Now(), Payload:nil}
}

func Disconnected(cid ConnId) {
	Events <- &Event{Cid:cid, Kind:Disconnect, Stamp:time.Now(), Payload:nil}
}

func Sent(cid ConnId, payload []byte) {
	Events <- &Event{Cid:cid, Kind:Send, Stamp:time.Now(), Payload:payload}
}

func Received(cid ConnId, payload []byte) {
	Events <- &Event{Cid:cid, Kind:Receive, Stamp:time.Now(), Payload:payload}
}
