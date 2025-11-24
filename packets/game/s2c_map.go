package game

//
//type MapSliceMsg struct {
//	Direction uint8 // 0=N, 1=E, 2=S, 3=W
//	Tiles     []TileData
//}
//
//type TileData struct {
//	Position types.Position
//	Ground   uint16 // ID of the ground
//	Items    []ItemData
//}
//
//type ItemData struct {
//	ID    uint16
//	Count uint8 // or SubType
//}
//
//// MapDescriptionMsg contains the visible world data.
//type MapDescriptionMsg struct {
//	// We flatten the map into a slice for simpler access,
//	// or you can use a 3D map/array.
//	// Here we store a list of non-empty tiles.
//	Tiles []TileData
//}
//
//const (
//	MapWidth  = 18
//	MapHeight = 14
//)
//
//// ParseMapDescriptionMsg reads the compressed map stream.
//// clientPos is the player's current coordinate (needed to calculate offsets).
//func ParseMapDescriptionMsg(pr *protocol.PacketReader, clientPos types.Position) (*MapDescriptionMsg, error) {
//	msg := &MapDescriptionMsg{
//		Tiles: make([]TileData, 0, MapWidth*MapHeight),
//	}
//
//	// 1. Determine the Z-range exactly like the C++ code
//	var startz, endz, zstep int
//	z := int(clientPos.Z)
//
//	if z > 7 {
//		startz = z - 2
//		endz = z + 2
//		// Clamp to max layers (usually 15 in Tibia)
//		if endz > 15 {
//			endz = 15
//		}
//		zstep = 1
//	} else {
//		startz = 7
//		endz = 0
//		zstep = -1
//	}
//
//	// 2. Prepare the linear counters
//	// The C++ code carries 'skip' across floors, so we must treat
//	// the whole operation as filling a linear buffer of tiles.
//
//	currentFloor := startz
//	tilesProcessedOnFloor := 0
//	totalTilesOnFloor := MapWidth * MapHeight
//
//	// Loop until we have processed all floors
//	for {
//		// Calculate the offset for 3D perspective shifting
//		// C++: GetFloorDescription(..., z - nz, ...)
//		offset := z - currentFloor
//
//		// 3. Read the Token (Uint16 LittleEndian)
//		// Tibia writes 2 bytes.
//		// If Value >= 0xFF00, it is a SKIP (N tiles empty).
//		// If Value < 0xFF00, it is a TILE (Value is the Ground ID).
//		val := pr.ReadUint16()
//		if err := pr.Err(); err != nil {
//			return nil, err
//		}
//
//		// --- CASE A: SKIP (RLE) ---
//		if val >= 0xFF00 {
//			skipCount := int(val & 0xFF) // The lower byte is the count
//
//			// Advance our logical cursor by 'skipCount' empty tiles
//			for i := 0; i < skipCount; i++ {
//				tilesProcessedOnFloor++
//
//				// Handle Floor Wrapping
//				// If we skipped past the end of this floor, move to next floor
//				if tilesProcessedOnFloor >= totalTilesOnFloor {
//					currentFloor += zstep
//					tilesProcessedOnFloor = 0
//
//					// Check if we are done with all floors
//					if (zstep > 0 && currentFloor > endz) || (zstep < 0 && currentFloor < endz) {
//						return msg, nil
//					}
//				}
//			}
//		} else {
//			// --- CASE B: REAL TILE ---
//			// 'val' is the Ground ID.
//
//			// 1. Calculate actual X,Y coordinates
//			// The loop implies: nx goes 0..18, ny goes 0..14
//			nx := tilesProcessedOnFloor / MapHeight // Note: Tibia loops X then Y, or Y then X?
//			// C++: for nx (0..Width) { for ny (0..Height) }
//			// This means Y is the inner loop.
//			// So every increment moves Y. When Y fills Height, X increments.
//
//			realNx := tilesProcessedOnFloor / MapHeight
//			realNy := tilesProcessedOnFloor % MapHeight
//
//			// Apply the offset (Perspective shift)
//			tilePos := types.Position{
//				X: uint16(int(clientPos.X) + realNx + offset),
//				Y: uint16(int(clientPos.Y) + realNy + offset),
//				Z: uint8(currentFloor),
//			}
//
//			// 2. Parse the Tile Content
//			tile := ParseTile(pr, val, tilePos)
//			msg.Tiles = append(msg.Tiles, tile)
//
//			// 3. Advance Cursor (1 tile processed)
//			tilesProcessedOnFloor++
//			if tilesProcessedOnFloor >= totalTilesOnFloor {
//				currentFloor += zstep
//				tilesProcessedOnFloor = 0
//
//				if (zstep > 0 && currentFloor > endz) || (zstep < 0 && currentFloor < endz) {
//					return msg, nil
//				}
//			}
//		}
//	}
//}
//
//func ParseMapSlice(opcode uint8, pr *protocol.PacketReader, clientPos types.Position) (*MapSliceMsg, error) {
//	// These constants come from TFS/Client source (Map::maxClientViewportX)
//	// Usually 18x14 view.
//	var width, height int
//
//	switch opcode {
//	case S2CMapSliceNorth: // 0x65
//		width, height = 18, 1
//	case S2CMapSliceEast: // 0x66
//		width, height = 1, 14
//	case S2CMapSliceSouth: // 0x67
//		width, height = 18, 1
//	case S2CMapSliceWest: // 0x68
//		width, height = 1, 14
//	}
//
//	// Reuse your existing map logic, but force the specific width/height
//	// Note: You need to refactor ParseMapDescriptionMsg to accept explicit w/h
//	// or create a ParseRawMapDescription(pr, width, height) helper.
//
//	tiles, err := ParseRawMapDescription(pr, clientPos, width, height)
//	if err != nil {
//		return nil, err
//	}
//
//	return &MapSliceMsg{
//		Direction: opcode - 0x65, // Hacky way to map 0x65->0, 0x66->1...
//		Tiles:     tiles,
//	}, nil
//}
//
//func ParseTile(pr *protocol.PacketReader, groundID uint16, pos types.Position) TileData {
//	tile := TileData{
//		Position: pos,
//		Ground:   groundID,
//		Items:    []ItemData{},
//	}
//
//	// Loop reading items on the stack
//	for {
//		// Peek at the next 2 bytes to check if we are done with this tile.
//		// We use Peek because if it's a Skip Marker (>= 0xFF00),
//		// it belongs to the NEXT iteration of the main loop.
//
//		nextVal, err := pr.PeekUint16()
//		if err != nil {
//			break
//		} // EOF
//
//		// Heuristic: If the next value looks like a Skip Marker (>= 0xFF00),
//		// we are definitely done with this tile.
//		if nextVal >= 0xFF00 {
//			break
//		}
//
//		// DANGER ZONE:
//		// If 'nextVal' is a regular ID, is it an Item on THIS tile?
//		// Or the Ground ID of the NEXT tile?
//		// Without 'items.otb', we cannot know for sure.
//		//
//		// Standard Tibia Protocol Logic:
//		// MapDescription writes: Ground -> Items -> End.
//		// But it doesn't write an "End" delimiter for the tile specifically.
//		// It relies on the client knowing that the next ID is a Ground ID (0-4000 range usually?)
//		// or a Skip Marker.
//
//		// FOR A PROXY:
//		// Usually, we assume we keep reading until we hit a specific condition.
//		// Since we can't implement full OTB parsing here, we assume standard item structure.
//
//		itemID := pr.ReadUint16()
//		item := ItemData{ID: itemID}
//
//		// If you implement OTB checking:
//		// if ItemIsStackable(itemID) {
//		//     item.Count = pr.ReadByte()
//		// }
//
//		tile.Items = append(tile.Items, item)
//
//		// Safety break to prevent infinite loops in bad parsing
//		if len(tile.Items) > 10 {
//			break
//		}
//	}
//
//	return tile
//}
