package main

import (
	"github.com/hhhhhhhhhn/hexes"
)

type BlockType int

const (
	L = 0
	J = 1
	I = 2
	O = 3
	Z = 4
	S = 5
	T = 6
	Empty = 7
)

var Blocks = [][]Loc {
	{ // L
		{0, 2},
		{0, 0},
		{0, -2},
		{2, -2},
	},
	{ // J
		{0, 2},
		{0, 0},
		{0, -2},
		{-2, -2},
	},
	{ // I
		{1, 3},
		{1, 1},
		{1, -1},
		{1, -3},
	},
	{ // O
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	},
	{ // Z
		{0, 0},
		{-2, 0},
		{0, 2},
		{2, 2},
	},
	{ // S
		{0, 0},
		{2, 0},
		{0, 2},
		{-2, 2},
	},
	{ // T
		{0, 0},
		{2, 0},
		{-2, 0},
		{0, 2},
	},
}

var Attributes []hexes.Attribute = []hexes.Attribute{
	hexes.BG_WHITE,
	hexes.BG_BLUE,
	hexes.BG_CYAN,
	hexes.BG_YELLOW,
	hexes.BG_RED,
	hexes.BG_GREEN,
	hexes.BG_MAGENTA,
	hexes.NORMAL,
}
