package worker

import (
	"image"
	"image/draw"
	"math"
)

func squareDifference(x, y float64) float64 {
	d := x - y
	return d * d
}

func imgDiff(a, b *image.RGBA) float64 {
	var d float64
	for i := 0; i < len(a.Pix); i++ {
		d += squareDifference(float64(a.Pix[i]), float64(b.Pix[i]))
	}

	return math.Sqrt(d)
}

func rgbaImg(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Src)
	return rgba
}
