package metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

func TestNext(t *testing.T) {
	type pet struct {
		Gender string `xml:"gender,attr"`
		Fixed  bool   `xml:"fixed,attr"`
		Name   string `xml:"name,attr"`
	}

	var raw = `<?xml version="1.0" encoding="utf-8"?>
	<pets>
		<dog gender="male" fixed="true" name="buddy"/>
		<cat gender="female" fixed="true" name="snowball"/>
		<dog gender="female" fixed="true" name="josie"/>
		<cat gender="male" fixed="false" name="mr bojangles"/>
	</pets>`
	body := ioutil.NopCloser(bytes.NewReader([]byte(raw)))
	parser := xml.NewDecoder(body)

	nextTest := func(tag, name, gender string, fixed bool) func(*testing.T) {
		return func(tt *testing.T) {
			start, err := Next(parser, tag)
			testutils.Ok(tt, err)
			pet := pet{}
			err = parser.DecodeElement(&pet, &start)
			testutils.Ok(tt, err)
			testutils.Equals(tt, name, pet.Name)
			testutils.Equals(tt, gender, pet.Gender)
			testutils.Equals(tt, fixed, pet.Fixed)
		}
	}

	t.Run("buddy", nextTest("dog", "buddy", "male", true))
	t.Run("josie", nextTest("dog", "josie", "female", true))
	t.Run("mr bojangles", nextTest("cat", "mr bojangles", "male", false))
}
