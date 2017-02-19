package main

//go:generate stringer -type=slavePacketType
type slavePacketType uint8

const (
	SlavePing slavePacketType = iota
	SlaveAcknowledgement
	SlaveDataRequest
	SlaveDataResponse
	SlaveStatusRequest
	SlaveStatusResponse
	SlaveErrorResponse
)
