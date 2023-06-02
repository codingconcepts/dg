package main

import (
	"dg/internal/pkg/generator"
	"dg/internal/pkg/model"
	"dg/internal/pkg/ui"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

var (
	version string
)

func main() {
	log.SetFlags(0)

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

	tt := ui.TimeTracker(os.Stdout, realClock{}, 40)

	c, err := loadConfig(*configPath, tt)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	files, err := generateTables(c, tt)
	if err != nil {
		log.Fatalf("error generating tables: %v", err)
	}

	if err := writeFiles(*outputDir, files, tt); err != nil {
		log.Fatalf("error writing csv files: %v", err)
	}
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func loadConfig(filename string, tt ui.TimerFunc) (model.Config, error) {
	defer tt(time.Now(), "loaded config file")

	file, err := os.Open(filename)
	if err != nil {
		return model.Config{}, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	return model.LoadConfig(file)
}

func generateTables(c model.Config, tt ui.TimerFunc) (map[string]model.CSVFile, error) {
	defer tt(time.Now(), "generated all tables")

	files := make(map[string]model.CSVFile)
	for _, table := range c {
		if err := generateTable(table, files, tt); err != nil {
			return nil, fmt.Errorf("generating csv file for %q: %w", table.Name, err)
		}
	}

	return files, nil
}

func generateTable(t model.Table, files map[string]model.CSVFile, tt ui.TimerFunc) error {
	defer tt(time.Now(), fmt.Sprintf("generated table: %s", t.Name))

	// Create the Cartesian product of any each types first.
	if err := generator.GenerateEachColumns(t, files); err != nil {
		return fmt.Errorf("generating each columns: %w", err)
	}

	for _, col := range t.Columns {
		switch col.Type {
		case "ref":
			var p model.ProcessorTableColumn
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing ref process for %s.%s: %w", t.Name, col.Name, err)
			}

			if err := generator.GenerateRefColumn(t, col, p, files); err != nil {
				return fmt.Errorf("parsing ref process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "gen":
			var p model.ProcessorGenerator
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing each process for %s: %w", col.Name, err)
			}

			if err := generator.GenerateGenColumn(t, col, p, files); err != nil {
				return fmt.Errorf("parsing gen process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "set":
			var p model.ProcessorSet
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing set process for %s.%s: %w", t.Name, col.Name, err)
			}

			if err := generator.GenerateSetColumn(t, col, p, files); err != nil {
				return fmt.Errorf("parsing set process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "inc":
			var p model.ProcessorInc
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing each process for %s: %w", col.Name, err)
			}

			if err := generator.GenerateIncColumn(t, col, p, files); err != nil {
				return fmt.Errorf("parsing inc process for %s.%s: %w", t.Name, col.Name, err)
			}
		}
	}

	return nil
}

func writeFiles(outputDir string, cfs map[string]model.CSVFile, tt ui.TimerFunc) error {
	defer tt(time.Now(), "wrote all csvs")

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for name, file := range cfs {
		if err := writeFile(outputDir, name, file, tt); err != nil {
			return fmt.Errorf("writing file %q: %w", file.Name, err)
		}
	}

	return nil
}

func writeFile(outputDir, name string, cf model.CSVFile, tt ui.TimerFunc) error {
	defer tt(time.Now(), fmt.Sprintf("generated csv: %s", name))

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
