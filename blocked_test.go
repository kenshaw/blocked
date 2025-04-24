package blocked

import (
	"bytes"
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

func testWrite(t *testing.T, buf []byte) {
	t.Helper()
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		t.Log(string(bytes.TrimRightFunc(line, unicode.IsSpace)))
	}
}
