package gorets_metadata

import (
	"encoding/xml"
)

// advances the cursor to the named xml.StartElement
func AdvanceToStartElem(parser *xml.Decoder, start string) (xml.StartElement, error) {
	for {
		token, err := parser.Token()
		if err != nil {
			return xml.StartElement{}, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			// clear any accumulated data
			switch t.Name.Local {
			case start:
				return t, nil
			}
		}
	}
}
