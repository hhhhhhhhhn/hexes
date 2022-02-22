package hexes

import (
	"strconv"
)

type Attribute []byte

func Join(attributes... Attribute) Attribute {
	size := len(attributes)
	for _, attribute := range attributes {
		size += len(attribute)
	}

	out := make(Attribute, size)
	index := 0

	for _, attribute := range attributes {
		index += copy(out[index:], attribute)
	}

	return out
}

var (
	NORMAL     = Attribute("\033[0m")
	BOLD       = Attribute("\033[1m")
	FAINT      = Attribute("\033[2m")
	ITALIC     = Attribute("\033[3m")
	UNDERLINE  = Attribute("\033[4m")
	SLOW_BLINK = Attribute("\033[5m")
	FAST_BLINK = Attribute("\033[6m")
	REVERSE    = Attribute("\033[7m")
	STRIKE     = Attribute("\033[8m")

	BLACK   = Attribute("\033[30m")
	RED     = Attribute("\033[31m")
	GREEN   = Attribute("\033[32m")
	YELLOW  = Attribute("\033[33m")
	BLUE    = Attribute("\033[34m")
	MAGENTA = Attribute("\033[35m")
	CYAN    = Attribute("\033[36m")
	WHITE   = Attribute("\033[37m")

	BG_BLACK   = Attribute("\033[40m")
	BG_RED     = Attribute("\033[41m")
	BG_GREEN   = Attribute("\033[42m")
	BG_YELLOW  = Attribute("\033[43m")
	BG_BLUE    = Attribute("\033[44m")
	BG_MAGENTA = Attribute("\033[45m")
	BG_CYAN    = Attribute("\033[46m")
	BG_WHITE   = Attribute("\033[47m")
)

func TrueColor(red, green, blue int) Attribute {
	return Attribute("\033[38;2;" + strconv.Itoa(red) + ";" + strconv.Itoa(green) + ";" + strconv.Itoa(blue) + "m")
	// return fmt.Sprintf("\033[38;2;%v;%v;%vm", red, green, blue) // This version is way slower
}

func TrueColorBg(red, green, blue int) Attribute {
	return Attribute("\033[48;2;" + strconv.Itoa(red) + ";" + strconv.Itoa(green) + ";" + strconv.Itoa(blue) + "m")
	// return fmt.Sprintf("\033[48;2;%v;%v;%vm", red, green, blue)
}
