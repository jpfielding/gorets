package rets

import (
	"encoding/xml"
	"time"
)

var example = `
<RETS ReplyCode="0" ReplyText="Operation Successful">
<METADATA>
	<METADATA-SYSTEM Date="2016-08-15T14:47:13Z" Version="1.00.04717">
		<System/>
	</METADATA-SYSTEM>
	<METADATA-RESOURCE Date="2016-08-15T14:47:13Z" Version="1.00.03519">
		<Resource>
		<METADATA-CLASS Resource="Agent" Date="2016-07-11T18:22:38Z" Version="1.00.00045">
			<Class>
			<METADATA-TABLE Class="Agent" Resource="Agent" Date="2016-07-11T18:22:38Z" Version="1.00.00040">
				<Field/>
			</METADATA-TABLE>
			</Class>
		</METADATA-CLASS>
		<METADATA-OBJECT Resource="Agent" Date="2014-06-04T15:55:22Z" Version="1.00.00001">
			<Object/>
		</METADATA-OBJECT>
		<METADATA-LOOKUP Resource="Agent" Date="2016-08-08T16:00:15Z" Version="1.00.00110">
			<LookupType>
			<METADATA-LOOKUP_TYPE Lookup="AgentStatus" Resource="Agent" Date="2014-02-24T15:55:08Z" Version="1.00.00004">
				<Lookup/>
			</METADATA-LOOKUP_TYPE>
			</LookupType>
		</METADATA-LOOKUP>
		</Resource>
	</METADATA-RESOURCE>
</METADATA>
</RETS>`

// XMLMetadataVersion ...
type XMLMetadataVersion struct {
	Date    time.Time `xml:"Date,attr"`
	Version string    `xml:"Version,attr"`
}

// StandardXML ...
type StandardXML struct {
	XMLName   xml.Name    `xml:"RETS"`
	ReplyCode int         `xml:"ReplyCode"`
	ReplyText string      `xml:"ReplyText"`
	Metadata  XMLMetadata `xml:"METADATA"`
}

// XMLMetadata ...
type XMLMetadata struct {
	System XMLMetadataSystem `xml:"METADATA-SYSTEM"`
}

// XMLMetadataSystem ...
type XMLMetadataSystem struct {
	XMLMetadataVersion
	System MetadataSystem `xml:"System"`
}

// MetadataSystem ...
type MetadataSystem struct {
	ID                    string                `xml:"SystemID"`
	Description           string                `xml:"SystemDescription"`
	Comment               string                `xml:"Comment"`
	XMLMetadataForeignKey XMLMetadataForeignKey `xml:"METADATA-FOREIGN_KEY"`
	XMLMetadataResource   XMLMetadataResource   `xml:"METADATA-RESOURCE"`
}

// XMLMetadataForeignKey ...
type XMLMetadataForeignKey struct {
	ForeignKey []MetadataForeignKey `xml:"ForeignKey"`
}

// MetadataForeignKey ...
type MetadataForeignKey struct {
	// <ForeignKeyID>18</ForeignKeyID>
	ForeignKeyID string `xml:"ForeignKeyID"`
	// <ParentResourceID>Property</ParentResourceID>
	ParentResourceID string `xml:"ParentResourceID"`
	// <ParentClassID>Listing</ParentClassID>
	ParentClassID string `xml:"ParentClassID"`
	// <ParentSystemName>Matrix_Unique_ID</ParentSystemName>
	ParentSystemName string `xml:"ParentSystemName"`
	// <ChildResourceID>PropertySubTable</ChildResourceID>
	ChildResourceID string `xml:"ChildResourceID"`
	// <ChildClassID>Room</ChildClassID>
	ChildClassID string `xml:"ChildClassID"`
	// <ChildSystemName>Listing_MUI</ChildSystemName>
	ChildSystemName string `xml:"ChildSystemName"`
}

// XMLMetadataResource ...
type XMLMetadataResource struct {
	XMLMetadataVersion
	Resource []MetadataResource `xml:"Resource"`
}

// MetadataResource ...
type MetadataResource struct {
	// 	<ResourceID>Agent</ResourceID>
	ResourceID string `xml:"ResourceID"`
	// <StandardName>Agent</StandardName>
	StandardName string `xml:"StandardName"`
	// <VisibleName>Agent</VisibleName>
	VisibleName string `xml:"VisibleName"`
	// <Description>Agent</Description>
	Description string `xml:"Description"`
	// <KeyField>Matrix_Unique_ID</KeyField>
	KeyField string `xml:"KeyField"`
	// <ClassCount>1</ClassCount>
	ClassCount int `xml:"ClassCount"`
	// <ClassVersion>1.00.00045</ClassVersion>
	ClassVersion string `xml:"ClassVersion"`
	// <ClassDate>2016-07-11T18:22:38Z</ClassDate>
	ClassDate time.Time `xml:"ClassDate"`
	// <ObjectVersion>1.00.00001</ObjectVersion>
	ObjectVersion string `xml:"ObjectVersion"`
	// <ObjectDate>2014-06-04T15:55:22Z</ObjectDate>
	ObjectDate time.Time `xml:"ObjectDate"`
	// <SearchHelpVersion>1.00.00000</SearchHelpVersion>
	SearchHelpVersion string `xml:"SearchHelpVersion"`
	// <SearchHelpDate>2014-02-05T19:15:32Z</SearchHelpDate>
	SearchHelpDate time.Time `xml:"SearchHelpDate"`
	// <EditMaskVersion>1.00.00000</EditMaskVersion>
	EditMaskVersion string `xml:"EditMaskVersion"`
	// <EditMaskDate>2014-02-05T19:15:32Z</EditMaskDate>
	EditMaskDate time.Time `xml:"EditMaskDate"`
	// <LookupVersion>1.00.00110</LookupVersion>
	LookupVersion string `xml:"LookupVersion"`
	// <LookupDate>2016-08-08T16:00:15Z</LookupDate>
	LookupDate time.Time `xml:"LookupDate"`
	// <UpdateHelpVersion>1.00.00000</UpdateHelpVersion>
	UpdateHelpVersion string `xml:"UpdateHelpVersion"`
	// <UpdateHelpDate>2014-02-05T19:15:32Z</UpdateHelpDate>
	UpdateHelpDate time.Time `xml:"UpdateHelpDate"`
	// <ValidationExpressionVersion>1.00.00000</ValidationExpressionVersion>
	ValidationExpressionVersion string `xml:"ValidationExpressionVersion"`
	// <ValidationExpressionDate>2014-02-05T19:15:32Z</ValidationExpressionDate>
	ValidationExpressionDate time.Time `xml:"ValidationExpressionDate"`
	// <ValidationLookupVersion>1.00.00000</ValidationLookupVersion>
	ValidationLookupVersion string `xml:"ValidationLookupVersion"`
	// <ValidationLookupDate>2014-02-05T19:15:32Z</ValidationLookupDate>
	ValidationLookupDate time.Time `xml:"ValidationLookupDate"`
	// <ValidationExternalVersion>1.00.00000</ValidationExternalVersion>
	ValidationExternalVersion string `xml:"ValidationExternalVersion"`
	// <ValidationExternalDate>2014-02-05T19:15:32Z</ValidationExternalDate>
	ValidationExternalDate time.Time `xml:"ValidationExternalDate"`
	// <X-DisplayOrder/>
	XDisplayOrder int `xml:"X-DisplayOrder"`

	// the resource children
	MetadataClass  XMLMetadataClass  `xml:"METADATA-CLASS"`
	MetadataObject XMLMetadataObject `xml:"METADATA-OBJECT"`
	MetadataLookup XMLMetadataLookup `xml:"METADATA-LOOKUP"`
}

// XMLMetadataClass ...
type XMLMetadataClass struct {
	XMLMetadataVersion
	Resource string          `xml:"Resource,attr"`
	Class    []MetadataClass `xml:"Class"`
}

// MetadataClass ...
type MetadataClass struct {
	// <ClassName>Agent</ClassName>
	ClassName string `xml:"ClassName"`
	// <StandardName/>
	StandardName string `xml:"StandardName"`
	// <VisibleName>Agent</VisibleName>
	VisibleName string `xml:"VisibleName"`
	// <Description>Agent</Description>
	Description string `xml:"Description"`
	// <TableVersion>1.00.00040</TableVersion>
	TableVersion string `xml:"TableVersion"`
	// <TableDate>2016-07-11T18:22:38Z</TableDate>
	TableDate time.Time `xml:"TableDate"`
	// <UpdateVersion>1.00.00000</UpdateVersion>
	UpdateVersion string `xml:"UpdateVersion"`
	// <UpdateDate>2014-02-05T19:16:09Z</UpdateDate>
	UpdateDate string `xml:"UpdateDate"`
	// <X-DisplayOrder/>
	XDisplayOrder int `xml:"X-DisplayOrder"`
	// <ColumnGroupVersion/>
	ColumnGroupVersion string `xml:"ColumnGroupVersion"`
	// <ColumnGroupDate/>
	ColumnGroupDate time.Time `xml:"ColumnGroupDate"`
	// <ColumnGroupSetVersion/>
	ColumnGroupSetVersion string `xml:"ColumnGroupSetVersion"`
	// <ColumnGroupSetDate/>
	ColumnGroupSetDate time.Time `xml:"ColumnGroupSetDate"`
}

// XMLMetadataTable ...
type XMLMetadataTable struct {
	XMLMetadataVersion
	Resource string          `xml:"Resource,attr"`
	Class    string          `xml:"Class,attr"`
	Table    []MetadataTable `xml:"Table"`
}

// MetadataTable ...
type MetadataTable struct {
	// <SystemName>AgentLSCcode</SystemName>
	SystemName string `xml:"SystemName"`
	// <StandardName/>
	StandardName string `xml:"StandardName"`
	// <LongName>Agent LS Ccode</LongName>
	LongName string `xml:"LongName"`
	// <DBName>R574</DBName>
	DBName string `xml:"DBName"`
	// <ShortName>Agent LSC Code</ShortName>
	ShortName string `xml:"ShortName"`
	// <MaximumLength>75</MaximumLength>
	MaximumLength int `xml:"MaximumLength"`
	// <DataType>Character</DataType>
	DataType string `xml:"DataType"`
	// <Precision/>
	Precision string `xml:"Precision"`
	// <Searchable>1</Searchable>
	Searchable int `xml:"Searchable"`
	// <Interpretation>Lookup</Interpretation>
	Interpretation string `xml:"Interpretation"`
	// <Alignment>Left</Alignment>
	Alignment string `xml:"Alignment"`
	// <UseSeparator>0</UseSeparator>
	UseSeparator int `xml:"UseSeparator"`
	// <EditMaskID/>
	EditMaskID string `xml:"EditMaskID"`
	// <LookupName>BoardID</LookupName>
	LookupName string `xml:"LookupName"`
	// <MaxSelect>1</MaxSelect>
	MaxSelect int `xml:"MaxSelect"`
	// <Units/>
	Units string `xml:"Units"`
	// <Index>1</Index>
	Index int `xml:"Index"`
	// <Minimum/>
	Minimum string `xml:"Minimum"`
	// <Maximum/>
	Maximum string `xml:"Maxiumum"`
	// <Default/>
	Default string `xml:"Default"`
	// <Required>0</Required>
	Required int `xml:"Required"`
	// <SearchHelpID/>
	SearchHelpID string `xml:"SearchHelpID"`
	// <Unique>0</Unique>
	Unique int `xml:"Unique"`
}

// XMLMetadataObject ...
type XMLMetadataObject struct {
	XMLMetadataVersion
	Resource string           `xml:"Resource,attr"`
	Object   []MetadataObject `xml:"Object"`
}

// MetadataObject ...
type MetadataObject struct {
	// <ObjectType>AgentPhoto</ObjectType>
	ObjectType string `xml:"ObjectType"`
	// <StandardName/>
	StandardName string `xml:"ObjectType"`
	// <MimeType>image/jpeg</MimeType>
	MimeType string `xml:"MimeType"`
	// <Description>AgentPhoto</Description>
	Description string `xml:"Description"`
}

// XMLMetadataLookup ...
type XMLMetadataLookup struct {
	XMLMetadataVersion
	Resource   string               `xml:"Resource,attr"`
	LookupType []MetadataLookupType `xml:"LookupType"`
}

// MetadataLookupType ...
type MetadataLookupType struct {
	// <LookupName>AgentStatus</LookupName>
	LookupName string `xml:"LookupName"`
	// <VisibleName>Agent Status</VisibleName>
	VisibleName string `xml:"VisibleName"`
	// <LookupTypeVersion>1.00.00004</LookupTypeVersion>
	LookupTypeVersion string `xml:"LookupTypeVersion"`
	// <LookupTypeDate>2014-02-24T15:55:08Z</LookupTypeDate>
	LookupTypeDate time.Time `xml:"LookupTypeDate"`

	LookupType []XMLMetadataLookupType `xml:"METADATA-LOOKUP_TYPE"`
}

// XMLMetadataLookupType ...
type XMLMetadataLookupType struct {
	XMLMetadataVersion
	Resource string           `xml:"Resource,attr"`
	Lookup   string           `xml:"Resource,attr"`
	Lookups  []MetadataLookup `xml:"Lookup"`
}

// MetadataLookup ...
type MetadataLookup struct {
	// <LongValue>Active</LongValue>
	LongValue string `xml:"LongValue"`
	// <ShortValue>A</ShortValue>
	ShortValue string `xml:"ShortValue"`
	// <Value>A</Value>
	Value string `xml:"Value"`
	// <X-DisplayOrder/>
	XDisplayOrder int `xml:"X-DisplayOrder"`
}
