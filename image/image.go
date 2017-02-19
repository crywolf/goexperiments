package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
)

func main() {
	m := Image{180, 200}
	saveImage(m, "pic2.png")
}

type Image struct {
	w, h int
}

// ColorModel returns the Image's color model.
func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the domain for which At can return non-zero color.
// The bounds do not necessarily contain the point (0, 0).
func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.w, i.h)
}

// At returns the color of the pixel at (x, y).
// At(Bounds().Min.X, Bounds().Min.Y) returns the upper-left pixel of the grid.
// At(Bounds().Max.X-1, Bounds().Max.Y-1) returns the lower-right one.
func (i Image) At(x, y int) color.Color {
	v := countColor(x, y)
	return color.RGBA{v, v, 255, 255}
}

func countColor(x, y int) uint8 {
	v := uint8(x ^ y)
	return uint8(v)
}

func saveImage(m image.Image, filename string) {
	var imgBuf bytes.Buffer

	err := png.Encode(&imgBuf, m)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(f, &imgBuf)

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image was saved to '" + filename + "'")
}
