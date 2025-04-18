package blocked

import (
	"io"
)

// Encoder is a bitmap block encoder.
type Encoder struct {
	Type Type
}

// NewEncoder creates a new bitmap block encoder.
func NewEncoder(opts ...EncoderOption) Encoder {
	enc := new(Encoder)
	for _, o := range opts {
		o(enc)
	}
	return *enc
}

// Encode encodes the bitmap to the writer.
func (enc Encoder) Encode(w io.Writer, img Bitmap) error {
	return nil
}

// Decoder is a bitmap block decoder.
type Decoder struct {
	Type Type
}

// NewDecoder creates a bitmap block decoder.
func NewDecoder(opts ...DecoderOption) Decoder {
	dec := new(Decoder)
	for _, o := range opts {
		o(dec)
	}
	return *dec
}

// Decode decodes a bitmap from the reader.
func (dec Decoder) Decode(r io.Reader, img *Bitmap) error {
	// TODO: implement this ...
	return nil
}

// EncoderOption is a [Encoder] option.
type EncoderOption func(*Encoder)

// DecoderOption is a [Decoder] option.
type DecoderOption func(*Decoder)
