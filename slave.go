package main

import (
	"crypto/cipher"
	"net"
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

func newSlaveConn(t *net.TCPConn) (conn *slaveConn, err error) {
	conn = new(slaveConn)
	conn.tcpConn = t
	conn.cipherBlock, err = cipher.NewGCM(globalCipher)
	return
}

type slave struct {
	nonce      []byte
	connection *slaveConn
}

func newSlave(conn *net.TCPConn) (s *slave, err error) {
	s.connection, err = newSlaveConn(conn)
	return
}
