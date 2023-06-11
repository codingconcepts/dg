package generator

import (
	"github.com/codingconcepts/dg/internal/pkg/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddToFile(t *testing.T) {
	cases := []struct {
		name        string
		table       string
		column      string
		line        []string
		filesBefore map[string]model.CSVFile
		filesAfter  map[string]model.CSVFile
	}{
		{
			name:        "first column for table",
			table:       "person",
			column:      "id",
			line:        []string{"a", "b", "c"},
			filesBefore: map[string]model.CSVFile{},
			filesAfter: map[string]model.CSVFile{
				"person": {
					Name:   "person",
					Header: []string{"id"},
					Lines:  [][]string{{"a", "b", "c"}},
				},
			},
		},
		{
			name:   "second column for table",
			table:  "person",
			column: "name",
			line:   []string{"1", "2", "3"},
			filesBefore: map[string]model.CSVFile{
				"person": {
					Name:   "person",
					Header: []string{"id"},
					Lines:  [][]string{{"a", "b", "c"}},
				},
			},
			filesAfter: map[string]model.CSVFile{
				"person": {
					Name:   "person",
					Header: []string{"id", "name"},
					Lines:  [][]string{{"a", "b", "c"}, {"1", "2", "3"}},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			AddTable(c.table, c.column, c.line, c.filesBefore)

			assert.Equal(t, c.filesAfter[c.table].Header, c.filesBefore[c.table].Header)
			assert.Equal(t, c.filesAfter[c.table].Lines, c.filesBefore[c.table].Lines)
			assert.Equal(t, c.filesAfter[c.table].Name, c.filesBefore[c.table].Name)
		})
	}
}

func TestFormatValue(t *testing.T) {
	cases := []struct {
		name   string
		format string
		value  any
		exp    string
	}{
		{
			name:  "no format",
			value: 1,
			exp:   "1",
		},
		{
			name:   "int format",
			value:  1,
			format: "PREFIX_%d_SUFFIX",
			exp:    "PREFIX_1_SUFFIX",
		},
		{
			name:   "time format",
			value:  time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC),
			format: "2006-01-02T15:04:05Z07:00",
			exp:    "2023-01-02T03:04:05Z",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			formatter := model.ProcessorGenerator{Format: c.format}
			act := formatValue(formatter, c.value)

			assert.Equal(t, c.exp, act)
		})
	}
}
