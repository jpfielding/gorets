package metadata

import "encoding/xml"

// StandardXML is a mapping that can be used directly with an xml.Decoder to extract a full mapping
type StandardXML struct {
	XMLName   xml.Name    `xml:"RETS,omitempty"`
	ReplyCode int         `xml:"ReplyCode,attr,omitempty"`
	ReplyText string      `xml:"ReplyText,attr,omitempty"`
	Metadata  XMLMetadata `xml:"METADATA,omitempty"`
}

// XMLMetadata ...
type XMLMetadata struct {
	System MSystem `xml:"METADATA-SYSTEM,omitempty"`
}

// MSystem ...
type MSystem struct {
	Date    DateTime `xml:"Date,attr,omitempty"`
	Version Version  `xml:"Version,attr,omitempty"`
	System  System   `xml:"SYSTEM,omitempty"`
}

// System ...
type System struct {
	SystemID          string `xml:"SystemID,attr,omitempty"`
	SystemDescription string `xml:"SystemDescription,attr,omitempty"`
	TimeZoneOffset    string `xml:"TimeZoneOffset,attr,omitempty"`
	MetadataID        string `xml:"MetadataID,attr,omitempty"`

	Comments          string   `xml:"COMMENTS,omitempty"`
	ResourceVersion   Version  `xml:"ResourceVersion,omitempty"`
	ResourceDate      DateTime `xml:"ResourceDate,omitempty"`
	ForeignKeyVersion Version  `xml:"ForeignKeyVersion,omitempty"`
	ForeignKeyDate    DateTime `xml:"ForeignKeyDate,omitempty"`
	FilterVersion     Version  `xml:"FilterVerision,omitempty"`
	FilterDate        DateTime `xml:"FilterDate,omitempty"`

	MForeignKey MForeignKey `xml:"METADATA-FOREIGN_KEY,omitempty"`
	MResource   MResource   `xml:"METADATA-RESOURCE,omitempty"`
	MFilter     MFilter     `xml:"METADATA-FILTER,omitempty"`
}
