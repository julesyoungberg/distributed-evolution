package worker

import (
	"image/color"
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"github.com/fogleman/gg"

	"github.com/rickyfitts/distributed-evolution/util"
)

type Triangle struct {
	Bounds   util.Vector
	Color    color.RGBA
	Vertices [][]float64
}

// creates a random triangle
func createTriangle(rng *rand.Rand, bounds util.Vector) Triangle {
	clr := color.RGBA{
		uint8(rng.Intn(255)),
		uint8(rng.Intn(255)),
		uint8(rng.Intn(255)),
		255,
	}

	vrt := [][]float64{
		{rng.Float64() * bounds.X, rng.Float64() * bounds.Y},
		{rng.Float64() * bounds.X, rng.Float64() * bounds.Y},
		{rng.Float64() * bounds.X, rng.Float64() * bounds.Y},
	}

	return Triangle{Bounds: bounds, Color: clr, Vertices: vrt}
}

func (t *Triangle) Draw(dc *gg.Context, offset util.Vector) {
	// draw triangle
	dc.NewSubPath()
	dc.MoveTo(t.Vertices[0][0]+offset.X, t.Vertices[0][1]+offset.Y)
	dc.LineTo(t.Vertices[1][0]+offset.X, t.Vertices[1][1]+offset.Y)
	dc.LineTo(t.Vertices[2][0]+offset.X, t.Vertices[2][1]+offset.Y)
	dc.ClosePath()

	dc.SetColor(t.Color)
	dc.Fill()
}

// clamps the triangles position by the width and height of the target image
func (t Triangle) clampPosition() {
	for _, v := range t.Vertices {
		v[0] = util.ClampFloat64(v[0], 0.0, t.Bounds.X)
		v[1] = util.ClampFloat64(v[0], 0.0, t.Bounds.Y)
	}
}

// mutate the properties of the triangle based on the mutation rate
func (t Triangle) Mutate(rng *rand.Rand, mutationRate float64) {
	c := []uint8{t.Color.R, t.Color.G, t.Color.B}

	for i, x := range c {
		if rng.Float64() < mutationRate {
			c[i] = uint8(rng.NormFloat64() * float64(x))
		}
	}

	t.Color.R = c[0]
	t.Color.G = c[1]
	t.Color.B = c[2]

	for _, p := range t.Vertices {
		eaopt.MutNormalFloat64(p, 0.8, rng)
	}

	t.clampPosition()
}

// copy all the data without pointers
// TODO figure out how deep this actually needs to be
func (t *Triangle) Clone() Triangle {
	var clone Triangle
	clone.Color = color.RGBA{t.Color.R, t.Color.G, t.Color.B, t.Color.A}

	for _, p := range t.Vertices {
		clone.Vertices = append(clone.Vertices, []float64{p[0], p[1]})
	}

	return clone
}
