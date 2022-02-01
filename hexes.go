package hexes

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"golang.org/x/text/width"
)

type Renderer struct {
	Lines            [][]string
	Attributes       [][]string
	Rows             int
	Cols             int
	CursorRow        int
	CursorCol        int
	CurrentAttribute string
	DefaultAttribute string
}

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

func New() *Renderer {
	return &Renderer{DefaultAttribute: NORMAL}
}

func commandWithStdin(name string, arg ...string) *exec.Cmd {
	command := exec.Command(name, arg...)
	command.Stdin = os.Stdin
	return command
}

func (r *Renderer) Start() {
	command := commandWithStdin("stty", "-icanon", "-echo")
	command.Stdout = os.Stdout
	command.Run()

	command = commandWithStdin("tput", "rmam", "civis")
	command.Stdout = os.Stdout
	command.Run()

	r.updateRowsAndCols()

	for i := 0; i < r.Rows; i++ {
		line := []string{}
		attributes := []string{}
		for j := 0; j < r.Cols; j++ {
			line = append(line, " ")
			attributes = append(attributes, r.DefaultAttribute)
		}
		r.Lines = append(r.Lines, line)
		r.Attributes = append(r.Attributes, attributes)
	}

	r.CurrentAttribute = r.DefaultAttribute

	r.Refresh()
	r.setupSignals()
}

func (r *Renderer) updateRowsAndCols() {
	rows, _ := commandWithStdin("tput", "lines").Output()
	cols, _ := commandWithStdin("tput", "cols").Output()
	r.Rows, _ = strconv.Atoi(string(rows[:len(rows)-1]))
	r.Cols, _ = strconv.Atoi(string(cols[:len(cols) - 1]))
}

func (r *Renderer) resizeLinesAndAttributes() {
	currentRows := len(r.Lines)
	newRows     := r.Rows
	currentCols := len(r.Lines[0])
	newCols     := r.Cols

	if newCols < currentCols {
		for i := 0; i < currentRows; i++ {
			r.Lines[i] = r.Lines[i][:newCols]
			r.Attributes[i] = r.Attributes[i][:newCols]
		}
	} else if newCols > currentCols {
		for i := 0; i < currentRows; i++ {
			for j := 0; j < newCols - currentCols; j++ {
				r.Lines[i] = append(r.Lines[i], " ")
				r.Attributes[i] = append(r.Attributes[i], r.DefaultAttribute)
			}
		}
	}

	if newRows < currentRows {
		r.Lines = r.Lines[:newRows]
		r.Attributes = r.Attributes[:newRows]
	} else if newRows > currentRows {
		for i := 0; i < newCols - currentRows; i++ {
			row := []string{}
			attributes := []string{}
			for j := 0; j < newCols; j++ {
				row = append(row, " ")
				attributes = append(attributes, r.DefaultAttribute)
			}
			r.Lines = append(r.Lines, row)
			r.Attributes = append(r.Attributes, attributes)
		}
	}
}

// Turn into refresh
func (r *Renderer) Refresh() {
	r.updateRowsAndCols()
	r.resizeLinesAndAttributes()
	r.SetAttribute(r.DefaultAttribute)

	fmt.Print("\033[H") // Move to top left corner
	r.CursorCol = 0
	r.CursorRow = 0
	fmt.Print("\033[J") // Clear to end of screen

	r.redraw()
}

func (r *Renderer) redraw() {
	for row := 0; row < r.Rows; row++ {
		for col := 0; col < r.Cols; col++ {
			r.MoveCursor(row, col)
			fmt.Print(r.Attributes[row][col])
			fmt.Print(r.Lines[row][col])
			r.CursorCol++
		}
	}
}

func (r *Renderer) MoveCursor(row, col int) {
	if r.CursorRow == row && r.CursorCol == col {
		return
	}
	r.CursorRow = row
	r.CursorCol = col
	fmt.Printf("\033[%v;%vH", row + 1, col + 1)
}

func (r *Renderer) setupSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		r.End()
		fmt.Fprintln(os.Stderr, "exiting")
		os.Exit(0)
	}()

	//// This makes sure WINCH signals aren't spammed
	//refresh := make(chan int)
	//refreshTimeouts := make(chan int)
	//go func() {
	//	latestRefreshId := 0
	//	for {
	//		select {
	//			case id := <- refresh:
	//				latestRefreshId = id
	//				go func() {
	//					time.Sleep(100 * time.Millisecond)
	//					refreshTimeouts <- id
	//				}()
	//			case id := <-refreshTimeouts:
	//				if id == latestRefreshId {
	//					r.Refresh()
	//				}
	//		}
	//	}
	//}()

	//w := make(chan os.Signal)
	//signal.Notify(w, syscall.SIGWINCH)
	//go func() {
	//	refreshId := 0
	//	for {
	//		<-w
	//		refresh <- (refreshId)
	//		refreshId++
	//		r.Refresh()
	//	}
	//}()
}

func (r *Renderer) End() {
	command := commandWithStdin("stty", "sane")
	command.Stdout = os.Stdout
	command.Run()

	command = commandWithStdin("tput", "smam", "cnorm")
	command.Stdout = os.Stdout
	command.Run()

	fmt.Print(NORMAL)
	fmt.Print("\033[H") // Move to top left corner
	fmt.Print("\033[J") // Clear to end of line
}

func (r *Renderer) Set(row, col int, value string) {
	if (row > r.Rows - 1 || col > r.Cols - 1 || (r.Lines[row][col] == value && r.Attributes[row][col] == r.CurrentAttribute)) {
		return
	}
	r.MoveCursor(row, col)
	r.Lines[row][col] = value
	r.Attributes[row][col] = r.CurrentAttribute
	fmt.Print(value)
	r.CursorCol++
}

func (r *Renderer) SetString(row, col int, value string) {
	value = width.Narrow.String(value)
	chrIndex := 0
	for _, chr := range value {
		r.Set(row, col + chrIndex, string(chr))
		chrIndex++
	}
}

func (r *Renderer) SetAttribute(attribute string) {
	if attribute == r.DefaultAttribute {
		r.CurrentAttribute = r.DefaultAttribute
	} else {
		r.CurrentAttribute = r.DefaultAttribute + attribute
	}
	fmt.Print(r.CurrentAttribute)
}

func (r *Renderer) SetDefaultAttribute(attribute string) {
	r.DefaultAttribute = attribute
}
