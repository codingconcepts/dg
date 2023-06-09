package generator

import (
	"dg/internal/pkg/model"
	"fmt"
)

// AddToFile adds a column to a table in the given files map.
func AddToFile(table, column string, fileType model.FileType, line []string, files map[string]model.CSVFile) {
	if _, ok := files[table]; !ok {
		files[table] = model.CSVFile{
			Name: table,
			Type: fileType,
		}
	}

	foundTable := files[table]
	foundTable.Header = append(foundTable.Header, column)
	foundTable.Lines = append(foundTable.Lines, line)
	files[table] = foundTable
}

func formatValue(fp model.FormatterProcessor, value any) string {
	format := fp.GetFormat()
	if format != "" {
		// Check if the value implements the formatter interface and use that first,
		// otherwise, just perform a simple string format.
		if f, ok := value.(model.Formatter); ok {
			return f.Format(format)
		} else {
			return fmt.Sprintf(format, value)
		}
	} else {
		return fmt.Sprintf("%v", value)
	}
}
