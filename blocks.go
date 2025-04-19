package blocked

import (
	"fmt"
	"io"
	"math/bits"
	"strings"
)

// SolidsRunes returns the runes for single block resolution bitmaps.
func SolidsRunes() []rune {
	return []rune{
		' ', '█',
	}
}

// BinariesRunes returns the runes for single block resolution bitmaps using
// '0', '1'.
func BinariesRunes() []rune {
	return []rune{
		'0', '1',
	}
}

// XXsRunes returns the runes for single block resolution bitmaps using ' ',
// 'X'.
func XXsRunes() []rune {
	return []rune{
		' ', 'X',
	}
}

// HalvesRunes returns the runes for double block resolution bitmaps.
func HalvesRunes() []rune {
	return []rune{
		' ', '▀', '▄', '█',
	}
}

// ASCIIsRunes returns the runes for double block resolution bitmaps.
func ASCIIsRunes() []rune {
	return []rune{
		' ', '^', 'v', '%',
	}
}

// QuadsRunes returns the runes for quadruple block resolution bitmaps.
func QuadsRunes() []rune {
	return []rune{
		' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛',
		'▗', '▚', '▐', '▜', '▄', '▙', '▟', '█',
	}
}

// QuadsSeparatedRunes returns the runes for quadruple block resolution images
// using separated quads.
func QuadsSeparatedRunes() []rune {
	return []rune{
		' ', '𜰡', '𜰢', '𜰣', '𜰤', '𜰥', '𜰦', '𜰧',
		'𜰨', '𜰩', '𜰪', '𜰫', '𜰬', '𜰭', '𜰮', '𜰯',
	}
}

// SextantsRunes returns the runes for sextuple block resolution images.
func SextantsRunes() []rune {
	return []rune{
		' ', '🬀', '🬁', '🬂', '🬃', '🬄', '🬅', '🬆',
		'🬇', '🬈', '🬉', '🬊', '🬋', '🬌', '🬍', '🬎',
		'🬏', '🬐', '🬑', '🬒', '🬓', '▌', '🬔', '🬕',
		'🬖', '🬗', '🬘', '🬙', '🬚', '🬛', '🬜', '🬝',
		'🬞', '🬟', '🬠', '🬡', '🬢', '🬣', '🬤', '🬥',
		'🬦', '🬧', '▐', '🬨', '🬩', '🬪', '🬫', '🬬',
		'🬭', '🬮', '🬯', '🬰', '🬱', '🬲', '🬳', '🬴',
		'🬵', '🬶', '🬷', '🬸', '🬹', '🬺', '🬻', '█',
	}
}

// SextantsSeparatedRunes returns the runes for sextuple block resolution
// images using separated sextants.
func SextantsSeparatedRunes() []rune {
	return []rune{
		' ', '𜹑', '𜹒', '𜹓', '𜹔', '𜹕', '𜹖', '𜹗',
		'𜹘', '𜹙', '𜹚', '𜹛', '𜹜', '𜹝', '𜹞', '𜹟',
		'𜹠', '𜹡', '𜹢', '𜹣', '𜹤', '𜹥', '𜹦', '𜹧',
		'𜹨', '𜹩', '𜹪', '𜹫', '𜹬', '𜹭', '𜹮', '𜹯',
		'𜹰', '𜹱', '𜹲', '𜹳', '𜹴', '𜹵', '𜹶', '𜹷',
		'𜹸', '𜹹', '𜹺', '𜹻', '𜹼', '𜹽', '𜹾', '𜹿',
		'𜺀', '𜺁', '𜺂', '𜺃', '𜺄', '𜺅', '𜺆', '𜺇',
		'𜺈', '𜺉', '𜺊', '𜺋', '𜺌', '𜺍', '𜺎', '𜺏',
	}
}

// OctantsRunes returns the runes for octuple block resolution images.
func OctantsRunes() []rune {
	return []rune{
		' ', '𜺨', '𜺫', '🮂', '𜴀', '▘', '𜴁', '𜴂',
		'𜴃', '𜴄', '▝', '𜴅', '𜴆', '𜴇', '𜴈', '▀',
		'𜴉', '𜴊', '𜴋', '𜴌', '🯦', '𜴍', '𜴎', '𜴏',
		'𜴐', '𜴑', '𜴒', '𜴓', '𜴔', '𜴕', '𜴖', '𜴗',
		'𜴘', '𜴙', '𜴚', '𜴛', '𜴜', '𜴝', '𜴞', '𜴟',
		'🯧', '𜴠', '𜴡', '𜴢', '𜴣', '𜴤', '𜴥', '𜴦',
		'𜴧', '𜴨', '𜴩', '𜴪', '𜴫', '𜴬', '𜴭', '𜴮',
		'𜴯', '𜴰', '𜴱', '𜴲', '𜴳', '𜴴', '𜴵', '🮅',
		'𜺣', '𜴶', '𜴷', '𜴸', '𜴹', '𜴺', '𜴻', '𜴼',
		'𜴽', '𜴾', '𜴿', '𜵀', '𜵁', '𜵂', '𜵃', '𜵄',
		'▖', '𜵅', '𜵆', '𜵇', '𜵈', '▌', '𜵉', '𜵊',
		'𜵋', '𜵌', '▞', '𜵍', '𜵎', '𜵏', '𜵐', '▛',
		'𜵑', '𜵒', '𜵓', '𜵔', '𜵕', '𜵖', '𜵗', '𜵘',
		'𜵙', '𜵚', '𜵛', '𜵜', '𜵝', '𜵞', '𜵟', '𜵠',
		'𜵡', '𜵢', '𜵣', '𜵤', '𜵥', '𜵦', '𜵧', '𜵨',
		'𜵩', '𜵪', '𜵫', '𜵬', '𜵭', '𜵮', '𜵯', '𜵰',
		'𜺠', '𜵱', '𜵲', '𜵳', '𜵴', '𜵵', '𜵶', '𜵷',
		'𜵸', '𜵹', '𜵺', '𜵻', '𜵼', '𜵽', '𜵾', '𜵿',
		'𜶀', '𜶁', '𜶂', '𜶃', '𜶄', '𜶅', '𜶆', '𜶇',
		'𜶈', '𜶉', '𜶊', '𜶋', '𜶌', '𜶍', '𜶎', '𜶏',
		'▗', '𜶐', '𜶑', '𜶒', '𜶓', '▚', '𜶔', '𜶕',
		'𜶖', '𜶗', '▐', '𜶘', '𜶙', '𜶚', '𜶛', '▜',
		'𜶜', '𜶝', '𜶞', '𜶟', '𜶠', '𜶡', '𜶢', '𜶣',
		'𜶤', '𜶥', '𜶦', '𜶧', '𜶨', '𜶩', '𜶪', '𜶫',
		'▂', '𜶬', '𜶭', '𜶮', '𜶯', '𜶰', '𜶱', '𜶲',
		'𜶳', '𜶴', '𜶵', '𜶶', '𜶷', '𜶸', '𜶹', '𜶺',
		'𜶻', '𜶼', '𜶽', '𜶾', '𜶿', '𜷀', '𜷁', '𜷂',
		'𜷃', '𜷄', '𜷅', '𜷆', '𜷇', '𜷈', '𜷉', '𜷊',
		'𜷋', '𜷌', '𜷍', '𜷎', '𜷏', '𜷐', '𜷑', '𜷒',
		'𜷓', '𜷔', '𜷕', '𜷖', '𜷗', '𜷘', '𜷙', '𜷚',
		'▄', '𜷛', '𜷜', '𜷝', '𜷞', '▙', '𜷟', '𜷠',
		'𜷡', '𜷢', '▟', '𜷣', '▆', '𜷤', '𜷥', '█',
	}
}

// BrailleRunes returns the runes for octuple block resolution images using
// braille.
func BrailleRunes() []rune {
	return []rune{
		'⠀', '⠁', '⠈', '⠉', '⠂', '⠃', '⠊', '⠋',
		'⠐', '⠑', '⠘', '⠙', '⠒', '⠓', '⠚', '⠛',
		'⠄', '⠅', '⠌', '⠍', '⠆', '⠇', '⠎', '⠏',
		'⠔', '⠕', '⠜', '⠝', '⠖', '⠗', '⠞', '⠟',
		'⠠', '⠡', '⠨', '⠩', '⠢', '⠣', '⠪', '⠫',
		'⠰', '⠱', '⠸', '⠹', '⠲', '⠳', '⠺', '⠻',
		'⠤', '⠥', '⠬', '⠭', '⠦', '⠧', '⠮', '⠯',
		'⠴', '⠵', '⠼', '⠽', '⠶', '⠷', '⠾', '⠿',
		'⡀', '⡁', '⡈', '⡉', '⡂', '⡃', '⡊', '⡋',
		'⡐', '⡑', '⡘', '⡙', '⡒', '⡓', '⡚', '⡛',
		'⡄', '⡅', '⡌', '⡍', '⡆', '⡇', '⡎', '⡏',
		'⡔', '⡕', '⡜', '⡝', '⡖', '⡗', '⡞', '⡟',
		'⡠', '⡡', '⡨', '⡩', '⡢', '⡣', '⡪', '⡫',
		'⡰', '⡱', '⡸', '⡹', '⡲', '⡳', '⡺', '⡻',
		'⡤', '⡥', '⡬', '⡭', '⡦', '⡧', '⡮', '⡯',
		'⡴', '⡵', '⡼', '⡽', '⡶', '⡷', '⡾', '⡿',
		'⢀', '⢁', '⢈', '⢉', '⢂', '⢃', '⢊', '⢋',
		'⢐', '⢑', '⢘', '⢙', '⢒', '⢓', '⢚', '⢛',
		'⢄', '⢅', '⢌', '⢍', '⢆', '⢇', '⢎', '⢏',
		'⢔', '⢕', '⢜', '⢝', '⢖', '⢗', '⢞', '⢟',
		'⢠', '⢡', '⢨', '⢩', '⢢', '⢣', '⢪', '⢫',
		'⢰', '⢱', '⢸', '⢹', '⢲', '⢳', '⢺', '⢻',
		'⢤', '⢥', '⢬', '⢭', '⢦', '⢧', '⢮', '⢯',
		'⢴', '⢵', '⢼', '⢽', '⢶', '⢷', '⢾', '⢿',
		'⣀', '⣁', '⣈', '⣉', '⣂', '⣃', '⣊', '⣋',
		'⣐', '⣑', '⣘', '⣙', '⣒', '⣓', '⣚', '⣛',
		'⣄', '⣅', '⣌', '⣍', '⣆', '⣇', '⣎', '⣏',
		'⣔', '⣕', '⣜', '⣝', '⣖', '⣗', '⣞', '⣟',
		'⣠', '⣡', '⣨', '⣩', '⣢', '⣣', '⣪', '⣫',
		'⣰', '⣱', '⣸', '⣹', '⣲', '⣳', '⣺', '⣻',
		'⣤', '⣥', '⣬', '⣭', '⣦', '⣧', '⣮', '⣯',
		'⣴', '⣵', '⣼', '⣽', '⣶', '⣷', '⣾', '⣿',
	}
}

// Dump dumps a ASCII drawing of the bitmask of the symobls to the writer.
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
// signficant bits of i. If width is less than or equal to 2, each line will be
// 1 rune, otherwise each line will be 2 runes.
func splitMask(i uint8, width int) []string {
	s := maskRepl.Replace(fmt.Sprintf("%0*b", width, bits.Reverse8(uint8(i))>>(8-width)))
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
