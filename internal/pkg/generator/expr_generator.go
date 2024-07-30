package generator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/samber/lo"
)

type ExprGenerator struct {
	Expression string `yaml:"expression"`
	Format     string `yaml:"format"`
}

type ExpressionContext struct {
	Files map[string]model.CSVFile
}

func (g ExprGenerator) Generate(t model.Table, c model.Column, files map[string]model.CSVFile) error {
	if g.Expression == "" {
		return fmt.Errorf("expression cannot be empty")
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].Lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	ctx := &ExpressionContext{Files: files}
	functions := map[string]govaluate.ExpressionFunction{
		"match": func(args ...interface{}) (interface{}, error) {
			if len(args) != 4 {
				return "", fmt.Errorf("match function expects 4 arguments")
			}
			sourceTable, sourceColumn, matchColumn := args[0].(string), args[1].(string), args[3].(string)
			sourceValue := fmt.Sprintf("%v", args[2])
			value, err := ctx.searchTable(sourceTable, sourceColumn, sourceValue, matchColumn)
			if err != nil {
				return nil, err
			}
			return coerce(value), nil
		},
		"add_date": func(args ...interface{}) (interface{}, error) {
			if len(args) != 4 {
				return "", fmt.Errorf("add_date function expects 5 arguments: add_date(years int, months int, days int, data string)")
			}
			years, err := strconv.Atoi(fmt.Sprintf("%v", args[0]))
			if err != nil {
				return "", fmt.Errorf("error parsing years: %w", err)
			}
			months, err := strconv.Atoi(fmt.Sprintf("%v", args[1]))
			if err != nil {
				return "", fmt.Errorf("error parsing months: %w", err)
			}
			days, err := strconv.Atoi(fmt.Sprintf("%v", args[2]))
			if err != nil {
				return "", fmt.Errorf("error parsing days: %w", err)
			}
			var data time.Time
			tipo := getType(args[3])
			switch tipo {
			case "time.Time":
				data = args[3].(time.Time)
				return data.AddDate(years, months, days), nil
			case "float64":
				float, _ := args[3].(float64)
				sec := int64(float)
				nano := int64((float - float64(sec)) * 1e9)
				data = time.Unix(sec, nano)
				return data.AddDate(years, months, days), nil
			case "int":
				sec := int64(args[3].(int))
				data = time.Unix(sec, 0)
				return data.AddDate(years, months, days), nil
			case "string":
				digits := regexp.MustCompile(`(\d)+`)
				match := digits.FindAllString(g.Format, -1)
				if len(match) >= 3 {
					data, err = time.Parse(g.Format, args[3].(string))
				} else {
					data, err = time.Parse("2006-01-02", args[3].(string))
				}
				if err != nil {
					return "", fmt.Errorf("error parsing date: %w", err)
				}
				return data.AddDate(years, months, days), nil
			}
			return "", fmt.Errorf("error parsing date")
		},
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(g.Expression, functions)
	if err != nil {
		return fmt.Errorf("error parsing expression: %w", err)
	}
	var lines []string
	for i := 0; i < t.Count; i++ {
		table := files[t.Name]
		columns := len(table.Header)
		parameters := make(map[string]interface{}, columns)
		for c := range columns {
			s := table.Lines[c][i]
			parameters[table.Header[c]] = coerce(s)
		}
		result, err := expression.Evaluate(parameters)
		if err != nil {
			return fmt.Errorf("error evaluating expression %w", err)
		}
		switch getType(result) {
		case "time.Time":
			if g.Format == "" {
				g.Format = "2006-01-02"
			}
			lines = append(lines, result.(time.Time).Format(g.Format))
		case "float64":
			if g.Format == "" {
				g.Format = "%g"
			}
			lines = append(lines, fmt.Sprintf(g.Format, result.(float64)))
		case "int":
			if g.Format == "" {
				g.Format = "%d"
			}
			lines = append(lines, fmt.Sprintf(g.Format, result.(int)))
		case "bool":
			if g.Format == "" {
				g.Format = "%t"
			}
			lines = append(lines, fmt.Sprintf(g.Format, result.(bool)))
		case "string":
			lines = append(lines, result.(string))
		default:
			lines = append(lines, fmt.Sprintf("%v", result))
		}
	}
	AddTable(t, c.Name, lines, files)
	return nil
}

func (tc *ExpressionContext) searchTable(sourceTable, sourceColumn, sourceValue, matchColumn string) (string, error) {
	csvFile, exists := tc.Files[sourceTable]
	if !exists {
		return "", fmt.Errorf("table not found: %s", sourceTable)
	}

	sourceColumnIndex := lo.IndexOf(csvFile.Header, sourceColumn)
	matchColumnIndex := lo.IndexOf(csvFile.Header, matchColumn)
	if sourceColumnIndex == -1 || matchColumnIndex == -1 {
		return "", fmt.Errorf("column not found: %s ou %s", sourceColumn, matchColumn)
	}
	_, index, found := lo.FindIndexOf(csvFile.Lines[sourceColumnIndex], func(item string) bool {
		return item == sourceValue
	})
	if found {
		return csvFile.Lines[matchColumnIndex][index], nil
	}

	return "", fmt.Errorf("value not found for %s in column %s", sourceValue, sourceColumn)
}

func coerce(value string) interface{} {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	dateFormats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"02/01/2006",
		"02-01-2006",
		"02/01/2006 15:04:05",
		"02-01-2006 15:04:05",
	}

	for _, format := range dateFormats {
		if t, err := time.Parse(format, value); err == nil {
			return t
		}
	}
	return value
}

func getType(value interface{}) string {
	return reflect.TypeOf(value).String()
}
