package metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

// err := p.DecodeElement(m, &t)

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

func TestAdvanceToStartElem(t *testing.T) {
	var raw = `<?xml version="1.0" encoding="utf-8"?>
	<pets>
		<dog gender="male" fixed="true" name="buddy"/>
		<cat gender="male" fixed="true" name="snowball"/>
		<dog gender="female" fixed="true" name="josie"/>
		<cat gender="female" fixed="true" name="mr bojangles"/>
	</pets>`
	body := ioutil.NopCloser(bytes.NewReader([]byte(raw)))

	parser := xml.NewDecoder(body)

	cat, err := AdvanceToStartElem(parser, "cat")
	Ok(t, err)
	Equals(t, "snowball", AttrsMap(cat.Attr)["name"])

	dog, err := AdvanceToStartElem(parser, "dog")
	Ok(t, err)
	Equals(t, "josie", AttrsMap(dog.Attr)["name"])

}
