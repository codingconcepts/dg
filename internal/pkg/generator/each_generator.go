package generator

import (
	"dg/internal/pkg/model"
	"fmt"

	"github.com/samber/lo"
)

// GenerateEachColumns looks for any each type columns for a table, and
// returns their Cartesian product back into the given files map.
func GenerateEachColumns(t model.Table, files map[string]model.CSVFile) error {
	cols := lo.Filter(t.Columns, func(c model.Column, _ int) bool {
		return c.Type == "each"
	})

	if len(cols) == 0 {
		return nil
	}

	var preCartesian [][]string
	for _, col := range cols {
		var ptc model.ProcessorTableColumn
		if err := col.Processor.UnmarshalFunc(&ptc); err != nil {
			return fmt.Errorf("parsing each process for %s.%s: %w", t.Name, col.Name, err)
		}

		srcTable := files[ptc.Table]
		srcColumn := ptc.Column
		srcColumnIndex := lo.IndexOf(srcTable.Header, srcColumn)

		preCartesian = append(preCartesian, srcTable.Lines[srcColumnIndex])
	}

	// Compute Cartesian product of all columns.
	cartesianColumns := Transpose(CartesianProduct(preCartesian...))

	// Add the header
	if _, ok := files[t.Name]; !ok {
		files[t.Name] = model.CSVFile{
			Name: t.Name,
		}
	}

	for i, col := range cartesianColumns {
		foundTable := files[t.Name]
		foundTable.Header = append(foundTable.Header, cols[i].Name)
		foundTable.Lines = append(foundTable.Lines, col)
		files[t.Name] = foundTable
	}

	return nil
}
