package rets

import (
	"encoding/xml"
	"io"

	"github.com/jpfielding/gominidom/minidom"

	"context"
)

// StandardXMLSearch if you set the wrong request Format you will get nothing back
func StandardXMLSearch(requester Requester, ctx context.Context, r SearchRequest) (*StandardXMLSearchResult, error) {
	body, err := SearchStream(requester, ctx, r)
	if err != nil {
		return nil, err
	}
	return NewStandardXMLSearchResult(body)
}

// StandardXMLSearchResult ...
type StandardXMLSearchResult struct {
	Response Response

	body   io.ReadCloser
	parser *xml.Decoder
}

// ForEach returns Count and MaxRows and any error that 'each' wont handle
func (c *StandardXMLSearchResult) ForEach(elem string, each minidom.EachDOM) (int, bool, error) {
	defer c.body.Close()
	count := 0
	maxrows := false
	md := minidom.MiniDom{
		StartFunc: func(start xml.StartElement) {
			switch start.Name.Local {
			case "COUNT":
				count, _ = countTag(start).Parse()
			case "MAXROWS":
				maxrows = true
			}
		},
		// quit on the the xml tag
		EndFunc: func(end xml.EndElement) bool {
			switch end.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				return true
			}
			return false
		},
	}
	err := md.Walk(c.parser, elem, each)
	return count, maxrows, err
}

// Close closesthe connection
func (c *StandardXMLSearchResult) Close() error {
	if c == nil || c.body == nil {
		return nil
	}
	return c.body.Close()
}

// NewStandardXMLSearchResult returns an XML search result handler to listen to elements
func NewStandardXMLSearchResult(body io.ReadCloser) (*StandardXMLSearchResult, error) {
	parser := DefaultXMLDecoder(body, false)
	result := &StandardXMLSearchResult{
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
			switch t.Name.Local {
			case XMLElemRETS, XMLElemRETSStatus:
				resp, err := ResponseTag(t).Parse()
				result.Response = *resp
				return result, err
			}
		}
	}
}
