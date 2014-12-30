package gorets_client

import (
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

/* the common wrapper details for each response */
type RetsResponse struct {
	ReplyCode int
	ReplyText string
}

func ParseRetsResponse(body io.ReadCloser) (*RetsResponse, error) {
	parser := xml.NewDecoder(body)
	for {
		token, err := parser.Token()
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			switch t.Name.Local {
			case "RETS":
				return ParseRetsResponseTag(t)
			}
		}
	}
}

func ParseRetsResponseTag(start xml.StartElement) (*RetsResponse, error) {
	rets := RetsResponse{}
	attrs := make(map[string]string)
	for _, v := range start.Attr {
		attrs[strings.ToLower(v.Name.Local)] = v.Value
	}
	code, err := strconv.ParseInt(attrs["replycode"], 10, 16)
	if err != nil {
		return nil, err
	}
	rets.ReplyCode = int(code)
	rets.ReplyText = attrs["replytext"]
	return &rets, nil
}

func ParseDelimiterTag(start xml.StartElement) (string, error) {
	del := start.Attr[0].Value
	pad := strings.Repeat("0", 2-len(del))
	decoded, err := hex.DecodeString(pad + del)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func ParseCountTag(count xml.StartElement) (int, error) {
	code, err := strconv.ParseInt(count.Attr[0].Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(code), nil
}

func ParseCompactRow(row, delim string) []string {
	split := strings.Split(row, delim)
	return split[1 : len(split)-1]
}

func OptionalStringValue(values url.Values) func(string, string) {
	return func(name, value string) {
		if value != "" {
			values.Add(name, value)
		}
	}
}
func OptionalIntValue(values url.Values) func(string, int) {
	return func(name string, value int) {
		if value >= 0 {
			values.Add(name, fmt.Sprintf("%d", value))
		}
	}
}

/* the common compact decoded structure */
type CompactData struct {
	Id, Date, Version string
	Columns           []string
	Rows              [][]string
}

/** cached lookup */
type Indexer func(col string, row int) string

/** create the cache */
func (m *CompactData) Indexer() Indexer {
	index := make(map[string]int)
	for i, c := range m.Columns {
		index[c] = i
	}
	return func(col string, row int) string {
		return m.Rows[row][index[col]]
	}
}

func ParseMetadataCompactDecoded(start xml.StartElement, parser *xml.Decoder, delim string) (*CompactData, error) {
	type XmlMetadataElement struct {
		Resource string `xml:"Resource,attr"`
		/* only valid for table */
		Class string `xml:"Class,attr"`
		/* only valid for lookup_type */
		Lookup  string   `xml:"Lookup,attr"`
		Version string   `xml:"Version,attr"`
		Date    string   `xml:"Date,attr"`
		Columns string   `xml:"COLUMNS"`
		Data    []string `xml:"DATA"`
	}
	xme := XmlMetadataElement{}
	err := parser.DecodeElement(&xme, &start)
	if err != nil {
		fmt.Println("failed to decode: ", err)
		return nil, err
	}
	if xme.Columns == "" {
		return nil, nil
	}
	data := *extractMap(xme.Columns, xme.Data, delim)
	data.Date = xme.Date
	data.Version = xme.Version
	data.Id = xme.Resource
	if xme.Class != "" {
		data.Id = xme.Resource + ":" + xme.Class
	}
	if xme.Lookup != "" {
		data.Id = xme.Resource + ":" + xme.Lookup
	}

	return &data, nil
}

/** extract a map of fields from columns and rows */
func extractMap(cols string, rows []string, delim string) *CompactData {
	data := CompactData{}
	// remove the first and last chars
	data.Columns = ParseCompactRow(cols, delim)
	data.Rows = make([][]string, len(rows))
	// create each
	for i, line := range rows {
		data.Rows[i] = ParseCompactRow(line, delim)
	}
	return &data
}
