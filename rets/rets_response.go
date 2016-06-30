package rets

import (
	"encoding/xml"
	"io"
	"strconv"
	"strings"
)

// RetsResponse is the common wrapper details for each response
type RetsResponse struct {
	ReplyCode int
	ReplyText string
}

// ParseRetsResponse ...
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
			case "RETS", "RETS-STATUS":
				return ParseRetsResponseTag(t)
			}
		}
	}
}

// ParseRetsResponseTag ...
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
