package worker

import (
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"

	"github.com/rickyfitts/distributed-evolution/util"
)

type Triangle struct {
	Bounds   util.Vector
	Color    color.RGBA
	Vertices []util.Vector
}

// creates a random triangle
func createTriangle(size float64, bounds util.Vector, rng *rand.Rand) Triangle {
	clr := color.RGBA{
		uint8(rng.Intn(255)),
		uint8(rng.Intn(255)),
		uint8(rng.Intn(255)),
		255,
	}

	offset := func() float64 {
		size := 100.0
		return rng.Float64()*size - (size / 2.0)
	}

	p1 := util.Vector{X: rng.Float64() * bounds.X, Y: rng.Float64() * bounds.Y}
	p2 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}
	p3 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}

	vrt := []util.Vector{p1, p2, p3}

	return Triangle{Bounds: bounds, Color: clr, Vertices: vrt}
}

func (t *Triangle) Draw(dc *gg.Context, offset util.Vector) {
	dc.NewSubPath()
	dc.MoveTo(t.Vertices[0].X+offset.X, t.Vertices[0].Y+offset.Y)
	dc.LineTo(t.Vertices[1].X+offset.X, t.Vertices[1].Y+offset.Y)
	dc.LineTo(t.Vertices[2].X+offset.X, t.Vertices[2].Y+offset.Y)
	dc.ClosePath()

	dc.SetColor(t.Color)
	dc.Fill()
}

// copy all the data without pointers
// TODO figure out how deep this actually needs to be
func (t *Triangle) Clone() Triangle {
	clone := Triangle{
		Bounds:   t.Bounds,
		Color:    color.RGBA{t.Color.R, t.Color.G, t.Color.B, t.Color.A},
		Vertices: make([]util.Vector, len(t.Vertices)),
	}

	for i, p := range t.Vertices {
		clone.Vertices[i] = util.Vector{X: p.X, Y: p.Y}
	}

	return clone
}
