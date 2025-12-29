package packets

import (
	"z07/internal/assets"
	"z07/internal/game/domain"
	"z07/internal/protocol"
)

func writePosition(pw *protocol.PacketWriter, position domain.Position) {
	pw.WriteUint16(position.X)
	pw.WriteUint16(position.Y)
	pw.WriteUint8(position.Z)
}

func readPosition(pr *protocol.PacketReader) domain.Position {
	return domain.Position{
		X: pr.ReadUint16(),
		Y: pr.ReadUint16(),
		Z: pr.ReadUint8(),
	}
}

func readItem(pr *protocol.PacketReader) domain.Item {
	id := pr.ReadUint16()
	item := domain.Item{ID: id}
	thing := assets.Get(id)

	if thing.IsStackable || thing.IsFluid {
		item.Count = pr.ReadUint8()
		item.HasCount = true
	}

	return item
}
