package generator

import (
	"testing"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorRelativeDateColumn(t *testing.T) {
	now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	cases := []struct {
		name       string
		unit       string
		after      int
		before     int
		format     string
		date       string
		exp_before time.Time
		exp_after  time.Time
	}{
		{
			name:       "adding days",
			unit:       day,
			after:      -3,
			before:     3,
			exp_before: now.AddDate(0, 4, 0),
			exp_after:  now.AddDate(0, -4, 0),
		},
		{
			name:       "adding months",
			unit:       month,
			after:      -3,
			before:     3,
			exp_before: now.AddDate(0, 4, 0),
			exp_after:  now.AddDate(0, -4, 0),
		},
		{
			name:       "adding years",
			unit:       year,
			after:      -3,
			before:     3,
			exp_before: now.AddDate(4, 0, 0),
			exp_after:  now.AddDate(-4, 0, 0),
		},
		{
			name:       "from now",
			unit:       day,
			after:      -3,
			before:     3,
			date:       "now",
			exp_before: now.AddDate(0, 0, 4),
			exp_after:  now.AddDate(0, 0, -4),
		},
		{
			name:       "from date string",
			unit:       day,
			after:      -3,
			before:     3,
			date:       "25/12/2020",
			format:     "02/01/2006",
			exp_before: time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
			exp_after:  time.Date(2020, 12, 21, 0, 0, 0, 0, time.Local),
		},
		{
			name:       "from another date column",
			unit:       day,
			after:      -3,
			before:     3,
			date:       "order_date",
			exp_before: time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
			exp_after:  time.Date(2020, 12, 21, 0, 0, 0, 0, time.Local),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := model.Table{
				Name:    "table",
				Columns: []model.Column{{Name: "id"}, {Name: "order_date"}},
				Count:   1,
			}
			column := model.Column{
				Name: "due_date",
			}
			files := map[string]model.CSVFile{
				"table": {
					Name:   "table",
					Header: []string{"id", "order_date"},
					Lines: [][]string{
						{"1", "2"},
						{"2020-12-25", "2024-12-25"},
					},
				},
			}
			g := RelDateGenerator{
				Date:   c.date,
				Format: c.format,
				After:  c.after,
				Before: c.before,
				Unit:   c.unit,
			}
			err := g.Generate(table, column, files)
			assert.Nil(t, err)
			last_line, ok := lo.Last(files["table"].Lines)
			assert.True(t, ok)
			last, ok := lo.Last(last_line)
			assert.True(t, ok)
			if c.format == "" {
				c.format = "2006-01-02"
			}
			due_date, err := time.Parse(c.format, last)
			assert.Nil(t, err)
			assert.NotNil(t, due_date)
			assert.True(t, due_date.Before(c.exp_before))
			assert.True(t, due_date.After(c.exp_after))
		})
	}
}
