package generator

import (
	"fmt"
	"math/rand/v2"
	"regexp"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

const (
	day   = "Day"
	month = "Month"
	year  = "Year"
)

type RelDateGenerator struct {
	Date   string `yaml:"date"`
	Unit   string `yaml:"unit"`
	After  int    `yaml:"after"`
	Before int    `yaml:"before"`
	Format string `yaml:"format"`
}

func findColumnIndex(t model.Table, name string) int {
	for i, column := range t.Columns {
		if column.Name == name {
			return i
		}
	}
	return -1
}

func (g RelDateGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if g.Format == "" {
		g.Format = "2006-01-02"
	}
	ref_date := time.Now()
	ref_column := -1
	if g.Date != "" && g.Date != "now" {
		var err error
		matched, _ := regexp.MatchString(`^[a-zA-Z]\w+$`, g.Date)
		if matched {
			ref_column = findColumnIndex(t, g.Date)
		} else {
			ref_date, err = time.Parse(g.Format, g.Date)
			if err != nil {
				return fmt.Errorf("error parsing date: %w", err)
			}
		}
	}

	if g.Unit != day && g.Unit != month && g.Unit != year {
		g.Unit = day
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}
	var lines []string
	for i := 0; i < t.Count; i++ {
		if ref_column > -1 {
			var err error
			ref_date, err = time.Parse(g.Format, files[t.Name].Lines[ref_column][i])
			if err != nil {
				return fmt.Errorf("error parsing date: %w", err)
			}
		}
		s := g.generate(ref_date)
		lines = append(lines, s)
	}
	AddTable(t, c.Name, lines, files)
	return nil
}

func (g RelDateGenerator) generate(reference time.Time) string {
	if g.After > g.Before {
		g.After, g.Before = g.Before, g.After
	}
	offset := rand.IntN(g.Before-g.After+1) + g.After
	switch g.Unit {
	case day:
		return reference.AddDate(0, 0, offset).Format(g.Format)
	case month:
		return reference.AddDate(0, offset, 0).Format(g.Format)
	case year:
		return reference.AddDate(offset, 0, 0).Format(g.Format)
	}
	return fmt.Errorf("invalid unit %s. unit must be 'day', 'month' or 'year'", g.Unit).Error()
}
