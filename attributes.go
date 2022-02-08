package hexes

import (
	"strconv"
)

const (
	NORMAL     = "\033[0m"
	BOLD       = "\033[1m"
	FAINT      = "\033[2m"
	ITALIC     = "\033[3m"
	UNDERLINE  = "\033[4m"
	SLOW_BLINK = "\033[5m"
	FAST_BLINK = "\033[6m"
	REVERSE    = "\033[7m"
	STRIKE     = "\033[8m"

	BLACK   = "\033[30m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	WHITE   = "\033[37m"

	BG_BLACK   = "\033[40m"
	BG_RED     = "\033[41m"
	BG_GREEN   = "\033[42m"
	BG_YELLOW  = "\033[43m"
	BG_BLUE    = "\033[44m"
	BG_MAGENTA = "\033[45m"
	BG_CYAN    = "\033[46m"
	BG_WHITE   = "\033[47m"
)

func TrueColor(red, green, blue int) string {
	return "\033[38;2;" + strconv.Itoa(red) + ";" + strconv.Itoa(green) + ";" + strconv.Itoa(blue) + "m"
	// return fmt.Sprintf("\033[38;2;%v;%v;%vm", red, green, blue) // This version is way slower
}

func TrueColorBg(red, green, blue int) string {
	return "\033[48;2;" + strconv.Itoa(red) + ";" + strconv.Itoa(green) + ";" + strconv.Itoa(blue) + "m"
	// return fmt.Sprintf("\033[48;2;%v;%v;%vm", red, green, blue)
}
