package generator

import (
	"dg/internal/pkg/model"
	"fmt"

	"github.com/samber/lo"
)

// GenerateRangeColumn generates sequential data between a given start and end range.
func GenerateRangeColumn(t model.Table, c model.Column, pr model.ProcessorRange, files map[string]model.CSVFile) error {
	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	switch pr.Type {
	case "date":
		lines, err := generateDateSlice(pr, count)
		if err != nil {
			return fmt.Errorf("generating date slice: %w", err)
		}

		AddToFile(t.Name, c.Name, model.FileTypeOutput, lines, files)
		return nil
	default:
		return fmt.Errorf("%q is not a valid range type", pr.Type)
	}
}
