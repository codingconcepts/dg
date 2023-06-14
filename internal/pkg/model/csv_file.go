package model

import (
	"github.com/samber/lo"
)

// CSVFile represents the content of a CSV file.
type CSVFile struct {
	Name          string
	Header        []string
	Lines         [][]string
	UniqueColumns []string
	Output        bool
}

// Unique removes any duplicates from the CSVFile's lines.
func (c *CSVFile) Unique() [][]string {
	uniqueColumnIndexes := uniqueIndexes(c.Header, c.UniqueColumns)

	uniqueValues := map[string]struct{}{}
	var uniqueLines [][]string

	for i := 0; i < len(c.Lines); i++ {
		key := uniqueKey(uniqueColumnIndexes, c.Lines[i])

		if _, ok := uniqueValues[key]; !ok {
			uniqueLines = append(uniqueLines, c.Lines[i])
			uniqueValues[key] = struct{}{}
		}
	}

	return uniqueLines
}

func uniqueIndexes(header, uniqueColumns []string) []int {
	indexes := []int{}

	for i, h := range header {
		if lo.Contains(uniqueColumns, h) {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func uniqueKey(indexes []int, line []string) string {
	output := ""

	for i, col := range line {
		if lo.Contains(indexes, i) {
			output += col
		} else {
			output += "-"
		}
	}

	return output
}
