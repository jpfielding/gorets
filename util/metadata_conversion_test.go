package util

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jpfielding/gorets/rets"
	"github.com/stretchr/testify/assert"
)

var retsStart = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">`
var retsEnd = `</RETS>`

var system = `<METADATA-SYSTEM Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<SYSTEM SystemID="SIRM" SystemDescription="MLS System"/>
<COMMENTS>The System is provided to you by Systems.</COMMENTS>
</METADATA-SYSTEM>`

var resource = `<METADATA-RESOURCE Version="1.12.30" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ResourceID	StandardName	VisibleName	Description	KeyField	ClassCount	ClassVersion	ClassDate	ObjectVersion	ObjectDate	SearchHelpVersion	SearchHelpDate	EditMaskVersion	EditMaskDate	LookupVersion	LookupDate	UpdateHelpVersion	UpdateHelpDate	ValidationExpressionVersion	ValidationExpressionDate	ValidationLookupVersion	ValidationLookupDate	ValidationExternalVersion	ValidationExternalDate	</COLUMNS>
<DATA>	Property	Property	Property	Property	ListingKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	Agent	Agent	Agent	Agent	AgentKey	1	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-RESOURCE>`

var class = `<METADATA-CLASS Resource="Agent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	ClassName	StandardName	VisibleName	Description	TableVersion	TableDate	UpdateVersion	UpdateDate	</COLUMNS>
<DATA>	ActiveAgent	SIRM Commercial	SIRM Commercial	SIRM_COM	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
<DATA>	Agent	SIRM Lot Land	SIRM Lot Land	SIRM_LOT	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT			</DATA>
</METADATA-CLASS>
`

var table = `<METADATA-TABLE Resource="Agent" Class="ActiveAgent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	SystemName	StandardName	LongName	DBName	ShortName	MaximumLength	DataType	Precision	Searchable	Interpretation	Alignment	UseSeparator	EditMaskID	LookupName	MaxSelect	Units	Index	Minimum	Maximum	Default	Required	SearchHelpID	Unique	</COLUMNS>
<DATA>	AgentListingServiceName		ListingServiceName	X49076033	ListingServiceName	4000	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentKey		AgentKey	X74130	AgentKey	15	Long	0	1	Number	Right	0					1			0			1	</DATA>
<DATA>	AgentID		AgentID	X74134	AgentID	30	Character		1		Left	0					1			0			0	</DATA>
<DATA>	AgentNationalID		NationalID	X74146	NationalID	20	Character		1		Left	0					1			0			0	</DATA>
</METADATA-TABLE>
`
var lookup = `<METADATA-LOOKUP Resource="Agent" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	LookupName	VisibleName	Version	Date	</COLUMNS>
<DATA>	COUNTIES_OR_REGIONS	Counties or Regions	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	TAX_HISTORIC_DESIGNATION_TYPES	Tax Historic Designation Types	1.12.6	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	SYSTEM_LOCALES	System Locales	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
<DATA>	SUB_SYSTEM_LOCALES	Sub System Locales	1.12.29	Tue, 3 Sep 2013 00:00:00 GMT	</DATA>
</METADATA-LOOKUP>
`

var lookupType = `<METADATA-LOOKUP_TYPE Resource="Agent" Lookup="COUNTIES_OR_REGIONS" Version="1.12.29" Date="Tue, 3 Sep 2013 00:00:00 GMT">
<COLUMNS>	LongValue	ShortValue	Value	</COLUMNS>
<DATA>	ALEUTIANS WEST-AK	ALEUTIANS WEST	85014594158	</DATA>
<DATA>	ATASCOSA-TX	ATASCOSA	85014594154	</DATA>
<DATA>	BROOMFIELD-CO	BROOMFIELD	85014594156	</DATA>
<DATA>	CLARK-CA	CLARK	50041199774	</DATA>
</METADATA-LOOKUP_TYPE>
`

func TestConvertCompactMetadata(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(retsStart + system + resource + class + table + lookup + lookupType + retsEnd))
	compact, err := rets.ParseMetadataCompactResult(body)
	assert.Nil(t, err)

	msystem, err := AsStandard(*compact).Convert()
	assert.Nil(t, err)

	assert.Equal(t, "MLS System", msystem.System.Description)
	assert.Equal(t, "1.12.30", string(msystem.Version))
	assert.Equal(t, "The System is provided to you by Systems.", msystem.System.Comments)

	mresource := msystem.System.MResource
	assert.Equal(t, "1.12.30", string(mresource.Version))
	assert.Equal(t, 2, len(mresource.Resource))

	mlookup := mresource.Resource[1].MLookup
	assert.Equal(t, "1.12.29", string(mlookup.Version))
	assert.Equal(t, 4, len(mlookup.Lookup))

	mlookupType := mlookup.Lookup[0].MLookupType
	assert.Equal(t, "1.12.29", string(mlookupType.Version))
	assert.Equal(t, 4, len(mlookupType.LookupType))

	agent := mresource.Resource[1]
	assert.Equal(t, 2, len(agent.MClass.Class))
	assert.Equal(t, "Agent", string(agent.ResourceID))
	assert.Equal(t, "1.12.29", string(agent.MClass.Version))

}
