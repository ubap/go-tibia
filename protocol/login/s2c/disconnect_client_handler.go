package s2c

import (
	"goTibia/protocol"
	"goTibia/protocol/login"
	"io"
	"log"
)

func init() {
	login.S2CHandlers.Register(ServerOpcodeDisconnectClient, &DisconnectClientHandler{})
}

type DisconnectClientHandler struct {
}

func (h *DisconnectClientHandler) Handle(r io.Reader) error {

	readString, err := protocol.ReadString(r)
	if err != nil {
		return err

	}
	log.Print("DisconnectClientHandler: " + readString)
	return nil
}
