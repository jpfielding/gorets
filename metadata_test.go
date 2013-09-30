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
	if ms.Version != "1.12.30" {
		t.Errorf("wrong version: %s ", ms.Version)
	}
	if ms.Comments != "The System is provided to you by Systems." {
		t.Errorf("wrong comments: %s ", ms.Comments)
	}
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
	if ms.Version != "1.12.30" {
		t.Errorf("wrong version: %s ", ms.Version)
	}
	if len(ms.MResources) != 2 {
		t.Errorf("wrong number of resources: %s ", len(ms.MResources))
	}
	if ms.MResources[0].Fields["ResourceID"] != "ActiveAgent" {
		t.Errorf("wrong field value: %s ", ms.MResources[0].Fields["ResourceID"])
	}
	if ms.MResources[1].Fields["ValidationExternalDate"] != "Tue, 3 Sep 2013 00:00:00 GMT" {
		t.Errorf("wrong field value: %s ", ms.MResources[1].Fields["ValidationExternalDate"])
	}
}


var class string =
	`<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<METADATA-CLASS Resource="ActiveAgent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ClassName	StandardName	VisibleName	Description	TableVersion	TableDate	UpdateVersion	UpdateDate	</COLUMNS>
<DATA>	ActiveAgent	MRIS Active Agents	MRIS Active Agents	MRIS_ActiveAgent	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
</METADATA-CLASS>
</RETS>
`

func TestParseClass(t *testing.T) {
	ms, err := parseMClasses([]byte(class))
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	if ms.Version != "1.12.29" {
		t.Errorf("wrong version: %s ", ms.Version)
	}
	if len(ms.MClasses) != 1 {
		t.Errorf("wrong number of resources: %s ", len(ms.MClasses))
	}
	if ms.MClasses[0].Fields["ClassName"] != "ActiveAgent" {
		t.Errorf("wrong field value: %s ", ms.MClasses[0].Fields["ClassName"])
	}
	if ms.MClasses[0].Fields["TableDate"] != "Tue, 3 Sep 2013 00:00:00 GMT" {
		t.Errorf("wrong field value: %s ", ms.MClasses[0].Fields["TableDate"])
	}
	if ms.MClasses[0].Fields["UpdateDate"] != "" {
		t.Errorf("wrong field value: %s ", ms.MClasses[0].Fields["UpdateDate"])
	}
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
}


