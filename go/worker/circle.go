package worker

import (
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"

	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Circle struct {
	Color    color.RGBA
	Position util.Vector
	Radius   float64
}

func createCircle(radius float64, bounds util.Vector, rng *rand.Rand) Shape {
	return Circle{
		Color:    util.RandomColor(rng),
		Position: util.RandomVector(rng, bounds),
		Radius:   (rng.Float64()*radius)/2.0 + radius/2.0, // from radius/2 to radius
	}
}

func (c Circle) Draw(dc *gg.Context, offset util.Vector) {
	dc.DrawCircle(c.Position.X+offset.X, c.Position.Y+offset.Y, c.Radius)
	dc.SetColor(c.Color)
	dc.Fill()
}
