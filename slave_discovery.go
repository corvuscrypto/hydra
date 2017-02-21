package main

import (
	"encoding/gob"
	"net"
)

func handleDiscovery(conn *net.TCPConn) {
	//attempt to handle the offer
	decoder := gob.NewDecoder(conn)
	discoveryPacket := new(slaveDiscoveryRequest)
	err := decoder.Decode(&discoveryPacket)
	if err != nil {
		conn.Close()
		return
	}
	//If it's a successful decode then we continue with the discovery sequence
	//Create a new slave representation
	ns, err := newSlave(conn, discoveryPacket.SlaveID, discoveryPacket.Resources)
	if err != nil {
		conn.Close()
		return
	}

	//Send Challenge
	challengePacket := slaveDiscoveryChallenge{
		newPacket(SlaveDiscoveryChallenge),
		ns.nonce,
	}
	err = ns.encoder.Encode(challengePacket)
	if err != nil {
		conn.Close()
		return
	}
}
