package rets

import (
	"encoding/hex"
	"encoding/xml"
	"strconv"
	"strings"
)

var MetadataLookup = map[string]string{
	"METADATA-SYSTEM":                     "SYSTEM",
	"METADATA-FOREIGN_KEY":                "ForeignKey",
	"METADATA-FILTER":                     "Filter",
	"METADATA-FILTER_TYPE":                "FilterType",
	"METADATA-RESOURCE":                   "Resource",
	"METADATA-CLASS":                      "Class",
	"METADATA-TABLE":                      "Field",
	"METADATA-UPDATE":                     "Update",
	"METADATA-UPDATE_TYPE":                "UpdateType",
	"METADATA-COLUMN_GROUP_SET":           "ColumnGroupSet",
	"METADATA-COLUMN_GROUP":               "ColumnGroup",
	"METADATA-COLUMN_GROUP_CONTROL":       "ColumnGroupControl",
	"METADATA-COLUMN_GROUP_TABLE":         "ColumnGroupTable",
	"METADATA-COLUMN_GROUP_NORMALIZATION": "ColumnGroupNormalization",

	"METADATA-OBJECT":                 "Object",
	"METADATA-SEARCHHELP":             "SearchHelp",
	"METADATA-EDITMASK":               "EditMask",
	"METADATA-UPDATEHELP":             "UpdateHelp",
	"METADATA-LOOKUP":                 "Lookup",
	"METADATA-LOOKUP_TYPE":            "LookupType",
	"METADATA-VALIDATION_LOOKUP":      "ValidationLookup",
	"METADATA-VALIDATION_LOOKUP_TYPE": "ValidationLookupType",

	"METADATA-VALIDATION_EXTERNAL":      "ValidationExternal",
	"METADATA-VALIDATION_EXTERNAL_TYPE": "ValidationExternalType",
	"METADATA-VALIDATION_EXPRESSION":    "ValidationExpression",
}

// NewCompactData parses a CompactData from a start element
func NewCompactData(start xml.StartElement, decoder *xml.Decoder, delim string) (CompactData, error) {
	cd := CompactData{}
	cd.ID = start.Name.Local
	cd.Delimiter = delim
	cd.Attr = map[string]string{}
	for _, a := range start.Attr {
		cd.Attr[a.Name.Local] = a.Value
	}
	return cd, decoder.DecodeElement(&cd, &start)
}

// CompactData is the common compact decoded structure
type CompactData struct {
	ID        string
	Delimiter string
	Attr      map[string]string
	// parse these values out with decode
	CompactColumns CompactRow   `xml:"COLUMNS"`
	CompactRows    []CompactRow `xml:"DATA"`
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

// Indexer provices cached lookup for CompactData
type Indexer func(col string, row Row) string

// Indexer create the cache
func (cd *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range cd.Columns() {
		index[c] = i
	}
	return func(col string, row Row) string {
		return row[index[col]]
	}
}

// Row is a string slice typedef for convenience
type Row []string

// CompactRow ...
type CompactRow string

// Parse ...
func (cr CompactRow) Parse(delim string) Row {
	split := strings.Split(string(cr), delim)
	return split[1 : len(split)-1]
}

// DelimiterTag holds the seperator for compact data
type DelimiterTag xml.StartElement

// Parse ...
func (dt DelimiterTag) Parse() (string, error) {
	del := dt.Attr[0].Value
	pad := strings.Repeat("0", 2-len(del))
	decoded, err := hex.DecodeString(pad + del)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// CountTag ...
type CountTag xml.StartElement

// Parse ...
func (ct CountTag) Parse() (int, error) {
	code, err := strconv.ParseInt(ct.Attr[0].Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(code), nil
}
