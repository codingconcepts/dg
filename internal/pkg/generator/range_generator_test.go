package generator

import (
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRangeColumn(t *testing.T) {
	cases := []struct {
		name     string
		files    map[string]model.CSVFile
		rtype    string
		count    int
		from     string
		to       string
		step     string
		format   string
		expLines []string
		expErr   error
	}{
		{
			name: "generates range for existing table",
			files: map[string]model.CSVFile{
				"table": {
					Lines: [][]string{
						{"a"},
						{"a", "b"},
						{"a", "b", "c"},
					},
				},
			},
			rtype:  "date",
			count:  5,
			from:   "2023-01-01",
			to:     "2023-02-01",
			step:   "24h",
			format: "2006-01-02",
			expLines: []string{
				"2023-01-01",
				"2023-01-11",
				"2023-01-21",
			},
		},
		{
			name:   "generates range for count",
			files:  map[string]model.CSVFile{},
			rtype:  "date",
			count:  4,
			from:   "2023-01-01",
			to:     "2023-02-01",
			step:   "24h",
			format: "2006-01-02",
			expLines: []string{
				"2023-01-01",
				"2023-01-08",
				"2023-01-16",
				"2023-01-24",
			},
		},
		{
			name:   "generates range for step",
			files:  map[string]model.CSVFile{},
			rtype:  "date",
			from:   "2023-01-01",
			to:     "2023-02-01",
			step:   "72h",
			format: "2006-01-02",
			expLines: []string{
				"2023-01-01",
				"2023-01-04",
				"2023-01-07",
				"2023-01-10",
				"2023-01-13",
				"2023-01-16",
				"2023-01-19",
				"2023-01-22",
				"2023-01-25",
				"2023-01-28",
				"2023-01-31",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := model.Table{
				Name:  "table",
				Count: c.count,
			}

			column := model.Column{
				Name: "col",
			}

			processor := model.ProcessorRange{
				Type:   c.rtype,
				From:   c.from,
				To:     c.to,
				Step:   c.step,
				Format: c.format,
			}

			files := c.files

			err := GenerateRangeColumn(table, column, processor, files)
			assert.Equal(t, c.expErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.expLines, files["table"].Lines[len(files["table"].Lines)-1])
		})
	}
}
