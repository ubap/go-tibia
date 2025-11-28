package game

import (
	"fmt"
	"goTibia/dat"
	"goTibia/protocol"
)

// ReadItem reads a full Item (ID + Optional Count/Subtype) from the stream.
// This is used for Inventory, Containers, and Tile Stacks.
func ReadItem(pr *protocol.PacketReader) Item {
	// 1. Read ID
	id := pr.ReadUint16()

	// 2. Setup Struct
	item := Item{ID: id}

	// 3. Check Flags for Extra Byte
	thing := dat.Get(id)

	fmt.Printf("    > Item ID: %d (Stackable: %v)\n", id, thing.IsStackable)

	// fmt.Printf("Thing Info: %+v\n", thing)

	if thing.IsStackable || thing.IsFluid {
		item.Count = pr.ReadByte()
		item.HasCount = true

		fmt.Printf("      > Count: %d\n", item.Count)
	}

	return item
}
