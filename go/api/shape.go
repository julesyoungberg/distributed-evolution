package api

import (
	"image"
	"image/color"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Shape interface {
	Draw(dc *gg.Context, offset util.Vector)
}

type ShapeOptions struct {
	Bounds       util.Vector
	MutationRate float64
	Palette      []color.RGBA
	PaletteType  string
	Quantization int
	Size         float64
	TargetImage  image.Image
}

func getShapeColor(rng *rand.Rand, opt ShapeOptions, position util.Vector) color.RGBA {
	if opt.PaletteType == "targetImage" {
		if rng.Float64() < opt.MutationRate {
			return util.RandomColor(rng)
		}

		c := opt.TargetImage.At(int(position.X), int(position.Y))
		r, g, b, _ := c.RGBA()
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8((rng.Intn(3) + 1) * 64)}
	}

	return util.RandomColorFromPalette(rng, opt.Palette)
}

//////////////////////////////////////////////////////////////////////////////
// CIRCLE
//////////////////////////////////////////////////////////////////////////////
type Circle struct {
	Color    color.RGBA
	Position util.Vector
	Radius   float64
}

func CreateCircle(opt ShapeOptions, rng *rand.Rand) Shape {
	position := util.RandomVector(rng, opt.Bounds, opt.Quantization)

	return Circle{
		Color:    getShapeColor(rng, opt, position),
		Position: position,
		Radius:   util.RandomRadius(rng, opt.Size),
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

func CreateLine(opt ShapeOptions, rng *rand.Rand) Shape {
	p1 := util.RandomVector(rng, opt.Bounds, opt.Quantization)
	p2 := util.Vector{
		X: p1.X + rng.Float64()*opt.Size*2,
		Y: p1.Y + rng.Float64()*opt.Size*2,
	}
	position := util.Vector{
		X: (p1.X + p2.X) / 2.0,
		Y: (p1.Y + p2.Y) / 2.0,
	}

	return Line{
		Color:    getShapeColor(rng, opt, position),
		Position: []util.Vector{p1, p2},
		Width:    float64(rng.Intn(16)),
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

func CreatePolygon(opt ShapeOptions, rng *rand.Rand) Shape {
	position := util.RandomVector(rng, opt.Bounds, opt.Quantization)

	return Polygon{
		Color:    getShapeColor(rng, opt, position),
		Position: position,
		Radius:   util.RandomRadius(rng, opt.Size),
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

func CreateTriangle(opt ShapeOptions, rng *rand.Rand) Shape {
	offset := func() float64 {
		return rng.Float64()*opt.Size*2.0 - opt.Size
	}

	p1 := util.RandomVector(rng, opt.Bounds, opt.Quantization)
	p2 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}
	p3 := util.Vector{X: p1.X + offset(), Y: p1.Y + offset()}
	position := util.Vector{
		X: (p1.X + p2.X + p3.X) / 3.0,
		Y: (p1.Y + p2.Y + p3.Y) / 3.0,
	}

	return Triangle{
		Color:    getShapeColor(rng, opt, position),
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
