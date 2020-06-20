package worker

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Polygon struct {
	Color    color.RGBA
	Position util.Vector
	Radius   float64
	Rotation float64
	Sides    int
}

func createPolygon(radius float64, bounds util.Vector, rng *rand.Rand) Shape {
	return Polygon{
		Color:    util.RandomColor(rng),
		Position: util.RandomVector(rng, bounds),
		Radius:   (rng.Float64()*radius)/2.0 + radius/2.0,
		Rotation: (rng.Float64() * 2 * math.Pi),
		Sides:    rng.Intn(8) + 1,
	}
}

func (p Polygon) Draw(dc *gg.Context, offset util.Vector) {
	dc.DrawRegularPolygon(p.Sides-1, p.Position.X+offset.X, p.Position.Y+offset.Y, p.Radius, p.Rotation)
	dc.SetColor(p.Color)
	dc.Fill()
}
