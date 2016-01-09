/**
parsing the 'login' action from RETS
*/
package client

import (
	"bytes"
	"io/ioutil"
	"testing"

	testutils "github.com/jpfielding/gorets/testutils"
)

var retsStart = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">`
var retsEnd = `</RETS>`

var system = `<METADATA-SYSTEM Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<SYSTEM SystemID="SIRM" SystemDescription="MLS System"/>
<COMMENTS>
The System is provided to you by Systems.
</COMMENTS>
</METADATA-SYSTEM>`

func TestSystem(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + system + retsEnd)))

	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)
	verifySystem(t, *ms)
}

func verifySystem(t *testing.T, ms Metadata) {
	testutils.Assert(t, ms.System.Version == "1.12.30", "bad version")
	testutils.Assert(t, ms.System.Comments == "The System is provided to you by Systems.", "bad comments")
}

var resource = `<METADATA-RESOURCE Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ResourceID	StandardName	VisibleName	Description	KeyField	ClassCount	ClassVersion	ClassDate	ObjectVersion	ObjectDate	SearchHelpVersion	SearchHelpDate	EditMaskVersion	EditMaskDate	LookupVersion	LookupDate	UpdateHelpVersion	UpdateHelpDate	ValidationExpressionVersion	ValidationExpressionDate	ValidationLookupVersion	ValidationLookupDate	ValidationExternalVersion	ValidationExternalDate	</COLUMNS>
<DATA>	ActiveAgent	ActiveAgent	Agent	ActiveAgent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	Agent	Agent	Agent	Agent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-RESOURCE>`

func TestParseResources(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + resource + retsEnd)))

	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)
	verifyParseResources(t, *ms)
}

func verifyParseResources(t *testing.T, ms Metadata) {
	resource := ms.find("METADATA-RESOURCE")[""]
	testutils.Equals(t, "1.12.30", resource.Version)
	testutils.Equals(t, len(resource.Rows), 2)

	indexer := resource.Indexer()

	testutils.Equals(t, "ActiveAgent", indexer("ResourceID", 0))
	testutils.Equals(t, "Tue, 3 Sep 2013 00:00:00 GMT", indexer("ValidationExternalDate", 1))
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
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + class + retsEnd)))

	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)
	verifyParseClass(t, *ms)
}

func verifyParseClass(t *testing.T, ms Metadata) {
	classes := ms.find("METADATA-CLASS")
	mdata := classes["Property"]

	testutils.Equals(t, mdata.Version, "1.12.29")
	testutils.Equals(t, len(mdata.Rows), 6)

	indexer := mdata.Indexer()

	testutils.Equals(t, "RESO_PROP_2012_05", indexer("ClassName", 5))
	testutils.Equals(t, "Tue, 3 Sep 2013 00:00:00 GMT", indexer("TableDate", 0))
	testutils.Equals(t, "MRIS Multi-Family", indexer("VisibleName", 2))
}

var table = `<METADATA-TABLE Resource="ActiveAgent" Class="ActiveAgent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	SystemName	StandardName	LongName	DBName	ShortName	MaximumLength	DataType	Precision	Searchable	Interpretation	Alignment	UseSeparator	EditMaskID	LookupName	MaxSelect	Units	Index	Minimum	Maximum	Default	Required	SearchHelpID	Unique	</COLUMNS>
<DATA>	AgentListingServiceName		ListingServiceName	X49076033	ListingServiceName	4000	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentKey		AgentKey	X74130	AgentKey	15	Long	0	1	Number	Right	0					1			0			1	</DATA>
<DATA>	AgentID		AgentID	X74134	AgentID	30	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentNationalID		NationalID	X74146	NationalID	20	Character		1		Left	0					1			0			0	</DATA>
</METADATA-TABLE>
`

func TestParseTable(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + table + retsEnd)))

	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)
	verifyParseTable(t, *ms)
}

func verifyParseTable(t *testing.T, ms Metadata) {
	tables := ms.find("METADATA-TABLE")
	mdata := tables["ActiveAgent:ActiveAgent"]

	testutils.Equals(t, "1.12.29", mdata.Version)
	testutils.Equals(t, len(mdata.Rows), 4)

	indexer := mdata.Indexer()

	testutils.Equals(t, "AgentListingServiceName", indexer("SystemName", 0))
	testutils.Equals(t, "0", indexer("Unique", 3))
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
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + lookup + retsEnd)))

	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)
	verifyParseLookup(t, *ms)
}

func verifyParseLookup(t *testing.T, ms Metadata) {
	lookups := ms.find("METADATA-LOOKUP")
	mdata := lookups["TaxHistoricalDesignation"]

	testutils.Equals(t, "1.12.29", mdata.Version)
	testutils.Equals(t, len(mdata.Rows), 4)

	indexer := mdata.Indexer()

	testutils.Equals(t, "COUNTIES_OR_REGIONS", indexer("LookupName", 0))
	testutils.Equals(t, "Tue, 3 Sep 2013 00:00:00 GMT", indexer("Date", 3))
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
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + lookupType + retsEnd)))

	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)
	verifyParseLookupType(t, *ms)
}
func verifyParseLookupType(t *testing.T, ms Metadata) {
	lookupTypes := ms.find("METADATA-LOOKUP_TYPE")
	mdata := lookupTypes["TaxHistoricalDesignation:COUNTIES_OR_REGIONS"]

	testutils.Equals(t, "1.12.29", mdata.Version)
	testutils.Equals(t, len(mdata.Rows), 4)

	indexer := mdata.Indexer()

	testutils.Equals(t, "BROOMFIELD-CO", indexer("LongValue", 2))
	testutils.Equals(t, "CLARK", indexer("ShortValue", 3))
}

func TestParseMetadata(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(retsStart + system + resource + class + table + lookup + lookupType + retsEnd)))
	ms, err := ParseMetadataCompactResult(body)
	testutils.Ok(t, err)

	verifySystem(t, *ms)
	verifyParseResources(t, *ms)
	verifyParseClass(t, *ms)
	verifyParseTable(t, *ms)
	verifyParseLookup(t, *ms)
	verifyParseLookupType(t, *ms)
}
