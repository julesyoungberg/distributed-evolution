package worker

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Shape interface {
	Draw(dc *gg.Context, offset util.Vector)
}

//////////////////////////////////////////////////////////////////////////////
// CIRCLE
//////////////////////////////////////////////////////////////////////////////
type Circle struct {
	Color    color.RGBA
	Position util.Vector
	Radius   float64
}

func createCircle(radius float64, bounds util.Vector, rng *rand.Rand) Shape {
	return Circle{
		Color:    randomColor(rng),
		Position: randomVector(rng, bounds),
		Radius:   (rng.Float64()*radius)/2.0 + radius/2.0, // from radius/2 to radius
	}
}

func (c Circle) Draw(dc *gg.Context, offset util.Vector) {
	dc.DrawCircle(c.Position.X+offset.X, c.Position.Y+offset.Y, c.Radius)
	dc.SetColor(c.Color)
	dc.Fill()
}

//////////////////////////////////////////////////////////////////////////////
// POLYGON
//////////////////////////////////////////////////////////////////////////////
type Polygon struct {
	Color    color.RGBA
	Position util.Vector
	Radius   float64
	Rotation float64
	Sides    int
}

func createPolygon(radius float64, bounds util.Vector, rng *rand.Rand) Shape {
	return Polygon{
		Color:    randomColor(rng),
		Position: randomVector(rng, bounds),
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

//////////////////////////////////////////////////////////////////////////////
// TRIANGLE
//////////////////////////////////////////////////////////////////////////////
type Triangle struct {
	Color    color.RGBA
	Vertices []util.Vector
}

func createTriangle(radius float64, bounds util.Vector, rng *rand.Rand) Shape {
	offset := func() float64 {
		return rng.Float64()*radius*2.0 - radius
	}

	p1 := randomVector(rng, bounds)
	p2 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}
	p3 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}

	return Triangle{
		Color:    randomColor(rng),
		Vertices: []util.Vector{p1, p2, p3},
	}
}

func (t Triangle) Draw(dc *gg.Context, offset util.Vector) {
	dc.NewSubPath()
	dc.MoveTo(t.Vertices[0].X+offset.X, t.Vertices[0].Y+offset.Y)
	dc.LineTo(t.Vertices[1].X+offset.X, t.Vertices[1].Y+offset.Y)
	dc.LineTo(t.Vertices[2].X+offset.X, t.Vertices[2].Y+offset.Y)
	dc.ClosePath()

	dc.SetColor(t.Color)
	dc.Fill()
}
