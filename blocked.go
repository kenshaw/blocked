// Package blocked provides a block encoder and decoder for bitmaps.
package blocked

import (
	"io"
	"maps"
	"sync"
)

// Type is a block type.
type Type rune

// Block types.
const (
	Auto              Type = 'v'
	Solids            Type = 'l'
	Binaries          Type = 'b'
	XXs               Type = 'L'
	Doubles           Type = 'D'
	Halves            Type = 'e'
	ASCIIs            Type = 'E'
	Quads             Type = 'q'
	QuadsSeparated    Type = 'Q'
	Sextants          Type = 'x'
	SextantsSeparated Type = 'X'
	Octants           Type = 'o'
	Braille           Type = 'O'
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

var (
	// blocks are the block maps.
	blocks = make(map[Type]map[uint8]rune)
	// blocksMu is the blocks mutex.
	blocksMu sync.Mutex
)
