package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/nrednav/cuid2"
	"github.com/samber/lo"
)

type Cuid2Generator struct {
	Length int `yaml:"length"`
}

func (g Cuid2Generator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if g.Length <= 0 {
		return fmt.Errorf("invalid length provided for cuid2 generator")
	}
	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	generate, err := cuid2.Init(cuid2.WithLength(g.Length))
	if err != nil {
		return fmt.Errorf("failed to initialize cuid2 generator: %w", err)
	}

	var lines []string
	for i := 0; i < count; i++ {
		lines = append(lines, generate())
	}

	AddTable(t, c.Name, lines, files)
	return nil
}
