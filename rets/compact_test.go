package rets

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

func TestCompactRowParsing(t *testing.T) {
	var col = `	A	B	C	D	E	F	`
	var row = `	1	2	3	4		6	`
	var delim = `	`
	headers := CompactRow(col).Parse(delim)
	values := CompactRow(row).Parse(delim)

	testutils.Equals(t, 6, int(len(headers)))
	testutils.Equals(t, 6, int(len(values)))
}

var compact = `<METADATA-ELEMENT Cat="Dog" Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ResourceID	StandardName	</COLUMNS>
<DATA>	ActiveAgent	ActiveAgent	</DATA>
<DATA>	Agent	Agent	</DATA>
</METADATA-ELEMENT>`

func TestParseCompactData(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(compact))
	decoder := DefaultXMLDecoder(body, false)
	token, err := decoder.Token()
	testutils.Ok(t, err)
	start, ok := token.(xml.StartElement)
	testutils.Assert(t, ok, "should be a start element")
	testutils.Equals(t, "METADATA-ELEMENT", start.Name.Local)
	cm, err := NewCompactData(start, decoder, "	")
	testutils.Ok(t, err)
	testutils.Equals(t, "METADATA-ELEMENT", cm.ID)
	testutils.Equals(t, "Dog", cm.Attr["Cat"])
	testutils.Equals(t, 2, len(cm.CompactRows))
	testutils.Equals(t, 2, len(cm.Columns()))
}
