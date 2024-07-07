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
				Generator: model.ToRawMessage(t, EachGenerator{
					Table:  "person",
					Column: "id",
				}),
			},
			{
				Name: "event_id",
				Type: "each",
				Generator: model.ToRawMessage(t, EachGenerator{
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

	g := EachGenerator{}

	err := g.Generate(table, files)
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

func TestGenerateEachColumnWithHigherCount(t *testing.T) {
	table := model.Table{
		Name:  "person_event",
		Count: 5,
		Columns: []model.Column{
			{
				Name: "person_id",
				Type: "each",
				Generator: model.ToRawMessage(t, EachGenerator{
					Table:  "person",
					Column: "id",
				}),
			},
			{
				Name: "event_id",
				Type: "each",
				Generator: model.ToRawMessage(t, EachGenerator{
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

	g := EachGenerator{}

	err := g.Generate(table, files)
	assert.Nil(t, err)

	exp := model.CSVFile{
		Name:   "person_event",
		Header: []string{"person_id", "event_id"},
		Lines: [][]string{
			{"p-i-1", "p-i-2", "p-i-1", "p-i-2", "p-i-1"},
			{"e-i-1", "e-i-1", "e-i-2", "e-i-2", "e-i-1"},
		},
		Output: true,
	}
	assert.Equal(t, exp, files["person_event"])
}

func TestGenerateEachColumnWithLowerCount(t *testing.T) {
	table := model.Table{
		Name:  "person_event",
		Count: 3,
		Columns: []model.Column{
			{
				Name: "person_id",
				Type: "each",
				Generator: model.ToRawMessage(t, EachGenerator{
					Table:  "person",
					Column: "id",
				}),
			},
			{
				Name: "event_id",
				Type: "each",
				Generator: model.ToRawMessage(t, EachGenerator{
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

	g := EachGenerator{}

	err := g.Generate(table, files)
	assert.Nil(t, err)

	exp := model.CSVFile{
		Name:   "person_event",
		Header: []string{"person_id", "event_id"},
		Lines: [][]string{
			{"p-i-1", "p-i-2", "p-i-1"},
			{"e-i-1", "e-i-1", "e-i-2"},
		},
		Output: true,
	}
	assert.Equal(t, exp, files["person_event"])
}
