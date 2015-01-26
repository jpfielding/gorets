package metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestAttrsMap(t *testing.T) {
	var raw = `<?xml version="1.0" encoding="utf-8"?>
	<pets>
		<dog gender="female" fixed="true" name="josie"/>
	</pets>`
	body := ioutil.NopCloser(bytes.NewReader([]byte(raw)))

	parser := xml.NewDecoder(body)
	start, err := AdvanceToStartElem(parser, "dog")
	Ok(t, err)
	attrs := AttrsMap(start.Attr)
	Equals(t, "josie", attrs["name"])
	Equals(t, "female", attrs["gender"])
	Equals(t, "true", attrs["fixed"])
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
