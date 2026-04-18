package lpstream

import (
	"fmt"
	"io"
)

// Decoder reads varint‑framed payloads from an io.Reader.
type Decoder struct {
	r    io.Reader
	buf []byte
	pos  int
}

// Creates a new frame decoder from an io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// ReadFrame returns the next complete frame.
// Blocks until the amount of bytes declared in the prefix have been read.
func (d *Decoder) ReadFrame() ([]byte, error) {
	// Use internal buffer if there's leftover data
	if d.pos < len(d.buf) {
		payload := d.buf[d.pos:]
		d.buf = nil
		d.pos = 0
		return payload, nil
	}

	d.buf = nil
	d.pos = 0

	// Read varint prefix byte‑by‑byte
	var (
		length uint64
		shift  uint
		b      [1]byte
	)
	for {
		if _, err := io.ReadFull(d.r, b[:]); err != nil {
			return nil, err
		}

		length |= uint64(b[0]&0x7F) << shift

		if (b[0] & 0x80) == 0 {
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

// Read implements io.Reader.
// It reads the next frame and copies as much as possible into p.
func (d *Decoder) Read(p []byte) (int, error) {
	// Refill buffer if empty
	if d.pos >= len(d.buf) {
		var err error
		d.buf, err = d.ReadFrame()
		if err != nil {
			return 0, err
		}
		d.pos = 0
	}

	// Copy as much as possible
	n := copy(p, d.buf[d.pos:])
	d.pos += n
	return n, nil
}
