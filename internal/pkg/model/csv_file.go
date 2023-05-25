package model

// CSVFile represents the content of a CSV file.
type CSVFile struct {
	Name   string
	Header []string
	Lines  [][]string
}
