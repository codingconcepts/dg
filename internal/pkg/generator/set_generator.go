package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/codingconcepts/dg/internal/pkg/random"
	"github.com/samber/lo"
)

// SetGenerator provides additional context to a set column.
type SetGenerator struct {
	Values  []string `yaml:"values"`
	Weights []int    `yaml:"weights"`
}

// Generate selects between a set of values for a given table.
func (ps SetGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if len(ps.Values) == 0 {
		return fmt.Errorf("no values provided for set generator")
	}

	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	var line []string
	if len(ps.Weights) > 0 {
		items, err := ps.buildWeightedItems()
		if err != nil {
			return fmt.Errorf("making weighted items collection: %w", err)
		}

		for i := 0; i < count; i++ {
			line = append(line, items.choose())
		}
	} else {
		for i := 0; i < count; i++ {
			line = append(line, ps.Values[random.Intn(len(ps.Values))])
		}
	}

	AddTable(t, c.Name, line, files)
	return nil
}

func (ps SetGenerator) buildWeightedItems() (weightedItems, error) {
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
