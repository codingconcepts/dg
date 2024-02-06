package model

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// Config represents the entire contents of a config file.
type Config struct {
	Tables []Table `yaml:"tables"`
	Inputs []Input `yaml:"inputs"`
}

// Table represents the instructions to create one CSV file.
type Table struct {
	Name          string   `yaml:"name"`
	Count         int      `yaml:"count"`
	Suppress      bool     `yaml:"suppress"`
	UniqueColumns []string `yaml:"unique_columns"`
	Columns       []Column `yaml:"columns"`
}

// Column represents the instructions to populate one CSV file column.
type Column struct {
	Name      string     `yaml:"name"`
	Type      string     `yaml:"type"`
	Suppress  bool       `yaml:"suppress"`
	Generator RawMessage `yaml:"processor"`
}

// Input represents a data source provided by the user.
type Input struct {
	Name   string     `yaml:"name"`
	Type   string     `yaml:"type"`
	Source RawMessage `yaml:"source"`
}

// Load config from a file
func LoadConfig(r io.Reader) (Config, error) {
	var c Config
	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("parsing file: %w", err)
	}

	return c, nil
}
