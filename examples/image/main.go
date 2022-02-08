package main

import (
	"bufio"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/hhhhhhhhhn/hexes"
)

func main() {
	img, err := loadImage()
	if err != nil {
		panic(err)
	}
	bounds := img.Bounds()

	out := bufio.NewWriterSize(os.Stdout, 4096)

	renderer := hexes.New(os.Stdin, out)
	renderer.Start()

	imageWidth   := bounds.Max.X - bounds.Min.X
	imageHeight  := bounds.Max.Y - bounds.Min.Y
	yOffset      := -bounds.Min.Y
	xOffset      := -bounds.Min.X

	for y := 0; y < renderer.Rows; y++ {
		for x := 0; x < renderer.Cols; x++ {
			imageX := x * imageWidth / renderer.Cols + xOffset
			imageY := y * imageHeight / renderer.Rows + yOffset

			r, g, b, _ := img.At(imageX, imageY).RGBA()
			r /= 256; g /= 256; b /= 256

			renderer.SetAttribute(hexes.TrueColorBg(int(r), int(g), int(b)))
			renderer.SetString(y, x, " ")
		}
	}

	out.Flush()
	
	os.Stdin.Read(make([]byte, 1))

	renderer.End()
	out.Flush()
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
