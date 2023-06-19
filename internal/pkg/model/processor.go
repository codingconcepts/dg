package model

// ProcessorTableColumn provides additional context to an each or ref column.
type ProcessorTableColumn struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}
