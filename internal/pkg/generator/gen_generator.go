package generator

import (
	"fmt"
	"strings"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/codingconcepts/dg/internal/pkg/random"
	"github.com/lucasjones/reggen"
	"github.com/samber/lo"
)

// GenGenerator provides additional context to a gen column.
type GenGenerator struct {
	Value          string `yaml:"value"`
	Pattern        string `yaml:"pattern"`
	NullPercentage int    `yaml:"null_percentage"`
	Format         string `yaml:"format"`

	patternGenerator *reggen.Generator
}

func (g GenGenerator) GetFormat() string {
	return g.Format
}

// Generate random data for a given column.
func (g GenGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if g.Value == "" && g.Pattern == "" {
		return fmt.Errorf("gen must have either 'value' or 'pattern'")
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	if g.Pattern != "" {
		var err error
		if g.patternGenerator, err = reggen.NewGenerator(g.Pattern); err != nil {
			return fmt.Errorf("creating regex generator: %w", err)
		}
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		s := g.generate()
		line = append(line, s)
	}

	AddTable(t, c.Name, line, files)
	return nil
}

func (pg GenGenerator) generate() string {
	r := random.Intn(100)
	if r < pg.NullPercentage {
		return ""
	}

	if pg.Pattern != "" {
		return pg.patternGenerator.Generate(255)
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
