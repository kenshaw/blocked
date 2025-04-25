package blocked

// SolidsRunes returns the runes for single block resolution bitmaps.
//
// See: https://www.amp-what.com/unicode/search/full%20block
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
//
// See: https://www.amp-what.com/unicode/search/half%20block
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
//
// See: https://www.amp-what.com/unicode/search/quarter%20block
func QuadsRunes() []rune {
	return []rune{
		' ', '▘', '▝', '▀', '▖', '▌', '▞', '▛',
		'▗', '▚', '▐', '▜', '▄', '▙', '▟', '█',
	}
}

// QuadsSeparatedRunes returns the runes for quadruple block resolution images
// using separated quads.
//
// See: https://www.amp-what.com/unicode/search/quad%20separated
func QuadsSeparatedRunes() []rune {
	return []rune{
		' ', '𜰡', '𜰢', '𜰣', '𜰤', '𜰥', '𜰦', '𜰧',
		'𜰨', '𜰩', '𜰪', '𜰫', '𜰬', '𜰭', '𜰮', '𜰯',
	}
}

// SextantsRunes returns the runes for sextuple block resolution images.
//
// See: https://www.amp-what.com/unicode/search/sextants
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
//
// See: https://www.amp-what.com/unicode/search/sextants%20separated
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
//
// See: https://www.amp-what.com/unicode/search/octants
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
//
// See: https://www.amp-what.com/unicode/search/braille
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
