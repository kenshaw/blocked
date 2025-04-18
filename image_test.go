package blocked

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func TestBitmapGetSet(t *testing.T) {
	const mask = 0b10101010
	for _, h := range []int{2, 4, 6, 8, 10, 12, 16, 24, 28, 34, 56} {
		for _, w := range []int{2, 4, 6, 8, 10, 12, 16, 24, 28, 34, 56} {
			t.Run(strconv.Itoa(h)+"/"+strconv.Itoa(w), func(t *testing.T) {
				n := w
				if n%8 != 0 {
					n = n/8 + 1
				} else {
					n /= 8
				}
				pix := make([]byte, n*h)
				for i := range len(pix) {
					pix[i] = mask
				}
				// clear upper bits
				if w%8 != 0 {
					for i := n - 1; i < len(pix); i += n {
						pix[i] &= 0xff >> (8 - (w % 8))
					}
				}
				img1 := Bitmap{
					Pix:    pix,
					Stride: n,
					Rect:   image.Rect(0, 0, w, h),
				}
				// t.Logf("pix: %b", img1.Pix)
				b1 := img1.Bounds()
				for y := range b1.Dy() {
					for x := range b1.Dx() {
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
				img2 := NewBitmap(image.Rect(0, 0, w, h))
				b2 := img2.Bounds()
				for y := range b2.Dy() {
					for x := range b2.Dx() {
						img2.Set(x, y, x%2 == 1)
					}
				}
				// compare
				if !slices.Equal(img1.Pix, img2.Pix) {
					t.Errorf("expected:\n%b\ngot:\n%b", img1.Pix, img2.Pix)
				}
				if !strings.Contains(os.Getenv("TESTS"), "image") {
					return
				}
				var buf bytes.Buffer
				if err := png.Encode(&buf, img2); err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if err := os.WriteFile(fmt.Sprintf("test_%02d_%02d.png", h, w), buf.Bytes(), 0o644); err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
			})
		}
	}
}

func TestNewBitmap(t *testing.T) {
	for n := 1330; n <= 1343; n++ {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			r := rand.New(rand.NewSource(int64(n)))
			img := NewBitmap(image.Rect(0, 0, 1+r.Intn(28), 1+r.Intn(28)))
			rect := img.Bounds()
			w, h, expTr, expFa := rect.Dx(), rect.Dy(), 0, 0
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
				t.Fatal("invalid test")
			}
			if n, exp := expTr+expFa, rect.Dx()*rect.Dy(); n != exp {
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
			var buf bytes.Buffer
			if err := png.Encode(&buf, img); err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			if err := os.WriteFile(fmt.Sprintf("test_%4d.png", n), buf.Bytes(), 0o644); err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			for _, typ := range Types() {
				t.Run(typ.String(), func(t *testing.T) {
					var buf bytes.Buffer
					t.Logf("%s:", typ)
					n := img.Width(typ)
					t.Logf("%s", strings.Repeat("_", n))
					if err := img.Encode(&buf, typ); err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					testWrite(t, buf.Bytes())
					t.Logf("%s", strings.Repeat("~", n))
				})
			}
		})
	}
}
