package worker

import (
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/go/util"
)

type Shape interface {
	Draw(dc *gg.Context, offset util.Vector)
}

type Shapes struct {
	Bounds  util.Vector
	Context *Worker
	Members []Shape
	Type    string
}

func getCreateShapeFunc(shapeType string) func(radius float64, bounds util.Vector, rng *rand.Rand) Shape {
	if shapeType == "triangles" {
		return createTriangle
	}

	if shapeType == "polygons" {
		return createPolygon
	}

	return createCircle
}

func createShapesFactory(ctx *Worker, shapeType string) func(rng *rand.Rand) eaopt.Genome {
	bounds := util.Vector{X: float64(ctx.TargetImage.Width), Y: float64(ctx.TargetImage.Height)}

	createShape := getCreateShapeFunc(shapeType)

	return func(rng *rand.Rand) eaopt.Genome {
		shapes := Shapes{
			Bounds:  bounds,
			Context: ctx,
			Members: make([]Shape, ctx.Job.NumShapes),
			Type:    shapeType,
		}

		for i := 0; i < int(ctx.Job.NumShapes); i++ {
			shapes.Members[i] = createShape(float64(ctx.Job.ShapeSize), bounds, rng)
		}

		return shapes
	}
}

func (s Shapes) Draw(dc *gg.Context, offset util.Vector) {
	for _, m := range s.Members {
		m.Draw(dc, offset)
	}
}

func (s Shapes) Evaluate() (float64, error) {
	dc := gg.NewContext(s.Context.TargetImage.Width, s.Context.TargetImage.Height)
	s.Draw(dc, util.Vector{X: 0, Y: 0})
	out := dc.Image()

	fitness := imgDiff(rgbaImg(out), rgbaImg(s.Context.TargetImage.Image))

	return fitness, nil
}

func (s Shapes) Mutate(rng *rand.Rand) {
	createShape := getCreateShapeFunc(s.Type)

	for i := range s.Members {
		if rng.Float64() < s.Context.Job.MutationRate {
			s.Members[i] = createShape(float64(s.Context.Job.ShapeSize), s.Bounds, rng)
		}
	}
}

func (s Shapes) Crossover(g eaopt.Genome, rng *rand.Rand) {
	o := g.(Shapes)

	for i := range s.Members {
		if rng.Float64() < 0.5 {
			s.Members[i] = o.Members[i]
		} else {
			o.Members[i] = s.Members[i]
		}
	}
}

// copy all the data without pointers
func (s Shapes) Clone() eaopt.Genome {
	return Shapes{
		Bounds:  s.Bounds,
		Context: s.Context,
		Members: append([]Shape{}, s.Members...),
		Type:    s.Type,
	}
}

func (s Shapes) CloneWithoutContext() eaopt.Genome {
	return Shapes{
		Bounds:  s.Bounds,
		Members: append([]Shape{}, s.Members...),
		Type:    s.Type,
	}
}

func (s Shapes) At(i int) interface{} {
	return s.Members[i]
}

func (s Shapes) Set(i int, v interface{}) {
	s.Members[i] = v.(Shape)
}

func (s Shapes) Len() int {
	return len(s.Members)
}

func (s Shapes) Swap(i, j int) {
	s.Members[i], s.Members[j] = s.Members[j], s.Members[i]
}

func (s Shapes) Slice(a, b int) eaopt.Slice {
	slice := s.Clone()
	s.Members = s.Members[a:b]
	return slice.(eaopt.Slice)
}

func (s Shapes) Split(k int) (eaopt.Slice, eaopt.Slice) {
	s1 := s.Clone().(Shapes)
	s1.Members = s.Members[:k]

	s2 := s.Clone().(Shapes)
	s2.Members = s.Members[k:]

	return s1, s2
}

func (s Shapes) Append(q eaopt.Slice) eaopt.Slice {
	new := s.Clone().(Shapes)
	new.Members = append(s.Members, q.(Shapes).Members...)
	return new
}

func (s Shapes) Replace(q eaopt.Slice) {
	copy(s.Members, q.(Shapes).Members)
}

func (s Shapes) Copy() eaopt.Slice {
	return s.Clone().(Shapes)
}
