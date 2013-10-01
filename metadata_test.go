/**
	parsing the 'login' action from RETS
 */
package gorets

import (
	"testing"
)


var metadataSystem string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-SYSTEM Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<SYSTEM SystemID="SIRM" SystemDescription="MLS System"/>
<COMMENTS>
The System is provided to you by Systems.
</COMMENTS>
</METADATA-SYSTEM>
</RETS>`

func TestParseSystem(t *testing.T) {
	ms, err := parseMSystem([]byte(metadataSystem))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEquals(t, "bad version", ms.Version, "1.12.30")
	AssertEquals(t, "bad comments", ms.Comments, "The System is provided to you by Systems.")
}

var resource string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-RESOURCE Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ResourceID	StandardName	VisibleName	Description	KeyField	ClassCount	ClassVersion	ClassDate	ObjectVersion	ObjectDate	SearchHelpVersion	SearchHelpDate	EditMaskVersion	EditMaskDate	LookupVersion	LookupDate	UpdateHelpVersion	UpdateHelpDate	ValidationExpressionVersion	ValidationExpressionDate	ValidationLookupVersion	ValidationLookupDate	ValidationExternalVersion	ValidationExternalDate	</COLUMNS>
<DATA>	ActiveAgent	ActiveAgent	Agent	ActiveAgent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	Agent	Agent	Agent	Agent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-RESOURCE>
	</RETS>
	`

func TestParseResources(t *testing.T) {
	ms, err := parseMResources([]byte(resource))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}

	AssertEquals(t, "bad version", "1.12.30", ms.MData.Version)
	AssertEqualsInt(t, "wrong number of resources", len(ms.MData.Rows), 2)

	indexer := ms.MData.Indexer()

	AssertEquals(t, "bad value", "ActiveAgent", indexer("ResourceID",0))
	AssertEquals(t, "bad value", "Tue, 3 Sep 2013 00:00:00 GMT", indexer("ValidationExternalDate",1))
}


var class string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-CLASS Resource="Property" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ClassName	StandardName	VisibleName	Description	TableVersion	TableDate	UpdateVersion	UpdateDate	</COLUMNS>
<DATA>	COM	MRIS Commercial	MRIS Commercial	MRIS_COM	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	LOT	MRIS Lot Land	MRIS Lot Land	MRIS_LOT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	MF	MRIS Multi-Family	MRIS Multi-Family	MRIS_MF	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	RES	MRIS Residential	MRIS Residential	MRIS_RES	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	ALL	MRIS All	MRIS All	MRIS_ALL	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	RESO_PROP_2012_05	RESO Prop 2012 05	RESO Prop 2012 05	Copyright 2012 RESO.  This software product includes software or other works developed by RESO and some of its contributors, subject to the RESO End User License published at www.reso.org	0.0.4	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
</METADATA-CLASS>
</RETS>

`

func TestParseClass(t *testing.T) {
	ms, err := parseMClasses([]byte(class))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEquals(t, "bad version", ms.MData.Version, "1.12.29")
	AssertEqualsInt(t, "wrong number of resources", len(ms.MData.Rows), 6)

	indexer := ms.MData.Indexer()

	AssertEquals(t, "bad value", "RESO_PROP_2012_05", indexer("ClassName",5))
	AssertEquals(t, "bad value", "Tue, 3 Sep 2013 00:00:00 GMT", indexer("TableDate",0))
	AssertEquals(t, "bad value", "MRIS Multi-Family", indexer("VisibleName",2))
}

var table string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-TABLE Resource="ActiveAgent" Class="ActiveAgent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	SystemName	StandardName	LongName	DBName	ShortName	MaximumLength	DataType	Precision	Searchable	Interpretation	Alignment	UseSeparator	EditMaskID	LookupName	MaxSelect	Units	Index	Minimum	Maximum	Default	Required	SearchHelpID	Unique	</COLUMNS>
<DATA>	AgentListingServiceName		ListingServiceName	X49076033	ListingServiceName	4000	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentKey		AgentKey	X74130	AgentKey	15	Long	0	1	Number	Right	0					1			0			1	</DATA>
<DATA>	AgentID		AgentID	X74134	AgentID	30	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentNationalID		NationalID	X74146	NationalID	20	Character		1		Left	0					1			0			0	</DATA>
</METADATA-TABLE>
</RETS>
`

func TestParseTable(t *testing.T) {
	ms, err := parseMTables([]byte(table))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEquals(t, "bad version", "1.12.29", ms.MData.Version)
	AssertEqualsInt(t, "wrong number of resources", len(ms.MData.Rows), 4)

	indexer := ms.MData.Indexer()

	AssertEquals(t, "bad value", "AgentListingServiceName", indexer("SystemName",0))
	AssertEquals(t, "bad value", "0", indexer("Unique",3))
}

var lookup string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-LOOKUP Resource="TaxHistoricalDesignation" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	LookupName	VisibleName	Version	Date	</COLUMNS>
<DATA>	COUNTIES_OR_REGIONS	Counties or Regions	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	TAX_HISTORIC_DESIGNATION_TYPES	Tax Historic Designation Types	1.12.6	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	SYSTEM_LOCALES	System Locales	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	SUB_SYSTEM_LOCALES	Sub System Locales	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-LOOKUP>
</RETS>
`

func TestParseLookup(t *testing.T) {
	ms, err := parseMLookups([]byte(lookup))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEquals(t, "bad version", "1.12.29", ms.MData.Version)
	AssertEqualsInt(t, "wrong number of resources", len(ms.MData.Rows), 4)

	indexer := ms.MData.Indexer()

	AssertEquals(t, "bad value", "COUNTIES_OR_REGIONS", indexer("LookupName",0))
	AssertEquals(t, "bad value", "Tue, 3 Sep 2013 00:00:00 GMT", indexer("Date",3))
}

var lookupType string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-LOOKUP_TYPE Resource="TaxHistoricalDesignation" Lookup="COUNTIES_OR_REGIONS" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	LongValue	ShortValue	Value	</COLUMNS>
<DATA>	ALEUTIANS WEST-AK	ALEUTIANS WEST	85014594158	</DATA>
<DATA>	ATASCOSA-TX	ATASCOSA	85014594154	</DATA>
<DATA>	BROOMFIELD-CO	BROOMFIELD	85014594156	</DATA>
<DATA>	CLARK-CA	CLARK	50041199774	</DATA>
</METADATA-LOOKUP_TYPE>
</RETS>
`

func TestParseLookupType(t *testing.T) {
	ms, err := parseMLookupTypes([]byte(lookupType))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEquals(t, "bad version", "1.12.29", ms.MData.Version)
	AssertEqualsInt(t, "wrong number of resources", len(ms.MData.Rows), 4)

	indexer := ms.MData.Indexer()

	AssertEquals(t, "bad value", "BROOMFIELD-CO", indexer("LongValue",2))
	AssertEquals(t, "bad value", "CLARK", indexer("ShortValue",3))
}


