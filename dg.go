package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"runtime/pprof"
	"strings"
	"text/template"
	"time"

	"github.com/codingconcepts/dg/internal/pkg/generator"
	"github.com/codingconcepts/dg/internal/pkg/model"
	"github.com/codingconcepts/dg/internal/pkg/source"
	"github.com/codingconcepts/dg/internal/pkg/ui"
)

var (
	version string
)

func main() {
	log.SetFlags(0)

	configPath := flag.String("c", "", "the absolute or relative path to the config file")
	outputDir := flag.String("o", ".", "the absolute or relative path to the output dir")
	createImports := flag.String("i", "", "write import statements to file")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	versionFlag := flag.Bool("version", false, "display the current version number")
	flag.Parse()

	if *cpuprofile != "" {
		defer launchProfiler(*cpuprofile)()
	}

	if *versionFlag {
		fmt.Println(version)
		return
	}

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	tt := ui.TimeTracker(os.Stdout, realClock{}, 40)
	defer tt(time.Now(), "done")

	c, err := loadConfig(*configPath, tt)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	files := make(map[string]model.CSVFile)

	if err = loadInputs(c, path.Dir(*configPath), tt, files); err != nil {
		log.Fatalf("error loading inputs: %v", err)
	}

	if err = generateTables(c, tt, files); err != nil {
		log.Fatalf("error generating tables: %v", err)
	}

	if err := writeFiles(*outputDir, files, tt); err != nil {
		log.Fatalf("error writing csv files: %v", err)
	}

	if *createImports != "" {
		if err := writeImports(*outputDir, *createImports, c, files, tt); err != nil {
			log.Fatalf("error writing import statements: %v", err)
		}
	}
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

func loadInputs(c model.Config, configDir string, tt ui.TimerFunc, files map[string]model.CSVFile) error {
	defer tt(time.Now(), "loaded data sources")

	for _, input := range c.Inputs {
		if err := loadInput(input, configDir, tt, files); err != nil {
			return fmt.Errorf("loading input for %q: %w", input.Name, err)
		}
	}

	return nil
}

func loadInput(input model.Input, configDir string, tt ui.TimerFunc, files map[string]model.CSVFile) error {
	defer tt(time.Now(), fmt.Sprintf("loaded data source: %s", input.Name))

	switch input.Type {
	case "csv":
		var s model.SourceCSV
		if err := input.Source.UnmarshalFunc(&s); err != nil {
			return fmt.Errorf("parsing csv source for %s: %w", input.Name, err)
		}

		if err := source.LoadCSVSource(input.Name, configDir, s, files); err != nil {
			return fmt.Errorf("loading csv for %s: %w", input.Name, err)
		}
	}

	return nil
}

func generateTables(c model.Config, tt ui.TimerFunc, files map[string]model.CSVFile) error {
	defer tt(time.Now(), "generated all tables")

	for _, table := range c.Tables {
		if err := generateTable(table, files, tt); err != nil {
			return fmt.Errorf("generating csv file for %q: %w", table.Name, err)
		}
	}

	return nil
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
				return fmt.Errorf("running ref process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "gen":
			var p model.ProcessorGenerator
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing each process for %s: %w", col.Name, err)
			}

			if err := generator.GenerateGenColumn(t, col, p, files); err != nil {
				return fmt.Errorf("running gen process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "set":
			var p model.ProcessorSet
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing set process for %s.%s: %w", t.Name, col.Name, err)
			}

			if err := generator.GenerateSetColumn(t, col, p, files); err != nil {
				return fmt.Errorf("running set process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "inc":
			var p model.ProcessorInc
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing each process for %s: %w", col.Name, err)
			}

			if err := generator.GenerateIncColumn(t, col, p, files); err != nil {
				return fmt.Errorf("running inc process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "range":
			var p model.ProcessorRange
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing range process for %s: %w", col.Name, err)
			}

			if err := generator.GenerateRangeColumn(t, col, p, files); err != nil {
				return fmt.Errorf("running range process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "match":
			var p model.ProcessorMatch
			if err := col.Processor.UnmarshalFunc(&p); err != nil {
				return fmt.Errorf("parsing match process for %s: %w", col.Name, err)
			}

			if err := generator.GenerateMatchColumn(t, col, p, files); err != nil {
				return fmt.Errorf("running match process for %s.%s: %w", t.Name, col.Name, err)
			}
		}
	}

	file := files[t.Name]
	if len(file.UniqueColumns) > 0 {
		file.Lines = generator.Transpose(file.Lines)
		file.Lines = file.Unique()
		file.Lines = generator.Transpose(file.Lines)
	}
	files[t.Name] = file

	return nil
}

func writeFiles(outputDir string, cfs map[string]model.CSVFile, tt ui.TimerFunc) error {
	defer tt(time.Now(), "wrote all csvs")

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for name, file := range cfs {
		if !file.Output {
			continue
		}

		if err := writeFile(outputDir, name, file, tt); err != nil {
			return fmt.Errorf("writing file %q: %w", file.Name, err)
		}
	}

	return nil
}

func writeFile(outputDir, name string, cf model.CSVFile, tt ui.TimerFunc) error {
	defer tt(time.Now(), fmt.Sprintf("wrote csv: %s", name))

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

	cf.Lines = generator.Transpose(cf.Lines)

	if err = writer.WriteAll(cf.Lines); err != nil {
		return fmt.Errorf("writing csv lines for %q: %w", name, err)
	}

	writer.Flush()
	return nil
}

func writeImports(outputDir, name string, c model.Config, files map[string]model.CSVFile, tt ui.TimerFunc) error {
	defer tt(time.Now(), fmt.Sprintf("wrote imports: %s", name))

	importTmpl := template.Must(template.New("import").
		Funcs(template.FuncMap{"join": strings.Join}).
		Parse(`IMPORT INTO {{.Name}} (
	{{ join .Header ", " }}
)
CSV DATA (
    '.../{{.Name}}.csv'
)
WITH skip='1', nullif = '', allow_quoted_null;

`),
	)

	fullPath := path.Join(outputDir, name)
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("creating csv file %q: %w", name, err)
	}
	defer file.Close()

	// Iterate through the tables in the config file, so the imports are in the right order.
	for _, table := range c.Tables {
		csv := files[table.Name]
		if !csv.Output {
			continue
		}

		if err := importTmpl.Execute(file, csv); err != nil {
			return fmt.Errorf("writing import statement for %q: %w", name, err)
		}
	}

	return nil
}

func launchProfiler(cpuprofile string) func() {
	f, err := os.Create(cpuprofile)
	if err != nil {
		log.Fatalf("creating file for profiler: %v", err)
	}
	pprof.StartCPUProfile(f)

	return func() {
		pprof.StopCPUProfile()
	}
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func (realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}
