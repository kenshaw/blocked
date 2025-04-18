// Package blocked provides a block encoder and decoder for bitmaps.
package blocked

import (
	"io"
	"maps"
)

// Type is a block type.
type Type int

// Block types.
const (
	Auto Type = iota
	Solids
	Binaries
	XXs
	Halves
	ASCIIs
	Quads
	QuadsSeparated
	Sextants
	SextantsSeparated
	Octants
	Braille
)

// Types returns all block types.
func Types() []Type {
	return []Type{
		Solids,
		Binaries,
		XXs,
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
	case Solids, Halves, Quads, Sextants, Octants:
		return true
	}
	return false
}

// RuneCount returns the number of runes for the block type.
func (typ Type) RuneCount() int {
	switch typ {
	case Solids, Binaries, XXs:
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
	case Solids, Binaries, XXs:
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

// Dump dumps a ASCII drawing of the bitmask the block type's symbols to the writer.
func (typ Type) Dump(w io.Writer) {
	if m := typ.runeMap(); m != nil {
		Dump(w, m)
	}
}

/*
// Encode encodes the bitmap to the writer.
func (typ Type) Encode(w io.Writer, img Bitmap, opts ...EncoderOption) error {
	return NewEncoder(opts...).Encode(w, img)
}

// Decode decodes a bitmap from the reader.
func (typ Type) Decode(r io.Reader, img *Bitmap, opts ...DecoderOption) error {
	return NewDecoder(opts...).Decode(r, img)
}
*/

// runeMap returns the rune map for the type.
func (typ Type) runeMap() map[uint8]rune {
	switch typ {
	case Solids:
		return solids
	case Binaries:
		return binaries
	case XXs:
		return xxs
	case Halves:
		return halves
	case ASCIIs:
		return asciis
	case Quads:
		return quads
	case QuadsSeparated:
		return quadsSeparated
	case Sextants:
		return sextants
	case SextantsSeparated:
		return sextantsSeparated
	case Octants:
		return octants
	case Braille:
		return braille
	}
	return nil
}

var (
	// solids is the map for single block resolution bitmaps.
	solids = toMap(SolidsRunes())
	// binaries is the map for single block resolution bitmaps using '0', '1'.
	binaries = toMap(BinariesRunes())
	// xxs is the map for single block resolution bitmaps using ' ', 'X'.
	xxs = toMap(XXsRunes())
	// halves is the map for double block resolution bitmaps.
	halves = toMap(HalvesRunes())
	// asciis is the map for double block resolution bitmaps.
	asciis = toMap(ASCIIsRunes())
	// quads is the map for quadruple block resolution bitmaps.
	quads = toMap(QuadsRunes())
	// quadsSeparated is the map for quadruple block resolution bitmaps using
	// separated quads.
	quadsSeparated = toMap(QuadsSeparatedRunes())
	// sextants is the map for sextuple block resolution bitmaps.
	sextants = toMap(SextantsRunes())
	// sextantsSeparated is the map for sextuple block resolution bitmaps using
	// separated sextants.
	sextantsSeparated = toMap(SextantsSeparatedRunes())
	// octants is the map for octuple block resolution bitmaps.
	octants = toMap(OctantsRunes())
	// braille is the map for octuple block resolution bitmaps using braille.
	braille = toMap(BrailleRunes())
)

// toMap converts a slice to a map.
func toMap(v []rune) map[uint8]rune {
	m := make(map[uint8]rune, len(v))
	for i, r := range v {
		m[uint8(i)] = r
	}
	return m
}

// Best returns the best (most compact) block type.
func Best(y int) Type {
	switch {
	case y == 1:
		return Solids
	case y == 2:
		return Halves
	case y <= 4:
		return Quads
	case y <= 6:
		return Sextants
	}
	return Octants
}
