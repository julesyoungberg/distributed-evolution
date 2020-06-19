package worker

import (
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Shapes struct {
	Context *Worker
	Members []Triangle
}

func createTrianglesFactory(ctx *Worker) func(rng *rand.Rand) eaopt.Genome {
	bounds := util.Vector{X: float64(ctx.TargetImage.Width), Y: float64(ctx.TargetImage.Height)}

	return func(rng *rand.Rand) eaopt.Genome {
		shapes := Shapes{
			Context: ctx,
			Members: make([]Triangle, ctx.NumShapes),
		}

		for i := 0; i < ctx.NumShapes; i++ {
			shapes.Members[i] = createTriangle(ctx.ShapeSize, bounds, rng)
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
	for i := range s.Members {
		if rng.Float64() < s.Context.MutationRate {
			s.Members[i] = createTriangle(s.Context.ShapeSize, s.Members[i].Bounds, rng)
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
	clone := Shapes{
		Context: s.Context,
		Members: make([]Triangle, len(s.Members)),
	}

	for i, m := range s.Members {
		clone.Members[i] = m.Clone()
	}

	return clone
}

func (s Shapes) CloneWithoutContext() eaopt.Genome {
	clone := Shapes{
		Members: make([]Triangle, len(s.Members)),
	}

	for i, m := range s.Members {
		clone.Members[i] = m.Clone()
	}

	return clone
}

func (s Shapes) At(i int) interface{} {
	return s.Members[i]
}

func (s Shapes) Set(i int, v interface{}) {
	s.Members[i] = v.(Triangle)
}

func (s Shapes) Len() int {
	return len(s.Members)
}

func (s Shapes) Swap(i, j int) {
	s.Members[i], s.Members[j] = s.Members[j], s.Members[i]
}

func (s Shapes) Slice(a, b int) eaopt.Slice {
	return Shapes{
		Context: s.Context,
		Members: s.Members[a:b],
	}
}

func (s Shapes) Split(k int) (eaopt.Slice, eaopt.Slice) {
	s1 := Shapes{
		Context: s.Context,
		Members: s.Members[:k],
	}

	s2 := Shapes{
		Context: s.Context,
		Members: s.Members[k:],
	}

	return s1, s2
}

func (s Shapes) Append(q eaopt.Slice) eaopt.Slice {
	return Shapes{
		Context: s.Context,
		Members: append(s.Members, q.(Shapes).Members...),
	}
}

func (s Shapes) Replace(q eaopt.Slice) {
	copy(s.Members, q.(Shapes).Members)
}

func (s Shapes) Copy() eaopt.Slice {
	clone := Shapes{Context: s.Context}
	copy(clone.Members, s.Members)
	return clone
}
