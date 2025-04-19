# blocked

Package `blocked` provides unicode block encoding for binary data.

[![Unit Tests][blocked-ci-status]][blocked-ci]
[![Go Reference][goref-blocked-status]][goref-blocked]
[![Releases][release-status]][Releases]
[![Discord Discussion][discord-status]][discord]

![blocked sample output](screenshot.png)

[blocked-ci]: https://github.com/kenshaw/blocked/actions/workflows/test.yml "Test CI"
[blocked-ci-status]: https://github.com/kenshaw/blocked/actions/workflows/test.yml/badge.svg "Test CI"
[goref-blocked]: https://pkg.go.dev/github.com/kenshaw/blocked "Go Reference"
[goref-blocked-status]: https://pkg.go.dev/badge/github.com/kenshaw/blocked.svg "Go Reference"
[release-status]: https://img.shields.io/github/v/release/kenshaw/blocked?display_name=tag&sort=semver "Latest Release"
[discord]: https://discord.gg/WDWAgXwJqN "Discord Discussion"
[discord-status]: https://img.shields.io/discord/829150509658013727.svg?label=Discord&logo=Discord&colorB=7289da&style=flat-square "Discord Discussion"
[releases]: https://github.com/kenshaw/blocked/releases "Releases"

## Overview

Enables displaying binary bit data using Unicode blocks such as [sextants][],
[octants][], [braille][], or other [block `Type`'s][b-type].

The following [block types][b-type] are currently supported:

|                               | Description                                |            Def            |
| :---------------------------- | :----------------------------------------- | :-----------------------: |
| **Singles (1x1 blocks)**      |                                            |                           |
| [`Solids`][b-type]            | [Full block][solids] (` `, `█`)            |       [≝][b-solids]       |
| [`Binaries`][b-type]          | Binary digits (`0`, `1`)                   |      [≝][b-binaries]      |
| [`XXs`][b-type]               | Binary mask characters (` `, `X`)          |        [≝][b-xxs]         |
|                               |                                            |                           |
| **Doubles (1x2 blocks)**      |                                            |                           |
| [`Halves`][b-type]            | [Half blocks][halves]                      |       [≝][b-halves]       |
| [`ASCIIs`][b-type]            | ASCII-safe characters (` `, `^`, `v`, `%`) |       [≝][b-asciis]       |
|                               |                                            |                           |
| **Quads (2x2 blocks)**        |                                            |                           |
| [`Quads`][b-type]             | [Quarter blocks][quads]                    |       [≝][b-quads]        |
| [`QuadsSeparated`][b-type]    | [Quarter blocks, separated][quads-sep]     |  [≝][b-quads-separated]   |
|                               |                                            |                           |
| **Sextants (2x3 blocks)**     |                                            |                           |
| [`Sextants`][b-type]          | [Sextant blocks][sextants]                 |      [≝][b-sextants]      |
| [`SextantsSeparated`][b-type] | [Sextant blocks, separated][sextants-sep]  | [≝][b-sextants-separated] |
|                               |                                            |                           |
| **Octants (2x4 blocks)**      |                                            |                           |
| [`Octants`][b-type]           | [Octant blocks (Unicode-16)][octants]      |      [≝][b-octants]       |
| [`Braille`][b-type]           | [Braille glyphs][braille]                  |      [≝][b-braille]       |

[solids]: https://www.amp-what.com/unicode/search/full%20block
[halves]: https://www.amp-what.com/unicode/search/half%20block
[quads]: https://www.amp-what.com/unicode/search/quarter%20block
[quads-sep]: https://www.amp-what.com/unicode/search/quad%20separated
[sextants]: https://www.amp-what.com/unicode/search/sextants
[sextants-sep]: https://www.amp-what.com/unicode/search/sextants%20separated
[octants]: https://www.amp-what.com/unicode/search/octants
[braille]: https://www.amp-what.com/unicode/search/braille
[b-type]: https://pkg.go.dev/github.com/kenshaw/blocked#Type
[b-solids]: https://pkg.go.dev/github.com/kenshaw/blocked#SolidsRunes
[b-binaries]: https://pkg.go.dev/github.com/kenshaw/blocked#BinariesRunes
[b-xxs]: https://pkg.go.dev/github.com/kenshaw/blocked#XXsRunes
[b-halves]: https://pkg.go.dev/github.com/kenshaw/blocked#HalvesRunes
[b-asciis]: https://pkg.go.dev/github.com/kenshaw/blocked#ASCIIsRunes
[b-quads]: https://pkg.go.dev/github.com/kenshaw/blocked#QuadsRunes
[b-quads-separated]: https://pkg.go.dev/github.com/kenshaw/blocked#QuadsSeparatedRunes
[b-sextants]: https://pkg.go.dev/github.com/kenshaw/blocked#SextantsRunes
[b-sextants-separated]: https://pkg.go.dev/github.com/kenshaw/blocked#SextantsSeparatedRunes
[b-octants]: https://pkg.go.dev/github.com/kenshaw/blocked#OctantsRunes
[b-braille]: https://pkg.go.dev/github.com/kenshaw/blocked#BrailleRunes

## Using

Install in the usual Go fashion:

```sh
$ go get -u github.com/kenshaw/blocked@latest
```

## Example

```go
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"

	"github.com/kenshaw/blocked"
)

func main() {
	for i, seed := range []int{
		999,
		4000,
		12555,
	} {
		if i != 0 {
			fmt.Println()
		}

		// create random data
		r := rand.New(rand.NewSource(int64(seed)))
		height := 1 + r.Intn(12)
		data := make([]uint64, height)
		for i := range height {
			data[i] = r.Uint64()
		}

		// create bitmap from data, using 64 bits per line (the width of
		// uint64)
		img, err := blocked.New(data, 64)
		if err != nil {
			panic(err)
		}

		// note: to interpret data as 46 bits wide (or another bit width) do
		// the following:
		//
		// img, err := blocked.New(data, 46)

		// encode as blocks and display
		for j, typ := range []blocked.Type{
			blocked.Solids,
			blocked.Halves,
			blocked.Sextants,
			blocked.Octants,
			blocked.Braille,
		} {
			if j != 0 {
				fmt.Println()
			}
			fmt.Printf("%d%c. %s:\n", i+1, 'a'+j, typ)
			var buf bytes.Buffer
			if err := img.Encode(&buf, typ); err != nil {
				panic(err)
			}

			// usually this would suffice:
			// fmt.Println(buf.String())

			// to work around with comment formatting issues in Go examples,
			// add pipes to surround the output:
			s := "|" + strings.ReplaceAll(buf.String(), "\n", "|\n|") + "|"
			fmt.Println(s)
		}
	}
}
```

Output:

```txt
1a. Solids:
|██  ██████  ██ █ █ ██  ██ █    █ ██ █  █  ██ █ █  █    ███ █████|

1b. Halves:
|▀▀  ▀▀▀▀▀▀  ▀▀ ▀ ▀ ▀▀  ▀▀ ▀    ▀ ▀▀ ▀  ▀  ▀▀ ▀ ▀  ▀    ▀▀▀ ▀▀▀▀▀|

1c. Sextants:
|🬂 🬂🬂🬂 🬂🬁🬁🬁🬀🬁🬀🬀 🬁🬁🬀🬀🬁 🬂🬁🬁 🬀 🬁🬂🬁🬂🬂|

1d. Octants:
|🮂 🮂🮂🮂 🮂𜺫𜺫𜺫𜺨𜺫𜺨𜺨 𜺫𜺫𜺨𜺨𜺫 🮂𜺫𜺫 𜺨 𜺫🮂𜺫🮂🮂|

1e. Braille:
|⠉⠀⠉⠉⠉⠀⠉⠈⠈⠈⠁⠈⠁⠁⠀⠈⠈⠁⠁⠈⠀⠉⠈⠈⠀⠁⠀⠈⠉⠈⠉⠉|

2a. Solids:
|█████  ████ █  █ █ █    █   █    █████ █  █ █ █   ███ ████ █ ███|
| ██ ██ ██   ██ █     █   █   ███    █   █ █  ███ ██ ██    ██   █|
|█  █ ███ ████ ██  █  █ ██ █ █    ██ █  ██   ██ ████ ██  █ █   █ |
|███    █  ██ ████   ███ █     █ █ ███ ███ █      ███  █     ███ |
| ████  ███  █ █  ██ █   █ █  ████ █ ███ █ █████      █ ████ ██ █|
|  █     █ █    ███ █  ██     █ █ ██ █████        █  ████  █ ██  |

2b. Halves:
|▀██▀█▄ ██▀▀ █▄ █ ▀ ▀ ▄  ▀▄  ▀▄▄▄ ▀▀▀█▀ ▀▄ █ ▀▄█▄ ▄█▀█▄▀▀▀▀▄█ ▀▀█|
|█▄▄▀ ▀▀█ ▀██▀▄██▄ ▀ ▄█▄▀█ ▀ ▀ ▄ ▄▀█▄█ ▄██ ▄ ▀▀ ▀▀██▄▀▀▄ ▀ ▀ ▄▄█ |
| ▀█▀▀  ▀█▀▄ ▀ ▀▄▄█▀▄▀ ▄▄▀ ▀  █▀█▀▄█ ███▄█ ▀▀▀▀▀  ▄  ▄█▄█▀▀█ ██ ▀|

2c. Sextants:
|🬙🬥🬪🬷🬥🬮🬛🬷🬁🬑🬦🬞🬗🬏🬗🬋🬠🬒🬕🬠🬓🬄🬶🬪🬵🬕🬺🬂🬒🬜🬁🬙|
|🬊🬛🬃🬉🬚🬒🬅🬥🬶🬢🬆🬮🬄🬃🬦🬪🬣🬕🬺🬴▌🬌🬋🬃🬠🬂🬵🬶🬋🬓█🬈|

2d. Octants:
|𜷆𜵘𜴤𜶦𜴟𜷛𜶍▟𜴷𜴋𜷓𜵑𜵌𜴉𜴑𜵁𜵓𜶾𜵊𜷍𜵈𜴺𜴰𜴤𜶤𜷂𜴴𜴸𜴌𜴖𜶭𜵍|
|𜺫𜴂𜺨𜺫𜴂𜴀𜺨𜴄𜴈𜴄𜺨𜴆𜺨𜺨▝𜴅𜴄▘▀𜴇▘🮂🮂𜺨𜴃 𜴈𜴈🮂▘▀𜺫|

2e. Braille:
|⣝⡫⠳⢼⠫⣥⢗⣼⡈⠌⣰⡠⡕⠄⠕⡒⡨⣍⡏⣨⡆⡃⠵⠳⢴⣏⠷⡉⠍⠞⣈⡝|
|⠈⠋⠁⠈⠋⠂⠁⠑⠚⠑⠁⠒⠁⠁⠘⠙⠑⠃⠛⠓⠃⠉⠉⠁⠐⠀⠚⠚⠉⠃⠛⠈|

3a. Solids:
|█   █  ██  ██    █ ██ █   █ ██ █ █     ██  █ █ ██ █ █  ████   ██|
|███ ██ █  █   ██ █  ███     █ █ ███    ████ ██ ███ ███  ██      |
| █ █ █ █  █  ██ █ ███ █  █   █ █ █  ██ █  ████████    ██ ███   █|
|█ ██   █ ███    ██ █      ████ ██  █ ███ █    ███  ██  ████ ███ |
| ██   ██ ██  ██ █  ██  █       ██ ██ █ █ ██  █   ████ █ █  █  █ |
|█ █ ██ █ ███  ██ █ ███ █ █ █ █ ████   ██     █ ██ █   ███  █ ███|
| █   ████  █  ███ █ ██    ██ ██ ████ █   ██  █   █ ███ ██      █|
| █  █████████ █    █  █  ██ █ █      █     ███ █  █ █  █ █    █ |
|█ █ █ █ ████ ██  █ ███ █  ██      █     ████   ██  ████ ██  ██  |
| █ ████  █ ████ █ ██ ███  █  █ █  ████  █       █ ██    ███  ███|
| █  █    █ █  █ █  █    █         ███  █  █    █    █  █   █ ██ |
|█    █ █ ██    █ ███ ███ █   ██      █  ████ █ █  ██ ███ ███ █ █|

3b. Halves:
|█▄▄ █▄ █▀ ▄▀▀ ▄▄ █ ▀█▄█   ▀ █▀▄▀▄█▄    ██▄▄▀▄█ ██▄▀▄█▄ ▀██▀   ▀▀|
|▄▀▄█ ▀ █ ▄█▄ ▀▀ █▄▀█▀ ▀  ▀▄▄▄█ █▄▀ ▄▀█▄█ ▄▀▀▀▀███▀ ▄▄ ▀█▄██▀▄▄▄▀|
|▄▀█ ▄▄▀█ ██▄ ▀█▄▀▄ ██▄ █ ▄ ▄ ▄ ██▄█▀ ▀▄█ ▀▀  █ ▄▄▀█▀▀ █▄█  █ ▄█▄|
| █  ▄████▄▄█▄ █▀▀ ▀▄▀▀▄  ▄█▀▄▀█ ▀▀▀▀ █   ▀▀▄▄█ ▄ ▀▄▀█▀ █▀▄    ▄▀|
|▀▄▀▄█▄█ ▀█▀█▄██ ▄▀▄█▀█▄█  █▀ ▄ ▄  █▄▄▄  █▀▀▀   ▀█ ▄█▀▀▀ ██▄ ▀█▄▄|
|▄▀  ▀▄ ▄ █▄▀  ▀▄▀▄▄█ ▄▄▄▀▄   ▄▄   ▀▀▀▄ ▀▄▄█▄ ▄ █  ▄▄▀▄▄█ ▄▄█ █▀▄|

3c. Sextants:
|🬪🬢🬪▐🬀🬔🬟🬚🬘🬯🬛▌🬞🬀🬥🬤🬫🬃🬭▐🬌🬳🬻🬷🬺🬈🬌🬯🬬🬮 🬡|
|🬗🬕🬭🬫▐🬴🬇🬱🬥▐🬱🬦🬞🬡🬡▐🬲🬜🬉🬸🬉🬃🬦🬡🬗🬜🬄🬳🬕🬧🬡🬲|
|🬘🬏🬜🬝🬺🬻🬢🬕🬟🬧🬰🬢🬇🬴🬅🬄🬂🬒🬉 🬯🬶🬍🬦🬑🬤🬴🬘🬶 🬭🬅|
|🬘🬁🬥🬟▐🬘🬂🬣🬣🬸🬠🬰🬢🬀🬠🬑 🬎🬥🬇🬮🬱🬞🬦🬀🬰🬢🬵🬡🬶▐🬥|

3d. Octants:
|𜵞𜷏𜴤▐𜵱𜷁𜴙𜴔𜷅𜶞𜴕𜴍𜴘𜶬𜷒𜶔𜵟𜵴𜶜𜷕𜵽𜴭𜴵▟𜵮𜵹𜵂𜶞𜷚𜵢▂𜵔|
|𜶔▘𜷗𜷣𜷅𜷘𜴷𜵮𜴑𜶊𜴴𜴿𜵸𜵩𜵙▞𜴴𜴮𜶑𜴈𜴚𜶁𜷕𜵸𜴞𜵘𜵢𜶚𜶅▝𜴃𜵞|
|𜵚𜴄𜶍𜵵𜶘𜵜𜴈𜶅𜶆𜷙𜵻𜶹𜶀𜴂𜵸𜴽 𜴴𜶌𜴘𜶲𜶾𜺠𜶑▘𜶹𜶃𜷌𜵿𜷏𜶘𜶌|

3e. Braille:
|⡳⣢⠳⢸⢁⣎⠡⠖⣜⢬⠗⠇⠠⣁⣫⢪⡺⢂⢤⣸⢓⠮⠾⣼⡷⢑⡓⢬⣻⡥⣀⡩|
|⢪⠃⣲⣽⣜⣳⡈⡷⠕⢜⠷⡘⢐⡴⡰⡜⠷⠯⢨⠚⠨⢅⣸⢐⠪⡫⡥⢳⢇⠘⠐⡳|
|⡱⠑⢗⢃⢹⡹⠚⢇⢎⣺⢙⣚⢄⠋⢐⡐⠀⠷⢖⠠⣋⣍⢀⢨⠃⣚⢍⣡⢛⣢⢹⢖|
```
