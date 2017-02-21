package main

//go:generate stringer -type=slavePacketType
type slavePacketType uint8

//The different packet types we expect
const (
	SlavePing slavePacketType = iota
	SlaveAcknowledgement
	SlaveDiscoveryRequest
	SlaveDiscoveryChallenge
	SlaveDiscoveryAccept
	SlaveDiscoveryReject
	SlaveDataRequest
	SlaveDataContinue
	SlaveDataResponse
	SlaveStatusRequest
	SlaveStatusResponse
	SlaveErrorResponse
)

func newPacket(t slavePacketType) slavePacket {
	basePacket := new(slavePacket)
	basePacket.Type = t
	return *basePacket
}

//Base struct to compose all other packets
type slavePacket struct {
	Type slavePacketType
}

type slavePingPacket struct {
	slavePacket
	ID packetID
}

type slaveAcknowledgementPacket struct {
	slavePacket
	ID packetID
}

type slaveKeyTransfer struct {
	slavePacket
	X     []byte
	Y     []byte
	XSign int
	YSign int
}

type slaveDiscoveryRequest struct {
	slavePacket
	SlaveID   string
	Resources []string
}

type slaveDiscoveryChallenge struct {
	slavePacket
	Nonce []byte
}

type slaveDiscoveryChallengeResponse struct {
	slavePacket
	Nonce  []byte
	Secret []byte
}

type slaveDiscoveryAccept struct {
	slavePacket
}

type slaveDiscoveryReject struct {
	slavePacket
	Reason string
}

type slaveDataRequestPacket struct {
	slavePacket
	ID       packetID
	Resource string
	Filter   string
	OrderBy  string
}

type slaveDataResponsePacket struct {
	slavePacket
	ID       packetID
	DataLeft int
	Data     []byte
}

type slaveStatusRequest struct {
	slavePacket
	ID     packetID
	Fields []string
}

type slaveStatusResponse struct {
	slavePacket
	ID   packetID
	Data []byte
}

type slaveErrorResponse struct {
	slavePacket
	ID    packetID
	Error string
}
