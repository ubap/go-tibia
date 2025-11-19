package protocol

import "testing"

func Test(t *testing.T) {

	packet := LoginPacket{
		Protocol:      1,
		ClientOS:      65535,
		ClientVersion: 1234,
		DatSignature:  7,
		SprSignature:  8,
		PicSignature:  9,
		XTEAKey:       [4]uint32{17, 18, 19, 20},
		AccountNumber: 42,
		Password:      "secret",
	}

	marshal, err := packet.Marshal(&keyForClientCommunication.PublicKey)
	if err != nil {
		t.Fatalf("Error marshalling public key: %v", err)
	}

	loginPacket, err := ParseLoginPacket(marshal)
	if err != nil {
		t.Fatalf("Error parsing public key: %v", err)
	}

	if *loginPacket != packet {
		t.Fatalf("Decrypted data does not match original. Got '%v', expected '%v'", loginPacket, packet)
	}
}
