package protocol

// Encodable represents anything that can write itself to a PacketWriter.
type Encodable interface {
	// Encode writes the packet data to the provided writer.
	// It does not return []byte.
	// It does not return error (errors are stored in the PacketWriter state).
	Encode(pw *PacketWriter)
}
