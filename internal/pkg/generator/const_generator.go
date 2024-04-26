package generator

import (
	"fmt"
	"sort"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

// ConstGenerator provides additional context to a const column.
type ConstGenerator struct {
	Values []string `yaml:"values"`
}

// Generate values for a column based on a series of provided values.
func (g ConstGenerator) Generate(t model.Table, files map[string]model.CSVFile) error {
	cols := lo.Filter(t.Columns, func(c model.Column, _ int) bool {
		return c.Type == "const"
	})

	sortColumns(cols)

	for _, c := range cols {
		var cg ConstGenerator
		if err := c.Generator.UnmarshalFunc(&cg); err != nil {
			return fmt.Errorf("parsing const process for %s.%s: %w", t.Name, c.Name, err)
		}
		if err := cg.generate(t, c, files); err != nil {
			return fmt.Errorf("generating const columns: %w", err)
		}
	}

	return nil
}

func sortColumns(cols []model.Column) {
	sort.Slice(cols, func(i, j int) bool {
		var g1 ConstGenerator
		if err := cols[i].Generator.UnmarshalFunc(&g1); err != nil {
			return false
		}

		var g2 ConstGenerator
		if err := cols[j].Generator.UnmarshalFunc(&g2); err != nil {
			return false
		}

		return len(g1.Values) > len(g2.Values)
	})
}

func (g ConstGenerator) generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if len(g.Values) == 0 {
		return fmt.Errorf("no values provided for const generator")
	}

	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	// Repeat the values until they equal the count.
	if count > len(g.Values) {
		for i := 0; len(g.Values) < count; i++ {
			g.Values = append(g.Values, g.Values[i%len(g.Values)])
		}
	}

	AddTable(t, c.Name, g.Values, files)
	return nil
}
