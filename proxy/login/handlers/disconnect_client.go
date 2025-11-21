package handlers

import (
	"goTibia/protocol"
	"io"
	"log"
)

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
