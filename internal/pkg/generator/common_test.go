package generator

import (
	"dg/internal/pkg/model"
	"testing"

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
			addToFile(c.table, c.column, c.line, c.filesBefore)
			assert.Equal(t, c.filesAfter, c.filesBefore)
		})
	}
}
