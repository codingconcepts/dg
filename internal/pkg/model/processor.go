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

// ProcessorSet provides additional context to a set column.
type ProcessorSet struct {
	Values  []string `yaml:"values"`
	Weights []int    `yaml:"weights"`
}

// ProcessorConst provides additional context to a const column.
type ProcessorConst struct {
	Values []string `yaml:"values"`
}

// ProcessorInc provides additional context to an inc column.
type ProcessorInc struct {
	Start  int    `yaml:"start"`
	Format string `yaml:"format"`
}

func (pi ProcessorInc) GetFormat() string {
	return pi.Format
}

// ProcessorMatch provides additional context to a match column.
type ProcessorMatch struct {
	SourceTable  string `yaml:"source_table"`
	SourceColumn string `yaml:"source_column"`
	SourceValue  string `yaml:"source_value"`
	MatchColumn  string `yaml:"match_column"`
}
