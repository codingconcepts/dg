package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

// RangeGenerator provides additional context to a range column.
type RangeGenerator struct {
	Type   string `yaml:"type"`
	From   string `yaml:"from"`
	To     string `yaml:"to"`
	Step   string `yaml:"step"`
	Format string `yaml:"format"`
}

// Generate generates sequential data between a given start and end range.
func (g RangeGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	switch g.Type {
	case "date":
		lines, err := generateDateSlice(g, count)
		if err != nil {
			return fmt.Errorf("generating date slice: %w", err)
		}

		AddTable(t, c.Name, lines, files)
		return nil
	default:
		return fmt.Errorf("%q is not a valid range type", g.Type)
	}
}
