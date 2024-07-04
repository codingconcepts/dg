package generator

import (
	"fmt"
	"testing"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorRelativeDateColumn(t *testing.T) {
	cases := []struct {
		name   string
		unit   string
		low    int
		high   int
		format string
		date   string
	}{
		{
			name: "adding days",
			unit: day,
			low:  -3,
			high: 3,
		},
		{
			name: "adding months",
			unit: month,
			low:  -3,
			high: 3,
		},
		{
			name: "adding years",
			unit: year,
			low:  -3,
			high: 3,
		},
		{
			name:   "from defined date",
			unit:   day,
			low:    -3,
			high:   3,
			date:   "25/12/2020",
			format: "02/01/2006",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := model.Table{
				Name:  "table",
				Count: 10,
			}
			column := model.Column{
				Name: "due_date",
			}

			g := RelativeDateGenerator{
				Date:   c.date,
				Format: c.format,
				Low:    -3,
				High:   3,
				Unit:   c.unit,
			}

			files := map[string]model.CSVFile{}
			err := g.Generate(table, column, files)
			assert.Nil(t, err)
			assert.Equal(t,
				[]string([]string{"due_date"}),
				files["table"].Header,
			)
			ref := time.Now()
			if c.format == "" {
				c.format = "2006-01-02"
			}
			if c.date != "" {
				ref, _ = time.Parse(c.format, c.date)
			}
			for _, line := range files["table"].Lines {
				for _, due_date_str := range line {
					due_date, err := time.Parse(c.format, due_date_str)
					assert.Nil(t, err)
					assert.NotNil(t, due_date)
					var before, after time.Time
					switch c.unit {
					case day:
						before = ref.AddDate(0, 0, c.high+1)
						after = ref.AddDate(0, 0, c.low-1)
					case month:
						before = ref.AddDate(0, c.high+1, 0)
						after = ref.AddDate(0, c.low-1, 0)
					case year:
						before = ref.AddDate(c.high+1, 0, 0)
						after = ref.AddDate(c.low-1, 0, 0)
					}
					assert.True(t, due_date.Before(before), fmt.Sprintf("due_date: %s not before: %s", due_date, before))
					assert.True(t, due_date.After(after), fmt.Sprintf("due_date: %s not after: %s", due_date, after))
				}
			}
		})
	}
}
