package rets

import (
	"bytes"
	"encoding/xml"
	"io"

	"context"
)

// StandardXMLSearch if you set the wrong request Format you will get nothing back
func StandardXMLSearch(requester Requester, ctx context.Context, r SearchRequest) (*CompactSearchResult, error) {
	body, err := SearchStream(requester, ctx, r)
	if err != nil {
		return nil, err
	}
	return NewCompactSearchResult(body)
}

// StandardXMLSearchResult ...
type StandardXMLSearchResult struct {
	Response Response
	Count    int
	XMLData  XMLData

	body   io.ReadCloser
	parser *xml.Decoder
	buf    bytes.Buffer
}

// EachEntry is the user hook to receive each element
type EachEntry func(row map[string]string, err error) error

// ForEach returns MaxRows and any error that 'each' wont handle
func (c *StandardXMLSearchResult) ForEach(each EachEntry) (bool, error) {
	defer c.body.Close()
	maxRows := false
	// need to capture maxrows
	c.XMLData.StartFunc = func(t xml.StartElement) {
		if t.Name.Local == "MAXROWS" {
			maxRows = true
		}
	}
	// make sure we exit cleanly
	c.XMLData.EndFunc = func(t xml.EndElement) bool {
		switch t.Name.Local {
		case XMLElemRETS, XMLElemRETSStatus:
			return true
		}
		return false
	}
	err := c.XMLData.Walk(c.parser, func(data map[string]string, err error) error {
		return each(data, err)
	})
	return maxRows, err
}

// Close closesthe connection
func (c *StandardXMLSearchResult) Close() error {
	if c == nil || c.body == nil {
		return nil
	}
	return c.body.Close()
}

// NewStandardXMLSearchResult returns an XML search result handler to listen to elements
func NewStandardXMLSearchResult(body io.ReadCloser, elem string, repeatElems ...string) (*StandardXMLSearchResult, error) {
	parser := DefaultXMLDecoder(body, false)
	result := &StandardXMLSearchResult{
		XMLData: XMLData{
			Prefix:      elem,
			RepeatElems: repeatElems,
		},
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
				result.Count, err = CountTag(t).Parse()
				return result, err
			}
		case xml.EndElement:
			switch t.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				// if there is only a RETS tag.. just exit
				return result, nil
			}
		}
	}
}
