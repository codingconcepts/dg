package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

// GenerateConstColumn generates values for a column based on a series of provided values.
func GenerateConstColumn(t model.Table, c model.Column, pc model.ProcessorConst, files map[string]model.CSVFile) error {
	if len(pc.Values) == 0 {
		return fmt.Errorf("no values provided for const generator")
	}

	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	if count > len(pc.Values) {
		return fmt.Errorf("wrong number of values provided for const generator (need %d, got %d)", count, len(pc.Values))
	}

	var line []string
	for _, value := range pc.Values {
		line = append(line, value)
	}

	AddTable(t, c.Name, line, files)
	return nil
}
