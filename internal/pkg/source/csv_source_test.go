package source

import (
	"os"
	"path"
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestLoadCSVSource(t *testing.T) {
	filePath := path.Join(t.TempDir(), "load_test.csv")
	assert.NoError(t, os.WriteFile(filePath, []byte("col_a,col_b,col_c\nA,B,C\n1,2,3"), os.ModePerm))

	table := "input"
	files := make(map[string]model.CSVFile)
	s := model.SourceCSV{FileName: "load_test.csv"}

	assert.NoError(t, LoadCSVSource(table, path.Dir(filePath), s, files))

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
