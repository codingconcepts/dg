package generator

import (
	"fmt"

	"github.com/codingconcepts/dg/internal/pkg/model"

	"github.com/samber/lo"
)

// MatchGenerator provides additional context to a match column.
type MatchGenerator struct {
	SourceTable  string `yaml:"source_table"`
	SourceColumn string `yaml:"source_column"`
	SourceValue  string `yaml:"source_value"`
	MatchColumn  string `yaml:"match_column"`
}

// Generate matches values from a previously generated table and inserts values
// into a new table where match is found.
func (g MatchGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	sourceTable, ok := files[g.SourceTable]
	if !ok {
		return fmt.Errorf("missing source table %q for match lookup", g.SourceTable)
	}

	sourceColumnIndex := lo.IndexOf(sourceTable.Header, g.SourceColumn)
	sourceColumn := sourceTable.Lines[sourceColumnIndex]

	valueColumnIndex := lo.IndexOf(sourceTable.Header, g.SourceValue)
	valueColumn := sourceTable.Lines[valueColumnIndex]

	sourceMap := map[string]string{}
	for i := 0; i < len(sourceColumn); i++ {
		sourceMap[sourceColumn[i]] = valueColumn[i]
	}

	matchTable, ok := files[t.Name]
	if !ok {
		return fmt.Errorf("missing destination table %q for match lookup", t.Name)
	}

	// Use the match table headers to determine index, as the each processor
	// will re-order columns.
	_, matchColumnIndex, ok := lo.FindIndexOf(matchTable.Header, func(c string) bool {
		return c == g.MatchColumn
	})
	if !ok {
		return fmt.Errorf("missing match column %q in current table", g.MatchColumn)
	}

	matchColumn := matchTable.Lines[matchColumnIndex]

	lines := make([]string, len(matchColumn))
	for i, matchC := range matchColumn {
		if sourceValue, ok := sourceMap[matchC]; ok {
			lines[i] = sourceValue
		}
	}

	AddTable(t, c.Name, lines, files)
	return nil
}
