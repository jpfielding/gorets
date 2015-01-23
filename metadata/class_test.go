package metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

var mclass MClass = MClass{
	Version:  "38.31.12348",
	Date:     "Thu, 21 Aug 2014 04:05:48 GMT",
	Resource: "Office",
	Class: []Class{
		Class{
			ClassName:        "6",
			StandardName:     "Office",
			VisibleName:      "Office",
			Description:      "Office",
			TableVersion:     "38.31.12348",
			TableDate:        "Thu, 21 Aug 2014 04:05:48 GMT",
			UpdateVersion:    "",
			UpdateDate:       "",
			ClassTimeStamp:   "",
			DeletedFlagField: "",
			DeletedFlagValue: "",
			HasKeyIndex:      "",
			OffsetSupport:    "",
		},
	},
}

var mclassCompact = `
<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA-CLASS Date="Thu, 21 Aug 2014 04:05:48 GMT" Resource="Office" Version="38.31.12348">
    <COLUMNS>	ClassName	StandardName	VisibleName	Description	TableVersion	TableDate	UpdateVersion	UpdateDate	</COLUMNS>
    <DATA>	6	Office	Office	Office	38.31.12348	Thu, 21 Aug 2014 04:05:48 GMT			</DATA>
  </METADATA-CLASS>
</RETS>
`

func TestClassCompact(t *testing.T) {

	body := ioutil.NopCloser(bytes.NewReader([]byte(mclassCompact)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-CLASS")
	Ok(t, err)

	found := MClass{}
	err = found.InitFromCompact(parser, start)

	Ok(t, err)

	Equals(t, mclass, found)
}

var mclassXml = `<RETS ReplyCode="0" ReplyText="Success">
  <METADATA>
    <METADATA-CLASS Date="Thu, 21 Aug 2014 04:05:48 GMT" Resource="Office" System="SUMRETS" Version="38.31.12348">
      <Class>
        <ClassName>6</ClassName>
        <StandardName>Office</StandardName>
        <VisibleName>Office</VisibleName>
        <Description>Office</Description>
        <TableVersion>38.31.12348</TableVersion>
        <TableDate>Thu, 21 Aug 2014 04:05:48 GMT</TableDate>
        <UpdateVersion></UpdateVersion>
        <UpdateDate></UpdateDate>
      </Class>
    </METADATA-CLASS>
  </METADATA>
</RETS>
`

func TestClassXml(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(mclassXml)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-CLASS")
	Ok(t, err)

	found := MClass{}
	err = found.InitFromXml(parser, start)

	Ok(t, err)

	Equals(t, mclass, found)
}
