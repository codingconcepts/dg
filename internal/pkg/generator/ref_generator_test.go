package generator

import (
	"dg/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRefColumn(t *testing.T) {
	table := model.Table{
		Name:  "pet",
		Count: 2,
	}

	column := model.Column{
		Name: "person_id",
	}

	processor := model.ProcessorTableColumn{
		Table:  "person",
		Column: "id",
	}

	files := map[string]model.CSVFile{
		"person": {
			Header: []string{"id"},
			Lines:  [][]string{{"ce9af887-37eb-4e08-9790-4f481b0fa594"}},
		},
	}
	err := GenerateRefColumn(table, column, processor, files)
	assert.Nil(t, err)
	assert.Equal(t, "ce9af887-37eb-4e08-9790-4f481b0fa594", files["pet"].Lines[0][0])
	assert.Equal(t, "ce9af887-37eb-4e08-9790-4f481b0fa594", files["pet"].Lines[0][1])
}
