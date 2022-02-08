package main

import (
	"bufio"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hhhhhhhhhn/hexes"
	"github.com/hhhhhhhhhn/hexes/input"
)

func main() {
	img, err := loadImage()
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()

	out := bufio.NewWriterSize(os.Stdout, 4096)

	listener := input.New(os.Stdin)
	listener.EnableMouseTracking(out)
	out.Flush()
	renderer := hexes.New(os.Stdin, out)
	renderer.OnEnd(func(*hexes.Renderer) {
		listener.DisableMouseTracking(out)
		out.Flush()
	})
	renderer.Start()

	imageWidth  := bounds.Max.X - bounds.Min.X
	imageHeight  := bounds.Max.Y - bounds.Min.Y
	yOffset      := -bounds.Min.Y + imageHeight * 50
	xOffset      := -bounds.Min.X + imageWidth * 50

	scale := 1000

	lastMouseX := 0
	lastMouseY := 0
	dragging := false

	for {
		for y := 0; y < renderer.Rows; y++ {
			for x := 0; x < renderer.Cols; x++ {
				imageX := ((x - renderer.Cols / 2) * imageHeight * 2 * scale / renderer.Cols / 1000 + xOffset / 100)
				imageY := ((y - renderer.Rows / 2) * imageHeight * scale / renderer.Rows / 1000 + yOffset / 100)

				//imageX := (x * imageHeight * 2 * scale / renderer.Cols / 1000 + xOffset)
				//imageY := (y * imageHeight * scale / renderer.Rows / 1000 + yOffset)

				r, g, b, _ := img.At(imageX, imageY).RGBA()
				r /= 256; g /= 256; b /= 256

				renderer.SetAttribute(hexes.TrueColorBg(int(r), int(g), int(b)))
				renderer.SetString(y, x, " ")
			}
		}
		out.Flush()

		for {
			event := listener.GetEvent()
			switch(event.EventType) {
			case input.ScrollDown:
				scale = scale * 10 / 9
				if scale < 10 {
					scale--
				}
			case input.ScrollUp:
				scale = scale * 9 / 10
				if scale < 10 {
					scale++
				}
			case input.MouseLeftClick:
				lastMouseX = event.X
				lastMouseY = event.Y
				dragging = true
			case input.MouseLeftRelease:
				lastMouseX = event.X
				lastMouseY = event.Y
				dragging = false
			case input.MouseMove:
				if dragging {
					xOffset += (lastMouseX - event.X) * imageHeight * scale * 2 * 100 / renderer.Cols / 1000
					yOffset += (lastMouseY - event.Y) * imageHeight * scale * 100 / renderer.Rows / 1000
					lastMouseX = event.X
					lastMouseY = event.Y
				} else {
					continue
				}
			default:
				continue
			}
			break
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
