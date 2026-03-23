package lpstream

import "io"

// Union of an Encoder and a Decoder. This is useful for bidirectional streams.
type FrameCodec struct {
	// The frame encoder
	*Encoder
	// The frame decoder
	*Decoder
}

// Create a new frame codec from a bidirectional io.ReadWriter.
func NewFrameCodec(readWriter io.ReadWriter) *FrameCodec {
	return &FrameCodec{
		Encoder: NewEncoder(readWriter),
		Decoder: NewDecoder(readWriter)}

}
