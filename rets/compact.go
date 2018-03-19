package rets

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// CompactDefaultDelim is the default field delimiter for COMPACT data
const CompactDefaultDelim = "\t"

// NewCompactData parses a CompactData from a start element.
// If delim is explicitly passed, it will override the DELIMITER element value, which defaults to \t. Pass
// empty string to automatically parse DELIMITER value or fallback to default of \t.
func NewCompactData(start xml.StartElement, decoder *xml.Decoder, delim string) (CompactData, error) {
	cd := CompactData{}
	cd.Element = start.Name.Local
	cd.Attr = map[string]string{}
	for _, a := range start.Attr {
		cd.Attr[a.Name.Local] = a.Value
	}
	err := decoder.DecodeElement(&cd, &start)
	if err != nil {
		return cd, err
	}
	if delim != "" {
		cd.Delimiter = delim
	} else if cd.CompactDelimiter.Value != "" {
		// explicitly specified delimiter, eg:  <DELIMITER value="2C"/>
		delimOrdinal, err := strconv.ParseInt(cd.CompactDelimiter.Value, 16, 8)
		if err != nil {
			return cd, fmt.Errorf("unable to parse DELIMITER value=%s", cd.CompactDelimiter.Value)
		}
		cd.Delimiter = string(rune(delimOrdinal))
	} else {
		cd.Delimiter = CompactDefaultDelim
	}
	return cd, nil
}

// CompactData is the common compact decoded structure
type CompactData struct {
	Element   string
	Delimiter string
	Attr      map[string]string
	// parse these values out with decode
	CompactDelimiter Delimiter    `xml:"DELIMITER"`
	CompactColumns   CompactRow   `xml:"COLUMNS"`
	CompactRows      []CompactRow `xml:"DATA"`
}

// Delimiter is the optional <DELIMITER value="2C"/> element
type Delimiter struct {
	Value string `xml:"value,attr"`
}

// Columns parses the compact values for the cols
func (cd CompactData) Columns() Row {
	return cd.CompactColumns.Parse(cd.Delimiter)
}

// Rows provides callback to access each row
func (cd CompactData) Rows(each func(i int, row Row)) {
	for i, row := range cd.CompactRows {
		each(i, row.Parse(cd.Delimiter))
	}
}

// CompactEntry ...
type CompactEntry map[string]string

// SetFields ...
func (ce CompactEntry) SetFields(target interface{}) {
	for key, value := range ce {
		e := reflect.ValueOf(target).Elem().FieldByNameFunc(func(n string) bool {
			// make this comparison case insensitive
			return strings.ToLower(n) == strings.ToLower(key)
		})
		if e.IsValid() {
			e.SetString(value)
		}
	}
}

// Entries turns all rows into maps
func (cd CompactData) Entries() []CompactEntry {
	index := cd.Indexer()
	cols := cd.Columns()
	var entries []CompactEntry
	cd.Rows(func(i int, r Row) {
		entry := CompactEntry{}
		for _, c := range cols {
			val, ok := index(c, r)
			if !ok {
				continue // declared column wasn't included in DATA row!
			}
			entry[c] = val
		}
		entries = append(entries, entry)
	})
	return entries
}

// Indexer provides cached lookup for CompactData
type Indexer func(col string, row Row) (val string, ok bool)

// Indexer create the cache
func (cd *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range cd.Columns() {
		index[c] = i
	}
	return func(col string, row Row) (val string, ok bool) {
		i, ok := index[col]
		if !ok || i >= len(row) {
			return "", false // non-existent column, or DATA row contains too few values
		}
		return row[i], true
	}
}

// Row is a string slice typedef for convenience
type Row []string

// CompactRow ...
type CompactRow string

// Parse ...
func (cr CompactRow) Parse(delim string) Row {
	asString := string(cr)
	if asString == "" {
		return []string{}
	}
	if delim == "" {
		delim = CompactDefaultDelim
	}
	split := strings.Split(asString, delim)
	start := 0
	if strings.HasPrefix(asString, delim) {
		start = 1
	}
	len := len(split)
	if strings.HasSuffix(asString, delim) {
		len = len - 1
	}
	return split[start:len]
}
