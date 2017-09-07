package rets

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompactEntry(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(compact))
	decoder := DefaultXMLDecoder(body, false)
	token, err := decoder.Token()
	assert.Nil(t, err)
	start, ok := token.(xml.StartElement)
	assert.Equal(t, true, ok, "should be a start element")
	cm, err := NewCompactData(start, decoder, "	")
	assert.Nil(t, err)
	type Test struct {
		ResourceID, Standardname string
	}
	row1 := Test{}
	maps := cm.Entries()
	maps[0].SetFields(&row1)
	assert.Equal(t, "ActiveAgent", row1.ResourceID)
	assert.Equal(t, "ActiveAgent", row1.Standardname)
}

func TestCompactRowParsing(t *testing.T) {
	var col = `	A	B	C	D	E	F	`
	var row = `	1	2	3	4		6	`
	var delim = `	`
	headers := CompactRow(col).Parse(delim)
	values := CompactRow(row).Parse(delim)

	assert.Equal(t, 6, int(len(headers)))
	assert.Equal(t, 6, int(len(values)))
}

func TestCompactRowParsingEmpty(t *testing.T) {
	var col = `	A	B	C	D	E	F	`
	var row = ``
	var delim = `	`

	headers := CompactRow(col).Parse(delim)
	values := CompactRow(row).Parse(delim)

	assert.Equal(t, 6, int(len(headers)))
	assert.Equal(t, 0, int(len(values)))
}

func TestCompactRowParsingFubar(t *testing.T) {
	var col = `	A	B	C	D	E	F	`
	var row = `	`
	var delim = `	`
	headers := CompactRow(col).Parse(delim)
	values := CompactRow(row).Parse(delim)

	assert.Equal(t, 6, int(len(headers)))
	assert.Equal(t, 0, int(len(values)))
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
	assert.Nil(t, err)
	start, ok := token.(xml.StartElement)
	assert.Equal(t, true, ok, "should be a start element")
	assert.Equal(t, "METADATA-ELEMENT", start.Name.Local)
	cm, err := NewCompactData(start, decoder, "	")
	assert.Nil(t, err)
	assert.Equal(t, "METADATA-ELEMENT", cm.Element)
	assert.Equal(t, "Dog", cm.Attr["Cat"])
	assert.Equal(t, 2, len(cm.CompactRows))
	assert.Equal(t, 2, len(cm.Columns()))
}
