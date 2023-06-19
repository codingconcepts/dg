package generator

import (
	"fmt"
	"time"

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
				// Create a new combination by appending the current element.
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
	max := lo.MaxBy(m, func(a, b []string) bool {
		return len(a) > len(b)
	})

	r := make([][]string, len(max))

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

func generateDateSlice(pr RangeGenerator, count int) ([]string, error) {
	// Validate that we have everything we need.
	if count == 0 && pr.Step == "" {
		return nil, fmt.Errorf("either a count or a step must be provided to a date range generator")
	}

	from, err := time.Parse(pr.Format, pr.From)
	if err != nil {
		return nil, fmt.Errorf("parsing from date: %w", err)
	}

	to, err := time.Parse(pr.Format, pr.To)
	if err != nil {
		return nil, fmt.Errorf("parsing to date: %w", err)
	}

	var step time.Duration
	if count > 0 {
		step = to.Sub(from) / time.Duration(count)
	} else {
		if step, err = time.ParseDuration(pr.Step); err != nil {
			return nil, fmt.Errorf("parsing step: %w", err)
		}
	}

	var s []string
	for i := from; i.Before(to); i = i.Add(step) {
		s = append(s, i.Format(pr.Format))
	}

	return s, nil
}
