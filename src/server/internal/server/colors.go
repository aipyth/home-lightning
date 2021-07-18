package server

import (
)

// TODO: use bytes instead of strings

var (
	black   = "30"
	red     = "31"
	green   = "32"
	yellow  = "33"
	blue    = "34"
	magenta = "35"
	cyan    = "36"
	white   = "37"

	bgBlack   = "40"
	bgRed     = "41"
	bgGreen   = "42"
	bgYellow  = "43"
	bgBlue    = "44"
	bgMagenta = "45"
	bgCyan    = "46"
	bgWhite   = "47"

	reset         = "0"
	boldBright    = "1"
	underline     = "4"
	inverse       = "7"
	boldBrightOff = "21"
	underlineOff  = "24"
	inverseOff    = "27"
)

func buildColor(colors ...string) string {
	out := "\033["
	for i := 0; i < len(colors); i++ {
		out += colors[i]
		if i != len(colors)-1 { out += ";" }
	}
	out += "m"
	return out
}

func getColorByStatus(code int) string {
	switch {
	case code < 200:
		return buildColor(boldBright, bgCyan)
	case code >= 200 && code < 300:
		return buildColor(boldBright, bgGreen)
	case code >= 300 && code < 400:
		return buildColor(boldBright, bgBlue)
	case code >= 400 && code < 500:
		return buildColor(boldBright, bgYellow)
	default:
		return buildColor(boldBright, bgRed)
	}
}

func wrapByStatus(code int, content string) string {
	return getColorByStatus(code) + content + buildColor(reset)
}
