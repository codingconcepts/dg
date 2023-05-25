package model

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

// Config represents the entire contents of a config file.
type Config []Table

// Table represents the instructions to create one CSV file.
type Table struct {
	Name    string   `yaml:"table"`
	Count   int      `yaml:"count"`
	Columns []Column `yaml:"columns"`
}

// Column represents the instructions to populate one CSV file column.
type Column struct {
	Name      string     `yaml:"name"`
	Type      string     `yaml:"type"`
	Processor RawMessage `yaml:"processor"`
}

// Load config from a file
func LoadConfig(r io.Reader) (Config, error) {
	var c Config
	if err := yaml.NewDecoder(r).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("parsing file: %w", err)
	}

	return c, nil
}
