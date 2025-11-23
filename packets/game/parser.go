package game

import (
	"errors"
	"goTibia/protocol"
)

// ErrUnknownOpcode is returned when we don't have a parser for this ID.
// The proxy uses this signal to just forward the raw bytes.
var ErrUnknownOpcode = errors.New("unknown opcode")

// IncomingPacket is a marker interface for any packet received from Client.
type S2CPacket interface {
	// We can add methods here later, e.g., Name() string
}

func ParseClientPacket(opcode uint8, pr *protocol.PacketReader) (S2CPacket, error) {
	switch opcode {
	case S2CLoginSuccessful:
		return ParseLoginResultMessage(pr)

	default:
		// We don't know this packet, let the proxy forward it raw.
		return nil, ErrUnknownOpcode
	}
}
