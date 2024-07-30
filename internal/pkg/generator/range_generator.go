package generator

import (
	"fmt"
	"strconv"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

// RangeGenerator provides additional context to a range column.
type RangeGenerator struct {
	Table  string `yaml:"table"`
	Type   string `yaml:"type"`
	From   string `yaml:"from"`
	To     string `yaml:"to"`
	Step   string `yaml:"step"`
	Format string `yaml:"format"`
}

// Generate sequential data between a given start and end range.
func (g RangeGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	if g.Table != "" {
		csvFile, ok := files[g.Table]
		if !ok {
			return fmt.Errorf("table %q not found", g.Table)
		}
		columnIndex := lo.IndexOf(csvFile.Header, c.Name)
		if columnIndex == -1 {
			return fmt.Errorf("column %q not found in table %q", c.Name, g.Table)
		}
		line := csvFile.Lines[columnIndex]
		size := len(line)
		if size == 0 {
			return fmt.Errorf("no data in column %q of table %q", c.Name, g.Table)
		}
		g.From = line[size-1]
		count += 1
	}

	switch g.Type {
	case "date":
		lines, err := g.generateDateSlice(count)
		if err != nil {
			return fmt.Errorf("generating date slice: %w", err)
		}
		if g.Table != "" && len(lines) > 0 {
			lines = lines[1:]
		}
		AddTable(t, c.Name, lines, files)
		return nil

	case "int":
		lines, err := g.generateIntSlice(count)
		if err != nil {
			return fmt.Errorf("generating int slice: %w", err)
		}
		if g.Table != "" && len(lines) > 0 {
			lines = lines[1:]
		}
		AddTable(t, c.Name, lines, files)
		return nil

	default:
		return fmt.Errorf("%q is not a valid range type", g.Type)
	}
}

func (g RangeGenerator) generateDateSlice(count int) ([]string, error) {
	// Validate that we have everything we need.
	if count == 0 && g.Step == "" {
		return nil, fmt.Errorf("either a count or a step must be provided to a date range generator")
	}

	from, err := time.Parse(g.Format, g.From)
	if err != nil {
		return nil, fmt.Errorf("parsing from date: %w", err)
	}

	to, err := time.Parse(g.Format, g.To)
	if err != nil {
		return nil, fmt.Errorf("parsing to date: %w", err)
	}

	var step time.Duration
	if count > 0 {
		step = to.Sub(from) / time.Duration(count)
	} else {
		if step, err = time.ParseDuration(g.Step); err != nil {
			return nil, fmt.Errorf("parsing step: %w", err)
		}
	}

	var s []string
	for i := from; i.Before(to); i = i.Add(step) {
		s = append(s, i.Format(g.Format))
	}

	return s, nil
}

func (g RangeGenerator) generateIntSlice(count int) ([]string, error) {
	// Validate that we have everything we need.
	if count == 0 && g.Step == "" {
		return nil, fmt.Errorf("either a count or a step must be provided to an int range generator")
	}

	from, err := strconv.Atoi(g.From)
	if err != nil {
		return nil, fmt.Errorf("parsing from number: %w", err)
	}

	var to int
	if g.To == "" {
		to = from + count - 1
	} else {
		if to, err = strconv.Atoi(g.To); err != nil {
			return nil, fmt.Errorf("parsing to number: %w", err)
		}
	}

	var step int
	if count > 0 {
		step = (to - from) / (count - 1)
	} else {
		if step, err = strconv.Atoi(g.Step); err != nil {
			return nil, fmt.Errorf("parsing step number: %w", err)
		}
	}

	var s []string
	for i := from; i <= to; i += step {
		s = append(s, strconv.Itoa(i))
	}

	return s, nil
}
