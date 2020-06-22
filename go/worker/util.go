package worker

import (
	"encoding/gob"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

func Register() {
	gob.Register(color.RGBA{})
	gob.Register(image.YCbCr{})

	gob.Register(Circle{})
	gob.Register(Polygon{})
	gob.Register(Triangle{})
	gob.Register(Shapes{})
}

func randomColor(rng *rand.Rand) color.RGBA {
	f := func() uint8 {
		return uint8(rng.Intn(64) * 4)
	}

	return color.RGBA{f(), f(), f(), f()}
}

func randomVector(rng *rand.Rand, bounds util.Vector) util.Vector {
	return util.Vector{X: rng.Float64() * bounds.X, Y: rng.Float64() * bounds.Y}
}

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
