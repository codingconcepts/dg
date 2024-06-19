package generator

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/lucasjones/reggen"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGenColumn(t *testing.T) {
	cases := []struct {
		name         string
		value        string
		pattern      string
		format       string
		expShapeFunc func(val string) bool
	}{
		{
			name:  "multiple space-delimited strings",
			value: "${first_name} ${last_name}",
			expShapeFunc: func(val string) bool {
				return len(strings.Split(val, " ")) == 2
			},
		},
		{
			name:   "formatted date string",
			value:  "${date}",
			format: "2006-01-02T15:04:05",
			expShapeFunc: func(val string) bool {
				_, err := time.Parse("2006-01-02T15:04:05", val)
				return err == nil
			},
		},
		{
			name:  "integer",
			value: "${int64}",
			expShapeFunc: func(val string) bool {
				_, err := strconv.Atoi(val)
				if err != nil {
					t.Fatal(err)
				}
				return err == nil
			},
		},
		{
			name:    "pattern",
			pattern: `[a-z]{3}-[A-Z]{3}-\d{3}`,
			expShapeFunc: func(val string) bool {
				re := regexp.MustCompile(`[a-z]{3}-[A-Z]{3}-\d{3}`)
				return re.MatchString(val)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := model.Table{
				Name:  "table",
				Count: 1,
			}

			column := model.Column{
				Name: "col",
			}

			g := GenGenerator{
				Value:   c.value,
				Pattern: c.pattern,
				Format:  c.format,
			}

			files := map[string]model.CSVFile{}
			err := g.Generate(table, column, files)
			assert.Nil(t, err)
			assert.True(t, c.expShapeFunc(files["table"].Lines[0][0]))
		})
	}
}

func BenchmarkGeneratePattern(b *testing.B) {
	pattern := `[a-z]{3}-[A-Z]{3}-\d{3}`
	patternGenerator, err := reggen.NewGenerator(pattern)
	assert.NoError(b, err)

	g := GenGenerator{
		Pattern:          pattern,
		patternGenerator: patternGenerator,
	}

	for i := 0; i < b.N; i++ {
		g.generate()
	}
}
