package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/codingconcepts/dg/internal/pkg/random"

	"github.com/samber/lo"
)

// RefGenerator provides additional context to a ref column.
type RefGenerator struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

// Generate looks to previously generated table data and references that when generating data
// for the given table.
func (g RefGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	table, ok := files[g.Table]
	if !ok {
		return fmt.Errorf("missing table %q for ref lookup", g.Table)
	}

	colIndex := lo.IndexOf(table.Header, g.Column)
	column := table.Lines[colIndex]

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, column[random.Intn(len(column))])
	}

	AddTable(t, c.Name, line, files)
	return nil
}
