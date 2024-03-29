package main

import (
	"bufio"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"time"

	"github.com/hhhhhhhhhn/hexes"
	"github.com/hhhhhhhhhn/hexes/input"
)

var quit = false
var render = true
var refresh = false

func main() {
	img, err := loadImage()
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()

	out := bufio.NewWriterSize(os.Stdout, 500000)

	listener := input.New(os.Stdin)
	listener.EnableMouseTracking(out)
	out.Flush()
	renderer := hexes.New(os.Stdin, out)
	renderer.Start()

	imageWidth   := bounds.Max.X - bounds.Min.X
	imageHeight  := bounds.Max.Y - bounds.Min.Y
	yOffset      := -bounds.Min.Y + imageHeight * 50
	xOffset      := -bounds.Min.X + imageWidth * 50

	scale := 1000

	lastMouseX := 0
	lastMouseY := 0
	dragging := false

	go func() {
		for {
			event := listener.GetEvent()
			switch(event.EventType) {
			case input.ScrollDown:
				scale = scale * 10 / 9
				if scale < 10 {
					scale--
				}
				render = true
				break
			case input.ScrollUp:
				scale = scale * 9 / 10
				if scale < 10 {
					scale++
				}
				render = true
				break
			case input.MouseLeftClick:
				lastMouseX = event.X
				lastMouseY = event.Y
				dragging = true
				break
			case input.MouseLeftRelease:
				lastMouseX = event.X
				lastMouseY = event.Y
				dragging = false
				break
			case input.MouseMove:
				if dragging {
					xOffset += (lastMouseX - event.X) * imageHeight * scale * 2 * 100 / renderer.Cols / 1000
					yOffset += (lastMouseY - event.Y) * imageHeight * scale * 100 / renderer.Rows / 1000
					lastMouseX = event.X
					lastMouseY = event.Y
					render = true
				}
				break
			case input.KeyPressed:
				switch event.Chr {
				case 'q':
					quit = true
				case 'r':
					refresh = true
				}
				break
			}
		}
	}()

	for {
		if quit {
			renderer.End()
			listener.DisableMouseTracking(out)
			out.Flush()
			os.Exit(0)
		}
		if refresh {
			refresh = false
			render = true
			renderer.Refresh()
		}
		if render {
			render = false
			for y := 0; y < renderer.Rows; y++ {
				for x := 0; x < renderer.Cols; x++ {
					imageX := ((x - renderer.Cols / 2) * imageHeight * scale / 2 / renderer.Rows / 1000 + (xOffset / 100))
					imageTopY := ((y - renderer.Rows / 2) * imageHeight * scale / renderer.Rows / 1000 + (yOffset / 100))
					imageBottomY := (((y - renderer.Rows / 2) * imageHeight * scale + imageHeight*scale/2) / renderer.Rows / 1000 + (yOffset / 100))

					tr, tg, tb, _ := img.At(imageX, imageTopY).RGBA()
					tr /= 256; tg /= 256; tb /= 256
					br, bg, bb, _ := img.At(imageX, imageBottomY).RGBA()
					br /= 256; bg /= 256; bb /= 256

					attribute := hexes.Join(
						hexes.TrueColor(int(tr), int(tg), int(tb)),
						hexes.TrueColorBg(int(br), int(bg), int(bb)),
					)
					renderer.SetAttribute(attribute)
					renderer.Set(y, x, '▀')
				}
			}
			out.Flush()
		} else {
			time.Sleep(time.Second / 60)
		}
	}
}

func loadImage() (image.Image, error) {
	if len(os.Args) < 2 {
		return nil, errors.New("Please provide image file")
	}

	reader, err := os.Open(os.Args[1])
	defer reader.Close()

	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(reader)

	if err != nil {
		return nil, err
	}

	return img, nil
}
