package rets

import (
	"bytes"
	"encoding/hex"
	"encoding/xml"
	"io"
	"strings"

	"context"
)

// SearchCompact if you set the wrong request Format you will get nothing back
func SearchCompact(ctx context.Context, requester Requester, r SearchRequest) (*CompactSearchResult, error) {
	body, err := SearchStream(ctx, requester, r)
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
	if c.body == nil {
		return false, nil
	}
	maxRows := false
	for {
		token, err := c.parser.Token()
		if err != nil {
			// dont catch io.EOF here since a clean read should exit at the </RETS> tag
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
			case XMLElemRETS, XMLElemRETSStatus:
				c.Close() // close the stream since we've read to the end
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
	tmp := c.body
	c.body = nil
	return tmp.Close()
}

// NewCompactSearchResult _always_ close this
func NewCompactSearchResult(body io.ReadCloser) (*CompactSearchResult, error) {
	parser := DefaultXMLDecoder(body, false)
	result := &CompactSearchResult{
		body:   body,
		parser: parser,
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
			case XMLElemRETS, XMLElemRETSStatus:
				resp, er := ResponseTag(t).Parse()
				if er != nil {
					return result, er
				}
				result.Response = *resp
			case "COUNT":
				result.Count, err = countTag(t).Parse()
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
			switch t.Name.Local {
			case "COLUMNS":
				result.Columns = CompactRow(result.buf.String()).Parse(result.Delimiter)
				return result, nil
			case XMLElemRETS, XMLElemRETSStatus:
				// if there is only a RETS tag.. close the stream and just exit
				result.Close()
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
