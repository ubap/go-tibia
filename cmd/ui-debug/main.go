package main

import (
	"goTibia/internal/bot"
	"goTibia/internal/game/state"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	gs := state.New()

	b := bot.NewBot(gs, nil, nil)
	b.StartUIOnly()

	wg.Wait()
}
