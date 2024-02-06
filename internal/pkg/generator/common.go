package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"
)

// AddTable adds a column to a table in the given files map.
func AddTable(table model.Table, column string, line []string, files map[string]model.CSVFile) {
	if _, ok := files[table.Name]; !ok {
		files[table.Name] = model.CSVFile{
			Name:          table.Name,
			Output:        !table.Suppress,
			UniqueColumns: table.UniqueColumns,
		}
	}

	add(files, table.Name, column, line)
}

// AddInput adds a column to a table in the given files map.
func AddInput(table, column string, line []string, files map[string]model.CSVFile) {
	if _, ok := files[table]; !ok {
		files[table] = model.CSVFile{
			Name:   table,
			Output: false,
		}
	}

	add(files, table, column, line)
}

func add(files map[string]model.CSVFile, table string, column string, line []string) {
	foundTable := files[table]
	foundTable.Header = append(foundTable.Header, column)
	foundTable.Lines = append(foundTable.Lines, line)
	files[table] = foundTable
}

func formatValue(fp FormatterProcessor, value any) string {
	format := fp.GetFormat()
	if format != "" {
		// Check if the value implements the formatter interface and use that first,
		// otherwise, just perform a simple string format.
		if f, ok := value.(Formatter); ok {
			return f.Format(format)
		} else {
			return fmt.Sprintf(format, value)
		}
	} else {
		return fmt.Sprintf("%v", value)
	}
}
