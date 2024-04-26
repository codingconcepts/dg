package generator

import (
	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

// IncGenerator provides additional context to an inc column.
type IncGenerator struct {
	Start  int    `yaml:"start"`
	Format string `yaml:"format"`
}

func (pi IncGenerator) GetFormat() string {
	return pi.Format
}

// Generate an incrementing number value for a column.
func (g IncGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, formatValue(g, g.Start+i))
	}

	AddTable(t, c.Name, line, files)
	return nil
}
