package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
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
			value := v()
			var valueStr string
			if pg.Format != "" {
				// Check if the value implements the formatter interface and use that first,
				// otherwise, just perform a simple string format.
				if f, ok := value.(formatter); ok {
					valueStr = f.Format(pg.Format)
				} else {
					valueStr = fmt.Sprintf(pg.Format, value)
				}
			} else {
				valueStr = fmt.Sprintf("%v", value)
			}
			s = strings.ReplaceAll(s, k, valueStr)
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
	Format         string `yaml:"format"`
}

type rawMessage struct {
	unmarshal func(interface{}) error
}

type csvFile struct {
	name   string
	header []string
	lines  [][]string
}

type formatter interface {
	Format(string) string
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

	replacements = map[string]func() any{
		"${ach_account}":                 func() any { return gofakeit.AchAccount() },
		"${ach_routing}":                 func() any { return gofakeit.AchRouting() },
		"${adjective_demonstrative}":     func() any { return gofakeit.AdjectiveDemonstrative() },
		"${adjective_descriptive}":       func() any { return gofakeit.AdjectiveDescriptive() },
		"${adjective_indefinite}":        func() any { return gofakeit.AdjectiveIndefinite() },
		"${adjective_interrogative}":     func() any { return gofakeit.AdjectiveInterrogative() },
		"${adjective_possessive}":        func() any { return gofakeit.AdjectivePossessive() },
		"${adjective_proper}":            func() any { return gofakeit.AdjectiveProper() },
		"${adjective_quantitative}":      func() any { return gofakeit.AdjectiveQuantitative() },
		"${adjective}":                   func() any { return gofakeit.Adjective() },
		"${adverb_degree}":               func() any { return gofakeit.AdverbDegree() },
		"${adverb_frequency_definite}":   func() any { return gofakeit.AdverbFrequencyDefinite() },
		"${adverb_frequency_indefinite}": func() any { return gofakeit.AdverbFrequencyIndefinite() },
		"${adverb_manner}":               func() any { return gofakeit.AdverbManner() },
		"${adverb_place}":                func() any { return gofakeit.AdverbPlace() },
		"${adverb_time_definite}":        func() any { return gofakeit.AdverbTimeDefinite() },
		"${adverb_time_indefinite}":      func() any { return gofakeit.AdverbTimeIndefinite() },
		"${adverb}":                      func() any { return gofakeit.Adverb() },
		"${animal_type}":                 func() any { return gofakeit.AnimalType() },
		"${animal}":                      func() any { return gofakeit.Animal() },
		"${app_author}":                  func() any { return gofakeit.AppAuthor() },
		"${app_name}":                    func() any { return gofakeit.AppName() },
		"${app_version}":                 func() any { return gofakeit.AppVersion() },
		"${bitcoin_address}":             func() any { return gofakeit.BitcoinAddress() },
		"${bitcoin_private_key}":         func() any { return gofakeit.BitcoinPrivateKey() },
		"${bool}":                        func() any { return gofakeit.Bool() },
		"${breakfast}":                   func() any { return gofakeit.Breakfast() },
		"${bs}":                          func() any { return gofakeit.BS() },
		"${car_fuel_type}":               func() any { return gofakeit.CarFuelType() },
		"${car_maker}":                   func() any { return gofakeit.CarMaker() },
		"${car_model}":                   func() any { return gofakeit.CarModel() },
		"${car_transmission_type}":       func() any { return gofakeit.CarTransmissionType() },
		"${car_type}":                    func() any { return gofakeit.CarType() },
		"${chrome_user_agent}":           func() any { return gofakeit.ChromeUserAgent() },
		"${city}":                        func() any { return gofakeit.City() },
		"${color}":                       func() any { return gofakeit.Color() },
		"${company_suffix}":              func() any { return gofakeit.CompanySuffix() },
		"${company}":                     func() any { return gofakeit.Company() },
		"${connective_casual}":           func() any { return gofakeit.ConnectiveCasual() },
		"${connective_complaint}":        func() any { return gofakeit.ConnectiveComplaint() },
		"${connective_examplify}":        func() any { return gofakeit.ConnectiveExamplify() },
		"${connective_listing}":          func() any { return gofakeit.ConnectiveListing() },
		"${connective_time}":             func() any { return gofakeit.ConnectiveTime() },
		"${connective}":                  func() any { return gofakeit.Connective() },
		"${country_abr}":                 func() any { return gofakeit.CountryAbr() },
		"${country}":                     func() any { return gofakeit.Country() },
		"${credit_card_cvv}":             func() any { return gofakeit.CreditCardCvv() },
		"${credit_card_exp}":             func() any { return gofakeit.CreditCardExp() },
		"${credit_card_type}":            func() any { return gofakeit.CreditCardType() },
		"${currency_long}":               func() any { return gofakeit.CurrencyLong() },
		"${currency_short}":              func() any { return gofakeit.CurrencyShort() },
		"${date}":                        func() any { return gofakeit.Date() },
		"${day}":                         func() any { return gofakeit.Day() },
		"${dessert}":                     func() any { return gofakeit.Dessert() },
		"${dinner}":                      func() any { return gofakeit.Dinner() },
		"${domain_name}":                 func() any { return gofakeit.DomainName() },
		"${domain_suffix}":               func() any { return gofakeit.DomainSuffix() },
		"${email}":                       func() any { return gofakeit.Email() },
		"${emoji}":                       func() any { return gofakeit.Emoji() },
		"${file_extension}":              func() any { return gofakeit.FileExtension() },
		"${file_mime_type}":              func() any { return gofakeit.FileMimeType() },
		"${firefox_user_agent}":          func() any { return gofakeit.FirefoxUserAgent() },
		"${first_name}":                  func() any { return gofakeit.FirstName() },
		"${flipacoin}":                   func() any { return gofakeit.FlipACoin() },
		"${float32}":                     func() any { return gofakeit.Float32() },
		"${float64}":                     func() any { return gofakeit.Float64() },
		"${fruit}":                       func() any { return gofakeit.Fruit() },
		"${gender}":                      func() any { return gofakeit.Gender() },
		"${hexcolor}":                    func() any { return gofakeit.HexColor() },
		"${hobby}":                       func() any { return gofakeit.Hobby() },
		"${hour}":                        func() any { return gofakeit.Hour() },
		"${http_method}":                 func() any { return gofakeit.HTTPMethod() },
		"${http_status_code_simple}":     func() any { return gofakeit.HTTPStatusCodeSimple() },
		"${http_status_code}":            func() any { return gofakeit.HTTPStatusCode() },
		"${http_version}":                func() any { return gofakeit.HTTPVersion() },
		"${int16}":                       func() any { return gofakeit.Int16() },
		"${int32}":                       func() any { return gofakeit.Int32() },
		"${int64}":                       func() any { return gofakeit.Int64() },
		"${int8}":                        func() any { return gofakeit.Int8() },
		"${ipv4_address}":                func() any { return gofakeit.IPv4Address() },
		"${ipv6_address}":                func() any { return gofakeit.IPv6Address() },
		"${job_descriptor}":              func() any { return gofakeit.JobDescriptor() },
		"${job_level}":                   func() any { return gofakeit.JobLevel() },
		"${job_title}":                   func() any { return gofakeit.JobTitle() },
		"${language_abbreviation}":       func() any { return gofakeit.LanguageAbbreviation() },
		"${language}":                    func() any { return gofakeit.Language() },
		"${last_name}":                   func() any { return gofakeit.LastName() },
		"${latitude}":                    func() any { return gofakeit.Latitude() },
		"${longitude}":                   func() any { return gofakeit.Longitude() },
		"${lunch}":                       func() any { return gofakeit.Lunch() },
		"${mac_address}":                 func() any { return gofakeit.MacAddress() },
		"${minute}":                      func() any { return gofakeit.Minute() },
		"${month_string}":                func() any { return gofakeit.MonthString() },
		"${month}":                       func() any { return gofakeit.Month() },
		"${name_prefix}":                 func() any { return gofakeit.NamePrefix() },
		"${name_suffix}":                 func() any { return gofakeit.NameSuffix() },
		"${name}":                        func() any { return gofakeit.Name() },
		"${nanosecond}":                  func() any { return gofakeit.NanoSecond() },
		"${nicecolors}":                  func() any { return gofakeit.NiceColors() },
		"${noun_abstract}":               func() any { return gofakeit.NounAbstract() },
		"${noun_collective_animal}":      func() any { return gofakeit.NounCollectiveAnimal() },
		"${noun_collective_people}":      func() any { return gofakeit.NounCollectivePeople() },
		"${noun_collective_thing}":       func() any { return gofakeit.NounCollectiveThing() },
		"${noun_common}":                 func() any { return gofakeit.NounCommon() },
		"${noun_concrete}":               func() any { return gofakeit.NounConcrete() },
		"${noun_countable}":              func() any { return gofakeit.NounCountable() },
		"${noun_uncountable}":            func() any { return gofakeit.NounUncountable() },
		"${noun}":                        func() any { return gofakeit.Noun() },
		"${opera_user_agent}":            func() any { return gofakeit.OperaUserAgent() },
		"${password}":                    func() any { return gofakeit.Password(true, true, true, true, true, 25) },
		"${pet_name}":                    func() any { return gofakeit.PetName() },
		"${phone_formatted}":             func() any { return gofakeit.PhoneFormatted() },
		"${phone}":                       func() any { return gofakeit.Phone() },
		"${phrase}":                      func() any { return gofakeit.Phrase() },
		"${preposition_compound}":        func() any { return gofakeit.PrepositionCompound() },
		"${preposition_double}":          func() any { return gofakeit.PrepositionDouble() },
		"${preposition_simple}":          func() any { return gofakeit.PrepositionSimple() },
		"${preposition}":                 func() any { return gofakeit.Preposition() },
		"${programming_language}":        func() any { return gofakeit.ProgrammingLanguage() },
		"${pronoun_demonstrative}":       func() any { return gofakeit.PronounDemonstrative() },
		"${pronoun_interrogative}":       func() any { return gofakeit.PronounInterrogative() },
		"${pronoun_object}":              func() any { return gofakeit.PronounObject() },
		"${pronoun_personal}":            func() any { return gofakeit.PronounPersonal() },
		"${pronoun_possessive}":          func() any { return gofakeit.PronounPossessive() },
		"${pronoun_reflective}":          func() any { return gofakeit.PronounReflective() },
		"${pronoun_relative}":            func() any { return gofakeit.PronounRelative() },
		"${pronoun}":                     func() any { return gofakeit.Pronoun() },
		"${quote}":                       func() any { return gofakeit.Quote() },
		"${rgbcolor}":                    func() any { return gofakeit.RGBColor() },
		"${safari_user_agent}":           func() any { return gofakeit.SafariUserAgent() },
		"${safecolor}":                   func() any { return gofakeit.SafeColor() },
		"${second}":                      func() any { return gofakeit.Second() },
		"${snack}":                       func() any { return gofakeit.Snack() },
		"${ssn}":                         func() any { return gofakeit.SSN() },
		"${state_abr}":                   func() any { return gofakeit.StateAbr() },
		"${state}":                       func() any { return gofakeit.State() },
		"${street_name}":                 func() any { return gofakeit.StreetName() },
		"${street_number}":               func() any { return gofakeit.StreetNumber() },
		"${street_prefix}":               func() any { return gofakeit.StreetPrefix() },
		"${street_suffix}":               func() any { return gofakeit.StreetSuffix() },
		"${street}":                      func() any { return gofakeit.Street() },
		"${time_zone_abv}":               func() any { return gofakeit.TimeZoneAbv() },
		"${time_zone_full}":              func() any { return gofakeit.TimeZoneFull() },
		"${time_zone_offset}":            func() any { return gofakeit.TimeZoneOffset() },
		"${time_zone_region}":            func() any { return gofakeit.TimeZoneRegion() },
		"${time_zone}":                   func() any { return gofakeit.TimeZone() },
		"${uint128_hex}":                 func() any { return gofakeit.HexUint128() },
		"${uint16_hex}":                  func() any { return gofakeit.HexUint16() },
		"${uint16}":                      func() any { return gofakeit.Uint16() },
		"${uint256_hex}":                 func() any { return gofakeit.HexUint256() },
		"${uint32_hex}":                  func() any { return gofakeit.HexUint32() },
		"${uint32}":                      func() any { return gofakeit.Uint32() },
		"${uint64_hex}":                  func() any { return gofakeit.HexUint64() },
		"${uint64}":                      func() any { return gofakeit.Uint64() },
		"${uint8_hex}":                   func() any { return gofakeit.HexUint8() },
		"${uint8}":                       func() any { return gofakeit.Uint8() },
		"${url}":                         func() any { return gofakeit.URL() },
		"${user_agent}":                  func() any { return gofakeit.UserAgent() },
		"${username}":                    func() any { return gofakeit.Username() },
		"${uuid}":                        func() any { return gofakeit.UUID() },
		"${vegetable}":                   func() any { return gofakeit.Vegetable() },
		"${verb_action}":                 func() any { return gofakeit.VerbAction() },
		"${verb_helping}":                func() any { return gofakeit.VerbHelping() },
		"${verb_linking}":                func() any { return gofakeit.VerbLinking() },
		"${verb}":                        func() any { return gofakeit.Verb() },
		"${weekday}":                     func() any { return gofakeit.WeekDay() },
		"${word}":                        func() any { return gofakeit.Word() },
		"${year}":                        func() any { return gofakeit.Year() },
		"${zip}":                         func() any { return gofakeit.Zip() },
	}
)
