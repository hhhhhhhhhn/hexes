package main

import (
	"os"

	"github.com/hhhhhhhhhn/hexes"
	"github.com/hhhhhhhhhn/hexes/input"
)

var renderer *hexes.Renderer
var listener *input.Listener
var xOffset int
var yOffset int

func initialise() {
	renderer = hexes.New(os.Stdin, os.Stdout)
	renderer.Start()
	yOffset = (renderer.Rows - height) / 2
	xOffset = (renderer.Cols - width*2) / 2
	listener = input.New(os.Stdin)
	renderBorder()
}

func exit() {
	renderer.End()
	os.Exit(0)
}

func render(grid [][]BlockType, piece Piece) {
	renderGrid(grid)
	renderGhost(piece, grid)
	renderPiece(piece)
}

func renderGrid(grid [][]BlockType) {
	for y, row := range grid {
		for x, block := range row {
			renderer.SetAttribute(Attributes[block])
			renderer.SetString(y + yOffset, x*2 + xOffset, "  ")
		}
	}
}

func renderPiece(piece Piece) {
	renderer.SetAttribute(Attributes[piece.Type])
	locs := pieceToLocs(piece)
	for _, loc := range locs {
		renderer.SetString(loc.Row + yOffset, loc.Col*2 + xOffset, "  ")
	}
}

func renderGhost(piece Piece, grid [][]BlockType) {
	for {
		moved := movePiece(piece, 1, 0)
		if isValid(moved, grid) {
			piece = moved
		} else {
			break
		}
	}
	renderer.SetAttribute(hexes.Join(Attributes[piece.Type], hexes.BLACK, hexes.REVERSE))
	locs := pieceToLocs(piece)
	for _, loc := range locs {
			renderer.SetString(loc.Row + yOffset, loc.Col*2 + xOffset, "░░")
	}
	renderer.SetAttribute(hexes.NORMAL)
}

func renderBorder() {
	renderer.SetAttribute(hexes.REVERSE)
	for y := yOffset; y < yOffset + height; y++ {
		renderer.SetString(y, xOffset - 2, "░░")
	}
	for x := xOffset - 2; x < xOffset + width*2; x++ {
		renderer.SetString(yOffset+height, x, "░░")
	}
	for y := yOffset; y < yOffset + height + 1; y++ {
		renderer.SetString(y, xOffset + width*2, "░░")
	}
	for y := yOffset; y < yOffset + height + 1; y++ {
		renderer.SetString(y, xOffset + width*2, "░░")
	}
	renderer.SetAttribute(hexes.NORMAL)
}

func keyboardListener(channel chan Event) {
	for {
		event := listener.GetEvent()
		for event.EventType != input.KeyPressed {
			event = listener.GetEvent()
		}
		switch event.Chr {
		case 'a':
			channel <- Left
			break
		case 'd':
			channel <- Right
			break
		case 'w':
			channel <- CCW
			break
		case 's':
			channel <- Soft
			break
		case ' ':
			channel <- Hard
			break
		case 'q':
			channel <- Quit
			break
		}
	}
}
