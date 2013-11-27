package event

import (
	"net/textproto"
	"time"
	"fmt"
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
	Global
)

func (k EventKind) String() string {
	switch k {
	case Connect:
		return " ["
	case Disconnect:
		return " ]"
	case Send:
		return "->"
	case Receive:
		return "<-"
	case Global:
		return "  "
	default:
		return "na"
	}
}


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

func PostGlobalf(s string, a ...interface{}) {
	PostGlobal(fmt.Sprintf(s, a...))
}

func PostGlobal(desc string) {
	Events <- &Event{
		Cid:     0,
		Kind:    Global,
		Stamp:   time.Now(),
		Payload: []byte(desc)}
}

func PostConnect(cid ConnId, from string) {
	Events <- &Event{
		Cid:     cid,
		Kind:    Connect,
		Stamp:   time.Now(),
		Payload: []byte(from)}
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
