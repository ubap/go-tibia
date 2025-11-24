package bot

import (
	"goTibia/game/state"
	"goTibia/packets/game"
	"goTibia/protocol"
	"log"
	"sync"
	"time"
)

type Bot struct {
	// Data Source (Read-Only mostly)
	State *state.GameState

	// Output Channels
	ServerConn *protocol.Connection // To send commands (Walk, Attack)
	ClientConn *protocol.Connection // To send feedback (Text overlay)

	// Lifecycle management
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func New(gameState *state.GameState, server, client *protocol.Connection) *Bot {
	return &Bot{
		State:      gameState,
		ServerConn: server,
		ClientConn: client,
		stopChan:   make(chan struct{}),
	}
}

func (b *Bot) Start() {
	log.Println("[Bot] Engine started")

	// Register your modules here:

	// 1. The Light Hack (For testing S2C injection)
	b.runModule("LightHack", b.loopLightHack)

	// 2. Auto Healer (For testing C2S injection - logic from previous steps)
	// b.runModule("AutoHealer", b.loopAutoHealer)
}

func (b *Bot) Stop() {
	close(b.stopChan) // Signal all loops to stop
	b.wg.Wait()       // Wait for them to finish
	log.Println("[Bot] Engine stopped")
}

// Helper to run a loop safely
func (b *Bot) runModule(name string, logic func()) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		log.Printf("[Bot] Module %s running", name)
		logic()
		log.Printf("[Bot] Module %s stopped", name)
	}()
}

// loopLightHack sends a World Light packet every 100ms.
func (b *Bot) loopLightHack() {
	// 1. Setup Ticker (100ms)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	log.Println("[Bot] LightHack started")

	for {
		select {
		// 2. Handle Stop Signal (When player disconnects)
		case <-b.stopChan:
			return

		// 3. Handle Tick
		case <-ticker.C:
			if b.State.Player.ID == 0 {
				continue // Player not logged in yet
			}
			pkt := &game.CreatureLightMsg{
				CreatureID: b.State.Player.ID,
				LightLevel: 0xFF,
				Color:      215,
			}

			// 4. Send to CLIENT (To make the game bright)
			// We use ClientConn because we want to cheat visuals for the user,
			// we don't want to tell the server anything.
			if err := b.ClientConn.SendPacket(pkt); err != nil {
				log.Printf("[Bot] Failed to send light: %v", err)
				return // Stop loop if connection dies
			}
		}
	}
}
