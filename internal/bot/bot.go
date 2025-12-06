package bot

import (
	"log"
	"sync"
	"time"
)

type Bot struct {
	reader  WorldStateReader
	actions ActionDispatcher
	client  ClientManipulator

	// lifecycle management
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func NewBot(r WorldStateReader, a ActionDispatcher, c ClientManipulator) *Bot {
	return &Bot{
		reader:  r,
		actions: a,
		client:  c,
	}
}

func (b *Bot) Start() {
	log.Println("[Bot] Engine started")

	// 1. The Light Hack (For testing S2C injection)
	b.runModule("LightHack", b.loopLightHack)
}

func (b *Bot) Stop() {
	close(b.stopChan) // Signal all loops to stop
	b.wg.Wait()       // Wait for them to finish
	log.Println("[Bot] Engine stopped")
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
			// 1. Abstract Query: Are we logged in?
			if b.reader.GetPlayerID() == 0 {
				continue
			}

			// 2. Abstract Action: Make it bright!
			// We don't care about Packet IDs or Byte construction here.
			b.client.SetLocalPlayerLight(0xFF, 215)
		}
	}
}
