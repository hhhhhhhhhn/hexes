package main

import (
	"os"
	"bufio"
	"runtime/pprof"
	"unicode"

	"github.com/hhhhhhhhhn/hexes"
)

func set() {
	cpuProf, _ := os.Create("setprofile")
	defer cpuProf.Close()
	pprof.StartCPUProfile(cpuProf)

	out := bufio.NewWriterSize(os.Stdout, 10000)

	// If you don't want buffering, you can use os.Stdout
	r := hexes.New(os.Stdin, out)
	r.SetDefaultAttribute(hexes.Join(hexes.NORMAL, hexes.BG_WHITE, hexes.GREEN))
	r.Start()

	var j rune
	for i := 0; i < 50 * r.Rows; i++ {
		j++
		for !unicode.IsGraphic(j) {
			j++
		}
		for row := 0; row < r.Rows; row++ {
			for col := 0; col < r.Cols; col++ {
				r.Set(row, col, j)
			}
		}
		out.Flush()
	}
	r.End()
	out.Flush()

	pprof.StopCPUProfile()
}
