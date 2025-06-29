package blocked

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func TestBlocks(t *testing.T) {
	t.Parallel()
	for _, typ := range Types() {
		t.Run(typ.String(), func(t *testing.T) {
			t.Parallel()
			m, exp := typ.runeMap(), typ.RuneCount()
			for i := range exp {
				if _, ok := m[uint8(i)]; !ok {
					t.Errorf("expected %08x to be defined in symbol map", i)
				}
			}
			if n := len(m); n != exp {
				t.Errorf("expected %d, got: %d", exp, n)
			}
			var buf bytes.Buffer
			typ.Dump(&buf)
			testWrite(t, buf.Bytes())
		})
	}
}

func TestNewImage(t *testing.T) {
	t.Parallel()
	const mask = 0b10101010
	for _, h := range []int{2, 4, 6, 8, 10, 12, 16, 24, 28, 34, 56} {
		for _, w := range []int{2, 4, 6, 8, 10, 12, 16, 24, 28, 34, 56} {
			t.Run(fmt.Sprintf("%02d_%02d", w, h), func(t *testing.T) {
				t.Parallel()
				n := (w*h + 7) / 8
				pix := make([]byte, n)
				for i := range pix {
					pix[i] = mask
				}
				if m := w * h % 8; m != 0 {
					pix[n-1] &= 0xff >> (8 - m)
				}
				img1 := Bitmap{
					Pix:         pix,
					Stride:      w,
					Rect:        image.Rect(0, 0, w, h),
					Opaque:      DefaultOpaque,
					Transparent: DefaultTransparent,
				}
				for y := range img1.Rect.Dy() {
					for x := range img1.Rect.Dx() {
						v := img1.Get(x, y)
						if x%2 != 1 && v {
							t.Errorf("expected false at (%d,%d)", x, y)
						}
						if x%2 == 1 && !v {
							t.Errorf("expected true at (%d,%d)", x, y)
						}
					}
				}
				// create
				img2 := NewImage(image.Rect(0, 0, w, h))
				for y := range img2.Rect.Dy() {
					for x := range img2.Rect.Dx() {
						img2.Set(x, y, x%2 == 1)
					}
				}
				// compare
				if !slices.Equal(img1.Pix, img2.Pix) {
					t.Errorf("expected:\n%b\ngot:\n%b", img1.Pix, img2.Pix)
				}
				/*
					var buf bytes.Buffer
					if err := png.Encode(&buf, img2); err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					name := filepath.Join("testdata", fmt.Sprintf("test_%02d_%02d.png", w, h))
					if err := os.WriteFile(name, buf.Bytes(), 0o644); err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
				*/
			})
		}
	}
}

func TestBitmap(t *testing.T) {
	t.Parallel()
	for seed := 1330; seed <= 1343; seed++ {
		t.Run(strconv.Itoa(seed), func(t *testing.T) {
			t.Parallel()
			r := rand.New(rand.NewSource(int64(seed)))
			img := NewImage(image.Rect(0, 0, 1+r.Intn(28), 1+r.Intn(28)))
			w, h, expTr, expFa := img.Rect.Dx(), img.Rect.Dy(), 0, 0
			for y := range h {
				for x := range w {
					exp := r.Intn(3) != 0
					if exp {
						expTr++
					} else {
						expFa++
					}
					img.Set(x, y, exp)
					if b := img.Get(x, y); b != exp {
						t.Errorf("(%d,%d) expected %t, got: %t", x, y, exp, b)
					}
				}
			}
			t.Logf("w: %d h: %d -- %d/%d", w, h, expTr, expFa)
			if expTr == 0 || expFa == 0 {
				t.Fatal("invalid test -- no significant bits")
			}
			if n, exp := expTr+expFa, img.Rect.Dx()*img.Rect.Dy(); n != exp {
				t.Fatalf("expected %d, got: %d", exp, n)
			}
			var tr, fa int
			for y := range h {
				for x := range w {
					if img.Get(x, y) {
						tr++
					} else {
						fa++
					}
				}
			}
			if tr != expTr {
				t.Errorf("expected %d trues, got: %d", expTr, tr)
			}
			if fa != expFa {
				t.Errorf("expected %d falses, got: %d", expFa, fa)
			}
			// export as png
			var buf bytes.Buffer
			if err := png.Encode(&buf, img); err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			name := filepath.Join("testdata", fmt.Sprintf("test_%4d.png", seed))
			if err := os.WriteFile(name, buf.Bytes(), 0o644); err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			// compare golden
			expBuf, err := os.ReadFile(name + ".golden")
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if !bytes.Equal(buf.Bytes(), expBuf) {
				t.Errorf("expected %s and %s to be the same", name+".golden", name)
			}
			for _, typ := range Types() {
				t.Run(typ.String(), func(t *testing.T) {
					var buf bytes.Buffer
					w := img.Width(typ)
					buf.WriteString(strings.Repeat("_", w) + "\n")
					if err := img.Encode(&buf, typ); err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					buf.WriteString("\n" + strings.Repeat("~", w))
					testWrite(t, buf.Bytes())
					name := filepath.Join("testdata", fmt.Sprintf("blocks_%4d_%s.txt", seed, typ))
					if err := os.WriteFile(name, append(buf.Bytes(), '\n'), 0o644); err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					expBuf, err := os.ReadFile(name + ".golden")
					if err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					if !bytes.Equal(bytes.TrimSpace(buf.Bytes()), bytes.TrimSpace(expBuf)) {
						t.Errorf("expected %s and %s to be the same", name+".golden", name)
					}
				})
			}
			t.Run("Auto", func(t *testing.T) {
				t.Logf("%s:", img.Best())
				testWrite(t, img.Bytes())
			})
		})
	}
}

func TestNewReader(t *testing.T) {
	t.Parallel()
	for y := 1; y < 135; y++ {
		for x := 1; x < 135; x++ {
			t.Run(fmt.Sprintf("%03d_%03d", x, y), func(t *testing.T) {
				t.Parallel()
				img, err := NewReader(newOneReader(x, y), x)
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				t.Logf("\n%s", img)
			})
		}
	}
}

type oneReader struct{}

func newOneReader(x, y int) io.Reader {
	n := x * y
	return io.LimitReader(&oneReader{}, int64((n+7)/8))
}

func (r *oneReader) Read(buf []byte) (int, error) {
	n := len(buf)
	for i := range n {
		buf[i] = 0xff
	}
	return n, nil
}

func testWrite(t *testing.T, buf []byte) {
	t.Helper()
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		t.Log(string(bytes.TrimRightFunc(line, unicode.IsSpace)))
	}
}
