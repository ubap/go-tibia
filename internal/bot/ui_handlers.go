package bot

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type BotSnapshot struct {
	FishingEnabled bool   `json:"fishingEnabled"`
	Name           string `json:"name"`
	X              uint16 `json:"x"`
	Y              uint16 `json:"y"`
	Z              uint8  `json:"z"`
}

func (b *Bot) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Stream updates every 100ms
	for {
		// Create a snapshot from your current bot fields
		snap := BotSnapshot{
			FishingEnabled: b.fishingEnabled,
			Name:           b.state.CaptureFrame().Player.Name,
			X:              b.state.CaptureFrame().Player.Pos.X,
			Y:              b.state.CaptureFrame().Player.Pos.Y,
			Z:              b.state.CaptureFrame().Player.Pos.Z,
		}

		payload, _ := json.Marshal(snap)
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}
