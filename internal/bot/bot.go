package bot

import (
	"goTibia/internal/game/packets"
	"goTibia/internal/game/state"
	"log"
	"sync"
	"time"
)

type Bot struct {
	Outbox      chan packets.C2SPacket
	State       *state.GameState
	UserActions chan packets.C2SPacket

	stopChan chan struct{}  // The broadcast channel
	wg       sync.WaitGroup // To wait for modules to finish
	stopOnce sync.Once      // To ensure we close the channel only once
}

func NewBot(state *state.GameState) *Bot {
	return &Bot{
		State: state,

		Outbox:      make(chan packets.C2SPacket, 100),
		UserActions: make(chan packets.C2SPacket, 100),

		stopChan: make(chan struct{}),
	}
}

func (b *Bot) Start() {
	log.Println("[Bot] Engine started")

	// 1. The Light Hack (For testing S2C injection)
	b.runModule("LightHack", b.loopLightHack)
	b.runModule("HandleUserAction", b.loopHandleUserAction)
}

func (b *Bot) Stop() {
	b.stopOnce.Do(func() {
		log.Println("[Bot] Stopping engine...")
		close(b.stopChan) // This broadcasts the signal to ALL loops instantly
	})

	b.wg.Wait()
	log.Println("[Bot] Engine stopped cleanly.")
}

func (b *Bot) runModule(name string, logic func()) {
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		log.Printf("[Bot] Module %s running", name)
		logic()
		log.Printf("[Bot] Module %s stopped", name)
	}()
}

func (b *Bot) loopLightHack() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	log.Println("[Bot] LightHack started")

	for {
		select {
		case <-b.stopChan:
			return

		case <-ticker.C:
			pId := b.State.CaptureFrame().Player.ID

			if pId == 0 {
				continue
			}

			//b.client.SetLocalPlayerLight(0xFF, 215)
		}
	}
}

func (b *Bot) loopHandleUserAction() {
	for {
		select {
		case <-b.stopChan:
			return
		case packet := <-b.UserActions:
			b.handleUserAction(packet)
		}
	}
}

func (b *Bot) handleUserAction(packet packets.C2SPacket) {
	switch p := packet.(type) {
	case *packets.LookRequest:
		log.Printf("User looked at item ID: %d", p.ItemId)
	}
}
