package main

import (
	"bufio"
	"os"
	"time"
	"strconv"

	"github.com/hhhhhhhhhn/hexes"
)
type cell bool
const (
	DEAD cell  = false
	ALIVE cell = true
)

var renderer   *hexes.Renderer
var grid       [][]cell
var out        *bufio.Writer
var lastSpeed  time.Duration = 5
var speed      time.Duration = 5
var mouseDown  bool
var mouseX     int
var mouseY     int

func main() {
	out = bufio.NewWriterSize(os.Stdin, 4096)
	os.Stdout.WriteString("\033[?1003;1006;1015h")
	renderer = hexes.New(os.Stdin, out)
	renderer.Start()
	renderer.OnEnd(func(*hexes.Renderer){
		out.Flush()
		os.Stdout.WriteString("\033[?1003;1006;1015l")
	})
	grid = createGrid(renderer.Rows, renderer.Cols / 2)

	go handleInput()

	lastTime := time.Now()

	for {
		renderGrid()
		if speed != 0 && (time.Since(lastTime) > time.Second / time.Duration(speed)) {
			step()
			lastTime = time.Now()
		}
	}
}

func handleInput() {
	in := bufio.NewReader(os.Stdin)
	for {
		chr, _, _ := in.ReadRune()
		switch(chr) {
		case 27:
			sequence := readEscapeSequence(in)
			parsedSequence := parseEscapeSequence(sequence)
			if parsedSequence[0] == "[<" {
				if parsedSequence[1] == "32" {
					mouseX, _ = strconv.Atoi(parsedSequence[3])
					if mouseX > 0 {
						mouseX--
					}
					mouseX /= 2
					mouseY, _ = strconv.Atoi(parsedSequence[5])
					mouseY--

					if mouseDown {
						grid[mouseY][mouseX] = !grid[mouseY][mouseX]
					}
				}
				if parsedSequence[1] == "0" {
					mouseX, _ = strconv.Atoi(parsedSequence[3])
					if mouseX > 0 {
						mouseX--
					}
					mouseX /= 2
					mouseY, _ = strconv.Atoi(parsedSequence[5])
					mouseY--
					if parsedSequence[6] == "M" {
						mouseDown = true
					} else {
						mouseDown = false
					}

					if mouseDown {
						grid[mouseY][mouseX] = !grid[mouseY][mouseX]
					}
				}
			}
			break
		case '+':
			speed++
			lastSpeed++
		case '-':
			if speed > 1 {
				speed--
				lastSpeed--
			}
		case ' ':
			if speed == 0 {
				speed = lastSpeed
			} else {
				speed = 0
			}
		case 'q':
			renderer.End()
			os.Exit(0)
		default:
			break
		}
	}
}

func parseEscapeSequence(escape string) (parsed []string) {
	inDigit := true
	for _, chr := range escape {
		if inDigit != isDigit(chr) {
			parsed = append(parsed, "")
			inDigit = !inDigit
		}
		parsed[len(parsed) - 1] += string(chr)
	}
	return parsed
}

func isDigit(chr rune) bool {
	return (chr >= '0' && chr <= '9')
}

func readEscapeSequence(in *bufio.Reader) (sequence string) {
	sequence = ""
	escapeStart := time.Now()
	for {
		chr, _, _ := in.ReadRune()
		currentTime := time.Now()
		if currentTime.Sub(escapeStart) > time.Millisecond * 20 {
			in.UnreadRune()
			return sequence
		}
		sequence += string(chr)
	}
}

func mod(a, b int) int {
	return (a % b + b) % b
}

func getCell(row, col int) cell {
	row = mod(row, len(grid))
	col = mod(col, len(grid[0]))
	return grid[row][col]
}

func neighbours(row, col int) (neighbours int) {
	if getCell(row - 1, col - 1) == ALIVE {
		neighbours++
	}
	if getCell(row - 1, col) == ALIVE {
		neighbours++
	}
	if getCell(row - 1, col + 1) == ALIVE {
		neighbours++
	}
	if getCell(row, col - 1) == ALIVE {
		neighbours++
	}
	if getCell(row, col + 1) == ALIVE {
		neighbours++
	}
	if getCell(row + 1, col - 1) == ALIVE {
		neighbours++
	}
	if getCell(row + 1, col) == ALIVE {
		neighbours++
	}
	if getCell(row + 1, col + 1) == ALIVE {
		neighbours++
	}
	return neighbours
}

func step() {
	newGrid := createGrid(len(grid), len(grid[0]))
	for i, row := range grid {
		for j, cell := range row {
			neighbourAmount := neighbours(i, j)
			if cell == ALIVE {
				if neighbourAmount == 2 || neighbourAmount == 3 {
					newGrid[i][j] = ALIVE
				} else {
					newGrid[i][j] = DEAD
				}
			} else { // cell == DEAD
				if neighbourAmount == 3 {
					newGrid[i][j] = ALIVE
				} else {
					newGrid[i][j] = DEAD
				}
			}
		}
	}
	grid = newGrid
}

func renderGrid() {
	for y, row := range grid {
		for x, cell := range row {
			switch(cell) {
			case DEAD:
				renderer.SetAttribute(renderer.DefaultAttribute)
				renderer.SetString(y, x*2, "  ")
				break
			case ALIVE:
				renderer.SetAttribute(hexes.REVERSE)
				renderer.SetString(y, x*2, "  ")
				break
			}
		}
	}
	out.Flush()
}

func createGrid(rows, cols int) (grid [][]cell) {
	for row := 0; row < rows; row++ {
		row := []cell{}
		for col := 0; col < cols; col++ {
			row = append(row, DEAD)
		}
		grid = append(grid, row)
	}
	return grid
}
