package generator

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

const (
	day   = "Day"
	month = "Month"
	year  = "Year"
)

type RelativeDateGenerator struct {
	Date   string `yaml:"date"`
	Unit   string `yaml:"unit"`
	Low    int    `yaml:"low"`
	High   int    `yaml:"high"`
	Format string `yaml:"format"`

	Ref_Date time.Time
}

func (g RelativeDateGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {

	if g.Format == "" {
		g.Format = "2006-01-02"
	}

	if g.Date == "" {
		g.Ref_Date = time.Now()
	} else {
		var err error
		g.Ref_Date, err = time.Parse(g.Format, g.Date)
		if err != nil {
			return fmt.Errorf("error parsing date: %w", err)
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
		s := g.generate()
		lines = append(lines, s)
	}
	AddTable(t, c.Name, lines, files)
	return nil
}

func (g RelativeDateGenerator) generate() string {
	if g.Low > g.High {
		g.Low, g.High = g.High, g.Low
	}

	offset := rand.IntN(g.High-g.Low+1) + g.Low
	switch g.Unit {
	case day:
		return g.Ref_Date.AddDate(0, 0, offset).Format(g.Format)
	case month:
		return g.Ref_Date.AddDate(0, offset, 0).Format(g.Format)
	case year:
		return g.Ref_Date.AddDate(offset, 0, 0).Format(g.Format)
	}
	return fmt.Errorf("invalid unit %s. unit must be 'day', 'month' or 'year'", g.Unit).Error()
}
