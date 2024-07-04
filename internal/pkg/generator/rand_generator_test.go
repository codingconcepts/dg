package generator

import (
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestGeneratorRand(t *testing.T) {
	cases := []struct {
		name     string
		Type     string
		Low      string
		High     string
		Format   string
		testFunc func(val string) bool
	}{
		{
			name:   "random_date_formatted",
			Type:   "date",
			Low:    "01-11-2023",
			High:   "30-12-2024",
			Format: "02-01-2006",
			testFunc: func(val string) bool {
				_, err := time.Parse("02-01-2006", val)
				return err == nil
			},
		},
		{
			name: "random_date_default_format",
			Type: "date",
			Low:  "2023-11-01",
			High: "2024-12-30",
			testFunc: func(val string) bool {
				_, err := time.Parse("2006-01-02", val)
				return err == nil
			},
		},
		{
			name:   "random_int_formatted",
			Type:   "int",
			Low:    "1",
			High:   "30",
			Format: "%05d",
			testFunc: func(val string) bool {
				_, err := strconv.Atoi(val)
				return err == nil
			},
		},
		{
			name:   "random_int_default_format",
			Type:   "int",
			Low:    "1",
			High:   "30",
			Format: "%05d",
			testFunc: func(val string) bool {
				_, err := strconv.Atoi(val)
				return err == nil
			},
		},
		{
			name:   "random_float64_formatted",
			Type:   "float64",
			Low:    "1.0",
			High:   "20.1234",
			Format: "%0.2f",
			testFunc: func(val string) bool {
				match, _ := regexp.MatchString(`\d+\.\d{2}`, val)
				_, err := strconv.ParseFloat(val, 64)
				return match && err == nil
			},
		},
		{
			name: "random_float64_default_format",
			Type: "float64",
			Low:  "1.0",
			High: "20.1234",
			testFunc: func(val string) bool {
				_, err := strconv.ParseFloat(val, 64)
				return err == nil
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := model.Table{
				Name:  "table",
				Count: 1,
			}
			column := model.Column{
				Name: "col",
			}

			files := map[string]model.CSVFile{}
			g := RandGenerator{
				Type:   c.Type,
				Low:    c.Low,
				High:   c.High,
				Format: c.Format,
			}

			err := g.Generate(table, column, files)
			assert.Nil(t, err)
			assert.True(t, c.testFunc(files["table"].Lines[0][0]))
		})
	}
}
