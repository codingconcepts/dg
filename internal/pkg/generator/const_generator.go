package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

// ConstGenerator provides additional context to a const column.
type ConstGenerator struct {
	Values []string `yaml:"values"`
}

// Generate generates values for a column based on a series of provided values.
func (g ConstGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if len(g.Values) == 0 {
		return fmt.Errorf("no values provided for const generator")
	}

	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	if count > len(g.Values) {
		return fmt.Errorf("wrong number of values provided for const generator (need %d, got %d)", count, len(g.Values))
	}

	var line []string
	for _, value := range g.Values {
		line = append(line, value)
	}

	AddTable(t, c.Name, line, files)
	return nil
}
