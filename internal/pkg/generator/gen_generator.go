package generator

import (
	"dg/internal/pkg/model"
	"fmt"
	"math/rand"
	"strings"

	"github.com/samber/lo"
)

// GenerateGenColumn generates random data for a given column.
func GenerateGenColumn(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	var pg model.ProcessorGenerator
	if err := c.Processor.UnmarshalFunc(&pg); err != nil {
		return fmt.Errorf("parsing each process for %s: %w", c.Name, err)
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, replacePlaceholders(pg))
	}

	addToFile(t.Name, c.Name, line, files)
	return nil
}

func replacePlaceholders(pg model.ProcessorGenerator) string {
	r := rand.Intn(100)
	if r < pg.NullPercentage {
		return ""
	}

	s := pg.Value
	for k, v := range replacements {
		if strings.Contains(s, k) {
			valueStr := formatValue(pg, v())
			s = strings.ReplaceAll(s, k, valueStr)
		}
	}

	return s
}
