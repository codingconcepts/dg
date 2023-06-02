package generator

import (
	"github.com/samber/lo"
)

// CartesianProduct returns the Cartesian product of a variable number of arrays.
func CartesianProduct(a ...[]string) [][]string {
	if len(a) == 0 {
		return [][]string{}
	}

	totalCombinations := lo.Reduce(a, func(agg int, item []string, index int) int {
		return agg * len(item)
	}, 1)

	// Preallocate the result slice with the correct capacity.
	result := make([][]string, 0, totalCombinations)
	result = append(result, []string{})

	// Generate the Cartesian products.
	for _, arr := range a {
		temp := make([][]string, 0, totalCombinations)
		for _, element := range arr {
			for _, combination := range result {
				// Create a new combination by appending the current element
				newCombination := make([]string, len(combination)+1)
				copy(newCombination, combination)
				newCombination[len(combination)] = element
				temp = append(temp, newCombination)
			}
		}
		result = temp
	}

	return result
}

// Transpose a multi-dimensional array.
func Transpose(m [][]string) [][]string {
	r := make([][]string, len(m[0]))

	for x := range r {
		r[x] = make([]string, len(m))
	}

	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}
