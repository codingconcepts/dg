package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCartesianProduct(t *testing.T) {
	cases := []struct {
		name   string
		input  [][]string
		output [][]string
	}{
		{
			name: "single input",
			input: [][]string{
				{"a", "b", "c"},
			},
			output: [][]string{
				{"a"}, {"b"}, {"c"},
			},
		},
		{
			name: "multiple input",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
			},
			output: [][]string{
				{"a", "d"},
				{"b", "d"},
				{"c", "d"},
				{"a", "e"},
				{"b", "e"},
				{"c", "e"},
				{"a", "f"},
				{"b", "f"},
				{"c", "f"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := CartesianProduct(c.input...)
			assert.Equal(t, c.output, actual)
		})
	}
}

func TestTranspose(t *testing.T) {
	cases := []struct {
		name   string
		input  [][]string
		output [][]string
	}{
		{
			name: "single input",
			input: [][]string{
				{"a", "b", "c"},
			},
			output: [][]string{
				{"a"}, {"b"}, {"c"},
			},
		},
		{
			name: "multiple input",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
			},
			output: [][]string{
				{"a", "d"},
				{"b", "e"},
				{"c", "f"},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := Transpose(c.input)
			assert.Equal(t, c.output, actual)
		})
	}
}
