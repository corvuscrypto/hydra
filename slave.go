package main

import (
	"crypto/cipher"
	"crypto/sha256"
	"encoding/gob"
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

func (s *slaveConn) Close() error {
	return s.tcpConn.Close()
}

func newSlaveConn(t *net.TCPConn) (conn *slaveConn, err error) {
	conn = new(slaveConn)
	conn.tcpConn = t
	//create a new private key for exchange
	priv, X, Y, err := createNewKey()
	if err != nil {
		return nil, err
	}
	conn.cipherBlock, err = createNewCipher(priv, X, Y)
	return
}

type slave struct {
	id         string
	nonce      []byte
	connection *slaveConn
	encoder    *gob.Encoder
	decoder    *gob.Decoder
	resources  []string
}

func (s *slave) sendChallenge() (err error) {
	challengePacket := slaveDiscoveryChallenge{
		newPacket(SlaveDiscoveryChallenge),
		s.nonce,
	}
	err = s.encoder.Encode(challengePacket)
	if err != nil {
		err2 := s.connection.tcpConn.Close()
		if err2 != nil {
			panic(err)
		}
	}
	return
}

func (s slave) getChallengeResponse() (ok bool, err error) {
	challengeResponse := new(slaveDiscoveryChallengeResponse)
	err = s.encoder.Encode(challengeResponse)
	if err != nil {
		err2 := s.connection.tcpConn.Close()
		if err2 != nil {
			panic(err)
		}
	}
	expectedHash := sha256.Sum256(append(s.nonce, secret...))
	ok = string(challengeResponse.Hash) == string(expectedHash[:])
	return
}

func newSlave(conn *net.TCPConn, id string, resources []string) (s *slave, err error) {
	s.connection, err = newSlaveConn(conn)
	s.id = id
	s.resources = resources
	s.decoder = gob.NewDecoder(s.connection)
	s.encoder = gob.NewEncoder(s.connection)
	return
}
