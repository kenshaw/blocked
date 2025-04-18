package blocked

import (
	"bufio"
	"errors"
	"image"
	"image/color"
	"io"
	"unicode/utf8"
)

var (
	// DefaultWidthScale is the default pixel width scale when bitmaps used as [image.Image].
	DefaultWidthScale uint = 24
	// DefaultHeightScale is the default pixel height scale when bitmaps used as [image.Image].
	DefaultHeightScale uint = 24
)

// Bitmap is a monotone bitmap image.
type Bitmap struct {
	Pix         []uint8
	Stride      int
	Rect        image.Rectangle
	WidthScale  uint
	HeightScale uint
	Opaque      color.Alpha16
	Transparent color.Alpha16
}

// NewBitmap creates a new bitmap.
func NewBitmap(rect image.Rectangle) Bitmap {
	stride := rect.Dx()
	if stride%8 != 0 {
		stride = stride/8 + 1
	} else {
		stride /= 8
	}
	return Bitmap{
		Pix:         make([]uint8, stride*rect.Dy()),
		Stride:      stride,
		Rect:        rect,
		WidthScale:  DefaultWidthScale,
		HeightScale: DefaultHeightScale,
		Opaque:      color.Opaque,
		Transparent: color.Transparent,
	}
}

// NewBitmapFromReader creates a new bitmap from the reader. Reads x pixels per line.
//
// TODO: allow compact byte streams.
func NewBitmapFromReader(r io.Reader, x int) (Bitmap, error) {
	var stride int
	if x%8 != 0 {
		stride = x/8 + 1
	} else {
		stride = x / 8
	}
	pix, buf := make([]byte, 0, 512), make([]byte, stride)
	var err error
loop:
	for br := bufio.NewReader(r); ; {
		switch _, err = br.Read(buf); {
		case errors.Is(err, io.EOF):
			break loop
		case err != nil:
			return Bitmap{}, err
		}
		pix = append(pix, buf...)
	}
	return Bitmap{
		Pix:         pix,
		Stride:      stride,
		WidthScale:  DefaultWidthScale,
		HeightScale: DefaultHeightScale,
		Rect:        image.Rect(0, 0, x, len(pix)/stride),
	}, nil
}

// Set sets the bit at x, y.
func (img Bitmap) Set(x, y int, v bool) {
	if v {
		img.Pix[y*img.Stride+x/8] |= 1 << (x % 8)
	} else {
		img.Pix[y*img.Stride+x/8] &= ^(1 << (x % 8))
	}
}

// Get returns the bit at x, y.
func (img Bitmap) Get(x, y int) bool {
	return img.Pix[y*img.Stride+x/8]&(1<<(x%8)) != 0
}

// Byte returns the byte containing x, y.
func (img Bitmap) Byte(x, y int) uint8 {
	return img.Pix[y*img.Stride+x/8]
}

// ColorModel satisfies the [image.Image] interface.
func (img Bitmap) ColorModel() color.Model {
	return color.Alpha16Model
}

// Bounds satisfies the [image.Image] interface.
func (img Bitmap) Bounds() image.Rectangle {
	w, h := img.scale()
	return image.Rect(0, 0, w*img.Rect.Dx(), h*img.Rect.Dy())
}

// At satisfies the [image.Image] interface.
func (img Bitmap) At(x, y int) color.Color {
	if w, h := img.scale(); img.Get(x/w, y/h) {
		return img.Opaque
	}
	return img.Transparent
}

// scale returns the width and height scale.
func (img Bitmap) scale() (int, int) {
	w, h := img.WidthScale, img.HeightScale
	if w == 0 {
		w = max(1, DefaultWidthScale)
	}
	if h == 0 {
		h = max(1, DefaultHeightScale)
	}
	return int(w), int(h)
}

// Encode encodes the bitmap to the writer using the block type.
func (img Bitmap) Encode(w io.Writer, typ Type) error {
	h := img.Rect.Dy()
	if typ == Auto {
		typ = Best(h)
	}
	var f func(io.Writer, []byte, int, int, int, map[uint8]rune) error
	// TODO: unify the func's below
	switch w, h := typ.Width(), typ.Height(); {
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
	return f(w, img.Pix, img.Rect.Dx(), h, img.Stride, typ.runeMap())
}

// Width returns the width for the block type.
func (img Bitmap) Width(typ Type) int {
	x, w := img.Rect.Dx(), typ.Width()
	if x%w != 0 {
		return x/w + 1
	}
	return x / w
}

// Height returns the height for the block type.
func (img Bitmap) Height(typ Type) int {
	y, h := img.Rect.Dy(), typ.Height()
	if y%h != 0 {
		return y/h + 1
	}
	return y / h
}

// enc1x1 encodes 1x1 blocks to the writer.
func enc1x1(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) (err error) {
	m, b, v := 0, uint8(0), make([]byte, 4)
	for y := range h {
		for x := range w {
			m = x % 8
			b = buf[y*n+x/8] & (1 << m) >> m
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
func enc1x2(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) (err error) {
	if h%2 != 0 {
		buf = append(buf, make([]byte, n)...)
	}
	d, m, b, v := 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 2 {
		for x := range w {
			d, m = x/8, x%8
			b = buf[y*n+d]&(1<<m)>>m |
				buf[(y+1)*n+d]&(1<<m)>>m<<1
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
func enc2x2(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) (err error) {
	if h%2 != 0 {
		buf = append(buf, make([]byte, n)...)
	}
	d, m, b, v := 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 2 {
		for x := 0; x < w; x += 2 {
			d, m = x/8, x%8
			b = buf[y*n+d]&(0b11<<m)>>m |
				buf[(y+1)*n+d]&(0b11<<m)>>m<<2
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
func enc2x3(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) (err error) {
	if x := h % 3; x != 0 {
		buf = append(buf, make([]byte, n*(3-x))...)
	}
	d, m, b, v := 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 3 {
		for x := 0; x < w; x += 2 {
			d, m = x/8, x%8
			b = buf[y*n+d]&(0b11<<m)>>m |
				buf[(y+1)*n+d]&(0b11<<m)>>m<<2 |
				buf[(y+2)*n+d]&(0b11<<m)>>m<<4
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
func enc2x4(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) (err error) {
	if x := h % 4; x != 0 {
		buf = append(buf, make([]byte, n*(4-x))...)
	}
	d, m, b, v := 0, 0, uint8(0), make([]byte, 4)
	for y := 0; y < h; y += 4 {
		for x := 0; x < w; x += 2 {
			d, m = x/8, x%8
			b = buf[(y)*n+d]&(0b11<<m)>>m |
				buf[(y+1)*n+d]&(0b11<<m)>>m<<2 |
				buf[(y+2)*n+d]&(0b11<<m)>>m<<4 |
				buf[(y+3)*n+d]&(0b11<<m)>>m<<6
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
