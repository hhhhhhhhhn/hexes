package input

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

type Listener struct {
	inChannel chan rune
	in        *bufio.Reader
}

func New(in io.Reader) *Listener {
	listener := &Listener{inChannel: make(chan rune)}
	listener.in = bufio.NewReader(in)
	go func() {
		defer close(listener.inChannel)
		for {
			chr, _, err := listener.in.ReadRune()
			if err == nil {
				listener.inChannel <- chr
			}
		}
	}()

	return listener
}

func (l *Listener) EnableMouseTracking(out io.Writer) {
	out.Write([]byte("\033[?1003;1006;1015h"))
}

func (l *Listener) DisableMouseTracking(out io.Writer) {
	out.Write([]byte("\033[?1003;1006;1015l"))
}

type Event struct {
	EventType eventType
	Chr rune
	X   int
	Y   int
}

func (l *Listener) GetEvent() *Event {
	chr := <- l.inChannel

	if chr == ESCAPE {
		return l.parseEscape()
	}
	return &Event{EventType: KeyPressed, Chr: chr}
}

func (l *Listener) parseEscape() *Event {
	timeoutChannel := make(chan bool)
	go func() {
		time.Sleep(time.Millisecond * 50)
		timeoutChannel <- true
	}()

	select {
	case <- timeoutChannel:
		return &Event{EventType: KeyPressed, Chr: ESCAPE}
	case chr := <- l.inChannel:
		l.in.UnreadRune()

		if chr == CSI { // Command Sequence Initiator, read ECMA-48 5.4
			return l.parseCommandSequence()
		}
		return &Event{EventType: KeyPressed, Chr: ESCAPE}
	}
}

func (l *Listener) parseCommandSequence() *Event {
	command := l.divideCommandSequence()
	subparameters := strings.Split(command[0], ";")

	// DEBUG
	//fmt.Fprintln(os.Stderr, "Command was divided into ", command)
	switch command[len(command) - 1] {
	case "A":
		return &Event{EventType: KeyPressed, Chr: KEY_UP}
	case "B":
		return &Event{EventType: KeyPressed, Chr: KEY_DOWN}
	case "C":
		return &Event{EventType: KeyPressed, Chr: KEY_RIGHT}
	case "D":
		return &Event{EventType: KeyPressed, Chr: KEY_LEFT}
	case "~":
		switch command[0] {
		case "11":
			return &Event{EventType: KeyPressed, Chr: KEY_F1}
		case "12":
			return &Event{EventType: KeyPressed, Chr: KEY_F2}
		case "13":
			return &Event{EventType: KeyPressed, Chr: KEY_F3}
		case "14":
			return &Event{EventType: KeyPressed, Chr: KEY_F4}
		case "15":
			return &Event{EventType: KeyPressed, Chr: KEY_F5}
		case "17":
			return &Event{EventType: KeyPressed, Chr: KEY_F6}
		case "18":
			return &Event{EventType: KeyPressed, Chr: KEY_F7}
		case "19":
			return &Event{EventType: KeyPressed, Chr: KEY_F8}
		case "20":
			return &Event{EventType: KeyPressed, Chr: KEY_F9}
		case "21":
			return &Event{EventType: KeyPressed, Chr: KEY_F10}
		case "23":
			return &Event{EventType: KeyPressed, Chr: KEY_F11}
		case "24":
			return &Event{EventType: KeyPressed, Chr: KEY_F12}
		case "2":
			return &Event{EventType: KeyPressed, Chr: KEY_INSERT}
		case "3":
			return &Event{EventType: KeyPressed, Chr: KEY_DELETE}
		case "5":
			return &Event{EventType: KeyPressed, Chr: KEY_PAGE_UP}
		case "6":
			return &Event{EventType: KeyPressed, Chr: KEY_PAGE_DOWN}
		case "7":
			return &Event{EventType: KeyPressed, Chr: KEY_HOME}
		case "8":
			return &Event{EventType: KeyPressed, Chr: KEY_END}
		}
		return &Event{EventType: KeyPressed, Chr: KEY_LEFT}
	case "M":
		switch(subparameters[0]) {
		// NOTE: Each mouse move is its corresponding button + 32
		// e.g. <32 is sent after a left click

		// NOTE: Row and Column positions are 0 indexed

		// FIXME: Add support for alt, shift and ctrl as read in
		// https://tintin.mudhalla.net/info/xterm/
		case "<32", "<33", "<34", "<96", "<97":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseMove, X: x - 1, Y: y - 1}
		case "<0":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseLeftClick, X: x - 1, Y: y - 1}
		case "<1":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseMiddleClick, X: x - 1, Y: y - 1}
		case "<2":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseRightClick, X: x - 1, Y: y - 1}
		case "<65":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: ScrollDown, X: x - 1, Y: y - 1}
		case "<64":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: ScrollUp, X: x - 1, Y: y - 1}
		}
	case "m":
		switch(subparameters[0]) {
		case "<0":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseLeftRelease, X: x - 1, Y: y - 1}
		case "<1":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseMiddleRelease, X: x - 1, Y: y - 1}
		case "<2":
			x, _ := strconv.Atoi(subparameters[1])
			y, _ := strconv.Atoi(subparameters[2])
			return &Event{EventType: MouseRightRelease, X: x - 1, Y: y - 1}
		}
	}
	return &Event{EventType: Unknown}
}

// Read ECMA-48 5.4
// Returned string is a list in [P, I, P, I, ..., F] format,
// where P are parameters (like numbers), I are intermediate bytes,
// and F is the final byte, indicating the command
func (l *Listener) divideCommandSequence() (parts []string) {
	<- l.inChannel // Drop CSI
	inParameter := true

	parts = []string{""}
	for {
		chr := <- l.inChannel
		if isParameterByte(chr) {
			if inParameter {
				parts[len(parts)-1] += string(chr)
			} else {
				inParameter = true
				parts = append(parts, string(chr))
			}
		} else if isIntermediateByte(chr) {
			if inParameter {
				inParameter = false
				parts = append(parts, string(chr))
			} else {
				parts[len(parts)-1] += string(chr)
			}
		} else { // Is final byte
			parts = append(parts, string(chr))
			return parts
		}
	}
}

// Read ECMA-48 5.4
// NOTE: AA/BB notation is read as AA * 16 + BB
func isParameterByte(chr rune) bool {
	return chr >= 48 && chr <= 63
}

// Read ECMA-48 5.4
// NOTE: AA/BB notation is read as AA * 16 + BB
func isIntermediateByte(chr rune) bool {
	return chr >= 32 && chr <= 47
}
