package main

import (
	"crypto/cipher"
	"encoding/gob"
	"math/big"
	"net"
	"time"
)

type slaveConn struct {
	tcpConn     *net.TCPConn
	cipherBlock cipher.AEAD
}

func (s slaveConn) Write(data []byte) (n int, err error) {
	var cipherText []byte
	s.cipherBlock.Seal(cipherText, nil, data, nil)
	return s.tcpConn.Write(cipherText)
}

func (s slaveConn) Read(dst []byte) (n int, err error) {
	var cipherText []byte
	n, err = s.tcpConn.Read(cipherText)
	s.cipherBlock.Open(dst, nil, cipherText, nil)
	return
}

func (s *slaveConn) Close() error {
	return s.tcpConn.Close()
}

func newSlaveConn(t *net.TCPConn) (conn *slaveConn, err error) {
	conn = new(slaveConn)
	conn.tcpConn = t

	//create a raw gob (d)e(n)coder
	decoder := gob.NewDecoder(t)
	encoder := gob.NewEncoder(t)

	//create a new private key for exchange
	priv, X, Y, err := createNewKey()
	if err != nil {
		return nil, err
	}

	keyTransferPacket := &slaveKeyTransfer{
		newPacket(SlaveKeyTransfer),
		X.Bytes(),
		Y.Bytes(),
		X.Sign(),
		Y.Sign(),
	}

	//receive the slave's public key
	slaveKeyPacket := new(slaveKeyTransfer)
	err = decoder.Decode(slaveKeyPacket)
	if err != nil {
		t.Close()
		return
	}
	//send the public key
	err = encoder.Encode(keyTransferPacket)
	if err != nil {
		t.Close()
		return
	}

	//reconstruct public key from the packet rx'd
	slaveX := big.NewInt(0).SetBytes(slaveKeyPacket.X)
	if slaveKeyPacket.XSign == -1 {
		slaveX = slaveX.Neg(slaveX)
	}
	slaveY := big.NewInt(0).SetBytes(slaveKeyPacket.Y)
	if slaveKeyPacket.YSign == -1 {
		slaveY = slaveY.Neg(slaveY)
	}

	conn.cipherBlock, err = createNewCipher(priv, slaveX, slaveY)
	return
}

type slave struct {
	id             string
	nonce          []byte
	connection     *slaveConn
	encoder        *gob.Encoder
	decoder        *gob.Decoder
	requestChannel chan interface{}
	resources      []string
}

func (s *slave) destroy() {
	s.connection.Close()
	return
}

func (s *slave) waitAndMaintain() {
	for {
		select {
		case <-s.requestChannel:
			//Handle a request
		case <-time.Tick(time.Duration(globalConfig.Server.HeartbeatInterval) * time.Second):
			// send a heartbeat just to check if things are okay
		}
	}
}

func newSlave(conn *net.TCPConn, id string, resources []string) (s *slave, err error) {
	s.connection, err = newSlaveConn(conn)
	s.id = id
	s.resources = resources
	s.decoder = gob.NewDecoder(s.connection)
	s.encoder = gob.NewEncoder(s.connection)
	return
}
