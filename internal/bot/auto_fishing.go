package bot

import (
	"goTibia/internal/game/domain"
	"log"
	"time"
)

func (b *Bot) loopFishing() {
	ticker := time.NewTicker(5000 * time.Millisecond)
	defer ticker.Stop()

	log.Println("[Bot] Auto fishing started")

	for {
		select {
		case <-b.stopChan:
			return

		case <-ticker.C:
			frame := b.state.CaptureFrame()
			pos := frame.Player.Pos

			for x := pos.X - 7; x <= pos.X+7; x++ {
				for y := pos.Y - 5; y <= pos.Y+5; y++ {
					tile := frame.WorldMap[domain.Position{X: x, Y: y, Z: pos.Z}]
					if tile.Ground.ID == 4598 { // Water with fish
						log.Printf("[Bot] Found water with tile at (%d, %d, %d)", x, y, pos.Z)
					}
				}

			}
		}
	}
}
