package main

import (
	"sync/atomic"
	"time"
)

//MonitorChanBufSize is the max buffer size for the channel used in the monitor
const MonitorChanBufSize = 1 << 16

var eventCount uint64

//go:generate stringer -type=eventType
type eventType uint8

//Event Types
const (
	EventIncomingRequest eventType = iota
	EventRequestForward
	EventSlaveOffer
	EventSlaveAccept
	EventSlaveTimeout
	EventSlaveError
	EventPanic
	EventRestart
	EventShutdown
)

type event struct {
	ID        uint64
	Type      eventType
	Timestamp time.Time
	Source    string
	Details   string
}

//NewEvent creates a new event
//NOTE: possible refactor
func NewEvent(t eventType, source, details string) (evt *event) {
	evt = new(event)
	evt.Type = t
	evt.Timestamp = time.Now()
	evt.ID = eventCount
	atomic.AddUint64(&eventCount, uint64(1))
	return
}

type monitor struct {
	channel     chan *event
	EventLog    []*event
	EventCounts map[eventType]uint
}

func (m *monitor) waitForEvents() {
	for {
		e := <-m.channel
		m.EventLog = append(m.EventLog, e)
		m.EventCounts[e.Type]++
	}
}

func (m monitor) addEvent(e *event) {
	m.channel <- e
}

var globalMonitor = new(monitor)

func init() {
	//init the vars for the global monitor instance
	globalMonitor.channel = make(chan *event, MonitorChanBufSize)
	globalMonitor.EventCounts = make(map[eventType]uint)
	go globalMonitor.waitForEvents()
}
