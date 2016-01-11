package client

import (
	"encoding/hex"
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
			case "RETS":
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
