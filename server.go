package main

import (
	"log"
	"net"
)

type server struct {
	listener *net.TCPListener
}

func (s server) listenDiscovery() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			log.Println(err)
		}
		defer conn.Close()
		globalMonitor.addEvent(NewEvent(EventSlaveOffer, conn.RemoteAddr().String(), "listenDiscovery"))
	}
}

func StartServer() {
	//Init the new server instance
	s := new(server)
	listener, _ := net.ListenTCP("net", nil)
	s.listener = listener
}
