package generator

import (
	"dg/internal/pkg/model"
	"fmt"
	"math/rand"
)

// GenerateSetColumn selects between a set of values for a given table.
func GenerateSetColumn(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	var ptc model.ProcessorSet
	if err := c.Processor.UnmarshalFunc(&ptc); err != nil {
		return fmt.Errorf("parsing set process for %s.%s: %w", t.Name, c.Name, err)
	}

	if len(ptc.Values) == 0 {
		return fmt.Errorf("no values provided for set generator")
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, ptc.Values[rand.Intn(len(ptc.Values))])
	}

	addToFile(t.Name, c.Name, line, files)

	return nil
}
