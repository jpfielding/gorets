package metadata

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
)

var mresource MResource = MResource{
	Version: "39.19.84596",
	Date:    "Tue, 02 Dec 2014 01:36:36 GMT",
	Resource: []Resource{
		Resource{
			ResourceID:      "Office",
			StandardName:    "Office",
			VisibleName:     "Office",
			Description:     "Office",
			KeyField:        "sysid",
			ClassCount:      1,
			ClassVersion:    "38.31.12348",
			ClassDate:       "Thu, 21 Aug 2014 04:05:48 GMT",
			ObjectVersion:   "28.92.81887",
			ObjectDate:      "Wed, 31 Aug 2011 04:04:47 GMT",
			EditMaskVersion: "38.31.12348",
			EditMaskDate:    "Thu, 21 Aug 2014 04:05:48 GMT",
			LookupVersion:   "24.03.55825",
			LookupDate:      "Thu, 11 Feb 2010 02:30:25 GMT",
		},
		Resource{
			ResourceID:      "Property",
			StandardName:    "Property",
			VisibleName:     "Property",
			Description:     "Property",
			KeyField:        "sysid",
			ClassCount:      4,
			ClassVersion:    "39.19.57282",
			ClassDate:       "Mon, 01 Dec 2014 13:01:22 GMT",
			ObjectVersion:   "37.33.45209",
			ObjectDate:      "Wed, 30 Apr 2014 03:00:09 GMT",
			EditMaskVersion: "39.19.57282",
			EditMaskDate:    "Mon, 01 Dec 2014 13:01:22 GMT",
			LookupVersion:   "39.11.47282",
			LookupDate:      "Sat, 22 Nov 2014 09:01:22 GMT",
		},
	},
}

var mresourceCompact = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA-RESOURCE Date="Tue, 02 Dec 2014 01:36:36 GMT" Version="39.19.84596">
    <COLUMNS>	ResourceID	StandardName	VisibleName	Description	KeyField	ClassCount	ClassVersion	ClassDate	ObjectVersion	ObjectDate	SearchHelpVersion	SearchHelpDate	EditMaskVersion	EditMaskDate	LookupVersion	LookupDate	UpdateHelpVersion	UpdateHelpDate	ValidationExpressionVersion	ValidationExpressionDate	ValidationLookupVersion	ValidationLookupDate	ValidationExternalVersion	ValidationExternalDate	</COLUMNS>
    <DATA>	Office	Office	Office	Office	sysid	1	38.31.12348	Thu, 21 Aug 2014 04:05:48 GMT	28.92.81887	Wed, 31 Aug 2011 04:04:47 GMT			38.31.12348	Thu, 21 Aug 2014 04:05:48 GMT	24.03.55825	Thu, 11 Feb 2010 02:30:25 GMT									</DATA>
    <DATA>	Property	Property	Property	Property	sysid	4	39.19.57282	Mon, 01 Dec 2014 13:01:22 GMT	37.33.45209	Wed, 30 Apr 2014 03:00:09 GMT			39.19.57282	Mon, 01 Dec 2014 13:01:22 GMT	39.11.47282	Sat, 22 Nov 2014 09:01:22 GMT									</DATA>
    <DATA>	User	Agent	User	User	sysid	1	38.31.12348	Thu, 21 Aug 2014 04:05:48 GMT	11.96.50294	Thu, 13 Feb 2003 10:24:54 GMT			38.31.12348	Thu, 21 Aug 2014 04:05:48 GMT	37.33.16423	Tue, 29 Apr 2014 23:00:23 GMT									</DATA>
  </METADATA-RESOURCE>
</RETS>
`

func TestResourceCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(mresourceCompact)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-RESOURCE")
	Ok(t, err)

	found := MResource{}
	err = found.InitFromCompact(parser, start)

	Ok(t, err)

	Equals(t, mresource, found)
}

var mresourceXml = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Success">
  <METADATA>
    <METADATA-RESOURCE Date="Tue, 02 Dec 2014 01:36:36 GMT" System="RETS" Version="39.19.84596">
      <Resource>
        <ResourceID>Office</ResourceID>
        <StandardName>Office</StandardName>
        <VisibleName>Office</VisibleName>
        <Description>Office</Description>
        <KeyField>sysid</KeyField>
        <ClassCount>1</ClassCount>
        <ClassVersion>38.31.12348</ClassVersion>
        <ClassDate>Thu, 21 Aug 2014 04:05:48 GMT</ClassDate>
        <ObjectVersion>28.92.81887</ObjectVersion>
        <ObjectDate>Wed, 31 Aug 2011 04:04:47 GMT</ObjectDate>
        <SearchHelpVersion></SearchHelpVersion>
        <SearchHelpDate></SearchHelpDate>
        <EditMaskVersion>38.31.12348</EditMaskVersion>
        <EditMaskDate>Thu, 21 Aug 2014 04:05:48 GMT</EditMaskDate>
        <LookupVersion>24.03.55825</LookupVersion>
        <LookupDate>Thu, 11 Feb 2010 02:30:25 GMT</LookupDate>
        <UpdateHelpVersion></UpdateHelpVersion>
        <UpdateHelpDate></UpdateHelpDate>
        <ValidationExpressionVersion></ValidationExpressionVersion>
        <ValidationExpressionDate></ValidationExpressionDate>
        <ValidationLookupVersion></ValidationLookupVersion>
        <ValidationLookupDate></ValidationLookupDate>
        <ValidationExternalVersion></ValidationExternalVersion>
        <ValidationExternalDate></ValidationExternalDate>
      </Resource>
      <Resource>
        <ResourceID>Property</ResourceID>
        <StandardName>Property</StandardName>
        <VisibleName>Property</VisibleName>
        <Description>Property</Description>
        <KeyField>sysid</KeyField>
        <ClassCount>4</ClassCount>
        <ClassVersion>39.19.57282</ClassVersion>
        <ClassDate>Mon, 01 Dec 2014 13:01:22 GMT</ClassDate>
        <ObjectVersion>37.33.45209</ObjectVersion>
        <ObjectDate>Wed, 30 Apr 2014 03:00:09 GMT</ObjectDate>
        <SearchHelpVersion></SearchHelpVersion>
        <SearchHelpDate></SearchHelpDate>
        <EditMaskVersion>39.19.57282</EditMaskVersion>
        <EditMaskDate>Mon, 01 Dec 2014 13:01:22 GMT</EditMaskDate>
        <LookupVersion>39.11.47282</LookupVersion>
        <LookupDate>Sat, 22 Nov 2014 09:01:22 GMT</LookupDate>
        <UpdateHelpVersion></UpdateHelpVersion>
        <UpdateHelpDate></UpdateHelpDate>
        <ValidationExpressionVersion></ValidationExpressionVersion>
        <ValidationExpressionDate></ValidationExpressionDate>
        <ValidationLookupVersion></ValidationLookupVersion>
        <ValidationLookupDate></ValidationLookupDate>
        <ValidationExternalVersion></ValidationExternalVersion>
        <ValidationExternalDate></ValidationExternalDate>
      </Resource>
    </METADATA-RESOURCE>
</RETS>
`

func TestResourceXml(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(mresourceXml)))

	parser := xml.NewDecoder(body)

	start, err := AdvanceToStartElem(parser, "METADATA-RESOURCE")
	Ok(t, err)

	found := MResource{}
	err = found.InitFromXml(parser, start)

	Ok(t, err)

	Equals(t, mresource, found)
}
