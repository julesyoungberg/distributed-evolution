package worker

import (
	"math"
	"math/rand"

	"github.com/MaxHalford/eaopt"
	"github.com/fogleman/gg"
	"github.com/rickyfitts/distributed-evolution/util"
)

type Triangles struct {
	Context *Worker
	Members []Triangle
}

// creates a triangles factory with a pointer to the worker object
func createTrianglesFactory(ctx *Worker) func(rng *rand.Rand) eaopt.Genome {
	bounds := util.Vector{X: float64(ctx.TargetImage.Width), Y: float64(ctx.TargetImage.Height)}

	return func(rng *rand.Rand) eaopt.Genome {
		nTriangles := 5

		triangles := Triangles{
			Context: ctx,
			Members: make([]Triangle, nTriangles),
		}

		for i := 0; i < nTriangles; i++ {
			triangles.Members[i] = createTriangle(rng, bounds)
		}

		return triangles
	}
}

func (t Triangles) Draw(dc *gg.Context, offset util.Vector) {
	for _, triangle := range t.Members {
		triangle.Draw(dc, offset)
	}
}

// determine how well the triangles match the pixels of the target image
func (t Triangles) Evaluate() (float64, error) {
	// draw the triangles
	dc := gg.NewContext(t.Context.TargetImage.Width, t.Context.TargetImage.Height)
	t.Draw(dc, util.Vector{X: 0, Y: 0})
	out := dc.Image()

	total := 0
	var fitness float64 = 0.0

	// calculates the fitness between two integers as an inverse ratio
	calcFit := func(a, b uint32) float64 {
		diff := math.Abs(float64(a - b))
		ratio := diff / 255.0
		return 1.0 - ratio
	}

	for y := 0; y < t.Context.TargetImage.Height; y++ {
		for x := 0; x < t.Context.TargetImage.Width; x++ {
			r1, g1, b1, _ := out.At(x, y).RGBA()
			r2, g2, b2, _ := t.Context.TargetImage.Image.At(x, y).RGBA()

			dr := calcFit(r1, r2)
			dg := calcFit(g1, g2)
			db := calcFit(b1, b2)

			fitness += (dr + dg + db) / 3.0
			total++
		}
	}

	return fitness, nil
}

func (t Triangles) Mutate(rng *rand.Rand) {
	for _, m := range t.Members {
		m.Mutate(rng, t.Context.MutationRate)
	}
}

// mix genes with another set of triangles
func (t Triangles) Crossover(g eaopt.Genome, rng *rand.Rand) {
	o := g.(Triangles)

	for i, m := range t.Members {
		if rng.Float32() < 0.5 {
			t.Members[i], o.Members[i] = o.Members[i], m
		}
	}
}

// copy all the data without pointers
func (t Triangles) Clone() eaopt.Genome {
	clone := Triangles{
		Context: t.Context,
		Members: make([]Triangle, len(t.Members)),
	}

	for i, m := range t.Members {
		clone.Members[i] = m.Clone()
	}

	return clone
}

func (t Triangles) CloneWithoutContext() eaopt.Genome {
	clone := Triangles{
		Members: make([]Triangle, len(t.Members)),
	}

	for i, m := range t.Members {
		clone.Members[i] = m.Clone()
	}

	return clone
}

func (t Triangles) At(i int) interface{} {
	return t.Members[i]
}

func (t Triangles) Set(i int, v interface{}) {
	t.Members[i] = v.(Triangle)
}

func (t Triangles) Len() int {
	return len(t.Members)
}

func (t Triangles) Swap(i, j int) {
	t.Members[i], t.Members[j] = t.Members[j], t.Members[i]
}

func (t Triangles) Slice(a, b int) eaopt.Slice {
	return Triangles{
		Context: t.Context,
		Members: t.Members[a:b],
	}
}

func (t Triangles) Split(k int) (eaopt.Slice, eaopt.Slice) {
	s1 := Triangles{
		Context: t.Context,
		Members: t.Members[:k],
	}

	s2 := Triangles{
		Context: t.Context,
		Members: t.Members[k:],
	}

	return s1, s2
}

func (t Triangles) Append(q eaopt.Slice) eaopt.Slice {
	return Triangles{
		Context: t.Context,
		Members: append(t.Members, q.(Triangles).Members...),
	}
}

func (t Triangles) Replace(q eaopt.Slice) {
	copy(t.Members, q.(Triangles).Members)
}

func (t Triangles) Copy() eaopt.Slice {
	clone := Triangles{Context: t.Context}
	copy(clone.Members, t.Members)
	return clone
}
