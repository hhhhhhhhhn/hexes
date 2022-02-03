package main

import (
	"time"
	"unicode"
	"fmt"
	"os"
	"bufio"

	"github.com/hhhhhhhhhn/hexes"
)

func main() {
	duration := 1 * time.Millisecond

	out := bufio.NewWriterSize(os.Stdout, 10000)

	// If you don't want buffering, you can use os.Stdout
	r := hexes.New(os.Stdin, out)
	r.SetDefaultAttribute(hexes.NORMAL + hexes.BG_WHITE + hexes.GREEN)
	r.Start()

	// Makes sure reset signals are sent
	r.OnEnd(func (*hexes.Renderer) {
		out.Flush()
	})

	for i := 0; i < 10000; i++ {
		for !unicode.IsGraphic(rune(i)) {
			i++
		}
		row := i % r.Rows
		col := (i * 4) % r.Cols

		r.SetString(row, col, fmt.Sprint(string(rune(i))))
		time.Sleep(duration)

		if i % 1000 == 0 {
			r.SetAttribute(r.DefaultAttribute)
		}
		if i % 1000 == 200 {
			r.SetAttribute(hexes.RED)
		}
		if i % 1000 == 400 {
			r.SetAttribute(hexes.BG_CYAN + hexes.WHITE + hexes.BOLD)
		}
		if i % 1000 == 600 {
			r.SetAttribute(hexes.BG_YELLOW + hexes.RED + hexes.BOLD + hexes.ITALIC)
		}
		if i % 1000 == 800 {
			r.SetAttribute(hexes.REVERSE)
		}

		out.Flush()
	}

	for row := 0; row < r.Rows; row++ {
		for col := 0; col < r.Cols; col++ {
			r.SetAttribute(
				hexes.TrueColorBg(row * 255 / r.Rows, col * 255 / r.Cols, 0))
			r.SetString(row, col, " ")
		}
	}
	out.Flush()
	time.Sleep(1000 * duration)

	colors := [][]string{}
	for row := 0; row < r.Rows; row++ {
		arr := []string{}
		for col := 0; col < r.Cols; col++ {
			arr = append(arr, hexes.TrueColorBg(row * 255 / r.Rows, col * 255 / r.Cols, 0))
		}
		colors = append(colors, arr)
	}

	for i := 0; i < 3 * r.Rows; i++ {
		for row := 0; row < r.Rows; row++ {
			for col := 0; col < r.Cols; col++ {
				r.SetAttribute(colors[(row+i) % r.Rows][(col+i) % r.Cols])
				r.SetString(row, col, fmt.Sprint(" "))
			}
		}
		out.Flush()
	}
	r.End()
}
