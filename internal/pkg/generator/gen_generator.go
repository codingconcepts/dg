package generator

import (
	"strings"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/codingconcepts/dg/internal/pkg/random"
	"github.com/samber/lo"
)

// GenGenerator provides additional context to a gen column.
type GenGenerator struct {
	Value          string `yaml:"value"`
	NullPercentage int    `yaml:"null_percentage"`
	Format         string `yaml:"format"`
}

func (g GenGenerator) GetFormat() string {
	return g.Format
}

// Generate generates random data for a given column.
func (g GenGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, g.replacePlaceholders())
	}

	AddTable(t, c.Name, line, files)
	return nil
}

func (pg GenGenerator) replacePlaceholders() string {
	r := random.Intn(100)
	if r < pg.NullPercentage {
		return ""
	}

	s := pg.Value

	// Look for quick single-replacements.
	if v, ok := replacements[s]; ok {
		return formatValue(pg, v())
	}

	// Process multipe-replacements.
	for k, v := range replacements {
		if strings.Contains(s, k) {
			valueStr := formatValue(pg, v())
			s = strings.ReplaceAll(s, k, valueStr)
		}
	}

	return s
}
