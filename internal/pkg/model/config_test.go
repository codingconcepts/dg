package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	y := `
inputs:
  - name: my_data
    type: csv
    source:
      file_name: my_data.csv

tables:
  - name: person
    count: 100
    columns:
      - name: id
        type: inc
        processor:
          start: 1
          format: "P%03d"
`

	config, err := LoadConfig(strings.NewReader(y))
	assert.Nil(t, err)

	exp := Config{
		Inputs: []Input{
			{
				Name: "my_data",
				Type: "csv",
				Source: ToRawMessage(t, SourceCSV{
					FileName: "my_data.csv",
				}),
			},
		},
		Tables: []Table{
			{
				Name:  "person",
				Count: 100,
				Columns: []Column{
					{
						Name: "id",
						Type: "inc",
						Generator: ToRawMessage(t, map[string]any{
							"start":  1,
							"format": "P%03d",
						}),
					},
				},
			},
		},
	}

	assert.Equal(t, exp.Inputs[0].Name, config.Inputs[0].Name)
	assert.Equal(t, exp.Inputs[0].Type, config.Inputs[0].Type)

	var expSource SourceCSV
	assert.Nil(t, exp.Inputs[0].Source.UnmarshalFunc(&expSource))

	var actSource SourceCSV
	assert.Nil(t, config.Inputs[0].Source.UnmarshalFunc(&actSource))

	assert.Equal(t, expSource, actSource)

	assert.Equal(t, exp.Tables[0].Name, config.Tables[0].Name)
	assert.Equal(t, exp.Tables[0].Count, config.Tables[0].Count)
	assert.Equal(t, exp.Tables[0].Columns[0].Name, config.Tables[0].Columns[0].Name)
	assert.Equal(t, exp.Tables[0].Columns[0].Type, config.Tables[0].Columns[0].Type)

	var expProcessor map[string]any
	assert.Nil(t, exp.Tables[0].Columns[0].Generator.UnmarshalFunc(&expProcessor))

	var actProcessor map[string]any
	assert.Nil(t, config.Tables[0].Columns[0].Generator.UnmarshalFunc(&actProcessor))

	assert.Equal(t, expProcessor, actProcessor)
}
