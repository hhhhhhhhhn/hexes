package main

import (
	"bufio"
	"os"
	"time"

	"github.com/hhhhhhhhhn/hexes"
	"github.com/hhhhhhhhhn/hexes/input"
)
type cell bool
const (
	DEAD cell  = false
	ALIVE cell = true
)

var renderer   *hexes.Renderer
var listener   *input.Listener
var grid       [][]cell
var out        *bufio.Writer
var lastSpeed  time.Duration = 5
var speed      time.Duration = 5
var mouseDown  bool
var mouseX     int
var mouseY     int
var quit       bool          = false
var render     bool          = true

func main() {
	out = bufio.NewWriterSize(os.Stdin, 4096)
	renderer = hexes.New(os.Stdin, out)
	renderer.Start()

	listener = input.New(os.Stdin)
	listener.EnableMouseTracking(out)

	out.Flush()

	renderer.OnEnd(func(*hexes.Renderer){
		listener.DisableMouseTracking(out)
		out.Flush()
	})
	grid = createGrid(renderer.Rows, renderer.Cols / 2)

	go handleInput()

	lastTime := time.Now()

	for {
		if quit {
			renderer.End()
			os.Exit(0)
			return
		}
		if render {
			render = false
			renderGrid()
		}
		if speed != 0 && (time.Since(lastTime) > time.Second / time.Duration(speed)) {
			render = true
			step()
			lastTime = time.Now()
		}
		time.Sleep(time.Second / 60)
	}
}

func handleInput() {
	for {
		event := listener.GetEvent()
		switch event.EventType {
		case input.KeyPressed:
			switch event.Chr {
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
				quit = true
			}
		case input.MouseMove:
			updateMouse(event.Y, event.X)

			if mouseDown {
				grid[mouseY][mouseX] = ALIVE
				render = true
			}
			break
		case input.MouseLeftClick:
			updateMouse(event.Y, event.X)
			mouseDown = true

			grid[mouseY][mouseX] = !grid[mouseY][mouseX]

			render = true
			break
		case input.MouseLeftRelease:
			updateMouse(event.Y, event.X)
			mouseDown = false
			break
		}
	}
}

func updateMouse(row, col int) {
	mouseX = col / 2
	if mouseX == len(grid[0]) {
		mouseX--
	}
	mouseY = row
	if mouseY == len(grid) {
		mouseY--
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
