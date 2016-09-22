package rets

import (
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"io"
	"strconv"
	"strings"

	"context"
)

// SearchCompact if you set the wrong request Format you will get nothing back
func SearchCompact(requester Requester, ctx context.Context, r SearchRequest) (*CompactSearchResult, error) {
	body, err := SearchStream(requester, ctx, r)
	if err != nil {
		return nil, err
	}
	return NewCompactSearchResult(body)
}

// CompactSearchResult ...
type CompactSearchResult struct {
	Response  Response
	Count     int
	Delimiter string
	Columns   Row

	body   io.ReadCloser
	parser *xml.Decoder
	buf    bytes.Buffer
}

// EachRow ...
type EachRow func(row Row, err error) error

// ForEach returns MaxRows and any error that 'each' wont handle
func (c *CompactSearchResult) ForEach(each EachRow) (bool, error) {
	defer c.body.Close()
	maxRows := false
	for {
		token, err := c.parser.Token()
		if err != nil {
			if err = each(nil, err); err != nil {
				return maxRows, err
			}
			continue
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			c.buf.Reset()
			// check tags
			switch t.Name.Local {
			case "MAXROWS":
				maxRows = true
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "DATA":
				err := each(CompactRow(c.buf.String()).Parse(c.Delimiter), nil)
				if err != nil {
					return maxRows, err
				}
			case "RETS":
				return maxRows, nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			c.buf.Write(bytes)
		}
	}
}

// Close ...
func (c *CompactSearchResult) Close() error {
	if c == nil || c.body == nil {
		return nil
	}
	return c.body.Close()
}

// NewCompactSearchResult _always_ close this
func NewCompactSearchResult(body io.ReadCloser) (*CompactSearchResult, error) {
	rets := Response{}
	parser := DefaultXMLDecoder(body, false)
	result := &CompactSearchResult{
		Response: rets,
		body:     body,
		parser:   parser,
	}
	// extract the basic content before delving into the data
	for {
		token, err := parser.Token()
		if err != nil {
			return result, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			result.buf.Reset()
			switch t.Name.Local {
			case "RETS", "RETS-STATUS":
				rets, err := ResponseTag(t).Parse()
				if err != nil {
					return result, err
				}
				result.Response = *rets
			case "COUNT":
				result.Count, err = CountTag(t).Parse()
				if err != nil {
					return result, err
				}
			case "DELIMITER":
				result.Delimiter, err = DelimiterTag(t).Parse()
				if err != nil {
					return result, err
				}
			}
		case xml.EndElement:
			elmt := xml.EndElement(t)
			name := elmt.Name.Local
			switch name {
			case "COLUMNS":
				result.Columns = CompactRow(result.buf.String()).Parse(result.Delimiter)
				return result, nil
			case "RETS", "RETS-STATUS":
				// if there is only a RETS tag.. just exit
				return result, nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			result.buf.Write(bytes)
		}
	}
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
