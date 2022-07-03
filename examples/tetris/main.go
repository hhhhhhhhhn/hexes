package main

import (
	"math/rand"
	"time"
)

var grid [][]BlockType
var interval = time.Second / 2

func main() {
	initialise()
	grid = newGrid()
	piece := randomPiece()
	lockTimer := lockWait

	eventChannel := make(chan Event, 5)
	go keyboardListener(eventChannel)
	go gravityListener(eventChannel)
	for true {
		event := <- eventChannel
		switch event {
		case Hard:
			for {
				moved := movePiece(piece, 1, 0)
				if isValid(moved, grid) {
					piece = moved
				} else {
					lockTimer = 0
					break
				}
			}
			break
		case Soft:
			moved := movePiece(piece, 1, 0)
			if isValid(moved, grid) {
				piece = moved
			} else {
				lockTimer -= 1
			}
			break
		case Left:
			moved := movePiece(piece, 0, -1)
			if isValid(moved, grid) {
				piece = moved
			}
			break
		case Right:
			moved := movePiece(piece, 0, 1)
			if isValid(moved, grid) {
				piece = moved
			}
			break
		case CCW:
			rotated := rotatePiece(piece)
			if isValid(rotated, grid) {
				piece = rotated
				lockTimer = lockWait
			}
			break
		case CW:
			rotated := rotatePiece(rotatePiece(rotatePiece(piece)))
			if isValid(rotated, grid) {
				piece = rotated
				lockTimer = lockWait
			}
			break
		case Quit:
			exit()
			break
		}
		if lockTimer <= 0 {
			lockTimer = lockWait
			piece = lockAndCreateNew(piece)
		}
		render(grid, piece)
	}
}

type Loc struct {
	Row int
	Col int
}

type Piece struct {
	Blocks   []Loc // Relative to piece origin
	Type     BlockType
	Position Loc
}

type Event int
const (
	CW Event = iota
	CCW
	Soft
	Hard
	Left
	Right
	Quit
)

const height = 25
const width = 10
const lockWait = 3

func gravityListener(channel chan Event) {
	for {
		channel <- Soft
		time.Sleep(interval)
	}
}

func newGrid() [][]BlockType {
	grid := make([][]BlockType, height)
	for i := 0; i < height; i++ {
		row := newRow()
		grid[i] = row
	}
	return grid
}


func isValid(piece Piece, grid [][]BlockType) bool {
	locs := pieceToLocs(piece)
	for _, loc := range locs {
		if loc.Row >= height || loc.Row < 0 || loc.Col >= width || loc.Col < 0 {
			return false
		}
		if grid[loc.Row][loc.Col] != Empty {
			return false
		}
	}
	return true
}

func lockAndCreateNew(piece Piece) Piece {
	pushToGrid(piece, grid)
	cleanFullLines(grid)
	piece = randomPiece()
	if !isValid(piece, grid) {
		exit()
	}
	interval = interval * 995 / 1000
	return piece
}

func pushToGrid(piece Piece, grid [][]BlockType) {
	locs := pieceToLocs(piece)
	for _, loc := range locs {
		grid[loc.Row][loc.Col] = piece.Type
	}

}

func randomPiece() Piece {
	blockType := rand.Int() % 7
	return Piece{Blocks: Blocks[blockType], Type: BlockType(blockType), Position: Loc{-1, 2}}
}

func cleanFullLines(grid [][]BlockType) {
	for y, row := range grid {
		if isFull(row) {
			moveDownFromRow(grid, y)
		}
	}
}

func isFull(row []BlockType) bool {
	for _, block := range row {
		if block == Empty {
			return false
		}
	}
	return true
}

func moveDownFromRow(grid [][]BlockType, row int) {
	for y := row; y > 0; y-- {
		grid[y] = grid[y-1]
	}
	grid[0] = newRow()
}


func newRow() []BlockType {
	row := make([]BlockType, width)
	for j := 0; j < width; j++ {
		row[j] = Empty
	}
	return row
}

func pieceToLocs(piece Piece) []Loc {
	var locs []Loc
	for _, loc := range piece.Blocks {
		locs = append(locs, Loc{
			Row: (loc.Row+4)/2 + piece.Position.Row,
			Col: (loc.Col+4)/2 + piece.Position.Col,
		})
	}
	return locs
}

func rotatePiece(piece Piece) Piece {
	rotated := piece
	rotated.Blocks = rotateLocs(piece.Blocks)
	return rotated
}

func movePiece(piece Piece, rows, cols int) Piece {
	moved := piece
	moved.Position.Row += rows
	moved.Position.Col += cols
	return moved
}

func cloneLocs(locs []Loc) []Loc {
	clone := make([]Loc, len(locs))
	copy(clone, locs)
	return clone
}

func rotateLocs(locs []Loc) []Loc {
	rotated := cloneLocs(locs)
	for i := range rotated {
		rotated[i].Col, rotated[i].Row = -rotated[i].Row, rotated[i].Col
	}
	return rotated
}
