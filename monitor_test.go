package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestBasicMonitorEventLogging(T *testing.T) {
	evt := NewEvent(EventIncomingRequest, "", "")
	globalMonitor.addEvent(evt)
	//wait a bit for it to consume the new event
	<-time.Tick(50 * time.Millisecond)
	if globalMonitor.EventLog[0] != evt || globalMonitor.EventCounts[EventIncomingRequest] != 1 {
		T.Error("Event not properly logged")
	}

	//Reset the counter to 0 before the next test
	globalMonitor.EventCounts[EventIncomingRequest] = 0
}

var expectedCountsMap = make(map[eventType]uint)
var globalMutex = new(sync.Mutex)

func eventGenerator(wg *sync.WaitGroup) {
	for i := 0; i < 50; i++ {
		whichEvent := eventType(rand.Intn(int(EventShutdown)))
		evt := NewEvent(whichEvent, "", "")
		globalMonitor.addEvent(evt)
		globalMutex.Lock()
		expectedCountsMap[whichEvent]++
		globalMutex.Unlock()
	}
	wg.Done()
}

func TestMultiEventLogging(T *testing.T) {
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go eventGenerator(wg)
	go eventGenerator(wg)
	wg.Wait()
	//wait a bit for it to consume the new events
	<-time.Tick(50 * time.Millisecond)
	for k, v := range expectedCountsMap {
		if globalMonitor.EventCounts[k] != v {
			T.Error("Event count mismatch", globalMonitor.EventCounts[k], v)
		}
	}
	fmt.Println(globalMonitor.EventCounts)
}
