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
		log.Println("received slave connection")
		if err != nil {
			log.Println(err)
		}
		go handleDiscovery(conn)
	}
}

func StartServer() {
	//Init the new server instance
	s := new(server)

	// make the bind address
	addr, err := net.ResolveTCPAddr("tcp", ":"+globalConfig.Server.DiscoveryPort)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener
	s.listenDiscovery()
}
