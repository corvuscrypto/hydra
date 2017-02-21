package main

import (
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
