package main

import (
	"encoding/gob"
	"sync/atomic"
	"time"
)

type packetID uint64

var packetPerturbator uint64

func newPacketID() packetID {
	perturbator := atomic.AddUint64(&packetPerturbator, uint64(1)) % 100
	timestamp := uint64(time.Now().Unix())
	return packetID((timestamp * 100) + perturbator)
}

func receivePacket(d *gob.Decoder, expectedPacket interface{}) error {
	return d.Decode(expectedPacket)
}

func sendPacket(e *gob.Encoder, packet interface{}) {
	e.Encode(packet)
}
