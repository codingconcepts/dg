package model

// FileType describes the various types of files dg can work with.
type FileType string

const (
	FileTypeOutput FileType = "output"
	FileTypeInput  FileType = "input"
)

// CSVFile represents the content of a CSV file.
type CSVFile struct {
	Name   string
	Type   FileType
	Header []string
	Lines  [][]string
}
