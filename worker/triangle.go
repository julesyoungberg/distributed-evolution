package worker

import (
	"image/color"
	"math/rand"

	"github.com/MaxHalford/eaopt"
)

type Triangle struct {
	Context  *Worker
	Color    color.RGBA
	Vertices [][]float64
}

// creates a triangle factory with a pointer to the worker object
func createTriangleFactory(ctx *Worker) func(rng *rand.Rand) eaopt.Genome {
	upper := Vector{float64(ctx.TargetImage.Width), float64(ctx.TargetImage.Height)}

	// creates a random triangle, used to generate the initial population
	return func(rng *rand.Rand) eaopt.Genome {
		clr := color.RGBA{
			uint8(rng.Intn(255)),
			uint8(rng.Intn(255)),
			uint8(rng.Intn(255)),
			255,
		}

		vrt := [][]float64{
			[]float64{rng.Float64() * upper.X, rng.Float64() * upper.Y},
			[]float64{rng.Float64() * upper.X, rng.Float64() * upper.Y},
			[]float64{rng.Float64() * upper.X, rng.Float64() * upper.Y},
		}

		return Triangle{
			Context:  ctx,
			Color:    clr,
			Vertices: vrt,
		}
	}
}

// checks if a point is contained in the triangle
func (t Triangle) contains(pt Vector) bool {
	v1 := t.Vertices[0]
	v2 := t.Vertices[1]
	v3 := t.Vertices[2]
	return pointInTriangle(pt, Vector{v1[0], v1[1]}, Vector{v2[0], v2[1]}, Vector{v3[0], v3[1]})
}

// clamps the triangles position by the width and height of the target image
func (t Triangle) clampPosition() {
	for _, v := range t.Vertices {
		v[0] = clampFloat64(v[0], 0.0, float64(t.Context.TargetImage.Width))
		v[1] = clampFloat64(v[0], 0.0, float64(t.Context.TargetImage.Height))
	}
}

// determine how well the triangle matches with the corresponding pixels of the target image
func (t Triangle) Evaluate() (float64, error) {
	totalInside := 0
	var fitness float64 = 0.0

	for y := 0; y < t.Context.TargetImage.Height; y++ {
		for x := 0; x < t.Context.TargetImage.Width; x++ {
			// if the pixel is within the triangle, check how well the color matches
			if t.contains(Vector{float64(x), float64(y)}) {
				totalInside++

				r, g, b, _ := t.Context.TargetImage.Image.At(x, y).RGBA()

				dr := uint32(t.Color.R) / r
				dg := uint32(t.Color.G) / g
				db := uint32(t.Color.B) / b

				fitness += float64(dr+dg+db) / 3.0
			}
		}
	}

	if totalInside > 0 {
		fitness /= float64(totalInside)
	}

	return fitness, nil
}

// mutate the properties of the triangle based on the mutation rate
func (t Triangle) Mutate(rng *rand.Rand) {
	c := []uint8{t.Color.R, t.Color.G, t.Color.B}

	for i, x := range c {
		if rng.Float64() < t.Context.MutationRate {
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

// mix genes with another triangle
func (t Triangle) Crossover(g eaopt.Genome, rng *rand.Rand) {
	o := g.(Triangle)

	c1 := []float64{float64(t.Color.R), float64(t.Color.G), float64(t.Color.B)}
	c2 := []float64{float64(o.Color.R), float64(o.Color.G), float64(o.Color.B)}

	eaopt.CrossUniformFloat64(c1, c2, rng)

	t.Color.R = uint8(c1[0])
	t.Color.G = uint8(c1[1])
	t.Color.B = uint8(c1[2])

	o.Color.R = uint8(c2[0])
	o.Color.G = uint8(c2[1])
	o.Color.B = uint8(c2[2])

	for i := range t.Vertices {
		eaopt.CrossUniformFloat64(t.Vertices[i], o.Vertices[i], rng)
	}

	t.clampPosition()
	o.clampPosition()
}

// copy all the data without pointers
// TODO figure out how deep this actually needs to be
func (t Triangle) Clone() eaopt.Genome {
	clone := Triangle{Context: t.Context}
	clone.Color = color.RGBA{t.Color.R, t.Color.G, t.Color.B, t.Color.A}

	for _, p := range t.Vertices {
		clone.Vertices = append(clone.Vertices, []float64{p[0], p[1]})
	}

	return clone
}
