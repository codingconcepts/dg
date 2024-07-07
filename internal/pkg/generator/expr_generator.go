package generator

import (
	"fmt"
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
	if g.Format == "" {
		g.Format = "%v"
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
		lines = append(lines, fmt.Sprintf(g.Format, result))
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
