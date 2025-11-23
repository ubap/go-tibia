package game

import (
	"goTibia/protocol"
	"goTibia/types"
)

func writePosition(pw *protocol.PacketWriter, position types.Position) {
	pw.WriteUint16(position.X)
	pw.WriteUint16(position.Y)
	pw.WriteByte(position.Z)
}

func readPosition(pr *protocol.PacketReader) types.Position {
	return types.Position{
		X: pr.ReadUint16(),
		Y: pr.ReadUint16(),
		Z: pr.ReadByte(),
	}
}
