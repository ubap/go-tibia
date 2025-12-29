package state

import "z07/internal/game/domain"

type WorldSnapshot struct {
	Player     domain.Player
	Equipment  [11]domain.Item
	Containers [16]*domain.Container
	WorldMap   map[domain.Position]*domain.Tile
}

type ItemInInventory struct {
	Item     domain.Item
	Position domain.Position
}

func (s WorldSnapshot) FindItemInEqAndOpenWindows(itemId uint16) *ItemInInventory {
	itemInEq := s.findItemInEq(itemId)
	if itemInEq != nil {
		return itemInEq
	}

	return s.findItemInContainers(itemId)
}

func (s WorldSnapshot) findItemInEq(itemId uint16) *ItemInInventory {
	for slot, item := range s.Equipment {
		if item.ID == itemId {
			equipmentSlot := domain.EquipmentSlot(slot)
			pos := domain.NewInventoryPosition(equipmentSlot)
			return &ItemInInventory{
				Item:     item,
				Position: pos,
			}
		}
	}
	return nil
}

func (s WorldSnapshot) findItemInContainers(itemId uint16) *ItemInInventory {
	for cid, container := range s.Containers {
		if container == nil {
			// the container is not open
			continue
		}
		for slot, item := range container.Items {
			if item.ID == itemId {
				pos := domain.NewContainerPosition(cid, slot)
				return &ItemInInventory{
					Item:     item,
					Position: pos,
				}
			}
		}
	}
	return nil
}

func (s WorldSnapshot) FindTileNearby(radiusX, radiusY int, criteria func(*domain.Tile) bool) (*domain.Position, *domain.Tile) {
	pos := s.Player.Pos

	for x := pos.X - uint16(radiusX); x <= pos.X+uint16(radiusX); x++ {
		for y := pos.Y - uint16(radiusY); y <= pos.Y+uint16(radiusY); y++ {
			currPos := domain.Position{X: x, Y: y, Z: pos.Z}

			if tile, ok := s.WorldMap[currPos]; ok {
				if criteria(tile) {
					return &currPos, tile
				}
			}
		}
	}
	return nil, nil
}
