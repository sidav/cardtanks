package calc

import "math/rand"

// May be extremely unefficient, depending on weightf imlementation.
func WeightedRand(min, max int, weightf func(index int) int) int {
	if min == max {
		return max
	}
	totalWeights := 0
	for i := min; i <= max; i++ {
		totalWeights += weightf(i)
	}
	roll := rand.Intn(totalWeights)
	for i := min; i <= max; i++ {
		wght := weightf(i)
		if wght == 0 {
			continue
		}
		if roll < wght {
			return i
		}
		roll -= wght
	}
	panic("weightedRand failed.")
}

