package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"

	"github.com/samber/lo"
)

// GenerateMatchColumn matches values from a previously generated table and inserts values
// into a new table where match is found.
func GenerateMatchColumn(t model.Table, c model.Column, ptc model.ProcessorMatch, files map[string]model.CSVFile) error {
	sourceTable, ok := files[ptc.SourceTable]
	if !ok {
		return fmt.Errorf("missing source table %q for match lookup", ptc.SourceTable)
	}

	sourceColumnIndex := lo.IndexOf(sourceTable.Header, ptc.SourceColumn)
	sourceColumn := sourceTable.Lines[sourceColumnIndex]

	valueColumnIndex := lo.IndexOf(sourceTable.Header, ptc.SourceValue)
	valueColumn := sourceTable.Lines[valueColumnIndex]

	matchTable, ok := files[t.Name]
	if !ok {
		return fmt.Errorf("missing destination table %q for match lookup", t.Name)
	}
	_, matchColumnIndex, ok := lo.FindIndexOf(t.Columns, func(c model.Column) bool {
		return c.Name == ptc.MatchColumn
	})
	if !ok {
		return fmt.Errorf("missing match column %q in current table", ptc.MatchColumn)
	}
	matchColumn := matchTable.Lines[matchColumnIndex]

	lines := make([]string, len(matchColumn))
	for sourceI, sourceC := range sourceColumn {
		if _, i, ok := lo.FindIndexOf(matchColumn, func(matchCol string) bool {
			return matchCol == sourceC
		}); ok {
			lines[i] = valueColumn[sourceI]
		}
	}

	AddTable(t.Name, c.Name, lines, files)
	return nil
}
