package game

import "goTibia/protocol"

type LoginResponse struct {
	ClientDisconnected       bool
	ClientDisconnectedReason string
}

func ParseLoginResultMessage(pr *protocol.PacketReader) (*LoginResponse, error) {
	pr.ReadAll()
	return nil, nil
}
