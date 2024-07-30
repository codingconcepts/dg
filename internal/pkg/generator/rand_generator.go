package generator

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

type RandGenerator struct {
	Type   string `yaml:"type"`
	Low    string `yaml:"low"`
	High   string `yaml:"high"`
	Format string `yaml:"format"`
}

func (g RandGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	count := len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
		return len(a) > len(b)
	}))

	if count == 0 {
		count = t.Count
	}

	switch g.Type {
	case "date":
		lines, err := g.generateDateRand(count)
		if err != nil {
			return fmt.Errorf("generating random date: %w", err)
		}

		AddTable(t, c.Name, lines, files)
		return nil

	case "int":
		lines, err := g.generateIntRand(count)
		if err != nil {
			return fmt.Errorf("generating random int: %w", err)
		}

		AddTable(t, c.Name, lines, files)
		return nil

	case "float64":
		lines, err := g.generateFloatRand(count)
		if err != nil {
			return fmt.Errorf("generating random float64: %w", err)
		}

		AddTable(t, c.Name, lines, files)
		return nil
	default:
		return fmt.Errorf("%q is not a valid random type", g.Type)
	}
}

func (g RandGenerator) generateIntRand(count int) ([]string, error) {
	if g.Low == "" && g.High == "" {
		return nil, fmt.Errorf("'low' and 'high' values must be provided to an int rand generator")
	}
	low, err := strconv.Atoi(g.Low)
	if err != nil {
		return nil, fmt.Errorf("parsing from number: %w", err)
	}
	high, err := strconv.Atoi(g.High)
	if err != nil {
		return nil, fmt.Errorf("parsing from number: %w", err)
	}
	if low > high {
		low, high = high, low
	}
	var lines []string
	for i := 0; i < count; i++ {
		value := rand.Intn(high-low) + low
		if g.Format == "" {
			g.Format = "%v"
		}
		lines = append(lines, fmt.Sprintf(g.Format, value))
	}

	return lines, nil
}

func (g RandGenerator) generateFloatRand(count int) ([]string, error) {
	if g.Low == "" && g.High == "" {
		return nil, fmt.Errorf("'low' and 'high' values must be provided to an float64 rand generator")
	}
	low, err := strconv.ParseFloat(g.Low, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing from number: %w", err)
	}
	high, err := strconv.ParseFloat(g.High, 64)
	if err != nil {
		return nil, fmt.Errorf("parsing from number: %w", err)
	}
	if low > high {
		low, high = high, low
	}
	var lines []string
	for i := 0; i < count; i++ {
		value := rand.Float64()*(high-low) + low
		if g.Format == "" {
			g.Format = "%g"
		}
		f := fmt.Sprintf(g.Format, value)
		lines = append(lines, f)
	}

	return lines, nil
}

func (g RandGenerator) generateDateRand(count int) ([]string, error) {
	if g.Low == "" || g.High == "" {
		return nil, fmt.Errorf("'low' and 'high' values must be provided to a date rand generator")
	}
	if g.Format == "" {
		g.Format = "2006-01-02"
	}

	low, err := time.Parse(g.Format, g.Low)
	if err != nil {
		return nil, fmt.Errorf("parsing low date: %w", err)
	}
	high, err := time.Parse(g.Format, g.High)
	if err != nil {
		return nil, fmt.Errorf("parsing high date: %w", err)
	}
	if low.After(high) {
		low, high = high, low
	}
	var lines []string
	diff := high.Unix() - low.Unix()
	if diff <= 0 {
		return nil, fmt.Errorf("no range found between low and high dates")
	}

	for i := 0; i < count; i++ {
		randomOffset := rand.Int63n(diff + 1) // +1 to include the high date in the range
		randomDate := low.Add(time.Duration(randomOffset) * time.Second).Format(g.Format)
		lines = append(lines, randomDate)
	}

	return lines, nil
}
