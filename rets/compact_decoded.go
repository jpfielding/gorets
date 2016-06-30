package rets

import (
	"encoding/hex"
	"encoding/xml"
	"strconv"
	"strings"
)

// CompactData is the common compact decoded structure
type CompactData struct {
	ID, Date, Version string
	Columns           []string
	Rows              [][]string
}

// Indexer provices cached lookup for CompactData
type Indexer func(col string, row int) string

// Indexer create the cache
func (m *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range m.Columns {
		index[c] = i
	}
	return func(col string, row int) string {
		return m.Rows[row][index[col]]
	}
}

// ParseDelimiterTag ...
func ParseDelimiterTag(start xml.StartElement) (string, error) {
	del := start.Attr[0].Value
	pad := strings.Repeat("0", 2-len(del))
	decoded, err := hex.DecodeString(pad + del)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// ParseCountTag ...
func ParseCountTag(count xml.StartElement) (int, error) {
	code, err := strconv.ParseInt(count.Attr[0].Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(code), nil
}

// ParseCompactRow ...
func ParseCompactRow(row, delim string) []string {
	split := strings.Split(row, delim)
	return split[1 : len(split)-1]
}
