package source

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/codingconcepts/dg/internal/pkg/generator"
	"github.com/codingconcepts/dg/internal/pkg/model"
)

// LoadCSVSource loads a CSV file from disk and adds it as a table to files.
func LoadCSVSource(table, configDir string, s model.SourceCSV, files map[string]model.CSVFile) (err error) {
	fullPath := path.Join(configDir, s.FileName)
	file, err := os.Open(fullPath)
	if err != nil {
		return fmt.Errorf("opening csv file: %w", err)
	}
	defer func() {
		if ferr := file.Close(); ferr != nil {
			err = ferr
		}
	}()

	return processCSVSource(file, table, files)
}

func processCSVSource(file io.Reader, table string, files map[string]model.CSVFile) error {
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("reading csv file: %w", err)
	}

	headers := rows[0]
	columns := generator.Transpose(rows[1:])

	for i, column := range columns {
		generator.AddInput(table, headers[i], column, files)
	}

	return nil
}
