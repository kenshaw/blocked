// Package blocked provides a block encoder and decoder for bitmaps.
package blocked

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"maps"
	"math/bits"
	"strings"
	"sync"
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

// Bytes returns the bitmap encoded using the [Bitmap.Best] block type.
func (img Bitmap) Bytes() []byte {
	var buf bytes.Buffer
	if err := img.Encode(&buf, img.Best()); err != nil {
		panic(err)
	}
	return buf.Bytes()
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

// Type is a block type.
type Type rune

// Block types.
const (
	// Auto uses [Best] to determine the best, contiguous block type to use for
	// a provided height.
	Auto Type = 'v'
	// Solids are single, 1x1 blocks using [SolidsRunes].
	Solids Type = 'l'
	// Binaries are single, 1x1 blocks using binary digits using [BinariesRunes].
	Binaries Type = 'b'
	// XXs are single, 1x1 blocks using [XXsRunes].
	XXs Type = 'L'
	// Doubles are single, 0.5x1 double wide blocks using [SolidsRunes].
	Doubles Type = 'D'
	// Halves are 0.5x1 double wide blocks using [HalvesRunes].
	Halves Type = 'e'
	// Halves are 1x2 double height blocks using ASCII-safe runes using
	// [ASCIIsRunes].
	ASCIIs Type = 'E'
	// Quads are 3x2 quarter blocks using [QuadsRunes].
	Quads Type = 'q'
	// QuadsSeparated are 2x2 quarter blocks using [QuadsSeparatedRunes].
	QuadsSeparated Type = 'Q'
	// Sextants are 2x3 blocks using [SextantsRunes].
	Sextants Type = 'x'
	// SextantsSeparated are 2x3 blocks using [SextantsSeparatedRunes].
	SextantsSeparated Type = 'X'
	// Octants are 2x4 blocks using [OctantsRunes].
	Octants Type = 'o'
	// Braille are 2x4 blocks using [BrailleRunes].
	Braille Type = 'O'
)

// Types returns all block types.
func Types() []Type {
	return []Type{
		Solids,
		Binaries,
		XXs,
		Doubles,
		Halves,
		ASCIIs,
		Quads,
		QuadsSeparated,
		Sextants,
		SextantsSeparated,
		Octants,
		Braille,
	}
}

// String satisfies the [fmt.Stringer] interface.
func (typ Type) String() string {
	switch typ {
	case Solids:
		return "Solids"
	case Binaries:
		return "Binaries"
	case XXs:
		return "XXs"
	case Doubles:
		return "Doubles"
	case Halves:
		return "Halves"
	case ASCIIs:
		return "ASCIIs"
	case Quads:
		return "Quads"
	case QuadsSeparated:
		return "QuadsSeparated"
	case Sextants:
		return "Sextants"
	case SextantsSeparated:
		return "SextantsSeparated"
	case Octants:
		return "Octants"
	case Braille:
		return "Braille"
	}
	return ""
}

// Contiguous returns true when the type is a contiguous block type.
func (typ Type) Contiguous() bool {
	switch typ {
	case Solids, Doubles, Halves, Quads, Sextants, Octants:
		return true
	}
	return false
}

// Rune returns the verb rune for the type.
func (typ Type) Rune() rune {
	return rune(typ)
}

// RuneCount returns the number of runes for the block type.
func (typ Type) RuneCount() int {
	switch typ {
	case Solids, Binaries, XXs, Doubles:
		return 2
	case Halves, ASCIIs:
		return 4
	case Quads, QuadsSeparated:
		return 16
	case Sextants, SextantsSeparated:
		return 64
	case Octants, Braille:
		return 256
	}
	return -1
}

// Width returns the width for the block type.
func (typ Type) Width() int {
	switch typ {
	case Doubles:
		return 0
	case Solids, Binaries, XXs, Halves, ASCIIs:
		return 1
	case Quads, QuadsSeparated, Sextants, SextantsSeparated, Octants, Braille:
		return 2
	}
	return -1
}

// Height returns the height for the block type.
func (typ Type) Height() int {
	switch typ {
	case Solids, Binaries, XXs, Doubles:
		return 1
	case Halves, ASCIIs, Quads, QuadsSeparated:
		return 2
	case Sextants, SextantsSeparated:
		return 3
	case Octants, Braille:
		return 4
	}
	return -1
}

// ToRune converts a byte to its block rune.
func (typ Type) ToRune(b uint8) rune {
	if m := typ.runeMap(); m != nil {
		return m[b]
	}
	return 0
}

// RuneMap returns the rune map for the block type.
func (typ Type) RuneMap() map[uint8]rune {
	return maps.Clone(typ.runeMap())
}

// Runes returns the runes for the block type.
func (typ Type) Runes() []rune {
	switch typ {
	case Solids:
		return SolidsRunes()
	case Binaries:
		return BinariesRunes()
	case XXs:
		return XXsRunes()
	case Doubles:
		return SolidsRunes()
	case Halves:
		return HalvesRunes()
	case ASCIIs:
		return ASCIIsRunes()
	case Quads:
		return QuadsRunes()
	case QuadsSeparated:
		return QuadsSeparatedRunes()
	case Sextants:
		return SextantsRunes()
	case SextantsSeparated:
		return SextantsSeparatedRunes()
	case Octants:
		return OctantsRunes()
	case Braille:
		return BrailleRunes()
	}
	return nil
}

// Dump dumps a ASCII drawing of the bitmask the block type's symbols to the
// writer.
func (typ Type) Dump(w io.Writer) {
	if m := typ.runeMap(); m != nil {
		Dump(w, m)
	}
}

// runeMap returns the rune map for the type.
func (typ Type) runeMap() map[uint8]rune {
	switch typ {
	case Solids, Binaries, XXs,
		Doubles,
		Halves, ASCIIs,
		Quads, QuadsSeparated,
		Sextants, SextantsSeparated,
		Octants, Braille:
		blocksMu.Lock()
		defer blocksMu.Unlock()
		b, ok := blocks[typ]
		if !ok {
			v := typ.Runes()
			b = make(map[uint8]rune, len(v))
			for i, r := range v {
				b[uint8(i)] = r
			}
			blocks[typ] = b
		}
		return b
	}
	return nil
}

// Best returns the best display block type for the height.
func Best(y int) Type {
	switch {
	case y == 1:
		return Solids
	case y <= 3:
		return Halves
	case y < 6:
		return Quads
	case y <= 24:
		return Sextants
	}
	return Octants
}

// Dump dumps a ASCII drawing of the bitmask of the symbols to the writer.
//
// Used to verify the symbols for different [Type]'s.
func Dump(w io.Writer, syms map[uint8]rune) {
	n := len(syms)
	for i := range n {
		if i%8 == 0 {
			fmt.Fprintf(w, "%*s|", 3, "")
		}
		fmt.Fprintf(w, "%c", syms[uint8(i)])
		if i%8 == 7 || i == n-1 {
			fmt.Fprintln(w, "|")
		}
	}
	fmt.Fprintln(w)
	width := 8 - bits.LeadingZeros8(uint8(len(syms)-1))
	for i := range n {
		if i != 0 {
			fmt.Fprintln(w)
		}
		v := splitMask(uint8(i), width)
		fmt.Fprintf(w, "%3d: |%s| │%c│\n", i, v[0], syms[uint8(i)])
		for j := 1; j < len(v); j++ {
			fmt.Fprintf(w, "%3s  |%s|\n", "", v[j])
		}
	}
}

// splitMask splits a mask into n lines of m runes per line, where n is the `8 - width`
// significant bits of i. If width is less than or equal to 2, each line will
// be 1 rune, otherwise each line will be 2 runes.
func splitMask(i uint8, width int) []string {
	s := maskRepl.Replace(fmt.Sprintf("%0*b", width, bits.Reverse8(i)>>(8-width)))
	var v []string
	n := 2
	if width <= 2 {
		n = 1
	}
	for i := 0; i < len(s); i += n {
		v = append(v, s[i:min(i+n, len(s))])
	}
	return v
}

// maskRepl replaces '0', '1' with ' ', 'X'.
var maskRepl = strings.NewReplacer(
	"0", " ",
	"1", "X",
)

var (
	// blocks are the block maps.
	blocks = make(map[Type]map[uint8]rune)
	// blocksMu is the blocks mutex.
	blocksMu sync.Mutex
	// nl is the newline.
	nl = []byte{'\n'}
)
