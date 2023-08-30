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
			name: "generates date range for existing table",
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
			name:   "generates date range for count",
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
			name:   "generates date range for step",
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
		{
			name: "generates int range for existing table",
			files: map[string]model.CSVFile{
				"table": {
					Lines: [][]string{
						{"a"},
						{"a", "b"},
						{"a", "b", "c"},
					},
				},
			},
			rtype: "int",
			count: 5,
			from:  "1",
			to:    "5",
			expLines: []string{
				"1",
				"3",
				"5",
			},
		},
		{
			name:  "generates int range for count",
			files: map[string]model.CSVFile{},
			rtype: "int",
			count: 4,
			from:  "10",
			to:    "40",
			step:  "10",
			expLines: []string{
				"10",
				"20",
				"30",
				"40",
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

			g := RangeGenerator{
				Type:   c.rtype,
				From:   c.from,
				To:     c.to,
				Step:   c.step,
				Format: c.format,
			}

			files := c.files

			err := g.Generate(table, column, files)
			assert.Equal(t, c.expErr, err)

			if err != nil {
				return
			}

			assert.Equal(t, c.expLines, files["table"].Lines[len(files["table"].Lines)-1])
		})
	}
}

func TestGenerateDateSlice(t *testing.T) {
	cases := []struct {
		name     string
		from     string
		to       string
		format   string
		count    int
		step     string
		expSlice []string
		expError string
	}{
		{
			name:     "no count or step",
			expError: "either a count or a step must be provided to a date range generator",
		},
		{
			name:   "count",
			count:  10,
			from:   "2023-01-01",
			to:     "2023-01-10",
			format: "2006-01-02",
			expSlice: []string{
				"2023-01-01", "2023-01-01", "2023-01-02", "2023-01-03", "2023-01-04", "2023-01-05", "2023-01-06", "2023-01-07", "2023-01-08", "2023-01-09",
			},
		},
		{
			name:   "step",
			step:   "24h",
			from:   "2023-01-10",
			to:     "2023-01-20",
			format: "2006-01-02",
			expSlice: []string{
				"2023-01-10", "2023-01-11", "2023-01-12", "2023-01-13", "2023-01-14", "2023-01-15", "2023-01-16", "2023-01-17", "2023-01-18", "2023-01-19",
			},
		},
		{
			name:     "invalid format",
			count:    10,
			from:     "2023-01-01",
			to:       "2023-01-10",
			format:   "abc",
			expError: `parsing from date: parsing time "2023-01-01" as "abc": cannot parse "2023-01-01" as "abc"`,
		},
		{
			name:   "invalid from date",
			count:  10,
			from:   "abc",
			format: "2006-01-02",

			to:       "2023-01-10",
			expError: `parsing from date: parsing time "abc" as "2006-01-02": cannot parse "abc" as "2006"`,
		},
		{
			name:     "invalid to date",
			count:    10,
			from:     "2023-01-01",
			to:       "abc",
			format:   "2006-01-02",
			expError: `parsing to date: parsing time "abc" as "2006-01-02": cannot parse "abc" as "2006"`,
		},
		{
			name:     "invalid step",
			step:     "abc",
			from:     "2023-01-01",
			to:       "2023-01-10",
			format:   "2006-01-02",
			expError: `parsing step: time: invalid duration "abc"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := RangeGenerator{
				From:   c.from,
				To:     c.to,
				Format: c.format,
				Step:   c.step,
			}

			actSlice, actErr := g.generateDateSlice(c.count)
			if c.expError != "" {
				assert.Equal(t, c.expError, actErr.Error())
				return
			}

			assert.Equal(t, c.expSlice, actSlice)
		})
	}
}
