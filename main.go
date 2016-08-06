package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/icrowley/fake"
)

var (
	// Debug logger
	Debug = log.New(os.Stderr, "DBG ", log.LstdFlags)

	// Info logger
	Info = log.New(os.Stderr, "INF ", log.LstdFlags)

	// Error logger
	Error = log.New(os.Stderr, "ERR ", log.LstdFlags)

	// Crit logger
	Crit = log.New(os.Stderr, "CRT ", log.LstdFlags)
)

// Field sets properties of a single attribute.
type Field struct {
	Type            string    `json:"type"`
	PctEmpty        float64   `json:"pct_empty"`
	Values          []string  `json:"values"`
	DateFormat      string    `json:"date_format"`
	DateLowBoundary time.Time `json:"date_low_boundary"`

	Mean   float64 `json:"mean"`
	StdDev float64 `json:"stddev"`

	// Key is the name of another field this field depends on.
	Key string `json:"key"`

	// KeyValues is a values for specific key.
	KeyValues map[string][]string `json:"key_values"`
}

// config sets structure for configuration.
var config = struct {
	Schema map[string]Field
	Keys   []string
	Rows   int
}{}

func init() {

	// reading command line flags
	schemaFile := flag.String("schema", "schema.json", "Data Schema")
	rows := flag.Int("rows", 1000, "Number of rows to be generated")
	lang := flag.String("lang", "en", "Language, en|ru")
	columns := flag.String("columns", "", "Output columns in order")
	flag.Parse()

	// setting language
	if err := fake.SetLang(*lang); err != nil {
		Crit.Fatalf("Can't set the language: %s", err)
	}

	// reading schema
	schemaData, err := ioutil.ReadFile(*schemaFile)
	if err != nil {
		Crit.Fatalf("Can't read schema file: %s", err)
	}

	schema := make(map[string]Field)
	// parse configuration
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		Crit.Fatalf("Can't unmarshal configuration: %s", err)
	}

	config.Rows = *rows
	config.Schema = schema
	config.Keys = strings.Split(*columns, ",")
	//Info.Printf("Configuration: %+v", config)
}

func emptier(val string, pct float64) string {
	return val
}

func generateValue(field Field, keyVal string) (string, error) {
	//Debug.Printf("Field generation started: key - %s, field - %+v", keyVal, field)
	var val string
	switch field.Type {
	case "enum":
		val = Enum(field.Values)
	case "uniq_int32_by_key":
		val = UniqInt32ByKey(keyVal)
	case "enum_by_key":
		val = Enum(field.KeyValues[keyVal])
	case "date":
		val = Date(field.DateFormat, field.DateLowBoundary)
	case "unix_timestamp":
		val = UnixTimestamp(field.DateLowBoundary)
	case "norm_int32":
		val = NormInt32(field.Mean, field.StdDev)
	case "norm_multiplier_key":
		val = NormMultiplierKey(keyVal, field.Mean, field.StdDev)
	case "uniq_email":
		val = UniqEmail()
	}
	return emptier(val, field.PctEmpty), nil
}

func generateRow(keys []string, schema map[string]Field) (record []string, err error) {
	// container for field results
	result := make(map[string]string)

	// going through independent fields
	for _, key := range keys {
		field := schema[key]
		if field.Key != "" {
			continue
		}
		result[key], err = generateValue(field, "")
		if err != nil {
			return nil, err
		}
	}

	// going through dependent fields
	for _, key := range keys {
		field := schema[key]
		if field.Key == "" {
			continue
		}
		result[key], err = generateValue(field, result[field.Key])
		if err != nil {
			return nil, err
		}
	}

	record = make([]string, len(result))
	for i, key := range keys {
		record[i] = result[key]
	}
	return record, nil
}

func main() {

	wr := csv.NewWriter(os.Stdout)
	wr.Write(config.Keys)
	for i := 0; i < config.Rows; i++ {
		row, err := generateRow(config.Keys, config.Schema)
		if err != nil {
			Crit.Fatalf("Row generation failed: %s", err)
		}
		if err := wr.Write(row); err != nil {
			Crit.Fatalf("Writing row failed: %s", err)
		}
	}
	wr.Flush()
}
