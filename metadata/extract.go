package metadata

import (
	"encoding/xml"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Extractor processes metadata elements
type Extractor struct {
	Body   io.ReadCloser
	parser *xml.Decoder
}

// Open a metadata stream and read in the RETS response
func (e *Extractor) Open() (RETSResponse, error) {
	// TODO extract common work from rets/rets_response.go
	rets := RETSResponse{}
	e.parser = xml.NewDecoder(e.Body)
	start, err := e.skipTo("(RETS|RETS-STATUS)")
	attrs := make(map[string]string)
	for _, v := range start.Attr {
		attrs[strings.ToLower(v.Name.Local)] = v.Value
	}
	code, err := strconv.ParseInt(attrs["replycode"], 10, 16)
	if err != nil {
		return rets, err
	}
	rets.ReplyCode = int(code)
	rets.ReplyText = attrs["replytext"]
	return rets, nil
}

// DecodeNext the provided elemment
func (e *Extractor) DecodeNext(match string, elem interface{}) error {
	next, err := e.skipTo(match)
	if err != nil {
		return err
	}
	return e.parser.DecodeElement(elem, &next)
}

// skipTo advances the cursor to the named xml.StartElement
func (e *Extractor) skipTo(match string) (xml.StartElement, error) {
	next, err := regexp.Compile(match)
	if err != nil {
		return xml.StartElement{}, err
	}
	for {
		token, err := e.parser.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if next.MatchString(t.Name.Local) {
				return t, nil
			}
		}
	}
}

// RETSResponse ...
type RETSResponse struct {
	// TODO extract common work from rets/rets_response.go
	ReplyCode int
	ReplyText string
}
