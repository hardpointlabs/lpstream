package lpstream

import (
	"fmt"
	"io"
)

// Reader reads varint‑framed payloads from an io.Reader.
type Decoder struct {
	r io.Reader
}

// NewReader wraps an io.Reader
func NewReader(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// ReadFrame returns the next complete frame.
func (d *Decoder) ReadFrame() ([]byte, error) {
	// Read varint prefix byte‑by‑byte
	var (
		length uint64
		shift  uint
		buf    [1]byte
	)
	for {
		if _, err := io.ReadFull(d.r, buf[:]); err != nil {
			return nil, err
		}

		b := buf[0]
		length |= uint64(b&0x7F) << shift

		if (b & 0x80) == 0 {
			break
		}
		shift += 7
		if shift >= 64 {
			return nil, fmt.Errorf("varint overflow")
		}
	}

	// Now read exactly `length` bytes of payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(d.r, payload); err != nil {
		return nil, err
	}

	return payload, nil
}
