package generator

import (
	"dg/internal/pkg/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIncColumn(t *testing.T) {
	table := model.Table{
		Name:  "table",
		Count: 10,
	}

	column := model.Column{
		Name: "id",
	}

	processor := model.ProcessorInc{
		Start: 100,
	}

	files := map[string]model.CSVFile{}

	err := GenerateIncColumn(table, column, processor, files)
	assert.Nil(t, err)
	assert.Equal(t,
		[]string([]string{"id"}),
		files["table"].Header,
	)
	assert.Equal(t,
		[][]string{{"100", "101", "102", "103", "104", "105", "106", "107", "108", "109"}},
		files["table"].Lines,
	)
}
