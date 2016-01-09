package client

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

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

// OptionalStringValue ...
func OptionalStringValue(values url.Values) func(string, string) {
	return func(name, value string) {
		if value != "" {
			values.Add(name, value)
		}
	}
}

// OptionalIntValue ...
func OptionalIntValue(values url.Values) func(string, int) {
	return func(name string, value int) {
		if value >= 0 {
			values.Add(name, fmt.Sprintf("%d", value))
		}
	}
}
