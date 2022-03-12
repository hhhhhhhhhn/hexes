package hexes

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"unicode/utf8"

	runeWidth "github.com/mattn/go-runewidth"
)

// Gives an abstraction to render text in any position in the terminal.
type Renderer struct {
	Lines            [][]rune      // (RO) The virtual output characters.
	Attributes       [][]Attribute // (RO) The virtual output Attributes.
	Rows             int           // (RO) The amount of rows in virtual output.
	Cols             int           // (RO) The amount of columns in virtual output.
	CursorRow        int           // (RO) The row the cursor is in the terminal.
	CursorCol        int           // (RO) The column the cursor is in the terminal.
	CurrentAttribute Attribute     // (RO) The terminal's current attribute.
	DefaultAttribute Attribute     // (RO) The preferred default attribute of the renderer.
	Out              io.Writer     // (RO) The default writer to send terminal data to.
	In               io.Reader     // (RO) The default reader to get data from.
}

// Creates a new Renderer using in (e.g. os.Stdin) as input and
// out (e.g. os.Stdout) as output.
func New(in io.Reader, out io.Writer) *Renderer {
	return &Renderer{DefaultAttribute: NORMAL, Out: out, In: in}
}

func (r *Renderer) commandWithStdin(name string, arg ...string) *exec.Cmd {
	command := exec.Command(name, arg...)
	command.Stdin = r.In
	return command
}

// Initializes and configures the Renderer and user's terminal.
func (r *Renderer) Start() {
	command := r.commandWithStdin("stty", "-icanon", "-echo")
	command.Stdout = r.Out
	command.Run()

	command = r.commandWithStdin("tput", "rmam")
	command.Stdout = r.Out
	command.Run()

	r.write([]byte("\033[?25l")) // Hide cursor
	r.updateRowsAndCols()

	for i := 0; i < r.Rows; i++ {
		line := []rune{}
		attributes := []Attribute{}
		for j := 0; j < r.Cols; j++ {
			line = append(line, ' ')
			attributes = append(attributes, r.DefaultAttribute)
		}
		r.Lines = append(r.Lines, line)
		r.Attributes = append(r.Attributes, attributes)
	}

	r.CurrentAttribute = r.DefaultAttribute

	r.Refresh()
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
				r.Lines[i] = append(r.Lines[i], ' ')
				r.Attributes[i] = append(r.Attributes[i], r.DefaultAttribute)
			}
		}
	}

	if newRows < currentRows {
		r.Lines = r.Lines[:newRows]
		r.Attributes = r.Attributes[:newRows]
	} else if newRows > currentRows {
		for i := 0; i < newCols - currentRows; i++ {
			row := []rune{}
			attributes := []Attribute{}
			for j := 0; j < newCols; j++ {
				row = append(row, ' ')
				attributes = append(attributes, r.DefaultAttribute)
			}
			r.Lines = append(r.Lines, row)
			r.Attributes = append(r.Attributes, attributes)
		}
	}
}

func (r *Renderer) write(data []byte) {
	r.Out.Write(data)
}

var tmpBuf = make([]byte, 4)
func (r *Renderer) writeRune(chr rune) {
	length := utf8.EncodeRune(tmpBuf, chr)
	r.Out.Write(tmpBuf[:length])
}

// Redraws virtual output to the terminal, handling resizes.
func (r *Renderer) Refresh() {
	r.updateRowsAndCols()
	r.resizeLinesAndAttributes()
	r.SetAttribute(r.DefaultAttribute)

	r.write([]byte("\033[H")) // Move to top left corner
	r.CursorCol = 0
	r.CursorRow = 0
	r.write([]byte("\033[J")) // Clear to end of screen

	r.redraw()
}

func (r *Renderer) redraw() {
	for row := 0; row < r.Rows; row++ {
		for col := 0; col < r.Cols; col++ {
			r.MoveCursor(row, col)
			r.write(r.Attributes[row][col])
			r.writeRune(r.Lines[row][col])
			r.CursorCol++
		}
	}
}

// Moves terminal cursor to a position (0 indexed).
func (r *Renderer) MoveCursor(row, col int) {
	// NOTE: This optimization doesn't always work, as some unicode characters
	// are 2 wide even if using the 'width' package to narrow them
	if r.CursorRow == row && r.CursorCol == col {
		return
	}
	r.CursorRow = row
	r.CursorCol = col
	fmt.Fprintf(r.Out, "\033[%v;%vH", row + 1, col + 1)
}

// Sets the cell at row, col (0 indexed) to the character given.
func (r *Renderer) Set(row, col int, value rune) {
	if (
		row > r.Rows - 1 ||
		col > r.Cols - 1 ||
		(r.Lines[row][col] == value && &r.Attributes[row][col][0] == &r.CurrentAttribute[0])) {
			return
	}
	var oldWidth int
	width := runeWidth.RuneWidth(value)
	if width == 2 {
		oldWidth = runeWidth.RuneWidth(r.Lines[row][col])
	}

	r.MoveCursor(row, col)
	r.Lines[row][col] = value
	r.Attributes[row][col] = r.CurrentAttribute
	r.writeRune(value)
	r.CursorCol += width

	if width < oldWidth && col < r.Cols - 1 {
		r.MoveCursor(row, col+1)
		r.SetAttribute(r.Attributes[row][col+1])
		r.writeRune(r.Lines[row][col+1])
	}
}

// Sets the cells starting at row, col (0 indexed) to value,
// accounting for wide characters.
func (r *Renderer) SetString(row, col int, value string) {
	for _, chr := range value {
		r.Set(row, col, chr)
		col += runeWidth.RuneWidth(chr)
	}
}

// Sets the Attribute of the text being written.
func (r *Renderer) SetAttribute(attribute Attribute) {
	r.CurrentAttribute = attribute
	r.write(r.CurrentAttribute)
}

// Sets the preferred Attribute with which to prepend new ones.
func (r *Renderer) SetDefaultAttribute(attribute Attribute) {
	r.DefaultAttribute = attribute
}

// Creates a new attribute based on the default attribute.
func (r *Renderer) NewAttribute(attributes... Attribute) Attribute {
	return Join(r.DefaultAttribute, Join(attributes...))
}

// Restores terminal to default state.
func (r *Renderer) End() {
	command := r.commandWithStdin("stty", "sane")
	command.Stdout = r.Out
	command.Run()

	command = r.commandWithStdin("tput", "smam")
	command.Stdout = r.Out
	command.Run()

	r.write([]byte("\033[?25h")) // Show cursor
	r.write(NORMAL)
	r.write([]byte("\033[H")) // Move to top left corner
	r.write([]byte("\033[J")) // Clear to end
}

