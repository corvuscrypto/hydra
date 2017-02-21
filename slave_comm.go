package main

//go:generate stringer -type=slavePacketType
type slavePacketType uint8

//The different packet types we expect
const (
	SlavePing slavePacketType = iota
	SlaveAcknowledgement
	SlaveDataRequest
	SlaveDataContinue
	SlaveDataResponse
	SlaveStatusRequest
	SlaveStatusResponse
	SlaveErrorResponse
)

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
