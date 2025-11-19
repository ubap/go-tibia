package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

// Connection is a wrapper around a raw network connection (net.Conn)
// that understands the Tibia protocol's message framing.
type Connection struct {
	conn net.Conn
}

// NewConnection creates a new protocol-aware connection wrapper.
func NewConnection(conn net.Conn) *Connection {
	return &Connection{conn: conn}
}

// ReadMessage reads a single, complete message from the stream.
// It handles the 2-byte length prefix and returns the message payload.
func (c *Connection) ReadMessage() ([]byte, error) {
	var length uint16
	// Read the 2-byte length prefix.
	if err := binary.Read(c.conn, binary.LittleEndian, &length); err != nil {
		// An io.EOF here is a clean disconnect.
		return nil, err
	}

	// Read the message body of the specified length.
	payload := make([]byte, length)
	if _, err := io.ReadFull(c.conn, payload); err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *Connection) WriteMessage(payload []byte) error {
	length := uint16(len(payload))

	// Write the 2-byte length prefix.
	if err := binary.Write(c.conn, binary.LittleEndian, length); err != nil {
		return err
	}

	// Write the actual message payload.
	_, err := c.conn.Write(payload)
	return err
}

// Close simply closes the underlying network connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// RemoteAddr returns the remote network address.
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
