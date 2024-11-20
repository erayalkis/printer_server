package main

import (
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/kenshaw/escpos"
	"github.com/nfnt/resize"
)

func printText(p *escpos.Escpos, body TextPrintPayload) {
	p.Init()

	p.SetSmooth(1)
	p.SetFont("A")

	parsedBody := fmt.Sprintf("%s\n", body.Text)
	p.Write(parsedBody)

	p.Formfeed()

	p.End()
}

func preprocessImage(file io.Reader, width int) ([]byte, error) {
	// Decode the image from the provided reader
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Resize the image to the printer's width
	resized := resize.Resize(uint(width), 0, img, resize.Lanczos3)

	// Convert to monochrome bitmap
	bounds := resized.Bounds()
	bitmap := make([]byte, (bounds.Dx()*bounds.Dy())/8)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(resized.At(x, y)).(color.Gray)
			byteIndex := (y*bounds.Dx() + x) / 8
			bitIndex := uint(x % 8)
			if gray.Y < 128 { // Threshold for black pixel
				bitmap[byteIndex] |= (1 << (7 - bitIndex))
			}
		}
	}

	return bitmap, nil
}
