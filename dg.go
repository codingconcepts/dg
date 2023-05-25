package main

import (
	"dg/internal/pkg/generator"
	"dg/internal/pkg/model"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
)

var (
	version string
)

func main() {
	configPath := flag.String("c", "", "the absolute or relative path to the config file")
	outputDir := flag.String("o", ".", "the absolute or relative path to the output dir")
	versionFlag := flag.Bool("version", false, "display the current version number")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		return
	}

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	c, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	files, err := generateTables(c)
	if err != nil {
		log.Fatalf("error generating tables: %v", err)
	}

	if err := writeFiles(*outputDir, files); err != nil {
		log.Fatalf("error writing csv files: %v", err)
	}
}

func loadConfig(filename string) (model.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return model.Config{}, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	return model.LoadConfig(file)
}

func generateTables(c model.Config) (map[string]model.CSVFile, error) {
	files := make(map[string]model.CSVFile)
	for _, table := range c {
		if err := generateTable(table, files); err != nil {
			return nil, fmt.Errorf("generating csv file for %q: %w", table.Name, err)
		}
	}

	return files, nil
}

func generateTable(t model.Table, files map[string]model.CSVFile) error {
	// Create the Cartesian product of any each types first.
	if err := generator.GenerateEachColumns(t, files); err != nil {
		return fmt.Errorf("generating each columns: %w", err)
	}

	for _, col := range t.Columns {
		switch col.Type {
		case "ref":
			if err := generator.GenerateRefColumn(t, col, files); err != nil {
				return fmt.Errorf("parsing ref process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "gen":
			if err := generator.GenerateGenColumn(t, col, files); err != nil {
				return fmt.Errorf("parsing gen process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "set":
			if err := generator.GenerateSetColumn(t, col, files); err != nil {
				return fmt.Errorf("parsing set process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "inc":
			if err := generator.GenerateIncColumn(t, col, files); err != nil {
				return fmt.Errorf("parsing inc process for %s.%s: %w", t.Name, col.Name, err)
			}
		}
	}

	return nil
}

func writeFiles(outputDir string, cfs map[string]model.CSVFile) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for name, file := range cfs {
		if err := writeFile(outputDir, name, file); err != nil {
			return fmt.Errorf("writing file %q: %w", file.Name, err)
		}
	}

	return nil
}

func writeFile(outputDir, name string, cf model.CSVFile) error {
	fullPath := path.Join(outputDir, fmt.Sprintf("%s.csv", name))
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("creating csv file %q: %w", name, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err = writer.Write(cf.Header); err != nil {
		return fmt.Errorf("writing csv header for %q: %w", name, err)
	}

	lines := generator.Transpose(cf.Lines)
	if err = writer.WriteAll(lines); err != nil {
		return fmt.Errorf("writing csv lines for %q: %w", name, err)
	}

	writer.Flush()
	return nil
}
