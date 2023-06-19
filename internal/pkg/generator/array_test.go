package generator

import (
	"strings"
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
		{
			name: "small array big array",
			input: [][]string{
				{"a", "b"},
				{"d", "e", "f"},
			},
			output: [][]string{
				{"a", "d"},
				{"b", "d"},
				{"a", "e"},
				{"b", "e"},
				{"a", "f"},
				{"b", "f"},
			},
		},
		{
			name: "big array small array",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e"},
			},
			output: [][]string{
				{"a", "d"},
				{"b", "d"},
				{"c", "d"},
				{"a", "e"},
				{"b", "e"},
				{"c", "e"},
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
		{
			name: "first input bigger than second",
			input: [][]string{
				{"a", "b", "c", "1"},
				{"d", "e", "f"},
			},
			output: [][]string{
				{"a", "d"},
				{"b", "e"},
				{"c", "f"},
				{"1", ""},
			},
		},
		{
			name: "second input bigger than first",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f", "2"},
			},
			output: [][]string{
				{"a", "d"},
				{"b", "e"},
				{"c", "f"},
				{"", "2"},
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

func BenchmarkCartesianProduct(b *testing.B) {
	cases := []struct {
		name  string
		input [][]string
	}{
		{
			name: "single array",
			input: [][]string{
				{"a", "b", "c"},
			},
		},
		{
			name: "two small arrays",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
			},
		},
		{
			name: "three small arrays",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
				{"g", "h", "i"},
			},
		},
		{
			name: "small array and big array",
			input: [][]string{
				{"a", "b", "c"},
				strings.Split(strings.Repeat("d", 1000), ""),
			},
		},
		{
			name: "big array and small array",
			input: [][]string{
				strings.Split(strings.Repeat("a", 1000), ""),
				{"d", "e", "f"},
			},
		},
		{
			name: "big arrays",
			input: [][]string{
				strings.Split(strings.Repeat("a", 1000), ""),
				strings.Split(strings.Repeat("d", 1000), ""),
			},
		},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CartesianProduct(c.input...)
			}
		})
	}
}

func BenchmarkTranspose(b *testing.B) {
	cases := []struct {
		name  string
		input [][]string
	}{
		{
			name: "single array",
			input: [][]string{
				{"a", "b", "c"},
			},
		},
		{
			name: "multiple small arrays",
			input: [][]string{
				{"a", "b", "c"},
				{"d", "e", "f"},
			},
		},
		{
			name: "small array and big array",
			input: [][]string{
				{"a", "b", "c"},
				strings.Split(strings.Repeat("d", 1000), ""),
			},
		},
		{
			name: "big array and small array",
			input: [][]string{
				strings.Split(strings.Repeat("a", 1000), ""),
				{"d", "e", "f"},
			},
		},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				CartesianProduct(c.input...)
			}
		})
	}
}
