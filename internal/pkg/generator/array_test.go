package generator

import (
	"strings"
	"testing"

	"github.com/codingconcepts/dg/internal/pkg/model"
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

func TestGenerateDateSlice(t *testing.T) {
	cases := []struct {
		name     string
		from     string
		to       string
		format   string
		count    int
		step     string
		expSlice []string
		expError string
	}{
		{
			name:     "no count or step",
			expError: "either a count or a step must be provided to a date range generator",
		},
		{
			name:   "count",
			count:  10,
			from:   "2023-01-01",
			to:     "2023-01-10",
			format: "2006-01-02",
			expSlice: []string{
				"2023-01-01", "2023-01-01", "2023-01-02", "2023-01-03", "2023-01-04", "2023-01-05", "2023-01-06", "2023-01-07", "2023-01-08", "2023-01-09",
			},
		},
		{
			name:   "step",
			step:   "24h",
			from:   "2023-01-10",
			to:     "2023-01-20",
			format: "2006-01-02",
			expSlice: []string{
				"2023-01-10", "2023-01-11", "2023-01-12", "2023-01-13", "2023-01-14", "2023-01-15", "2023-01-16", "2023-01-17", "2023-01-18", "2023-01-19",
			},
		},
		{
			name:     "invalid format",
			count:    10,
			from:     "2023-01-01",
			to:       "2023-01-10",
			format:   "abc",
			expError: `parsing from date: parsing time "2023-01-01" as "abc": cannot parse "2023-01-01" as "abc"`,
		},
		{
			name:     "invalid from date",
			count:    10,
			from:     "abc",
			format:   "2006-01-02",
			to:       "2023-01-10",
			expError: `parsing from date: parsing time "abc" as "2006-01-02": cannot parse "abc" as "2006"`,
		},
		{
			name:     "invalid to date",
			count:    10,
			from:     "2023-01-01",
			to:       "abc",
			format:   "2006-01-02",
			expError: `parsing to date: parsing time "abc" as "2006-01-02": cannot parse "abc" as "2006"`,
		},
		{
			name:     "invalid step",
			step:     "abc",
			from:     "2023-01-01",
			to:       "2023-01-10",
			format:   "2006-01-02",
			expError: `parsing step: time: invalid duration "abc"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pr := model.ProcessorRange{
				From:   c.from,
				To:     c.to,
				Format: c.format,
				Step:   c.step,
			}

			actSlice, actErr := generateDateSlice(pr, c.count)
			if c.expError != "" {
				assert.Equal(t, c.expError, actErr.Error())
				return
			}

			assert.Equal(t, c.expSlice, actSlice)
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
