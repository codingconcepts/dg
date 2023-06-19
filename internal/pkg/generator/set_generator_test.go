package generator

import (
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSetColumn(t *testing.T) {
	table := model.Table{
		Name:  "table",
		Count: 10,
	}

	column := model.Column{
		Name: "id",
	}

	g := SetGenerator{
		Values:  []string{"a", "b", "c"},
		Weights: []int{0, 1, 0},
	}

	files := map[string]model.CSVFile{}

	err := g.Generate(table, column, files)
	assert.Nil(t, err)
	assert.Equal(t,
		[]string([]string{"id"}),
		files["table"].Header,
	)
	assert.Equal(t,
		[][]string{{"b", "b", "b", "b", "b", "b", "b", "b", "b", "b"}},
		files["table"].Lines,
	)
}
