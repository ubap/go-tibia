package packets

import (
	"errors"
	"goTibia/internal/game/domain"
	"goTibia/internal/protocol"
	"io"
)

var ErrUnknownOpcode = errors.New("unknown opcode")

type C2SPacket interface {
}

// S2CPacket is a marker interface for any packet received from Server.
type S2CPacket interface {
	// We can add methods here later, e.g., Name() string
	// TODO consider adding Opcode() uint8 to identify the packet type
	// Opcode() uint8
}

type InjectablePacket interface {
	S2CPacket
	protocol.Encodable
}

func ReadAndParseS2C(reader *protocol.PacketReader, ctx ParsingContext) (S2CPacket, error) {
	if reader.Remaining() == 0 {
		return nil, io.EOF
	}
	opcode := S2COpcode(reader.ReadUint8())
	return ParseS2CPacket(opcode, reader, ctx)
}

func ReadAndParseC2S(reader *protocol.PacketReader) (C2SPacket, error) {
	if reader.Remaining() == 0 {
		return nil, io.EOF
	}
	opcode := C2SOpcode(reader.ReadUint8())
	return ParseC2SPacket(opcode, reader)
}

func ParseS2CPacket(opcode S2COpcode, pr *protocol.PacketReader, ctx ParsingContext) (S2CPacket, error) {
	switch opcode {
	case S2CLoginSuccessful:
		return ParseLoginResultMessage(pr)
	case S2CSLoginQueue:
		return ParseLoginQueueMsg(pr)
	case S2CMapDescription:
		return ParseMapDescriptionMsg(pr)
	case S2CMapSliceNorth:
		return ParseMove(pr, ctx, domain.North)
	case S2CMapSliceSouth:
		return ParseMove(pr, ctx, domain.South)
	case S2CMapSliceEast:
		return ParseMove(pr, ctx, domain.East)
	case S2CMapSliceWest:
		return ParseMove(pr, ctx, domain.West)
	case S2CMoveCreature:
		return ParseMoveCreature(pr)
	case S2CPing:
		return &PingMsg{}, nil
	case S2CMagicEffect:
		return ParseMagicEffect(pr)
	case S2CAddTileThing:
		return ParseAddTileThingMsg(pr)
	case S2CRemoveTileThing:
		return ParseRemoveTileThing(pr)
	case S2CAddInventoryItem:
		return ParseAddInventoryItemMsg(pr)
	case S2CRemoveInventoryItem:
		return ParseRemoveInventoryItemMsg(pr)
	case S2CWorldLight:
		return ParseWorldLight(pr)
	case S2CCreatureLight:
		return ParseCreatureLight(pr)
	case S2CCreatureHealth:
		return ParseCreatureHealth(pr)
	case S2CPlayerIcons:
		return ParsePlayerIcons(pr)
	case S2CServerClosed:
		return ParseServerClosedMsg(pr)
	case S2COpenContainer:
		return ParseOpenContainerMsg(pr)
	case S2CCloseContainer:
		return ParseCloseContainerMsg(pr)
	case S2CRemoveContainerItem:
		return ParseRemoveContainerItemMsg(pr)
	case S2CAddContainerItem:
		return ParseAddContainerItemMsg(pr)
	case S2CUpdateContainerItem:
		return ParseUpdateContainerItemMsg(pr)
	case S2CUpdateTileItem:
		return ParseUpdateTileItemMsg(pr)
	case S2CPlayerSkills:
		return ParsePlayerSkillMsg(pr)
	case S2CPlayerStats:
		return ParsePlayerStatsMsg(pr)

	default:
		return nil, ErrUnknownOpcode
	}
}

func ParseC2SPacket(opcode C2SOpcode, pr *protocol.PacketReader) (C2SPacket, error) {
	switch opcode {
	case C2SLookRequest:
		return ParseLookRequest(pr)
	case C2SUseItemWithCrosshair:
		return ParseUseItemWithCrosshairRequest(pr)
	default:
		return nil, ErrUnknownOpcode
	}
}
