package login

import (
	"bytes"
	"encoding/binary"
	"errors"
	"goTibia/protocol"
	"io"
	"log"
	"strconv"
	"strings"
)

// region LoginResultMessage

type LoginResultMessage struct {
	ClientDisconnected       bool
	ClientDisconnectedReason string
	Motd                     *Motd
	CharacterList            *CharacterList
}

func (lp *LoginResultMessage) Marshal() ([]byte, error) {
	buf := new(bytes.Buffer)

	// 2. Write a 2-byte (uint16) placeholder for the length. We'll overwrite it later.
	// We write two zero bytes for now.
	err := binary.Write(buf, binary.LittleEndian, uint16(0))
	if err != nil {
		return nil, err
	}

	// 3. Write the actual message payload.
	if lp.ClientDisconnected {
		buf.WriteByte(S2COpcodeDisconnectClient)
		protocol.WriteString(buf, lp.ClientDisconnectedReason)
	}

	if lp.CharacterList != nil {
		buf.WriteByte(S2COpcodeCharacterList)
		err := WriteCharacterList(buf, lp.CharacterList)
		if err != nil {
			return nil, err
		}
	}

	// 4. Get the final byte slice from the buffer.
	finalBytes := buf.Bytes()

	// 5. Calculate the length of the PAYLOAD (total length minus the 2 placeholder bytes).
	payloadLength := len(finalBytes) - 2

	binary.LittleEndian.PutUint16(finalBytes, uint16(payloadLength))

	// 7. Return the complete message with the correct length prefix.
	return finalBytes, nil
}

// endregion LoginResultMessage

// region CharacterList
type CharacterList struct {
	Characters  []*CharacterEntry
	PremiumDays uint16
}

type CharacterEntry struct {
	Name      string
	WorldName string
	WorldIp   uint32
	WorldPort uint16
}

func ReadCharacterList(r io.Reader) (*CharacterList, error) {
	entryCount, err := protocol.ReadByte(r)
	if err != nil {
		return nil, err
	}

	var characterEntries []*CharacterEntry
	for i := 0; i < int(entryCount); i++ {
		name, err := protocol.ReadString(r)
		if err != nil {
			return nil, err
		}

		worldName, err := protocol.ReadString(r)
		if err != nil {
			return nil, err
		}

		var worldIp uint32
		if err := binary.Read(r, binary.LittleEndian, &worldIp); err != nil {
			return nil, err
		}

		var worldPort uint16
		if err := binary.Read(r, binary.LittleEndian, &worldPort); err != nil {
			return nil, err
		}

		characterEntries = append(characterEntries, &CharacterEntry{Name: name, WorldName: worldName, WorldIp: worldIp, WorldPort: worldPort})
	}

	var premiumDays uint16
	if err := binary.Read(r, binary.LittleEndian, &premiumDays); err != nil {
		return nil, err
	}

	return &CharacterList{Characters: characterEntries, PremiumDays: premiumDays}, nil
}

func WriteCharacterEntry(w io.Writer, entry *CharacterEntry) error {
	if err := protocol.WriteString(w, entry.Name); err != nil {
		return err
	}
	if err := protocol.WriteString(w, entry.WorldName); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, entry.WorldIp); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, entry.WorldPort); err != nil {
		return err
	}
	return nil
}

func WriteCharacterList(w io.Writer, charList *CharacterList) error {
	entryCount := uint8(len(charList.Characters))
	if err := binary.Write(w, binary.LittleEndian, entryCount); err != nil {
		return err
	}

	for _, entry := range charList.Characters {
		if err := WriteCharacterEntry(w, entry); err != nil {
			return err

		}
	}
	if err := binary.Write(w, binary.LittleEndian, charList.PremiumDays); err != nil {
		return err
	}
	return nil
}

// endregion CharacterList

// region MOTD

type Motd struct {
	MotdId  int
	Message string
}

func ReadMotd(r io.Reader) (*Motd, error) {
	data, err := protocol.ReadString(r)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(data, "\n", 2)

	if len(parts) != 2 {
		return nil, errors.New("invalid format")
	}

	motdId, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.New("failed to parse MOTD ID")
	}

	message := parts[1]

	log.Printf("MOTDID :%d, MOTD: %s", motdId, message)
	return &Motd{MotdId: motdId, Message: message}, nil
}

// endregion MOTD
