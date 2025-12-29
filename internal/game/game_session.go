package game

import (
	"z07/internal/bot"
	"z07/internal/game/state"
	"z07/internal/protocol"
)

type GameSession struct {
	ID         string
	State      *state.GameState
	Bot        *bot.Bot
	ClientConn protocol.Connection
	ServerConn protocol.Connection
	ErrChan    chan error
}

func newGameSession(client protocol.Connection, server protocol.Connection, gameState *state.GameState) *GameSession {
	return &GameSession{
		ID:         client.RemoteAddr().String(),
		State:      gameState,
		ClientConn: client,
		ServerConn: server,
		ErrChan:    make(chan error, 100),
		Bot:        bot.NewBot(gameState, client, server),
	}
}
