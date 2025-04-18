package blocked

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
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

// NewBitmapFromReader creates a new bitmap from the reader. Reads x pixel per line.
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

// At satisfies the [image.Image] interface.
func (img Bitmap) At(x, y int) color.Color {
	if img.Get(
		x/max(int(DefaultWidthScale), int(img.WidthScale)),
		y/max(int(DefaultHeightScale), int(img.HeightScale)),
	) {
		return img.Opaque
	}
	return img.Transparent
}

// ColorModel satisfies the [image.Image] interface.
func (img Bitmap) ColorModel() color.Model {
	return color.Alpha16Model
}

// Bounds satisfies the [image.Image] interface.
func (img Bitmap) Bounds() image.Rectangle {
	return image.Rect(
		0, 0,
		img.Rect.Dx()*max(int(DefaultWidthScale), int(img.WidthScale)),
		img.Rect.Dy()*max(int(DefaultHeightScale), int(img.HeightScale)),
	)
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
func enc1x1(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) error {
	for y := range h {
		for x := range w {
			fmt.Fprintf(wr, "%c", syms[buf[y*n+x/8]&(1<<(x%8))>>(x%8)])
		}
		if y != h-1 {
			fmt.Fprintln(wr)
		}
	}
	return nil
}

// enc1x2 encodes 1x2 blocks to the writer.
func enc1x2(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) error {
	if h%2 != 0 {
		buf = append(buf, make([]byte, n)...)
	}
	for y := 0; y < h; y += 2 {
		for x := range w {
			b := buf[y*n+x/8]&(1<<(x%8))>>(x%8) |
				buf[(y+1)*n+x/8]&(1<<(x%8))>>(x%8)<<1
			fmt.Fprintf(wr, "%c", syms[b])
		}
		if y < h-2 {
			fmt.Fprintln(wr)
		}
	}
	return nil
}

// enc2x2 encodes 2x2 blocks to the writer.
func enc2x2(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) error {
	if h%2 != 0 {
		buf = append(buf, make([]byte, n)...)
	}
	for y := 0; y < h; y += 2 {
		for x := 0; x < w; x += 2 {
			b := buf[y*n+x/8]&(0b11<<(x%8))>>(x%8) |
				buf[(y+1)*n+x/8]&(0b11<<(x%8))>>(x%8)<<2
			fmt.Fprintf(wr, "%c", syms[b])
		}
		if y < h-2 {
			fmt.Fprintln(wr)
		}
	}
	return nil
}

// enc2x3 encodes 2x3 blocks to the writer.
func enc2x3(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) error {
	if x := h % 3; x != 0 {
		buf = append(buf, make([]byte, n*(3-x))...)
	}
	for y := 0; y < h; y += 3 {
		for x := 0; x < w; x += 2 {
			b := buf[y*n+x/8]&(0b11<<(x%8))>>(x%8) |
				buf[(y+1)*n+x/8]&(0b11<<(x%8))>>(x%8)<<2 |
				buf[(y+2)*n+x/8]&(0b11<<(x%8))>>(x%8)<<4
			fmt.Fprintf(wr, "%c", syms[b])
		}
		if y < h-3 {
			fmt.Fprintln(wr)
		}
	}
	return nil
}

// enc2x4 encodes 2x4 blocks to the writer.
func enc2x4(wr io.Writer, buf []byte, w, h, n int, syms map[uint8]rune) error {
	if x := h % 4; x != 0 {
		buf = append(buf, make([]byte, n*(4-x))...)
	}
	for y := 0; y < h; y += 4 {
		for x := 0; x < w; x += 2 {
			b := buf[(y)*n+x/8]&(0b11<<(x%8))>>(x%8) |
				buf[(y+1)*n+x/8]&(0b11<<(x%8))>>(x%8)<<2 |
				buf[(y+2)*n+x/8]&(0b11<<(x%8))>>(x%8)<<4 |
				buf[(y+3)*n+x/8]&(0b11<<(x%8))>>(x%8)<<6
			fmt.Fprintf(wr, "%c", syms[b])
		}
		if y < h-4 {
			fmt.Fprintln(wr)
		}
	}
	return nil
}
