package model

// ProcessorTableColumn provides additional context to an each or ref column.
type ProcessorTableColumn struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

// ProcessorRange provides additional context to a range column.
type ProcessorRange struct {
	Type   string `yaml:"type"`
	From   string `yaml:"from"`
	To     string `yaml:"to"`
	Step   string `yaml:"step"`
	Format string `yaml:"format"`
}

// ProcessorMatch provides additional context to a match column.
type ProcessorMatch struct {
	SourceTable  string `yaml:"source_table"`
	SourceColumn string `yaml:"source_column"`
	SourceValue  string `yaml:"source_value"`
	MatchColumn  string `yaml:"match_column"`
}
