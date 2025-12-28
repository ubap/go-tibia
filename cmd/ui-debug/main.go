package main

import (
	"goTibia/internal/bot"
	"goTibia/internal/game/domain"
	"goTibia/internal/game/state"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	gs := state.New()
	gs.SetPlayerName("JohnDoe")
	gs.SetPlayerPos(domain.Position{X: 5, Y: 6, Z: 7})

	b := bot.NewBot(gs, nil, nil)
	b.StartUIOnly()

	wg.Wait()
}
