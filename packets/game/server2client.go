package game

import (
	"goTibia/protocol"
	"goTibia/types"
)

type LoginResponse struct {
	ClientDisconnected bool
	PlayerId           uint32
	BeatDuration       uint16
	CanReportBugs      bool
}

type MapDescription struct {
	Pos types.Position
}

func ParseLoginResultMessage(pr *protocol.PacketReader) (*LoginResponse, error) {
	lr := &LoginResponse{}

	lr.PlayerId = pr.ReadUint32()
	lr.BeatDuration = pr.ReadUint16()
	lr.CanReportBugs = pr.ReadBool()

	return lr, pr.Err()
}

func ParseMapDescription(pr *protocol.PacketReader) (*MapDescription, error) {
	return &MapDescription{Pos: readPosition(pr)}, ErrNotFullyImplemented
}

func ParsePlayerStats(pr *protocol.PacketReader) (*MapDescription, error) {
	return nil, ErrUnknownOpcode
}
