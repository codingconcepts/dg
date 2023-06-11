package generator

import (
	"github.com/codingconcepts/dg/v1/internal/pkg/model"

	"github.com/samber/lo"
)

// GenerateIncColumn generates an incrementing number value for a column.
func GenerateIncColumn(t model.Table, c model.Column, pi model.ProcessorInc, files map[string]model.CSVFile) error {
	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, formatValue(pi, pi.Start+i))
	}

	AddTable(t.Name, c.Name, line, files)
	return nil
}
