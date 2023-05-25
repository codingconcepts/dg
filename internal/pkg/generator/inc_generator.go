package generator

import (
	"dg/internal/pkg/model"
	"fmt"

	"github.com/samber/lo"
)

// GenerateIncColumn generates an incrementing number value for a column.
func GenerateIncColumn(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	var pi model.ProcessorInc
	if err := c.Processor.UnmarshalFunc(&pi); err != nil {
		return fmt.Errorf("parsing each process for %s: %w", c.Name, err)
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, formatValue(pi, pi.Start+i))
	}

	addToFile(t.Name, c.Name, line, files)
	return nil
}
