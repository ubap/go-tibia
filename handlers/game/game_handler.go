package game

import (
	"goTibia/packets/game"
	"goTibia/protocol"
	"goTibia/proxy"
	"log"
)

type GameHandler struct {
	TargetAddr string
	// You could add "DB *sql.DB" here later!
}

func (h *GameHandler) Handle(client *protocol.Connection) {
	log.Printf("[Game] New Connection: %s", client.RemoteAddr())

	packetReader, err := client.ReadMessage()
	if err != nil {
		log.Printf("Game: error reading message from %s: %v", client.RemoteAddr(), err)
		return
	}

	loginRequest, err := game.ParseLoginRequest(packetReader)
	if err != nil {
		log.Printf("Game: Failed to parse login packet: %v", err)
		return
	}

	protoServerConn, err := proxy.ConnectToBackend(h.TargetAddr)
	if err != nil {
		log.Printf("Game: Failed to connect to %s: %v", client.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	if err := protoServerConn.SendPacket(loginRequest); err != nil {
		log.Printf("Game: Failed to forward credentials to backend: %v", err)
		return
	}

	log.Println("Game: LoginRequest forwarded to backend.")

	protoServerConn.EnableXTEA(loginRequest.XTEAKey)
	client.EnableXTEA(loginRequest.XTEAKey)

	message, err := protoServerConn.ReadMessage()
	if err != nil {
		log.Printf("Game: Failed to read server response for %s: %v", client.RemoteAddr(), err)
		return
	}

	_, err = game.ParseLoginResultMessage(message)
	if err != nil {
		log.Printf("Game: Failed to receive login result message for %s: %v", client.RemoteAddr(), err)
		return
	}

}
