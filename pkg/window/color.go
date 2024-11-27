package window

import "strconv"

type ColorCode int

const FgNone ColorCode = 0

// Foreground text colors
const (
	FgBlack ColorCode = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	FgHiBlack ColorCode = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack ColorCode = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite

	BgHiBlack ColorCode = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)

const noColor = "\033[0m"

type color struct {
	code ColorCode
}

func (c color) paint() []byte {
	// return none
	if c.code == FgNone {
		return []byte(noColor)
	}

	code := strconv.Itoa(int(c.code))
	return []byte("\033[0;" + code + "m")
}
