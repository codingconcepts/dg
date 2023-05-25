package generator

import (
	"dg/internal/pkg/model"
	"fmt"
	"math/rand"
)

// GenerateSetColumn selects between a set of values for a given table.
func GenerateSetColumn(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	var ps model.ProcessorSet
	if err := c.Processor.UnmarshalFunc(&ps); err != nil {
		return fmt.Errorf("parsing set process for %s.%s: %w", t.Name, c.Name, err)
	}

	if len(ps.Values) == 0 {
		return fmt.Errorf("no values provided for set generator")
	}

	var line []string

	if len(ps.Weights) > 0 {
		items, err := buildWeightedItems(ps)
		if err != nil {
			return fmt.Errorf("making weighted items collection: %w", err)
		}

		for i := 0; i < t.Count; i++ {
			line = append(line, items.choose())
		}
	} else {
		for i := 0; i < t.Count; i++ {
			line = append(line, ps.Values[rand.Intn(len(ps.Values))])
		}
	}

	addToFile(t.Name, c.Name, line, files)

	return nil
}

func buildWeightedItems(ps model.ProcessorSet) (weightedItems, error) {
	if len(ps.Values) != len(ps.Weights) {
		return weightedItems{}, fmt.Errorf("set values and weights need to be the same")
	}

	weightedItems := make([]weightedItem, len(ps.Values))
	for i, v := range ps.Values {
		weightedItems = append(weightedItems, weightedItem{
			Value:  v,
			Weight: ps.Weights[i],
		})
	}

	return makeWeightedItems(weightedItems), nil
}
