package bot

import (
	"z07/internal/game/domain"
	"z07/internal/game/packets"
	"z07/internal/game/state"
)

func (b *Bot) UseItemFromInventoryOnTile(item state.ItemInInventory, to domain.Tile) {
	pkt := packets.UseItemWithCrosshairRequest{
		FromPos:      item.Position,
		FromItemId:   item.Item.ID,
		FromStackPos: 0, // stack pos is always 0 for inventory items

		ToPos:      to.Position,
		ToItemId:   to.TopItem().ID,
		ToStackPos: uint8(len(to.Items) - 1),
	}
	b.serverConn.SendPacket(&pkt)
}
