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
	FishingEnabled   bool       `json:"fishingEnabled"`
	LighthackEnabled bool       `json:"lighthackEnabled"`
	Name             string     `json:"name"`
	X                uint16     `json:"x"`
	Y                uint16     `json:"y"`
	Z                uint8      `json:"z"`
	Waypoints        []Waypoint `json:"waypoints"`
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
	// This ensures that when this function exits, the socket is closed.
	// Closing the socket will also force the "Command Reader" goroutine to exit.
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
				} else if cmd.Type == "TOGGLE_LIGHTHACK" {
					b.lighthackEnabled = !b.lighthackEnabled
				}
			}
		}
	}()

	// --- 2. THE STATE WRITER (Go -> Browser) ---
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		// EXIT if the Bot is stopped via Stop()
		case <-b.stopChan:
			return

		// EXECUTE update every tick
		case <-ticker.C:
			snap := BotSnapshot{
				FishingEnabled:   b.fishingEnabled,
				LighthackEnabled: b.lighthackEnabled,
				Name:             b.state.CaptureFrame().Player.Name,
				X:                b.state.CaptureFrame().Player.Pos.X,
				Y:                b.state.CaptureFrame().Player.Pos.Y,
				Z:                b.state.CaptureFrame().Player.Pos.Z,
				Waypoints: []Waypoint{
					{ID: "wp-1", Type: "Walk", X: 32345, Y: 32222, Z: 7},
					{ID: "wp-2", Type: "Walk", X: 32350, Y: 32230, Z: 7},
					{ID: "wp-3", Type: "Rope", X: 32350, Y: 32230, Z: 7},
					{ID: "wp-4", Type: "Walk", X: 32352, Y: 32235, Z: 6},
				},
			}

			// We use WriteJSON directly to simplify the code
			if err := conn.WriteJSON(snap); err != nil {
				// If the browser tab is closed, this will error out and exit the loop
				return
			}
		}
	}
}
