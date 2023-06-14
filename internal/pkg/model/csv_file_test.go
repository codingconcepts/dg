package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	cases := []struct {
		name          string
		uniqueColumns []string
		exp           [][]string
	}{
		{
			name:          "1 column unique",
			uniqueColumns: []string{"col_1"},
			exp: [][]string{
				{"a", "d", "g"},
				{"b", "d", "g"},
				{"c", "d", "g"},
			},
		},
		{
			name:          "2 column unique",
			uniqueColumns: []string{"col_1", "col_2"},
			exp: [][]string{
				{"a", "d", "g"},
				{"b", "d", "g"},
				{"c", "d", "g"},
				{"a", "e", "g"},
				{"b", "e", "g"},
				{"c", "e", "g"},
				{"a", "f", "g"},
				{"b", "f", "g"},
				{"c", "f", "g"},
			},
		},
		{
			name:          "3 column unique",
			uniqueColumns: []string{"col_1", "col_2", "col_3"},
			exp: [][]string{
				{"a", "d", "g"},
				{"b", "d", "g"},
				{"c", "d", "g"},
				{"a", "e", "g"},
				{"b", "e", "g"},
				{"c", "e", "g"},
				{"a", "f", "g"},
				{"b", "f", "g"},
				{"c", "f", "g"},
				{"a", "d", "h"},
				{"b", "d", "h"},
				{"c", "d", "h"},
				{"a", "e", "h"},
				{"b", "e", "h"},
				{"c", "e", "h"},
				{"a", "f", "h"},
				{"b", "f", "h"},
				{"c", "f", "h"},
				{"a", "d", "i"},
				{"b", "d", "i"},
				{"c", "d", "i"},
				{"a", "e", "i"},
				{"b", "e", "i"},
				{"c", "e", "i"},
				{"a", "f", "i"},
				{"b", "f", "i"},
				{"c", "f", "i"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			file := CSVFile{
				Header:        []string{"col_1", "col_2", "col_3"},
				UniqueColumns: c.uniqueColumns,
				Lines: [][]string{
					{"a", "d", "g"},
					{"b", "d", "g"},
					{"c", "d", "g"},
					{"a", "e", "g"},
					{"b", "e", "g"},
					{"c", "e", "g"},
					{"a", "f", "g"},
					{"b", "f", "g"},
					{"c", "f", "g"},
					{"a", "d", "h"},
					{"b", "d", "h"},
					{"c", "d", "h"},
					{"a", "e", "h"},
					{"b", "e", "h"},
					{"c", "e", "h"},
					{"a", "f", "h"},
					{"b", "f", "h"},
					{"c", "f", "h"},
					{"a", "d", "i"},
					{"b", "d", "i"},
					{"c", "d", "i"},
					{"a", "e", "i"},
					{"b", "e", "i"},
					{"c", "e", "i"},
					{"a", "f", "i"},
					{"b", "f", "i"},
					{"c", "f", "i"},
				},
			}

			act := file.Unique()

			assert.Equal(t, c.exp, act)
		})
	}
}
