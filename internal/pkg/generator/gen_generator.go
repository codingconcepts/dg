package generator

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/codingconcepts/dg/internal/pkg/random"
	"github.com/lucasjones/reggen"
	"github.com/martinusso/go-docs/cnpj"
	"github.com/martinusso/go-docs/cpf"
	"github.com/samber/lo"
)

// GenGenerator provides additional context to a gen column.
type GenGenerator struct {
	Value          string `yaml:"value"`
	Pattern        string `yaml:"pattern"`
	NullPercentage int    `yaml:"null_percentage"`
	Format         string `yaml:"format"`
	Template       string `yaml:"template"`

	patternGenerator *reggen.Generator
	templateOptions  gofakeit.TemplateOptions
}

func (g GenGenerator) GetFormat() string {
	return g.Format
}

// Generate random data for a given column.
func (g GenGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if g.Value == "" && g.Pattern == "" && g.Template == "" {
		return fmt.Errorf("gen must have either 'value', 'pattern' or 'template'")
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

	if g.Template != "" {
		var err error
		g.templateOptions = gofakeit.TemplateOptions{
			Funcs: template.FuncMap{
				"cpf":  cpf.Generate,
				"Cpf":  cpf.Generate,
				"CPF":  cpf.Generate,
				"cnpj": cnpj.Generate,
				"Cnpj": cnpj.Generate,
				"CNPJ": cnpj.Generate,
			},
		}
		if _, err = gofakeit.Template(g.Template, &g.templateOptions); err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}
	}

	var lines []string
	for i := 0; i < t.Count; i++ {
		s := g.generate()
		lines = append(lines, s)
	}

	AddTable(t, c.Name, lines, files)
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

	if pg.Template != "" {
		value, err := gofakeit.Template(pg.Template, &pg.templateOptions)
		if err != nil {
			return fmt.Errorf("generating template: %w", err).Error()
		}
		return value
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
