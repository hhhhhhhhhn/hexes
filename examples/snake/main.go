package main

import (
	"bufio"
	"os"
	"time"
	"math/rand"

	"github.com/hhhhhhhhhn/hexes"
)
type cell int
const (
	EMPTY cell = iota
	SNAKE
	FRUIT
)

type direction int
const (
	LEFT  = -1
	RIGHT = 1
	UP    = -2
	DOWN  = 2
)

var renderer   *hexes.Renderer
var grid       [][]cell
var out        *bufio.Writer
var snake      [][]int
var snakeDir   direction
var wantedDir  []direction // A queue
var difficulty int = 20

func main() {
	out = bufio.NewWriterSize(os.Stdin, 4096)
	renderer = hexes.New(os.Stdin, out)
	renderer.Start()
	renderer.OnEnd(func(*hexes.Renderer){
		out.Flush()
	})
	grid = createGrid(renderer.Rows, renderer.Cols / 2)

	grid[2][3] = FRUIT
	snakeDir = LEFT
	wantedDir = []direction{LEFT}

	snake = [][]int{{renderer.Rows/2, renderer.Cols/4}, {renderer.Rows/2, renderer.Cols/4 + 1}, {renderer.Rows/2 + 2, renderer.Cols/4}}

	go handleInput()

	for {
		moveSnake()
		renderGrid()
		time.Sleep(time.Second / time.Duration(difficulty))
	}
}

func changeDir(dir direction) {
	if len(wantedDir) > 0 {
		if wantedDir[len(wantedDir) - 1] == -dir || wantedDir[len(wantedDir) - 1] == dir {
			return
		}
	} else if snakeDir == -dir {
		return
	}
	wantedDir = append(wantedDir, dir)
}

func handleInput() {
	in := bufio.NewReader(os.Stdin)
	for {
		chr, _, _ := in.ReadRune()
		switch(chr) {
		case 'h':
			changeDir(LEFT)
			break
		case 'l':
			changeDir(RIGHT)
			break
		case 'j':
			changeDir(DOWN)
			break
		case 'k':
			changeDir(UP)
			break
		case 'r':
			renderer.Refresh()
			break
		case '+':
			difficulty++
			break
		case '-':
			difficulty--
			if difficulty < 1 {
				difficulty = 1
			}
			break
		case 'q':
			renderer.End()
			os.Exit(0)
			break
		}
	}
}

func mod(a, b int) int {
	return (a % b + b) % b
}

func moveSnake() {
	if len(wantedDir) > 0 {
		snakeDir = wantedDir[0]
		wantedDir = wantedDir[1:]
	}
	grid[snake[0][0]][snake[0][1]] = EMPTY
	head := snake[len(snake)-1]

	var newHead []int
	switch(snakeDir) {
	case UP:
		newHead = []int{head[0] - 1, head[1]}
		break
	case DOWN:
		newHead = []int{head[0] + 1, head[1]}
		break
	case LEFT:
		newHead = []int{head[0], head[1] - 1}
		break
	case RIGHT:
		newHead = []int{head[0], head[1] + 1}
		break
	}

	newHead[0] = mod(newHead[0], len(grid))
	newHead[1] = mod(newHead[1], len(grid[0]))

	switch (grid[newHead[0]][newHead[1]]) {
	case FRUIT:
		spawnFruit()
	case EMPTY:
		snake = snake[1:]
	case SNAKE:
		renderer.End()
		os.Exit(0)
	}

	grid[newHead[0]][newHead[1]] = SNAKE
	snake = append(snake, newHead)
}

func spawnFruit() {
	for i := 0; i < 10000; i++ {
		row := rand.Intn(len(grid))
		col := rand.Intn(len(grid[0]))
		if grid[row][col] == EMPTY {
			grid[row][col] = FRUIT
			return
		}
	}
	renderer.End()
	os.Exit(0)
}

func renderGrid() {
	for y, row := range grid {
		for x, cell := range row {
			switch(cell) {
			case EMPTY:
				renderer.SetAttribute(renderer.DefaultAttribute)
				renderer.SetString(y, x*2, "  ")
				break
			case SNAKE:
				renderer.SetAttribute(hexes.REVERSE)
				renderer.SetString(y, x*2, "  ")
				break
			case FRUIT:
				renderer.SetAttribute(hexes.BOLD + hexes.RED)
				renderer.SetString(y, x*2, "()")
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
			row = append(row, EMPTY)
		}
		grid = append(grid, row)
	}
	return grid
}
