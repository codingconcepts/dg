package generator

import "dg/internal/pkg/model"

func addToFile(table, column string, line []string, files map[string]model.CSVFile) {
	if _, ok := files[table]; !ok {
		files[table] = model.CSVFile{
			Name: table,
		}
	}

	foundTable := files[table]
	foundTable.Header = append(foundTable.Header, column)
	foundTable.Lines = append(foundTable.Lines, line)
	files[table] = foundTable
}
