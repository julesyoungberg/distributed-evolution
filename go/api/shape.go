package api

import (
	"image/color"
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

func CreateCircle(radius float64, bounds util.Vector, rng *rand.Rand, palette []color.RGBA) Shape {
	return Circle{
		Color:    util.RandomColorFromPalette(rng, palette),
		Position: util.RandomVector(rng, bounds),
		Radius:   util.RandomRadius(rng, radius),
	}
}

func (c Circle) Draw(dc *gg.Context, offset util.Vector) {
	dc.DrawCircle(c.Position.X+offset.X, c.Position.Y+offset.Y, c.Radius)
	dc.SetColor(c.Color)
	dc.Fill()
}

//////////////////////////////////////////////////////////////////////////////
// LINE
//////////////////////////////////////////////////////////////////////////////
type Line struct {
	Color    color.RGBA
	Position []util.Vector
	Width    float64
}

func CreateLine(length float64, bounds util.Vector, rng *rand.Rand, palette []color.RGBA) Shape {
	p1 := util.RandomVector(rng, bounds)
	p2 := util.Vector{X: p1.X + rng.Float64()*length*2, Y: p1.Y + rng.Float64()*length*2}

	return Line{
		Color:    util.RandomColorFromPalette(rng, palette),
		Position: []util.Vector{p1, p2},
		Width:    float64(rng.Intn(8)) * 3.0,
	}
}

func (l Line) Draw(dc *gg.Context, offset util.Vector) {
	dc.DrawLine(l.Position[0].X, l.Position[0].Y, l.Position[1].X, l.Position[1].Y)
	dc.SetColor(l.Color)
	dc.SetLineWidth(l.Width)
	dc.Stroke()
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

func CreatePolygon(radius float64, bounds util.Vector, rng *rand.Rand, palette []color.RGBA) Shape {
	return Polygon{
		Color:    util.RandomColorFromPalette(rng, palette),
		Position: util.RandomVector(rng, bounds),
		Radius:   util.RandomRadius(rng, radius),
		Rotation: util.RandomRotation(rng),
		Sides:    rng.Intn(5) + 3,
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

func CreateTriangle(radius float64, bounds util.Vector, rng *rand.Rand, palette []color.RGBA) Shape {
	offset := func() float64 {
		return rng.Float64()*radius*2.0 - radius
	}

	p1 := util.RandomVector(rng, bounds)
	p2 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}
	p3 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}

	return Triangle{
		Color:    util.RandomColorFromPalette(rng, palette),
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
