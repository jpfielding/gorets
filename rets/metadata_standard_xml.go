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
	ForeignKeyID     string `xml:"ForeignKeyID"`
	ParentResourceID string `xml:"ParentResourceID"`
	ParentClassID    string `xml:"ParentClassID"`
	ParentSystemName string `xml:"ParentSystemName"`
	ChildResourceID  string `xml:"ChildResourceID"`
	ChildClassID     string `xml:"ChildClassID"`
	ChildSystemName  string `xml:"ChildSystemName"`
}

// XMLMetadataResource ...
type XMLMetadataResource struct {
	XMLMetadataVersion
	Resource []MetadataResource `xml:"Resource"`
}

// MetadataResource ...
type MetadataResource struct {
	ResourceID                  string    `xml:"ResourceID"`
	StandardName                string    `xml:"StandardName"`
	VisibleName                 string    `xml:"VisibleName"`
	Description                 string    `xml:"Description"`
	KeyField                    string    `xml:"KeyField"`
	ClassCount                  int       `xml:"ClassCount"`
	ClassVersion                string    `xml:"ClassVersion"`
	ClassDate                   time.Time `xml:"ClassDate"`
	ObjectVersion               string    `xml:"ObjectVersion"`
	ObjectDate                  time.Time `xml:"ObjectDate"`
	SearchHelpVersion           string    `xml:"SearchHelpVersion"`
	SearchHelpDate              time.Time `xml:"SearchHelpDate"`
	EditMaskVersion             string    `xml:"EditMaskVersion"`
	EditMaskDate                time.Time `xml:"EditMaskDate"`
	LookupVersion               string    `xml:"LookupVersion"`
	LookupDate                  time.Time `xml:"LookupDate"`
	UpdateHelpVersion           string    `xml:"UpdateHelpVersion"`
	UpdateHelpDate              time.Time `xml:"UpdateHelpDate"`
	ValidationExpressionVersion string    `xml:"ValidationExpressionVersion"`
	ValidationExpressionDate    time.Time `xml:"ValidationExpressionDate"`
	ValidationLookupVersion     string    `xml:"ValidationLookupVersion"`
	ValidationLookupDate        time.Time `xml:"ValidationLookupDate"`
	ValidationExternalVersion   string    `xml:"ValidationExternalVersion"`
	ValidationExternalDate      time.Time `xml:"ValidationExternalDate"`
	XDisplayOrder               int       `xml:"X-DisplayOrder"`

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
	ClassName             string    `xml:"ClassName"`
	StandardName          string    `xml:"StandardName"`
	VisibleName           string    `xml:"VisibleName"`
	Description           string    `xml:"Description"`
	TableVersion          string    `xml:"TableVersion"`
	TableDate             time.Time `xml:"TableDate"`
	UpdateVersion         string    `xml:"UpdateVersion"`
	UpdateDate            string    `xml:"UpdateDate"`
	XDisplayOrder         int       `xml:"X-DisplayOrder"`
	ColumnGroupVersion    string    `xml:"ColumnGroupVersion"`
	ColumnGroupDate       time.Time `xml:"ColumnGroupDate"`
	ColumnGroupSetVersion string    `xml:"ColumnGroupSetVersion"`
	ColumnGroupSetDate    time.Time `xml:"ColumnGroupSetDate"`
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
	SystemName     string `xml:"SystemName"`
	StandardName   string `xml:"StandardName"`
	LongName       string `xml:"LongName"`
	DBName         string `xml:"DBName"`
	ShortName      string `xml:"ShortName"`
	MaximumLength  int    `xml:"MaximumLength"`
	DataType       string `xml:"DataType"`
	Precision      string `xml:"Precision"`
	Searchable     int    `xml:"Searchable"`
	Interpretation string `xml:"Interpretation"`
	Alignment      string `xml:"Alignment"`
	UseSeparator   int    `xml:"UseSeparator"`
	EditMaskID     string `xml:"EditMaskID"`
	LookupName     string `xml:"LookupName"`
	MaxSelect      int    `xml:"MaxSelect"`
	Units          string `xml:"Units"`
	Index          int    `xml:"Index"`
	Minimum        string `xml:"Minimum"`
	Maximum        string `xml:"Maxiumum"`
	Default        string `xml:"Default"`
	Required       int    `xml:"Required"`
	SearchHelpID   string `xml:"SearchHelpID"`
	Unique         int    `xml:"Unique"`
}

// XMLMetadataObject ...
type XMLMetadataObject struct {
	XMLMetadataVersion
	Resource string           `xml:"Resource,attr"`
	Object   []MetadataObject `xml:"Object"`
}

// MetadataObject ...
type MetadataObject struct {
	ObjectType   string `xml:"ObjectType"`
	StandardName string `xml:"ObjectType"`
	MimeType     string `xml:"MimeType"`
	Description  string `xml:"Description"`
}

// XMLMetadataLookup ...
type XMLMetadataLookup struct {
	XMLMetadataVersion
	Resource   string               `xml:"Resource,attr"`
	LookupType []MetadataLookupType `xml:"LookupType"`
}

// MetadataLookupType ...
type MetadataLookupType struct {
	LookupName        string    `xml:"LookupName"`
	VisibleName       string    `xml:"VisibleName"`
	LookupTypeVersion string    `xml:"LookupTypeVersion"`
	LookupTypeDate    time.Time `xml:"LookupTypeDate"`

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
	LongValue     string `xml:"LongValue"`
	ShortValue    string `xml:"ShortValue"`
	Value         string `xml:"Value"`
	XDisplayOrder int    `xml:"X-DisplayOrder"`
}
