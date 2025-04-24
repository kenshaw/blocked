package blocked

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"unicode/utf8"
)

var (
	// DefaultScaleWidth is the pixel width scale for bitmaps used as
	// [image.Image].
	DefaultScaleWidth uint = 24
	// DefaultScaleHeight is the pixel height scale for bitmaps used
	// as [image.Image].
	DefaultScaleHeight uint = 24
	// DefaultOpaque is the default opaque color.
	DefaultOpaque = color.Opaque
	// DefaultTransparent is the default transparent color.
	DefaultTransparent = color.Transparent
)

// Bitmap is a monotone bitmap image.
type Bitmap struct {
	Pix         []uint8
	Stride      int
	Rect        image.Rectangle
	ScaleWidth  uint
	ScaleHeight uint
	Opaque      color.Alpha16
	Transparent color.Alpha16
}

// New creates a new bitmap from data with width x. Data can be any type that
// works with [binary.Write].
func New(data any, x int) (Bitmap, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.NativeEndian, data); err != nil {
		return Bitmap{}, err
	}
	return NewReader(&buf, x)
}

// NewReader creates a new bitmap from the reader with bit width x.
func NewReader(r io.Reader, x int) (Bitmap, error) {
	data, buf := make([]byte, 0, 512), make([]byte, 512)
	var err error
	var c int
loop: // read
	for {
		switch c, err = r.Read(buf); {
		case errors.Is(err, io.EOF):
			break loop
		case err != nil:
			return Bitmap{}, err
		}
		data = append(data, buf[:c]...)
	}
	return NewBytes(data, x, (len(data)*8+7)/x)
}

// NewBytes creates a new bitmap from for the unaligned bytes in data with
// width x, height y.
func NewBytes(data []byte, x, y int) (Bitmap, error) {
	pix := make([]byte, (x*y+7)/8)
	copy(pix, data)
	if m := x * y % 8; m != 0 {
		// clear upper bits of last byte
		pix[len(pix)-1] &= 0xff >> (8 - m)
	}
	return Bitmap{
		Pix:         pix,
		Stride:      x,
		Rect:        image.Rect(0, 0, x, y),
		Opaque:      DefaultOpaque,
		Transparent: DefaultTransparent,
	}, nil
}

// NewImage creates a blank bitmap image with dimensions in rect.
func NewImage(rect image.Rectangle) Bitmap {
	x := rect.Dx()
	return Bitmap{
		Pix:         make([]uint8, (x*rect.Dy()+7)/8),
		Stride:      x,
		Rect:        rect,
		Opaque:      DefaultOpaque,
		Transparent: DefaultTransparent,
	}
}

// Set sets the bit at x, y.
func (img Bitmap) Set(x, y int, b bool) {
	if i := y*img.Stride + x; b {
		img.Pix[i/8] |= 1 << (i % 8)
	} else {
		img.Pix[i/8] &= ^(1 << (i % 8))
	}
}

// Get returns the bit at x, y.
func (img Bitmap) Get(x, y int) bool {
	i := y*img.Stride + x
	return img.Pix[i/8]&(1<<(i%8)) != 0
}

// ColorModel satisfies the [image.Image] interface.
func (img Bitmap) ColorModel() color.Model {
	return color.Alpha16Model
}

// Bounds satisfies the [image.Image] interface.
func (img Bitmap) Bounds() image.Rectangle {
	w, h := img.Scale()
	return image.Rect(0, 0, w*img.Rect.Dx(), h*img.Rect.Dy())
}

// At satisfies the [image.Image] interface.
func (img Bitmap) At(x, y int) color.Color {
	if w, h := img.Scale(); img.Get(x/w, y/h) {
		return img.Opaque
	}
	return img.Transparent
}

// Width returns the width for the block type.
func (img Bitmap) Width(typ Type) int {
	w := typ.Width()
	if w == 0 {
		return img.Stride * 2
	}
	return (img.Stride + w - 1) / w
}

// Height returns the height for the block type.
func (img Bitmap) Height(typ Type) int {
	h := typ.Height()
	return (img.Rect.Dy() + h - 1) / h
}

// Format satisfies the [fmt.Formatter] interface.
func (img Bitmap) Format(f fmt.State, verb rune) {
	switch typ := Type(verb); typ {
	case Auto, 's':
		typ = img.Best()
		fallthrough
	case Solids, Binaries, XXs,
		Doubles,
		Halves, ASCIIs,
		Quads, QuadsSeparated,
		Sextants, SextantsSeparated,
		Octants, Braille:
		if err := img.Encode(f, typ); err != nil {
			fmt.Fprintf(f, "%%!%c(ERROR: %v)", verb, err)
		}
	default:
		fmt.Fprintf(f, "%%!%c(BAD VERB)", verb)
	}
}

// Encode encodes the bitmap to the writer using the block type.
func (img Bitmap) Encode(w io.Writer, typ Type) error {
	if typ == Auto {
		typ = img.Best()
	}
	var f func(io.Writer, []byte, int, int, map[uint8]rune) error
	switch w, h := typ.Width(), typ.Height(); {
	case w == 0:
		f = enc0_5x1
	case w == 1 && h == 1:
		f = enc1x1
	case w == 1 && h == 2:
		f = enc1x2
	case h == 2:
		f = enc2x2
	case h == 3:
		f = enc2x3
	case h == 4:
		f = enc2x4
	}
	return f(w, img.Pix, img.Stride, img.Rect.Dy(), typ.runeMap())
}

// Best returns the [Best] block type for the image.
func (img Bitmap) Best() Type {
	return Best(img.Rect.Dy())
}

// Scale returns the image width and height scale factors.
func (img Bitmap) Scale() (int, int) {
	w, h := img.ScaleWidth, img.ScaleHeight
	if w == 0 {
		w = max(1, DefaultScaleWidth)
	}
	if h == 0 {
		h = max(1, DefaultScaleHeight)
	}
	return int(w), int(h)
}

// enc0_5x1 encodes 0.5x1 blocks to the writer, used when width == 0 for the
// [Type].
func enc0_5x1(wr io.Writer, buf []byte, w, h int, syms map[uint8]rune) (err error) {
	i, m, b, v, o := 0, 0, uint8(0), make([]byte, 8), 0
	for y := range h {
		for x := range w {
			i = y*w + x
			m = i % 8
			b = buf[i/8] & (1 << m) >> m
			o = utf8.EncodeRune(v, syms[b])
			copy(v[o:], v[:o])
			if _, err = wr.Write(v[:2*o]); err != nil {
				return err
			}
		}
		if y < h-1 {
			if _, err = wr.Write(nl); err != nil {
				return err
			}
		}
	}
	return nil
}

// enc1x1 encodes 1x1 blocks to the writer.
func enc1x1(wr io.Writer, buf []byte, w, h int, syms map[uint8]rune) (err error) {
	i, m, b, v := 0, 0, uint8(0), make([]byte, 4)
	for y := range h {
		for x := range w {
			i = y*w + x
			m = i % 8
			b = buf[i/8] & (1 << m) >> m
			if _, err = wr.Write(v[:utf8.EncodeRune(v, syms[b])]); err != nil {
				return err
			}
		}
		if y < h-1 {
			if _, err = wr.Write(nl); err != nil {
				return err
			}
		}
	}
	return nil
}

// enc1x2 encodes 1x2 blocks to the writer.
func enc1x2(wr io.Writer, buf []byte, w, h int, syms map[uint8]rune) (err error) {
	buf = append(buf, make([]byte, (w+7)/8*(h+1)/2)...)
	i, d, m, b, v := 0, 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 2 {
		for x := range w {
			i = y*w + x
			d, m = i/8, i%8
			b = buf[d] & (1 << m) >> m
			d, m = (i+w)/8, (i+w)%8
			b |= buf[d] & (1 << m) >> m << 1
			if _, err = wr.Write(v[:utf8.EncodeRune(v, syms[b])]); err != nil {
				return err
			}
		}
		if y < h-2 {
			if _, err = wr.Write(nl); err != nil {
				return err
			}
		}
	}
	return nil
}

// enc2x2 encodes 2x2 blocks to the writer.
func enc2x2(wr io.Writer, buf []byte, w, h int, syms map[uint8]rune) (err error) {
	buf = append(buf, make([]byte, (w+7)/8*(h+3)/2)...)
	i, d, m, b, v := 0, 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 2 {
		for x := 0; x < w; x += 2 {
			i = y*w + x
			d, m = i/8, i%8
			b = buf[d] & (1 << m) >> m
			d, m = (i+w)/8, (i+w)%8
			b |= buf[d] & (1 << m) >> m << 2
			if x+2 <= w {
				d, m = (i+1)/8, (i+1)%8
				b |= buf[d] & (1 << m) >> m << 1
				d, m = (i+w+1)/8, (i+w+1)%8
				b |= buf[d] & (1 << m) >> m << 3
			}
			if _, err = wr.Write(v[:utf8.EncodeRune(v, syms[b])]); err != nil {
				return err
			}
		}
		if y < h-2 {
			if _, err = wr.Write(nl); err != nil {
				return err
			}
		}
	}
	return nil
}

// enc2x3 encodes 2x3 blocks to the writer.
func enc2x3(wr io.Writer, buf []byte, w, h int, syms map[uint8]rune) (err error) {
	buf = append(buf, make([]byte, (w+7)/8*(h+5)/3)...)
	i, d, m, b, v := 0, 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 3 {
		for x := 0; x < w; x += 2 {
			i = y*w + x
			d, m = i/8, i%8
			b = buf[d] & (1 << m) >> m
			d, m = (i+w)/8, (i+w)%8
			b |= buf[d] & (1 << m) >> m << 2
			d, m = (i+2*w)/8, (i+2*w)%8
			b |= buf[d] & (1 << m) >> m << 4
			if x+2 <= w {
				d, m = (i+1)/8, (i+1)%8
				b |= buf[d] & (1 << m) >> m << 1
				d, m = (i+w+1)/8, (i+w+1)%8
				b |= buf[d] & (1 << m) >> m << 3
				d, m = (i+2*w+1)/8, (i+2*w+1)%8
				b |= buf[d] & (1 << m) >> m << 5
			}
			if _, err = wr.Write(v[:utf8.EncodeRune(v, syms[b])]); err != nil {
				return err
			}
		}
		if y < h-3 {
			if _, err = wr.Write(nl); err != nil {
				return err
			}
		}
	}
	return nil
}

// enc2x4 encodes 2x4 blocks to the writer.
func enc2x4(wr io.Writer, buf []byte, w, h int, syms map[uint8]rune) (err error) {
	buf = append(buf, make([]byte, (w+7)/8*(h+11)/4)...)
	i, d, m, b, v := 0, 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 4 {
		for x := 0; x < w; x += 2 {
			i = y*w + x
			d, m = i/8, i%8
			b = buf[d] & (1 << m) >> m
			d, m = (i+w)/8, (i+w)%8
			b |= buf[d] & (1 << m) >> m << 2
			d, m = (i+2*w)/8, (i+2*w)%8
			b |= buf[d] & (1 << m) >> m << 4
			d, m = (i+3*w)/8, (i+3*w)%8
			b |= buf[d] & (1 << m) >> m << 6
			if x+2 <= w {
				d, m = (i+1)/8, (i+1)%8
				b |= buf[d] & (1 << m) >> m << 1
				d, m = (i+w+1)/8, (i+w+1)%8
				b |= buf[d] & (1 << m) >> m << 3
				d, m = (i+2*w+1)/8, (i+2*w+1)%8
				b |= buf[d] & (1 << m) >> m << 5
				d, m = (i+3*w+1)/8, (i+3*w+1)%8
				b |= buf[d] & (1 << m) >> m << 7
			}
			if _, err = wr.Write(v[:utf8.EncodeRune(v, syms[b])]); err != nil {
				return err
			}
		}
		if y < h-4 {
			if _, err = wr.Write(nl); err != nil {
				return err
			}
		}
	}
	return nil
}

// nl is the newline.
var nl = []byte{'\n'}
