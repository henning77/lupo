package event

import (
	"net/textproto"
	"time"
)

// All events from all connections go here
var Events = make(chan interface{}, 1000)

type ConnId int
type EventKind byte

const (
	_                 = iota
	Connect EventKind = iota
	Disconnect
	Send
	Receive
)

// An event on a specific connection.
type Event struct {
	Cid     ConnId
	Kind    EventKind
	Stamp   time.Time
	Payload []byte
}

type HttpEvent struct {
	Event
	Start   []byte
	Headers textproto.MIMEHeader
	Body    []byte
}

func PostConnect(cid ConnId) {
	Events <- &Event{
		Cid:     cid,
		Kind:    Connect,
		Stamp:   time.Now(),
		Payload: nil}
}

func PostDisconnect(cid ConnId) {
	Events <- &Event{
		Cid:     cid,
		Kind:    Disconnect,
		Stamp:   time.Now(),
		Payload: nil}
}

func Post(cid ConnId, kind EventKind, stamp time.Time, payload []byte) {
	Events <- &Event{
		Cid:     cid,
		Kind:    kind,
		Stamp:   stamp,
		Payload: payload}
}

func PostHttp(cid ConnId, kind EventKind, stamp time.Time, payload []byte, start []byte, headers textproto.MIMEHeader, body []byte) {
	Events <- &HttpEvent{
		Event:   Event{Cid: cid, Kind: kind, Stamp: stamp, Payload: payload},
		Start:   start,
		Headers: headers,
		Body:    body}
}
