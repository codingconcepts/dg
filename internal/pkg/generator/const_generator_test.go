package generator

import (
	"fmt"
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateConstColumn(t *testing.T) {
	cases := []struct {
		name       string
		tableCount int
		files      map[string]model.CSVFile
		values     []string
		expErr     error
	}{
		{
			name:       "first column in table",
			tableCount: 3,
			files:      map[string]model.CSVFile{},
			values:     []string{"a", "b", "c"},
		},
		{
			name: "less than current table size",
			files: map[string]model.CSVFile{
				"table": {
					Name:   "table",
					Header: []string{"col_a", "col_b", "col_c"},
					Lines: [][]string{
						{"val_1", "val_2", "val_3"},
						{"val_1", "val_2", "val_3"},
					},
				},
			},
			values: []string{"a"},
			expErr: fmt.Errorf("wrong number of values provided for const generator (need 3, got 1)"),
		},
		{
			name:       "less than current table size with table count",
			tableCount: 10,
			files: map[string]model.CSVFile{
				"table": {
					Name:   "table",
					Header: []string{"col_a", "col_b", "col_c"},
					Lines: [][]string{
						{"val_1", "val_2", "val_3"},
						{"val_1", "val_2", "val_3"},
					},
				},
			},
			values: []string{"a"},
			expErr: fmt.Errorf("wrong number of values provided for const generator (need 3, got 1)"),
		},
		{
			name: "same as current table size",
			files: map[string]model.CSVFile{
				"table": {
					Name:   "table",
					Header: []string{"col_a", "col_b", "col_c"},
					Lines: [][]string{
						{"val_1", "val_2", "val_3"},
						{"val_1", "val_2", "val_3"},
					},
				},
			},
			values: []string{"a", "b", "c"},
		},
		{
			name: "more than current table size",
			files: map[string]model.CSVFile{
				"table": {
					Name:   "table",
					Header: []string{"col_a", "col_b", "col_c"},
					Lines: [][]string{
						{"val_1", "val_2", "val_3"},
						{"val_1", "val_2", "val_3"},
					},
				},
			},
			values: []string{"a", "b", "c", "d", "e"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pc := model.ProcessorConst{
				Values: c.values,
			}

			table := model.Table{
				Name:  "table",
				Count: c.tableCount,
			}

			column := model.Column{
				Name: "col",
			}

			actErr := GenerateConstColumn(table, column, pc, c.files)
			assert.Equal(t, c.expErr, actErr)
			if actErr != nil {
				return
			}

			assert.Equal(t, c.values, c.files["table"].Lines[len(c.files["table"].Lines)-1])
		})
	}
}
