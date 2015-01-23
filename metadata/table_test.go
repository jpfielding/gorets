package metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

var mtable MTable = MTable{
	Version:  "38.31.12348",
	Date:     "Thu, 21 Aug 2014 04:05:48 GMT",
	Resource: "Office",
	Class:    "6",
	Table: []Table{
		Table{
			SystemName:    "sysid",
			LongName:      "sysid",
			DBName:        "sysid",
			ShortName:     "sysid",
			MaximumLength: "10",
			DataType:      "Int",
			Precision:     "0",
			Searchable:    "1",
			UseSeparator:  "0",
			MaxSelect:     "0",
			Index:         "1",
			Default:       "0",
			Required:      "0",
			Unique:        "1",
		},
		Table{
			SystemName:    "1724",
			StandardName:  "ModificationTimestamp",
			LongName:      "Last Trans Date",
			DBName:        "LastT_1724",
			ShortName:     "Last Trans Date",
			MaximumLength: "10",
			DataType:      "DateTime",
			Precision:     "0",
			Searchable:    "1",
			UseSeparator:  "0",
			MaxSelect:     "0",
			Index:         "0",
			Default:       "0",
			Required:      "0",
			Unique:        "0",
		},
	},
}

func TestTableCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(mtableCompact)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-TABLE")
	Ok(t, err)

	found := MTable{}
	err = found.InitFromCompact(parser, start)

	Ok(t, err)

	Equals(t, mtable, found)
}

func TestTableXml(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(mtableXml)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-TABLE")
	Ok(t, err)

	found := MTable{}
	err = found.InitFromXml(parser, start)

	Ok(t, err)

	Equals(t, mtable, found)
}

var mtableCompact = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA-TABLE Class="6" Date="Thu, 21 Aug 2014 04:05:48 GMT" Resource="Office" Version="38.31.12348">
    <COLUMNS>	SystemName	StandardName	LongName	DBName	ShortName	MaximumLength	DataType	Precision	Searchable	Interpretation	Alignment	UseSeparator	EditMaskID	LookupName	MaxSelect	Units	Index	Minimum	Maximum	Default	Required	SearchHelpID	Unique	</COLUMNS>
    <DATA>	sysid		sysid	sysid	sysid	10	Int	0	1			0			0		1			0	0		1	</DATA>
    <DATA>	1724	ModificationTimestamp	Last Trans Date	LastT_1724	Last Trans Date	10	DateTime	0	1			0			0		0			0	0		0	</DATA>
  </METADATA-TABLE>
</RETS>
`
var mtableXml = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA>
    <METADATA-TABLE Class="6" Date="Thu, 21 Aug 2014 04:05:48 GMT" Resource="Office" System="SUMRETS" Version="38.31.12348">
      <Field>
        <SystemName>sysid</SystemName>
        <StandardName></StandardName>
        <LongName>sysid</LongName>
        <DBName>sysid</DBName>
        <ShortName>sysid</ShortName>
        <MaximumLength>10</MaximumLength>
        <DataType>Int</DataType>
        <Precision>0</Precision>
        <Searchable>1</Searchable>
        <Interpretation></Interpretation>
        <Alignment></Alignment>
        <UseSeparator>0</UseSeparator>
        <EditMaskID></EditMaskID>
        <LookupName></LookupName>
        <MaxSelect>0</MaxSelect>
        <Units></Units>
        <Index>1</Index>
        <Minimum></Minimum>
        <Maximum></Maximum>
        <Default>0</Default>
        <Required>0</Required>
        <SearchHelpID></SearchHelpID>
        <Unique>1</Unique>
      </Field>
      <Field>
        <SystemName>1724</SystemName>
        <StandardName>ModificationTimestamp</StandardName>
        <LongName>Last Trans Date</LongName>
        <DBName>LastT_1724</DBName>
        <ShortName>Last Trans Date</ShortName>
        <MaximumLength>10</MaximumLength>
        <DataType>DateTime</DataType>
        <Precision>0</Precision>
        <Searchable>1</Searchable>
        <Interpretation></Interpretation>
        <Alignment></Alignment>
        <UseSeparator>0</UseSeparator>
        <EditMaskID></EditMaskID>
        <LookupName></LookupName>
        <MaxSelect>0</MaxSelect>
        <Units></Units>
        <Index>0</Index>
        <Minimum></Minimum>
        <Maximum></Maximum>
        <Default>0</Default>
        <Required>0</Required>
        <SearchHelpID></SearchHelpID>
        <Unique>0</Unique>
      </Field>
    </METADATA-TABLE>
  </METADATA>
</RETS>
`
