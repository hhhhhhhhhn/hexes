package main

import (
	"os"
	"bufio"
	"runtime/pprof"

	"github.com/hhhhhhhhhn/hexes"
)

func colors() {
	cpuProf, _ := os.Create("cpuprofile")
	defer cpuProf.Close()
	pprof.StartCPUProfile(cpuProf)

	out := bufio.NewWriterSize(os.Stdout, 10000)

	// If you don't want buffering, you can use os.Stdout
	r := hexes.New(os.Stdin, out)
	r.SetDefaultAttribute(hexes.Join(hexes.NORMAL, hexes.BG_WHITE, hexes.GREEN))
	r.Start()

	colors := [][]hexes.Attribute{}
	for row := 0; row < r.Rows; row++ {
		arr := []hexes.Attribute{}
		for col := 0; col < r.Cols; col++ {
			arr = append(arr, hexes.TrueColorBg(row * 255 / r.Rows, col * 255 / r.Cols, 0))
		}
		colors = append(colors, arr)
	}

	for i := 0; i < 50 * r.Rows; i++ {
		for row := 0; row < r.Rows; row++ {
			for col := 0; col < r.Cols; col++ {
				r.SetAttribute(colors[(row+i) % r.Rows][(col+i) % r.Cols])
				r.Set(row, col, ' ')
			}
		}
		out.Flush()
	}
	r.End()
	out.Flush()

	pprof.StopCPUProfile()
}

func main() {
	set()
	colors()
}
