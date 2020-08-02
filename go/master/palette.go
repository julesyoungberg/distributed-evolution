package master

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"gocv.io/x/gocv"
)

func getPalette(src image.Image, k int) ([]color.RGBA, error) {
	n := src.Bounds().Dx() * src.Bounds().Dy()

	data, err := gocv.ImageToMatRGBA(src)
	if err != nil {
		return []color.RGBA{}, fmt.Errorf("getting palette: %v", err)
	}
	defer data.Close()

	data.Reshape(1, n)
	data.ConvertTo(&data, gocv.MatTypeCV32F)

	labels := gocv.NewMat()
	criteria := gocv.NewTermCriteria(gocv.MaxIter+gocv.EPS, 10, 1.0)
	colors := gocv.NewMat()

	gocv.KMeans(data, k, &labels, criteria, 1, gocv.KMeansPPCenters, &colors)

	out, err := colors.ToImage()
	if err != nil {
		return []color.RGBA{}, fmt.Errorf("getting palette: %v", err)
	}

	b := out.Bounds()
	dx := b.Dx()
	dy := b.Dy()

	palette := make([]color.RGBA, dx*dy)
	for i := range palette {
		x := i % dy
		y := i / dy

		r, g, b, a := out.At(x, y).RGBA()
		palette[i] = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	}

	return palette, nil
}

func (m *Master) preparePalette() {
	palette, err := getPalette(m.TargetImage.Image, 8)
	if err != nil {
		log.Fatal(err)
	}

	err = m.db.SetPalette(palette)
	if err != nil {
		log.Fatalf("error setting palette: %v", err)
	}
}
