package generator

import (
	"testing"

	"github.com/codingconcepts/dg/v1/internal/pkg/model"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIncColumn(t *testing.T) {
	cases := []struct {
		name  string
		count int
		start int
		files map[string]model.CSVFile
		exp   [][]string
	}{
		{
			name:  "with count generates as many as specified by count",
			count: 10,
			start: 100,
			files: map[string]model.CSVFile{},
			exp: [][]string{
				{"100", "101", "102", "103", "104", "105", "106", "107", "108", "109"},
			},
		},
		{
			name:  "without count generates as many as the max line",
			start: 200,
			files: map[string]model.CSVFile{
				"table": {
					Lines: [][]string{
						{"a", "b", "c"},
						{"a", "b", "c", "d", "e"},
					},
				},
			},
			exp: [][]string{
				{"a", "b", "c"},
				{"a", "b", "c", "d", "e"},
				{"200", "201", "202", "203", "204"},
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
				Name: "id",
			}

			processor := model.ProcessorInc{
				Start: c.start,
			}

			err := GenerateIncColumn(table, column, processor, c.files)
			assert.Nil(t, err)
			assert.Equal(t,
				[]string([]string{"id"}),
				c.files["table"].Header,
			)
			assert.Equal(t,
				c.exp,
				c.files["table"].Lines,
			)
		})
	}
}
