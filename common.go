package gorets


import (
	"encoding/xml"
	"encoding/hex"
	"fmt"
	"strings"
	"strconv"
	"net/url"
	"io"
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
			elmt := xml.StartElement(t)
			name := elmt.Name.Local
			switch name {
			case "RETS":
				return ParseRetsResponseTag(elmt)
			}
		}
	}
}

func ParseRetsResponseTag(start xml.StartElement) (*RetsResponse, error) {
	rets := RetsResponse{}
	attrs := make(map[string]string)
	for _,v := range start.Attr {
		attrs[strings.ToLower(v.Name.Local)] = v.Value
	}
	code,err := strconv.ParseInt(attrs["replycode"],10,16)
	if err != nil {
		return nil, err
	}
	rets.ReplyCode = int(code)
	rets.ReplyText = attrs["replytext"]
	return &rets, nil
}

func ParseDelimiterTag(start xml.StartElement) (string, error) {
	del := start.Attr[0].Value
	pad := strings.Repeat("0",2-len(del))
	decoded,err := hex.DecodeString(pad+del)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func ParseCountTag(count xml.StartElement) (int, error) {
	code,err := strconv.ParseInt(count.Attr[0].Value, 10, 64)
	if err != nil {
		return -1, err
	}
	return int(code), nil
}

func SplitRowByDelim(row, delim string) ([]string) {
	return strings.Split(strings.Trim(row, delim), delim)
}

func OptionalStringValue(values url.Values) (func (string, string)) {
	return func (name, value string) {
		if value != "" {
			values.Add(name, value)
		}
	}
}
func OptionalIntValue(values url.Values) (func (string, int)) {
	return func (name string, value int) {
		if value >= 0 {
			values.Add(name, fmt.Sprintf("%d",value))
		}
	}
}

