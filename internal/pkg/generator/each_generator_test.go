package generator

import (
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateEachColumn(t *testing.T) {
	table := model.Table{
		Name: "person_event",
		Columns: []model.Column{
			{
				Name: "person_id",
				Type: "each",
				Processor: toRawMessage(t, model.ProcessorTableColumn{
					Table:  "person",
					Column: "id",
				}),
			},
			{
				Name: "event_id",
				Type: "each",
				Processor: toRawMessage(t, model.ProcessorTableColumn{
					Table:  "event",
					Column: "id",
				}),
			},
		},
	}

	files := map[string]model.CSVFile{
		"person": {
			Name:   "person",
			Header: []string{"id", "name"},
			Lines: [][]string{
				{"p-i-1", "p-i-2"},
				{"p-one", "p-two"},
			},
		},
		"event": {
			Name:   "event",
			Header: []string{"id", "name"},
			Lines: [][]string{
				{"e-i-1", "e-i-2"},
				{"e-one", "e-two"},
			},
		},
	}

	err := GenerateEachColumns(table, files)
	assert.Nil(t, err)

	exp := model.CSVFile{
		Name:   "person_event",
		Header: []string{"person_id", "event_id"},
		Lines: [][]string{
			{"p-i-1", "p-i-2", "p-i-1", "p-i-2"},
			{"e-i-1", "e-i-1", "e-i-2", "e-i-2"},
		},
		Output: true,
	}
	assert.Equal(t, exp, files["person_event"])
}
