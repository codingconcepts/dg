package generator

import (
	"github.com/codingconcepts/dg/internal/pkg/model"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGenerateMatchColumn(t *testing.T) {
	cases := []struct {
		name          string
		srcTable      *model.CSVFile
		srcTableName  string
		srcColumnName string
		srcValueName  string
		dstTable      *model.CSVFile
		dstColumns    []model.Column
		dstColumn     model.Column
		matchColumn   string
		expColumn     []string
		expError      error
	}{
		{
			name: "generates matching columns",
			srcTable: &model.CSVFile{
				Name:   "significant_events",
				Header: []string{"date", "event"},
				Lines: [][]string{
					{"2023-01-01", "2023-01-03"},
					{"abc", "def"},
				},
			},
			srcTableName:  "significant_events",
			srcColumnName: "date",
			srcValueName:  "event",
			dstTable: &model.CSVFile{
				Name:   "timeline",
				Header: []string{"timeline_date"},
				Lines: [][]string{
					{"2023-01-01", "2023-01-02", "2023-01-03"},
				},
			},
			dstColumns: []model.Column{
				{Name: "timeline_date"},
			},
			dstColumn: model.Column{
				Name: "timeline_event",
			},
			matchColumn: "timeline_date",
			expColumn:   []string{"abc", "", "def"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			table := model.Table{
				Name:    c.dstTable.Name,
				Columns: c.dstColumns,
			}

			column := c.dstColumn

			processor := model.ProcessorMatch{
				SourceTable:  c.srcTableName,
				SourceColumn: c.srcColumnName,
				SourceValue:  c.srcValueName,
				MatchColumn:  c.matchColumn,
			}

			files := map[string]model.CSVFile{}
			if c.srcTable != nil {
				files[c.srcTable.Name] = *c.srcTable
			}
			if c.dstTable != nil {
				files[c.dstTable.Name] = *c.dstTable
			}

			err := GenerateMatchColumn(table, column, processor, files)
			assert.Equal(t, c.expError, err)
			if err != nil {
				return
			}

			actColumnIndex := lo.IndexOf(files[c.dstTable.Name].Header, c.dstColumn.Name)
			assert.Equal(t, c.expColumn, files[c.dstTable.Name].Lines[actColumnIndex])
		})
	}
}

/*
	source table:

	date, event
	2023-01-01, abc
	2023-01-03, def


	dest table:

	timeline_date, timeline_event
	2023-01-01
	2023-01-02
	2023-01-03


	outcome:
	timeline_date, timeline_event
	2023-01-01, abc
	2023-01-02
	2023-01-03, def
*/
