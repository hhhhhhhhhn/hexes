package hexes

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	runeWidth "github.com/mattn/go-runewidth"
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
	Out              io.Writer
	In               io.Reader
	onEnd            func(*Renderer)
}

func New(in io.Reader, out io.Writer) *Renderer {
	return &Renderer{DefaultAttribute: NORMAL, Out: out, In: in}
}

func (r *Renderer) commandWithStdin(name string, arg ...string) *exec.Cmd {
	command := exec.Command(name, arg...)
	command.Stdin = r.In
	return command
}

func (r *Renderer) Start() {
	command := r.commandWithStdin("stty", "-icanon", "-echo")
	command.Stdout = r.Out
	command.Run()

	command = r.commandWithStdin("tput", "rmam", "civis")
	command.Stdout = r.Out
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
	rows, _ := r.commandWithStdin("tput", "lines").Output()
	cols, _ := r.commandWithStdin("tput", "cols").Output()
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

func (r *Renderer) print(str string) {
	r.Out.Write([]byte(str))
}

// Turn into refresh
func (r *Renderer) Refresh() {
	r.updateRowsAndCols()
	r.resizeLinesAndAttributes()
	r.SetAttribute(r.DefaultAttribute)

	r.print("\033[H") // Move to top left corner
	r.CursorCol = 0
	r.CursorRow = 0
	r.print("\033[J") // Clear to end of screen

	r.redraw()
}

func (r *Renderer) redraw() {
	for row := 0; row < r.Rows; row++ {
		for col := 0; col < r.Cols; col++ {
			r.MoveCursor(row, col)
			r.print(r.Attributes[row][col])
			r.print(r.Lines[row][col])
			r.CursorCol++
		}
	}
}

func (r *Renderer) MoveCursor(row, col int) {
	// NOTE: This optimization doesn't always work, as some unicode characters
	// are 2 wide even if using the 'width' package to narrow them
	if r.CursorRow == row && r.CursorCol == col {
		return
	}
	r.CursorRow = row
	r.CursorCol = col
	r.print(fmt.Sprintf("\033[%v;%vH", row + 1, col + 1))
}

func (r *Renderer) setupSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		r.End()
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
	command := r.commandWithStdin("stty", "sane")
	command.Stdout = r.Out
	command.Run()

	command = r.commandWithStdin("tput", "smam", "cnorm")
	command.Stdout = r.Out
	command.Run()

	r.print(NORMAL)
	r.print("\033[H") // Move to top left corner
	r.print("\033[J") // Clear to end of line

	if (r.onEnd != nil)  {
		r.onEnd(r)
	}
}

func (r *Renderer) Set(row, col int, value string) {
	if (row > r.Rows - 1 || col > r.Cols - 1 || (r.Lines[row][col] == value && r.Attributes[row][col] == r.CurrentAttribute)) {
		return
	}
	oldWidth := runeWidth.StringWidth(r.Lines[row][col])
	width := runeWidth.StringWidth(value)
	r.MoveCursor(row, col)
	r.Lines[row][col] = value
	r.Attributes[row][col] = r.CurrentAttribute
	r.print(value)
	r.CursorCol += width

	if width < oldWidth && col < r.Cols - 1 {
		r.MoveCursor(row, col+1)
		r.SetAttribute(r.Attributes[row][col+1])
		r.print(r.Lines[row][col+1])
	}
}

func (r *Renderer) SetString(row, col int, value string) {
	for _, chr := range value {
		r.Set(row, col, string(chr))
		col += runeWidth.RuneWidth(chr)
	}
}

func (r *Renderer) SetAttribute(attribute string) {
	if attribute == r.DefaultAttribute {
		r.CurrentAttribute = r.DefaultAttribute
	} else {
		r.CurrentAttribute = r.DefaultAttribute + attribute
	}
	r.print(r.CurrentAttribute)
}

func (r *Renderer) SetDefaultAttribute(attribute string) {
	r.DefaultAttribute = attribute
}

func (r *Renderer) OnEnd(f func(r *Renderer)) {
	r.onEnd = f
}
