package metadata

import (
	"encoding/xml"
	"io"
)

// err := p.DecodeElement(m, &t)

// MClasses extracts classes from a metadata stream
func MClasses(body io.ReadCloser, err error) ([]MClass, error) {
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Next advances the cursor to the named xml.StartElement
func Next(parser *xml.Decoder, start string) (xml.StartElement, error) {
	for {
		token, err := parser.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case start:
				return t, nil
			}
		}
	}
}
