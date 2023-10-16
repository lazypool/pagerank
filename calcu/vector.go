package calcu

import "math"

// add a parse vector and a number
func VecAdd(v *map[int]float64, a float64) {
	for idx := range *v {
		(*v)[idx] += a
	}
}

// multiply a parse vector and a number
func VecMult(v *map[int]float64, a float64) {
	for idx := range *v {
		(*v)[idx] *= a
	}
}

// compare 2 vectors and return the difference
func VecComp(v1 *map[int]float64, v2 *map[int]float64) float64 {
	sum := 0.0
	for idx := range *v1 {
		sum += math.Abs((*v1)[idx] - (*v2)[idx])
	}
	return sum
}
