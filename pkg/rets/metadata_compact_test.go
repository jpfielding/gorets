/**
parsing the 'login' action from RETS
*/
package rets

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var retsStart = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">`
var retsEnd = `</RETS>`

var system = `<METADATA-SYSTEM Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<SYSTEM SystemID="SIRM" SystemDescription="MLS System"/>
<COMMENTS>The System is provided to you by Systems.</COMMENTS>
</METADATA-SYSTEM>`

func TestSystem(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + system + retsEnd))

	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)
	verifySystem(t, *ms)
}

func verifySystem(t *testing.T, cm CompactMetadata) {
	assert.Equal(t, "MLS System", cm.MSystem.System.Description)
	assert.Equal(t, "1.12.30", cm.MSystem.Version)
	assert.Equal(t, "The System is provided to you by Systems.", cm.MSystem.Comments)
}

var resource = `<METADATA-RESOURCE Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ResourceID	StandardName	VisibleName	Description	KeyField	ClassCount	ClassVersion	ClassDate	ObjectVersion	ObjectDate	SearchHelpVersion	SearchHelpDate	EditMaskVersion	EditMaskDate	LookupVersion	LookupDate	UpdateHelpVersion	UpdateHelpDate	ValidationExpressionVersion	ValidationExpressionDate	ValidationLookupVersion	ValidationLookupDate	ValidationExternalVersion	ValidationExternalDate	</COLUMNS>
<DATA>	ActiveAgent	ActiveAgent	Agent	ActiveAgent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	Agent	Agent	Agent	Agent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-RESOURCE>`

func TestParseResources(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + resource + retsEnd))

	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)
	verifyParseResources(t, *ms)
}

func verifyParseResources(t *testing.T, cm CompactMetadata) {
	resource := cm.Elements["METADATA-RESOURCE"][0]

	assert.Equal(t, "1.12.30", resource.Attr["Version"])
	assert.Equal(t, len(resource.CompactRows), 2)

	indexer := resource.Indexer()
	var rows []Row
	resource.Rows(func(i int, r Row) {
		rows = append(rows, r)
	})

	val, ok := indexer("ResourceID", rows[0])
	assert.True(t, ok, "ResourceID not found in row")
	assert.Equal(t, "ActiveAgent", val)

	val, ok = indexer("ValidationExternalDate", rows[1])
	assert.True(t, ok, "ValidationExternalDate not found in row")
	assert.Equal(t, "Tue, 3 Sep 2013 00:00:00 GMT", val)
}

var class = `<METADATA-CLASS Resource="Property" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ClassName	StandardName	VisibleName	Description	TableVersion	TableDate	UpdateVersion	UpdateDate	</COLUMNS>
<DATA>	COM	MRIS Commercial	MRIS Commercial	MRIS_COM	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	LOT	MRIS Lot Land	MRIS Lot Land	MRIS_LOT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	MF	MRIS Multi-Family	MRIS Multi-Family	MRIS_MF	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	RES	MRIS Residential	MRIS Residential	MRIS_RES	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	ALL	MRIS All	MRIS All	MRIS_ALL	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	RESO_PROP_2012_05	RESO Prop 2012 05	RESO Prop 2012 05	Copyright 2012 RESO.  This software product includes software or other works developed by RESO and some of its contributors, subject to the RESO End User License published at www.reso.org	0.0.4	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
</METADATA-CLASS>
`

func TestParseClass(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + class + retsEnd))

	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)
	verifyParseClass(t, *ms)
}

func verifyParseClass(t *testing.T, cm CompactMetadata) {
	mdata := cm.Elements["METADATA-CLASS"][0]

	assert.Equal(t, "Property", mdata.Attr["Resource"])
	assert.Equal(t, mdata.Attr["Version"], "1.12.29")
	assert.Equal(t, len(mdata.CompactRows), 6)

	indexer := mdata.Indexer()
	var row []Row
	mdata.Rows(func(i int, r Row) {
		row = append(row, r)
	})

	val, ok := indexer("ClassName", row[5])
	assert.True(t, ok, "ClassName not found in row")
	assert.Equal(t, "RESO_PROP_2012_05", val)

	val, ok = indexer("TableDate", row[0])
	assert.True(t, ok, "TableDate not found in row")
	assert.Equal(t, "Tue, 3 Sep 2013 00:00:00 GMT", val)

	val, ok = indexer("VisibleName", row[2])
	assert.True(t, ok, "VisibleName not found in row")
	assert.Equal(t, "MRIS Multi-Family", val)
}

var table = `<METADATA-TABLE Resource="Agent" Class="ActiveAgent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	SystemName	StandardName	LongName	DBName	ShortName	MaximumLength	DataType	Precision	Searchable	Interpretation	Alignment	UseSeparator	EditMaskID	LookupName	MaxSelect	Units	Index	Minimum	Maximum	Default	Required	SearchHelpID	Unique	</COLUMNS>
<DATA>	AgentListingServiceName		ListingServiceName	X49076033	ListingServiceName	4000	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentKey		AgentKey	X74130	AgentKey	15	Long	0	1	Number	Right	0					1			0			1	</DATA>
<DATA>	AgentID		AgentID	X74134	AgentID	30	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentNationalID		NationalID	X74146	NationalID	20	Character		1		Left	0					1			0			0	</DATA>
</METADATA-TABLE>
`

func TestParseTable(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + table + retsEnd))

	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)
	verifyParseTable(t, *ms)
}

func verifyParseTable(t *testing.T, cm CompactMetadata) {
	mdata := cm.Elements["METADATA-TABLE"][0]

	assert.Equal(t, "ActiveAgent", mdata.Attr["Class"])
	assert.Equal(t, "Agent", mdata.Attr["Resource"])
	assert.Equal(t, "1.12.29", mdata.Attr["Version"])
	assert.Equal(t, len(mdata.CompactRows), 4)

	indexer := mdata.Indexer()
	var row []Row
	mdata.Rows(func(i int, r Row) {
		row = append(row, r)
	})

	val, ok := indexer("SystemName", row[0])
	assert.True(t, ok, "SystemName not found in row")
	assert.Equal(t, "AgentListingServiceName", val)

	val, ok = indexer("Unique", row[3])
	assert.True(t, ok, "Unique not found in row")
	assert.Equal(t, "0", val)
}

var lookup = `<METADATA-LOOKUP Resource="TaxHistoricalDesignation" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	LookupName	VisibleName	Version	Date	</COLUMNS>
<DATA>	COUNTIES_OR_REGIONS	Counties or Regions	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	TAX_HISTORIC_DESIGNATION_TYPES	Tax Historic Designation Types	1.12.6	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	SYSTEM_LOCALES	System Locales	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	SUB_SYSTEM_LOCALES	Sub System Locales	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-LOOKUP>
`

func TestParseLookup(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + lookup + retsEnd))

	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)
	verifyParseLookup(t, *ms)
}

func verifyParseLookup(t *testing.T, cm CompactMetadata) {
	mdata := cm.Elements["METADATA-LOOKUP"][0]

	assert.Equal(t, "TaxHistoricalDesignation", mdata.Attr["Resource"])
	assert.Equal(t, "1.12.29", mdata.Attr["Version"])
	assert.Equal(t, len(mdata.CompactRows), 4)

	indexer := mdata.Indexer()
	var row []Row
	mdata.Rows(func(i int, r Row) {
		row = append(row, r)
	})

	val, ok := indexer("LookupName", row[0])
	assert.True(t, ok, "LookupName not found in row")
	assert.Equal(t, "COUNTIES_OR_REGIONS", val)

	val, ok = indexer("Date", row[3])
	assert.True(t, ok, "Date not found in row")
	assert.Equal(t, "Tue, 3 Sep 2013 00:00:00 GMT", val)
}

var lookupType = `<METADATA-LOOKUP_TYPE Resource="TaxHistoricalDesignation" Lookup="COUNTIES_OR_REGIONS" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	LongValue	ShortValue	Value	</COLUMNS>
<DATA>	ALEUTIANS WEST-AK	ALEUTIANS WEST	85014594158	</DATA>
<DATA>	ATASCOSA-TX	ATASCOSA	85014594154	</DATA>
<DATA>	BROOMFIELD-CO	BROOMFIELD	85014594156	</DATA>
<DATA>	CLARK-CA	CLARK	50041199774	</DATA>
</METADATA-LOOKUP_TYPE>
`

func TestParseLookupType(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + lookupType + retsEnd))

	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)
	verifyParseLookupType(t, *ms)
}

func verifyParseLookupType(t *testing.T, cm CompactMetadata) {
	mdata := cm.Elements["METADATA-LOOKUP_TYPE"][0]

	assert.Equal(t, "TaxHistoricalDesignation", mdata.Attr["Resource"])
	assert.Equal(t, "COUNTIES_OR_REGIONS", mdata.Attr["Lookup"])

	assert.Equal(t, "1.12.29", mdata.Attr["Version"])
	assert.Equal(t, len(mdata.CompactRows), 4)

	indexer := mdata.Indexer()
	var row []Row
	mdata.Rows(func(i int, r Row) {
		row = append(row, r)
	})

	val, ok := indexer("LongValue", row[2])
	assert.True(t, ok, "LongValue not found in row")
	assert.Equal(t, "BROOMFIELD-CO", val)

	val, ok = indexer("ShortValue", row[3])
	assert.True(t, ok, "ShortValue not found in row")
	assert.Equal(t, "CLARK", val)
}

func TestParseMetadata(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + system + resource + class + table + lookup + lookupType + retsEnd))
	ms, err := ParseMetadataCompactResult(body)
	assert.Nil(t, err)

	verifySystem(t, *ms)
	verifyParseResources(t, *ms)
	verifyParseClass(t, *ms)
	verifyParseTable(t, *ms)
	verifyParseLookup(t, *ms)
	verifyParseLookupType(t, *ms)
}
