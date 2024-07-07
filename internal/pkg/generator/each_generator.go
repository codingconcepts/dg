package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"

	"github.com/samber/lo"
)

// EachGenerator provides additional context to an each or ref column.
type EachGenerator struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

// Generate looks for any each type columns for a table, and
// returns their Cartesian product back into the given files map.
func (g EachGenerator) Generate(t model.Table, files map[string]model.CSVFile) error {
	cols := lo.Filter(t.Columns, func(c model.Column, _ int) bool {
		return c.Type == "each"
	})

	if len(cols) == 0 {
		return nil
	}

	var preCartesian [][]string
	for _, col := range cols {
		var gCol EachGenerator
		if err := col.Generator.UnmarshalFunc(&gCol); err != nil {
			return fmt.Errorf("parsing each process for %s.%s: %w", t.Name, col.Name, err)
		}

		srcTable := files[gCol.Table]
		srcColumn := gCol.Column
		srcColumnIndex := lo.IndexOf(srcTable.Header, srcColumn)

		if len(srcTable.Lines)-1 < srcColumnIndex {
			return fmt.Errorf("column %q out of bounds for table %q", srcColumn, srcTable.Name)
		}

		preCartesian = append(preCartesian, srcTable.Lines[srcColumnIndex])
	}

	// Compute Cartesian product of all columns.
	cartesianColumns := Transpose(CartesianProduct(preCartesian...))

	// if count is set adjust the cartesian product to fit the count
	if t.Count > 0 {
		newCartesianColumns := make([][]string, len(cartesianColumns))
		size := len(cartesianColumns[0])
		for c := range cartesianColumns {
			for i := 0; i < t.Count; i++ {
				newCartesianColumns[c] = append(newCartesianColumns[c], cartesianColumns[c][i%size])
			}
		}
		cartesianColumns = newCartesianColumns
	}
	// Add the header
	for i, col := range cartesianColumns {
		AddTable(t, cols[i].Name, col, files)
	}

	return nil
}
