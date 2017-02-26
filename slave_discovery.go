package main

import (
	"crypto/sha256"
	"encoding/gob"
	"net"
)

func abortConnection(s *slave) {
	err := s.connection.tcpConn.Close()
	if err != nil {
		panic(err)
	}
}

//sendChallenge sends the challenge packet containing the encrypted nonce
func sendChallenge(s *slave) (err error) {
	challengePacket := &slaveDiscoveryChallenge{
		newPacket(SlaveDiscoveryChallenge),
		s.nonce,
	}
	err = s.encoder.Encode(challengePacket)
	return
}

//sendChallenge sends the challenge packet containing the encrypted nonce
func receiveChallengeResponse(s *slave) (verified bool, err error) {
	challengePacket := new(slaveDiscoveryChallengeResponse)
	// this is to allow hashing even if the packet reception fails (to avoid a possible timing attack)
	challengePacket.Hash = []byte{}

	err = s.decoder.Decode(challengePacket)
	// check that hashes are the same that we expect
	verificationSignature := sha256.Sum256(append(s.nonce, []byte(globalConfig.Security.SecretKey)...))
	verified = string(verificationSignature[:]) == string(challengePacket.Hash)

	return
}

func sendDiscoveryRejection(s *slave) (err error) {
	rejectionPacket := &slaveDiscoveryReject{
		newPacket(SlaveDiscoveryReject),
		"Verification signature did not match!",
	}
	err = s.encoder.Encode(rejectionPacket)
	return
}

func sendDiscoveryAcceptance(s *slave) (err error) {
	acceptPacket := &slaveDiscoveryAccept{
		newPacket(SlaveDiscoveryAccept),
	}
	err = s.encoder.Encode(acceptPacket)
	return
}

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
	err = sendChallenge(ns)
	if err != nil {
		abortConnection(ns)
		return
	}

	//Receive Challenge Response
	verified, err := receiveChallengeResponse(ns)
	if err != nil {
		abortConnection(ns)
		return
	}

	if !verified {
		//We send a rejection response
		err = sendDiscoveryRejection(ns)
		if err != nil {
			abortConnection(ns)
		}
		return
	}

	//We send an acceptance if everything is okay and add this slave to our registry
	err = sendDiscoveryAcceptance(ns)
	if err != nil {
		abortConnection(ns)
		return
	}

}
