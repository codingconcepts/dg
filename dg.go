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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

var (
	refExtractor = regexp.MustCompile(`ref (\w+) (\w+)`)

	replacements = map[string]func() string{
		"${gen latitude}":               func() string { return strconv.FormatFloat(faker.Latitude(), 'f', -1, 64) },
		"${gen longitude}":              func() string { return strconv.FormatFloat(faker.Longitude(), 'f', -1, 64) },
		"${gen address}":                func() string { return faker.GetRealAddress().Address },
		"${gen city}":                   func() string { return faker.GetRealAddress().City },
		"${gen state}":                  func() string { return faker.GetRealAddress().State },
		"${gen postcode}":               func() string { return faker.GetRealAddress().PostalCode },
		"${gen unix_time}":              func() string { return strconv.FormatInt(faker.UnixTime(), 10) },
		"${gen date}":                   func() string { return faker.Date() },
		"${gen time}":                   func() string { return faker.TimeString() },
		"${gen month_name}":             func() string { return faker.MonthName() },
		"${gen year}":                   func() string { return faker.YearString() },
		"${gen day_of_week}":            func() string { return faker.DayOfWeek() },
		"${gen day_of_month}":           func() string { return faker.DayOfMonth() },
		"${gen timestamp}":              func() string { return faker.Timestamp() },
		"${gen century}":                func() string { return faker.Century() },
		"${gen timezone}":               func() string { return faker.Timezone() },
		"${gen timeperiod}":             func() string { return faker.Timeperiod() },
		"${gen email}":                  func() string { return faker.Email() },
		"${gen mac_address}":            func() string { return faker.MacAddress() },
		"${gen domain_name}":            func() string { return faker.DomainName() },
		"${gen url}":                    func() string { return faker.URL() },
		"${gen user_name}":              func() string { return faker.Username() },
		"${gen ipv4}":                   func() string { return faker.IPv4() },
		"${gen ipv6}":                   func() string { return faker.IPv6() },
		"${gen password}":               func() string { return faker.Password() },
		"${gen word}":                   func() string { return faker.Word() },
		"${gen sentence}":               func() string { return faker.Sentence() },
		"${gen paragraph}":              func() string { return faker.Paragraph() },
		"${gen cc_type}":                func() string { return faker.CCType() },
		"${gen cc_number}":              func() string { return faker.CCNumber() },
		"${gen currency}":               func() string { return faker.Currency() },
		"${gen amount_with_currency}":   func() string { return faker.AmountWithCurrency() },
		"${gen title_male}":             func() string { return faker.TitleMale() },
		"${gen title_female}":           func() string { return faker.TitleFemale() },
		"${gen first_name}":             func() string { return faker.FirstName() },
		"${gen first_name_male}":        func() string { return faker.FirstNameMale() },
		"${gen first_name_female}":      func() string { return faker.FirstNameFemale() },
		"${gen last_name}":              func() string { return faker.LastName() },
		"${gen name}":                   func() string { return faker.Name() },
		"${gen phone_number}":           func() string { return faker.Phonenumber() },
		"${gen toll_free_phone_number}": func() string { return faker.TollFreePhoneNumber() },
		"${gen e164_phone_number}":      func() string { return faker.E164PhoneNumber() },
		"${gen uuid_hyphen}":            func() string { return faker.UUIDHyphenated() },
		"${gen uuid}":                   func() string { return faker.UUIDDigit() },
		"${gen bool}":                   func() string { return strconv.FormatBool(rand.Int()%2 == 0) },
		"${gen int8}":                   func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt8), 10) },
		"${gen int16}":                  func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt16), 10) },
		"${gen int32}":                  func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt32), 10) },
		"${gen int64}":                  func() string { return strconv.FormatInt(rand.Int63n(math.MaxInt64), 10) },
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	configPath := flag.String("c", "", "absolute or relative path to config file")
	outputDir := flag.String("o", ".", "absolute or relative path of output directory")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		os.Exit(2)
	}

	c, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	csvFiles, err := generateTables(c)
	if err != nil {
		log.Fatalf("error generating data for tables: %v", err)
	}

	if err = writeFiles(*outputDir, csvFiles); err != nil {
		log.Fatalf("error writing csv files: %v", err)
	}
}

type config []table

type table struct {
	Name    string   `yaml:"table"`
	Count   int      `yaml:"count"`
	Columns []column `yaml:"columns"`
	Foreach string   `yaml:"foreach"`
}

type column struct {
	Name           string `yaml:"name"`
	Value          string `yaml:"value"`
	NullPercentage int    `yaml:"null_percentage"`
}

type csvFile struct {
	name   string
	header []string
	lines  [][]string
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
		file, err := generateTable(table, files)
		if err != nil {
			return nil, fmt.Errorf("generating csv file for %q: %w", table.Name, err)
		}

		files[table.Name] = file
	}

	return files, nil
}

func generateTable(t table, files map[string]csvFile) (csvFile, error) {
	if t.Foreach != "" {
		return generateForeachTable(t, files)
	}

	var lines [][]string
	for i := 0; i < t.Count; i++ {
		lines = append(lines, generateRow(t.Columns, files))
	}

	file := csvFile{
		name: t.Name,
		header: lo.Map(t.Columns, func(c column, i int) string {
			return c.Name
		}),
		lines: lines,
	}
	return file, nil
}

func generateForeachTable(t table, files map[string]csvFile) (csvFile, error) {
	var lines [][]string

	for i := 0; i < t.Count; i++ {
		ref, ok := files[t.Foreach]
		if !ok {
			return csvFile{}, fmt.Errorf("no reference table called %q, make sure you've generated it first", t.Foreach)
		}

		for _, line := range ref.lines {
			lines = append(lines, generateRefRow(t.Columns, ref.header, line))
		}
	}

	file := csvFile{
		name: t.Name,
		header: lo.Map(t.Columns, func(c column, i int) string {
			return c.Name
		}),
		lines: lines,
	}
	return file, nil
}

func generateRefRow(columns []column, refHeader []string, refLine []string) []string {
	var line []string

	for _, c := range columns {
		if value, ok := getRefColumnValue(c.Value, refHeader, refLine); ok {
			line = append(line, value)
		} else {
			line = append(line, replacePlaceholders(c))
		}
	}

	return line
}

func generateRow(columns []column, files map[string]csvFile) []string {
	var line []string

	for _, c := range columns {
		if value, ok := getRefTableColumnValue(c.Value, files); ok {
			line = append(line, value)
		} else {
			line = append(line, replacePlaceholders(c))
		}
	}

	return line
}

func replacePlaceholders(c column) string {
	r := rand.Intn(100)
	if r < c.NullPercentage {
		return ""
	}

	s := c.Value
	for k, v := range replacements {
		if strings.Contains(s, k) {
			s = strings.ReplaceAll(s, k, v())
		}
	}

	return s
}

func writeFiles(outputDir string, cfs map[string]csvFile) error {
	for name, file := range cfs {
		if err := writeFile(outputDir, name, file); err != nil {
			return fmt.Errorf("writing file %q: %w", file.name, err)
		}
	}

	return nil
}

func writeFile(outputDir, name string, cf csvFile) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

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

	if err = writer.WriteAll(cf.lines); err != nil {
		return fmt.Errorf("writing csv lines for %q: %w", name, err)
	}

	writer.Flush()

	return nil
}

func getRefProp(refHeader []string, refLine []string, prop string) string {
	index := lo.IndexOf(refHeader, prop)

	return refLine[index]
}

func getRefTableColumnValue(s string, files map[string]csvFile) (string, bool) {
	matches := refExtractor.FindStringSubmatch(s)
	if len(matches) == 0 {
		return "", false
	}

	// TODO: This should be an error.
	table, ok := files[matches[1]]
	if !ok {
		return "", false
	}

	// Determine the property's index in the header.
	index := lo.IndexOf(table.header, matches[2])

	// Pick a random line.
	line := table.lines[rand.Intn(len(table.lines))]

	return line[index], true
}

func getRefColumnValue(s string, refHeader, refLine []string) (string, bool) {
	matches := refExtractor.FindStringSubmatch(s)
	if len(matches) == 0 {
		return "", false
	}

	// Determine the property's index in the header.
	index := lo.IndexOf(refHeader, matches[2])

	return refLine[index], true
}
