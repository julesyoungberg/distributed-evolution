package main

import (
	"image/color"
	"math/rand"

	"github.com/MaxHalford/eaopt"
)

type Vector struct {
	x float64
	y float64
}

type Triangle struct {
	context  *Worker
	Color    color.RGBA
	Vertices []Vector
}

// creates a triangle factory with a pointer to the worker object
func createTriangleFactory(ctx *Worker) func(rng *rand.Rand) eaopt.Genome {
	upper := Vector{float64(ctx.targetImage.width), float64(ctx.targetImage.height)}

	return func(rng *rand.Rand) eaopt.Genome {
		clr := color.RGBA{
			uint8(rng.Intn(255)),
			uint8(rng.Intn(255)),
			uint8(rng.Intn(255)),
			255,
		}

		pos := []Vector{
			Vector{rng.Float64() * upper.x, rng.Float64() * upper.y},
			Vector{rng.Float64() * upper.x, rng.Float64() * upper.y},
			Vector{rng.Float64() * upper.x, rng.Float64() * upper.y},
		}

		return Triangle{
			context:  ctx,
			Color:    clr,
			Vertices: pos,
		}
	}
}

// determine how well the triangle matches with the containing pixels of the target image
func (t Triangle) Evaluate() (float64, error) {
	totalInside := 0
	var fitness float64 = 0.0

	for y := 0; y < t.context.targetImage.height; y++ {
		for x := 0; x < t.context.targetImage.width; x++ {
			if pointInTriangle(Vector{float64(x), float64(y)}, t.Vertices[0], t.Vertices[1], t.Vertices[2]) {
				totalInside++

				r, g, b, _ := t.context.targetImage.image.At(x, y).RGBA()

				dr := uint32(t.Color.R) / r
				dg := uint32(t.Color.G) / g
				db := uint32(t.Color.B) / b

				fitness += float64(dr+dg+db) / 3.0
			}
		}
	}

	fitness /= float64(totalInside)

	return fitness, nil
}

func (t Triangle) Mutate(rng *rand.Rand) {
	c := []uint8{t.Color.R, t.Color.G, t.Color.B}

	for i, x := range c {
		if rng.Float64() < t.context.mutationRate {
			c[i] = uint8(rng.NormFloat64() * float64(x))
		}
	}

	t.Color.R = c[0]
	t.Color.G = c[1]
	t.Color.B = c[2]

	for i, p := range t.Vertices {
		q := []float64{p.x, p.y}
		eaopt.MutNormalFloat64(q, 0.8, rng)
		t.Vertices[i].x = q[0]
		t.Vertices[i].y = q[1]
	}
}

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
		p := []float64{t.Vertices[i].x, t.Vertices[i].y}
		q := []float64{o.Vertices[i].x, o.Vertices[i].y}

		eaopt.CrossUniformFloat64(p, q, rng)

		t.Vertices[i].x = p[0]
		t.Vertices[i].y = p[1]

		o.Vertices[i].x = q[0]
		o.Vertices[i].y = q[1]
	}
}

func (t Triangle) Clone() eaopt.Genome {
	clone := Triangle{context: t.context}
	clone.Color = color.RGBA{t.Color.R, t.Color.G, t.Color.B, t.Color.A}

	for _, p := range t.Vertices {
		clone.Vertices = append(clone.Vertices, Vector{p.x, p.y})
	}

	return clone
}
