package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
)

type PacketReader struct {
	reader *bytes.Reader
	err    error
}

func NewPacketReader(data []byte) *PacketReader {
	return &PacketReader{
		reader: bytes.NewReader(data),
		err:    nil,
	}
}

// Err returns the first error encountered during reading.
func (pr *PacketReader) Err() error {
	return pr.err
}

// ReadByte reads a single byte. Returns 0 if an error occurred previously.
func (pr *PacketReader) ReadByte() uint8 {
	if pr.err != nil {
		return 0
	}

	b, err := pr.reader.ReadByte()
	if err != nil {
		pr.err = err
		return 0
	}
	return b
}

func (pr *PacketReader) ReadBool() bool {
	if pr.err != nil {
		return false
	}

	return pr.ReadByte() != 0
}

// ReadUint16 reads a uint16 (LittleEndian). Returns 0 on error.
func (pr *PacketReader) ReadUint16() uint16 {
	if pr.err != nil {
		return 0
	}

	var b [2]byte
	// Use ReadFull to ensure we get all bytes
	if _, err := io.ReadFull(pr.reader, b[:]); err != nil {
		pr.err = err
		return 0
	}
	return binary.LittleEndian.Uint16(b[:])
}

// ReadUint32 reads a uint32 (LittleEndian). Returns 0 on error.
func (pr *PacketReader) ReadUint32() uint32 {
	if pr.err != nil {
		return 0
	}

	var b [4]byte
	if _, err := io.ReadFull(pr.reader, b[:]); err != nil {
		pr.err = err
		return 0
	}
	return binary.LittleEndian.Uint32(b[:])
}

// ReadString reads a uint16 length, followed by the string bytes.
func (pr *PacketReader) ReadString() string {
	if pr.err != nil {
		return ""
	}

	// 1. Read Length using our own internal method
	length := pr.ReadUint16()

	// If ReadUint16 failed, pr.err is set. Return immediately.
	if pr.err != nil {
		return ""
	}

	if length == 0 {
		return ""
	}

	// 2. Read Body
	buf := make([]byte, length)
	if _, err := io.ReadFull(pr.reader, buf); err != nil {
		pr.err = err
		return ""
	}

	return string(buf)
}

// Skip advances the reader by n bytes.
func (pr *PacketReader) Skip(n int) {
	if pr.err != nil {
		return
	}
	_, err := pr.reader.Seek(int64(n), io.SeekCurrent)
	if err != nil {
		pr.err = err
	}
}

func (pr *PacketReader) ReadAll() []byte {
	if pr.err != nil {
		return nil
	}

	// Optimization: bytes.Reader knows exactly how many bytes are left.
	// We can allocate the exact size once, avoiding array resizing.
	remaining := pr.reader.Len()

	if remaining == 0 {
		return []byte{}
	}

	buf := make([]byte, remaining)

	// We use ReadFull, but since we sized 'buf' to exactly 'remaining',
	// this effectively reads everything until EOF.
	_, err := io.ReadFull(pr.reader, buf)
	if err != nil {
		pr.err = err
		return nil
	}

	return buf
}

func (pr *PacketReader) PeekUint16() (uint16, error) {
	// 1. Snapshot position
	currentPos, err := pr.reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// 2. Read forward
	var b [2]byte
	_, err = pr.reader.Read(b[:])
	if err != nil {
		return 0, err
	}

	// 3. Rewind back
	_, err = pr.reader.Seek(currentPos, 0) // 0 = io.SeekStart
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint16(b[:]), nil
}

func (pr *PacketReader) PeekBytes(n int) ([]byte, error) {
	// 1. Save current position
	// 1 = io.SeekCurrent
	currentPos, err := pr.reader.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	// 2. Read the bytes
	buf := make([]byte, n)
	readBytes, err := pr.reader.Read(buf)
	// Ignore EOF here, we just want to see what's available
	if err != nil && readBytes == 0 {
		return nil, err
	}

	// 3. Rewind to original position
	// 0 = io.SeekStart (Absolute offset)
	_, err = pr.reader.Seek(currentPos, 0)
	if err != nil {
		return nil, err
	}

	// Return only the bytes we actually managed to read
	return buf[:readBytes], nil
}

func (pr *PacketReader) Remaining() int {
	return pr.reader.Len()
}
