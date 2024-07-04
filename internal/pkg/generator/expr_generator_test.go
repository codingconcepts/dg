package generator

import (
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorExpMatchColumn(t *testing.T) {

	table := model.Table{
		Name:  "table",
		Count: 1,
	}

	column := model.Column{
		Name: "column",
	}

	files := map[string]model.CSVFile{
		"products": {
			Name:   "products",
			Header: []string{"product_id", "product_name", "product_price"},
			Lines: [][]string{
				{"1", "2", "3"},
				{"Apple", "bananas", "carrots"},
				{"3.00", "5.00", "2.0"},
			},
		},
		"table": {
			Name:   "table",
			Header: []string{"id"},
			Lines: [][]string{
				{"2", "1", "3"},
			},
		},
	}
	g := ExprGenerator{
		Expression: "match('products','product_id', id,'product_price') / 2.0",
	}
	err := g.Generate(table, column, files)
	assert.Nil(t, err)
	assert.Equal(t, files["table"].Lines[1][0], "2.5")
}

func TestGeneratorColumnValues(t *testing.T) {

	table := model.Table{
		Name:  "table",
		Count: 3,
	}

	column := model.Column{
		Name: "column",
	}

	files := map[string]model.CSVFile{
		"table": {
			Name:   "table",
			Header: []string{"name", "rate", "months"},
			Lines: [][]string{
				{"jhon", "jack", "joe"},
				{"3.00", "5.00", "2.0"},
				{"2", "3", "5"},
			},
		},
	}
	g := ExprGenerator{
		Expression: "rate * months",
		Format:     "%.4f",
	}
	err := g.Generate(table, column, files)
	assert.Nil(t, err)
	assert.Equal(t, files["table"].Lines[3][0], "6.0000")
}
