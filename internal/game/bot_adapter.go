package game

import (
	"goTibia/internal/bot"
	"goTibia/internal/game/domain"
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"goTibia/internal/protocol"
	"log"
)

type BotAdapter struct {
	State      *state.GameState
	ServerConn *protocol.Connection
	ClientConn *protocol.Connection
}

// region Implementing WorldStateReader

func (ba *BotAdapter) GetPlayerPosition() domain.Coordinate {
	return ba.State.CaptureFrame().Player.Pos
}

func (ba *BotAdapter) GetInventoryItem(slot domain.EquipmentSlot) domain.Item {
	return ba.State.CaptureFrame().Equipment[slot]
}

func (ba *BotAdapter) GetPlayerID() uint32 {
	return ba.State.CaptureFrame().Player.ID
}

// endregion Implementing WorldStateReader

// region Implementing ActionDispatcher

func (ba *BotAdapter) Say(text string) {
	log.Printf("[Game] Say: %s", text)
}

// endregion Implementing ActionDispatcher

// region Implementing ClientManipulator

var _ bot.ClientManipulator = (*BotAdapter)(nil)

func (ba *BotAdapter) SetLocalPlayerLight(lightLevel uint8, color uint8) {
	id := ba.State.CaptureFrame().Player.ID

	if id == 0 {
		// Player is not logged in yet
		return
	}

	// 2. Construct the specific packet
	pkt := &packets.CreatureLightMsg{
		CreatureID: id,
		LightLevel: lightLevel,
		Color:      color,
	}

	// 3. Send it to the CLIENT (Cheating visually)
	if err := ba.ClientConn.SendPacket(pkt); err != nil {
		// You might log errors here or handle disconnects
	}
}

// endregion Implementing ClientManipulator
