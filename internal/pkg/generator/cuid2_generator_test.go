package generator

import (
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/nrednav/cuid2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCuid2Column(t *testing.T) {
	table := model.Table{
		Name:  "table",
		Count: 10,
	}

	column := model.Column{
		Name: "id",
	}

	g := Cuid2Generator{
		Length: 14,
	}

	files := map[string]model.CSVFile{}

	err := g.Generate(table, column, files)
	assert.Nil(t, err)

	assert.Equal(t,
		[]string([]string{"id"}),
		files["table"].Header,
	)
	for _, line := range files["table"].Lines {
		for _, cuid := range line {
			assert.True(t, cuid2.IsCuid(cuid))
		}
	}
}
