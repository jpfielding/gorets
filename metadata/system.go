package metadata

import "encoding/xml"

// StandardXML ...
type StandardXML struct {
	XMLName   xml.Name    `xml:"RETS"`
	ReplyCode int         `xml:"ReplyCode,attr"`
	ReplyText string      `xml:"ReplyText,attr"`
	Metadata  XMLMetadata `xml:"METADATA"`
}

// XMLMetadata ...
type XMLMetadata struct {
	System MSystem `xml:"METADATA-SYSTEM"`
}

// MSystem ...
type MSystem struct {
	Date    DateTime `xml:"Date,attr"`
	Version Version  `xml:"Version,attr"`
	System  System   `xml:"SYSTEM"`
}

// System ...
type System struct {
	ID             string `xml:"SystemID,attr"`
	Description    string `xml:"SystemDescription,attr"`
	TimeZoneOffset string `xml:"TimeZoneOffset,attr"`
	MetadataID     string `xml:"MetadataID,attr"`

	Comments          string   `xml:"COMMENTS"`
	ResourceVersion   Version  `xml:"ResourceVersion"`
	ResourceDate      DateTime `xml:ResourceDate"`
	ForeignKeyVersion Version  `xml:"ForeignKeyVersion"`
	ForeignKeyDate    DateTime `xml:"ForeignKeyDate"`
	FilterVersion     Version  `xml:"FilterVerision"`
	FilterDate        DateTime `xml:"FilterDate"`

	MForeignKey MForeignKey `xml:"METADATA-FOREIGN_KEY"`
	MResource   MResource   `xml:"METADATA-RESOURCE"`
	MFilter     MFilter     `xml:"METADATA-FILTER"`
}
