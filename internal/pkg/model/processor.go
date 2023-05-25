package model

// ProcessorTableColumn provides additional context to an each or ref column.
type ProcessorTableColumn struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

// ProcessorGenerator provides additional context to a gen column.
type ProcessorGenerator struct {
	Value          string `yaml:"value"`
	NullPercentage int    `yaml:"null_percentage"`
	Format         string `yaml:"format"`
}

func (pg ProcessorGenerator) GetFormat() string {
	return pg.Format
}

// ProcessorSet provides additional context to a set column.
type ProcessorSet struct {
	Values  []string `yaml:"values"`
	Weights []int    `yaml:"weights"`
}

// ProcessorInc provides additional context to an inc column.
type ProcessorInc struct {
	Start  int    `yaml:"start"`
	Format string `yaml:"format"`
}

func (pi ProcessorInc) GetFormat() string {
	return pi.Format
}
