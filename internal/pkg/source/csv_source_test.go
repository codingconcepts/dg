package source

import (
	"strings"
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestProcessCSVSource(t *testing.T) {
	file := strings.NewReader("col_a,col_b,col_c\nA,B,C\n1,2,3")

	table := "input"

	files := make(map[string]model.CSVFile)

	err := processCSVSource(file, table, files)
	assert.Nil(t, err)

	expCSVFile := model.CSVFile{
		Name:   "input",
		Header: []string{"col_a", "col_b", "col_c"},
		Lines: [][]string{
			{"A", "1"},
			{"B", "2"},
			{"C", "3"}},
		Output: false}
	assert.Equal(t, expCSVFile, files["input"])
}
