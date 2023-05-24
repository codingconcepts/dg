package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
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

func loadConfig(filename string) (config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return config{}, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var c config
	if err := yaml.NewDecoder(file).Decode(&c); err != nil {
		return config{}, fmt.Errorf("parsing file: %w", err)
	}

	return c, nil
}

func generateTables(c config) (map[string]csvFile, error) {
	files := make(map[string]csvFile)
	for _, table := range c {
		if err := generateTable(table, files); err != nil {
			return nil, fmt.Errorf("generating csv file for %q: %w", table.Name, err)
		}
	}

	return files, nil
}

func generateTable(t table, files map[string]csvFile) error {
	// Create the Cartesian product of any each types first.
	if err := generateEachColumns(t, files); err != nil {
		return fmt.Errorf("generating each columns: %w", err)
	}

	for _, col := range t.Columns {
		switch col.Type {
		case "ref":
			if err := generateRefColumn(t, col, files); err != nil {
				return fmt.Errorf("parsing ref process for %s.%s: %w", t.Name, col.Name, err)
			}

		case "gen":
			if err := generateGenColumn(t, col, files); err != nil {
				return fmt.Errorf("parsing gen process for %s.%s: %w", t.Name, col.Name, err)
			}
		}
	}

	return nil
}

func generateEachColumns(t table, files map[string]csvFile) error {
	cols := lo.Filter(t.Columns, func(c column, _ int) bool {
		return c.Type == "each"
	})

	if len(cols) == 0 {
		return nil
	}

	var preCartesian [][]string
	for _, col := range cols {
		var ptc processorTableColumn
		if err := col.Processor.unmarshal(&ptc); err != nil {
			return fmt.Errorf("parsing each process for %s.%s: %w", t.Name, col.Name, err)
		}

		srcTable := files[ptc.Table]
		srcColumn := ptc.Column
		srcColumnIndex := lo.IndexOf(srcTable.header, srcColumn)

		preCartesian = append(preCartesian, srcTable.lines[srcColumnIndex])
	}

	// Compute Cartesian product of all columns.
	cartesianColumns := transpose(cartesianProduct(preCartesian...))

	// Add the header
	if _, ok := files[t.Name]; !ok {
		files[t.Name] = csvFile{
			name: t.Name,
		}
	}

	for i, col := range cartesianColumns {
		foundTable := files[t.Name]
		foundTable.header = append(foundTable.header, cols[i].Name)
		foundTable.lines = append(foundTable.lines, col)
		files[t.Name] = foundTable
	}

	return nil
}

func generateRefColumn(t table, c column, files map[string]csvFile) error {
	var ptc processorTableColumn
	if err := c.Processor.unmarshal(&ptc); err != nil {
		return fmt.Errorf("parsing ref process for %s.%s: %w", t.Name, c.Name, err)
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	table, ok := files[ptc.Table]
	if !ok {
		return fmt.Errorf("missing table %q for ref lookup", ptc.Table)
	}

	colIndex := lo.IndexOf(table.header, ptc.Column)
	column := table.lines[colIndex]

	var lines []string
	for i := 0; i < t.Count; i++ {
		lines = append(lines, column[rand.Intn(len(column))])
	}

	// Add the header
	if _, ok := files[t.Name]; !ok {
		files[t.Name] = csvFile{
			name: t.Name,
		}
	}

	foundTable := files[t.Name]
	foundTable.header = append(foundTable.header, c.Name)
	foundTable.lines = append(foundTable.lines, lines)
	files[t.Name] = foundTable

	return nil
}

func generateGenColumn(t table, c column, files map[string]csvFile) error {
	var pg processorGenerator
	if err := c.Processor.unmarshal(&pg); err != nil {
		return fmt.Errorf("parsing each process for %s: %w", c.Name, err)
	}

	if t.Count == 0 {
		t.Count = len(lo.MaxBy(files[t.Name].lines, func(a, b []string) bool {
			return len(a) > len(b)
		}))
	}

	var line []string
	for i := 0; i < t.Count; i++ {
		line = append(line, replacePlaceholders(pg))
	}

	// Add the header
	if _, ok := files[t.Name]; !ok {
		files[t.Name] = csvFile{
			name: t.Name,
		}
	}

	foundTable := files[t.Name]
	foundTable.header = append(foundTable.header, c.Name)
	foundTable.lines = append(foundTable.lines, line)
	files[t.Name] = foundTable

	return nil
}

func replacePlaceholders(pg processorGenerator) string {
	r := rand.Intn(100)
	if r < pg.NullPercentage {
		return ""
	}

	s := pg.Value
	for k, v := range replacements {
		if strings.Contains(s, k) {
			s = strings.ReplaceAll(s, k, v())
		}
	}

	return s
}

func writeFiles(outputDir string, cfs map[string]csvFile) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for name, file := range cfs {
		if err := writeFile(outputDir, name, file); err != nil {
			return fmt.Errorf("writing file %q: %w", file.name, err)
		}
	}

	return nil
}

func writeFile(outputDir, name string, cf csvFile) error {
	fullPath := path.Join(outputDir, fmt.Sprintf("%s.csv", name))
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("creating csv file %q: %w", name, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err = writer.Write(cf.header); err != nil {
		return fmt.Errorf("writing csv header for %q: %w", name, err)
	}

	lines := transpose(cf.lines)
	if err = writer.WriteAll(lines); err != nil {
		return fmt.Errorf("writing csv lines for %q: %w", name, err)
	}

	writer.Flush()
	return nil
}

func cartesianProduct(a ...[]string) (c [][]string) {
	if len(a) == 0 {
		return [][]string{nil}
	}
	last := len(a) - 1
	l := cartesianProduct(a[:last]...)
	for _, e := range a[last] {
		for _, p := range l {
			c = append(c, append(p, e))
		}
	}
	return
}

func transpose(m [][]string) [][]string {
	r := make([][]string, len(m[0]))
	for x := range r {
		r[x] = make([]string, len(m))
	}
	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}

type config []table

type table struct {
	Name    string   `yaml:"table"`
	Count   int      `yaml:"count"`
	Columns []column `yaml:"columns"`
}

type column struct {
	Name      string     `yaml:"name"`
	Type      string     `yaml:"type"`
	Processor rawMessage `yaml:"processor"`
}

type processorTableColumn struct {
	Table  string `yaml:"table"`
	Column string `yaml:"column"`
}

type processorGenerator struct {
	Value          string `yaml:"value"`
	NullPercentage int    `yaml:"null_percentage"`
}

type rawMessage struct {
	unmarshal func(interface{}) error
}

type csvFile struct {
	name   string
	header []string
	lines  [][]string
}

func (msg *rawMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}

func (msg *rawMessage) Unmarshal(v any) error {
	return msg.unmarshal(v)
}

var (
	version string

	replacements = map[string]func() string{
		"latitude":               func() string { return strconv.FormatFloat(faker.Latitude(), 'f', -1, 64) },
		"longitude":              func() string { return strconv.FormatFloat(faker.Longitude(), 'f', -1, 64) },
		"address":                func() string { return faker.GetRealAddress().Address },
		"city":                   func() string { return faker.GetRealAddress().City },
		"state":                  func() string { return faker.GetRealAddress().State },
		"postcode":               func() string { return faker.GetRealAddress().PostalCode },
		"unix_time":              func() string { return strconv.FormatInt(faker.UnixTime(), 10) },
		"date":                   func() string { return faker.Date() },
		"time":                   func() string { return faker.TimeString() },
		"month_name":             func() string { return faker.MonthName() },
		"year":                   func() string { return faker.YearString() },
		"day_of_week":            func() string { return faker.DayOfWeek() },
		"day_of_month":           func() string { return faker.DayOfMonth() },
		"timestamp":              func() string { return faker.Timestamp() },
		"century":                func() string { return faker.Century() },
		"timezone":               func() string { return faker.Timezone() },
		"timeperiod":             func() string { return faker.Timeperiod() },
		"email":                  func() string { return faker.Email() },
		"mac_address":            func() string { return faker.MacAddress() },
		"domain_name":            func() string { return faker.DomainName() },
		"url":                    func() string { return faker.URL() },
		"user_name":              func() string { return faker.Username() },
		"ipv4":                   func() string { return faker.IPv4() },
		"ipv6":                   func() string { return faker.IPv6() },
		"password":               func() string { return faker.Password() },
		"word":                   func() string { return faker.Word() },
		"sentence":               func() string { return faker.Sentence() },
		"paragraph":              func() string { return faker.Paragraph() },
		"cc_type":                func() string { return faker.CCType() },
		"cc_number":              func() string { return faker.CCNumber() },
		"currency":               func() string { return faker.Currency() },
		"amount_with_currency":   func() string { return faker.AmountWithCurrency() },
		"title_male":             func() string { return faker.TitleMale() },
		"title_female":           func() string { return faker.TitleFemale() },
		"first_name":             func() string { return faker.FirstName() },
		"first_name_male":        func() string { return faker.FirstNameMale() },
		"first_name_female":      func() string { return faker.FirstNameFemale() },
		"last_name":              func() string { return faker.LastName() },
		"name":                   func() string { return faker.Name() },
		"phone_number":           func() string { return faker.Phonenumber() },
		"toll_free_phone_number": func() string { return faker.TollFreePhoneNumber() },
		"e164_phone_number":      func() string { return faker.E164PhoneNumber() },
		"uuid_hyphen":            func() string { return faker.UUIDHyphenated() },
		"uuid":                   func() string { return faker.UUIDDigit() },
		"bool":                   func() string { return strconv.FormatBool(rand.Int()%2 == 0) },
		"int8":                   func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt8), 10) },
		"int16":                  func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt16), 10) },
		"int32":                  func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt32), 10) },
		"int64":                  func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt64), 10) },
	}
)
