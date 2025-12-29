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
	FishingEnabled bool       `json:"fishingEnabled"`
	Name           string     `json:"name"`
	X              uint16     `json:"x"`
	Y              uint16     `json:"y"`
	Z              uint8      `json:"z"`
	Waypoints      []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	ID   string `json:"id"` // Required for DND reordering
	Type string `json:"type"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Z    int    `json:"z"`
}

func (b *Bot) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// --- 1. THE COMMAND READER (Browser -> Go) ---
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return // Connection closed
			}

			// Parse the command
			var cmd struct {
				Type string `json:"type"`
			}
			if err := json.Unmarshal(message, &cmd); err == nil {
				if cmd.Type == "TOGGLE_FISHING" {
					b.fishingEnabled = !b.fishingEnabled
				}
			}
		}
	}()

	// --- 2. THE STATE WRITER (Go -> Browser) ---
	// Stream updates every 100ms
	for {
		// Create a snapshot from your current bot fields
		snap := BotSnapshot{
			FishingEnabled: b.fishingEnabled,
			Name:           b.state.CaptureFrame().Player.Name,
			X:              b.state.CaptureFrame().Player.Pos.X,
			Y:              b.state.CaptureFrame().Player.Pos.Y,
			Z:              b.state.CaptureFrame().Player.Pos.Z,
			Waypoints: []Waypoint{
				{ID: "wp-1", Type: "Walk", X: 32345, Y: 32222, Z: 7},
				{ID: "wp-2", Type: "Walk", X: 32350, Y: 32230, Z: 7},
				{ID: "wp-3", Type: "Rope", X: 32350, Y: 32230, Z: 7},
				{ID: "wp-4", Type: "Walk", X: 32352, Y: 32235, Z: 6},
			},
		}

		payload, _ := json.Marshal(snap)
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}
