package worker

type Vector struct {
	X float64
	Y float64
}

// https://stackoverflow.com/questions/2049582/how-to-determine-if-a-point-is-in-a-2d-triangle
func sign(p1 Vector, p2 Vector, p3 Vector) float64 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

// checks if pt is contained within the triangle formed by points v1, v2, and v3
// https://stackoverflow.com/questions/2049582/how-to-determine-if-a-point-is-in-a-2d-triangle
func pointInTriangle(pt Vector, v1 Vector, v2 Vector, v3 Vector) bool {
	d1 := sign(pt, v1, v2)
	d2 := sign(pt, v2, v3)
	d3 := sign(pt, v3, v1)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}

func clampFloat64(n, min, max float64) float64 {
	if n < min {
		return min
	}

	if n > max {
		return max
	}

	return n
}