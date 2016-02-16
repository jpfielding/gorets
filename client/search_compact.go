package client

import (
	"bytes"
	"encoding/xml"
	"errors"
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

	body   io.ReadCloser
	parser *xml.Decoder
	buf    bytes.Buffer
}

// CompactRow ...
type CompactRow func(row []string, err error) error

// ForEach returns MaxRows and any error that 'each' wont handle
func (c *CompactSearchResult) ForEach(each CompactRow) (bool, error) {
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
				err := each(ParseCompactRow(c.buf.String(), c.Delimiter), nil)
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
	return maxRows, errors.New("invalid exit from RETS stream")
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
	rets := RetsResponse{}
	parser := DefaultXMLDecoder(body, false)
	result := &CompactSearchResult{
		RetsResponse: rets,
		body:         body,
		parser:       parser,
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
				rets, err := ParseRetsResponseTag(t)
				if err != nil {
					return result, err
				}
				result.RetsResponse = *rets
			case "COUNT":
				result.Count, err = ParseCountTag(t)
				if err != nil {
					return result, err
				}
			case "DELIMITER":
				result.Delimiter, err = ParseDelimiterTag(t)
				if err != nil {
					return result, err
				}
			}
		case xml.EndElement:
			elmt := xml.EndElement(t)
			name := elmt.Name.Local
			switch name {
			case "COLUMNS":
				result.Columns = ParseCompactRow(result.buf.String(), result.Delimiter)
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
