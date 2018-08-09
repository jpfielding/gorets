package rets

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

const (
	// XMLElemRETS tag for RETS responses
	XMLElemRETS = "RETS"
	// XMLElemRETSStatus tag for RETS responses
	XMLElemRETSStatus = "RETS-STATUS"
)

// Response is the common wrapper details for each response
type Response struct {
	Code int    `xml:"ReplyCode,attr"`
	Text string `xml:"ReplyText,attr"`
}

// ResponseTag holds the separator for compact data
type ResponseTag xml.StartElement

// Parse ...
func (start ResponseTag) Parse() (*Response, error) {
	attrs := make(map[string]string)
	for _, v := range start.Attr {
		attrs[strings.ToLower(v.Name.Local)] = v.Value
	}
	code, err := strconv.ParseInt(attrs["replycode"], 10, 16)
	if err != nil {
		return nil, err
	}
	return &Response{
		Code: int(code),
		Text: attrs["replytext"],
	}, nil
}

// ReadResponse  ...
func ReadResponse(body io.ReadCloser) (*Response, error) {
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
			case XMLElemRETS, XMLElemRETSStatus:
				return ResponseTag(t).Parse()
			}
		}
	}
}
