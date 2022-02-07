package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/hhhhhhhhhn/hexes"
	"github.com/hhhhhhhhhn/hexes/listener"
)

func main() {
	out := bufio.NewWriterSize(os.Stdin, 4096)
	renderer := hexes.New(os.Stdin, out)

	eventListener := listener.New(os.Stdin)
	eventListener.EnableMouseTracking(out)
	out.Flush()

	renderer.Start()
	renderer.OnEnd(func(*hexes.Renderer){
		eventListener.DisableMouseTracking(out)
		out.Flush()
	})


	for {
		event := eventListener.GetEvent()
		switch(event.EventType) {
		case listener.KeyPressed:
			if unicode.IsGraphic(event.Chr) {
				printLine(renderer, "Key Pressed: " + string(event.Chr))
			} else {
				printLine(renderer, "Key Pressed: " + fmt.Sprint(event.Chr))
			}
			out.Flush()
			if event.Chr == 'q'{
				renderer.End()
				return
			}
			break
		case listener.MouseMove:
			printLine(renderer, "Mouse Move: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.MouseLeftClick:
			printLine(renderer, "Mouse Left Click: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.MouseLeftRelease:
			printLine(renderer, "Mouse Left Release: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.MouseMiddleClick:
			printLine(renderer, "Mouse Middle Click: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.MouseMiddleRelease:
			printLine(renderer, "Mouse Middle Release: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.MouseRightClick:
			printLine(renderer, "Mouse Right Click: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.MouseRightRelease:
			printLine(renderer, "Mouse Right Release: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.ScrollDown:
			printLine(renderer, "Scroll Down: " + fmt.Sprint(event.X, event.Y))
			out.Flush()
			break
		case listener.ScrollUp:
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
