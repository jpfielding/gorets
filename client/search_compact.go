package client

import (
	"bytes"
	"encoding/xml"
	"io"

	"golang.org/x/net/context"
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
	RetsResponse RetsResponse
	Count        int
	Delimiter    string
	Columns      []string
	MaxRows      bool

	body   io.ReadCloser
	parser *xml.Decoder
	buf    bytes.Buffer
}

// CompactRow ...
type CompactRow func(row []string, err error) error

// Listen ...
func (c *CompactSearchResult) Listen(each CompactRow) error {
	defer c.body.Close()
	for {
		token, err := c.parser.Token()
		if err != nil {
			if err = each(nil, err); err != nil {
				return err
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
				c.MaxRows = true
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "DATA":
				err := each(ParseCompactRow(c.buf.String(), c.Delimiter), nil)
				if err != nil {
					return err
				}
			case "RETS":
				return nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			c.buf.Write(bytes)
		}
	}
	return nil
}

// Close ...
func (c *CompactSearchResult) Close() error {
	return c.Close()
}

// NewCompactSearchResult ...
func NewCompactSearchResult(body io.ReadCloser) (*CompactSearchResult, error) {
	rets := RetsResponse{}
	parser := DefaultXMLDecoder(body, false)
	result := &CompactSearchResult{
		RetsResponse: rets,
		MaxRows:      false,
		body:         body,
		parser:       parser,
	}

	// extract the basic content before delving into the data
	for {
		token, err := parser.Token()
		switch err {
		case nil:
			// nothing
		case io.EOF:
			return result, nil
		default:
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			result.buf.Reset()
			switch t.Name.Local {
			case "RETS", "RETS-STATUS":
				rets, err := ParseRetsResponseTag(t)
				if err != nil {
					return nil, err
				}
				result.RetsResponse = *rets
			case "COUNT":
				result.Count, err = ParseCountTag(t)
				if err != nil {
					return nil, err
				}
			case "DELIMITER":
				result.Delimiter, err = ParseDelimiterTag(t)
				if err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			elmt := xml.EndElement(t)
			name := elmt.Name.Local
			switch name {
			case "COLUMNS":
				result.Columns = ParseCompactRow(result.buf.String(), result.Delimiter)
				return result, nil
			}
		case xml.CharData:
			bytes := xml.CharData(t)
			result.buf.Write(bytes)
		}
	}
}
