package state

import "goTibia/internal/game/domain"

type WorldSnapshot struct {
	Player    domain.Player
	Equipment [11]domain.Item
}
