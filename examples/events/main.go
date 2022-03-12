package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/hhhhhhhhhn/hexes"
	"github.com/hhhhhhhhhn/hexes/input"
)

func main() {
	out := bufio.NewWriterSize(os.Stdin, 4096)
	renderer := hexes.New(os.Stdin, out)

	listener := input.New(os.Stdin)
	listener.EnableMouseTracking(out)
	out.Flush()

	renderer.Start()

	for {
		event := listener.GetEvent()
		switch(event.EventType) {
		case input.KeyPressed:
			if unicode.IsGraphic(event.Chr) {
				printLine(renderer, "Key Pressed: " + string(event.Chr))
			} else {
				printLine(renderer, "Key Pressed: " + fmt.Sprint(event.Chr))
			}
			out.Flush()
			if event.Chr == 'q'{
				renderer.End()
				listener.DisableMouseTracking(out)
				out.Flush()
				return
			}
			break
		case input.MouseMove:
			printLine(renderer, "Mouse Move: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.MouseLeftClick:
			printLine(renderer, "Mouse Left Click: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.MouseLeftRelease:
			printLine(renderer, "Mouse Left Release: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.MouseMiddleClick:
			printLine(renderer, "Mouse Middle Click: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.MouseMiddleRelease:
			printLine(renderer, "Mouse Middle Release: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.MouseRightClick:
			printLine(renderer, "Mouse Right Click: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.MouseRightRelease:
			printLine(renderer, "Mouse Right Release: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.ScrollDown:
			printLine(renderer, "Scroll Down: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case input.ScrollUp:
			printLine(renderer, "Scroll Up: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		}
	}
}

func printLine(r *hexes.Renderer, line string) {
	line += strings.Repeat(" ", r.Cols - len(line))
	r.SetString(0, 0, line)
}
