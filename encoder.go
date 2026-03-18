package lpstream

import (
	"encoding/binary"
	"io"
)

// Writer writes length‑prefixed frames using protobuf varints.
type Encoder struct {
	w io.Writer
}

// NewWriter wraps an io.Writer
func NewWriter(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// WriteFrame writes a single frame with varint length prefix.
func (e *Encoder) WriteFrame(payload []byte) error {
	// varint length prefix
	lenBuf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(lenBuf, uint64(len(payload)))

	// write prefix
	if _, err := e.w.Write(lenBuf[:n]); err != nil {
		return err
	}

	// write payload
	_, err := e.w.Write(payload)
	return err
}
