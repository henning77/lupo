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


// A generic event on a specific connection.
type Event struct {
	Cid     ConnId
	Kind    EventKind
	Stamp   time.Time
}

// Event where interesting data is sent
type DataEvent struct {
	Event
	Payload []byte
}

// HTTP protocol event
type HttpEvent struct {
	DataEvent
	Start   []byte
	Headers textproto.MIMEHeader
	Body    []byte
}

type MessageEvent struct {
	Event
	Message string
}

type ConnectEvent struct {
	Event
	From string
}

type DisconnectEvent struct {
	Event
	Initiator EventKind
}

func PostGlobalf(s string, a ...interface{}) {
	PostGlobal(fmt.Sprintf(s, a...))
}

func PostGlobal(desc string) {
	Events <- &MessageEvent{
		Event:   Event{Cid: 0, Kind: Global, Stamp: time.Now()},
		Message: desc}
}

func PostConnect(cid ConnId, from string) {
	Events <- &ConnectEvent{
		Event:   Event{Cid: cid, Kind: Connect, Stamp: time.Now()},
		From:    from}
}

// dir specifies which partner of the connection initiated the close.
func PostDisconnect(cid ConnId, dir EventKind) {
	Events <- &DisconnectEvent{
		Event:   Event{Cid: cid, Kind: Disconnect, Stamp: time.Now()},		
		Initiator: dir}
}

func PostData(cid ConnId, kind EventKind, stamp time.Time, payload []byte) {
	Events <- &DataEvent{
		Event:   Event{Cid: cid, Kind: kind, Stamp: stamp},
		Payload: payload}
}

func PostHttp(cid ConnId, kind EventKind, stamp time.Time, payload []byte, start []byte, headers textproto.MIMEHeader, body []byte) {
	Events <- &HttpEvent{
		DataEvent:   DataEvent{
			Event:   Event{Cid: cid, Kind: kind, Stamp: stamp},
			Payload: payload},
		Start:   start,
		Headers: headers,
		Body:    body}
}
