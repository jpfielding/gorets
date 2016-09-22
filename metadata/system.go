package metadata

import "encoding/xml"

// RETSResponseWrapper is a mapping that can be used directly with an xml.Decoder to extract a full mapping
type RETSResponseWrapper struct {
	XMLName   xml.Name `xml:"RETS,omitempty"`
	ReplyCode int      `xml:",attr,omitempty"`
	ReplyText string   `xml:",attr,omitempty"`
	Metadata  Standard `xml:"METADATA,omitempty"`
}

// Standard ...
type Standard struct {
	MSystem MSystem `xml:"METADATA-SYSTEM,omitempty" json:"METADATA-SYSTEM,omitempty"`
}

// MSystem ...
type MSystem struct {
	Date    DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version Version  `xml:",attr,omitempty" json:",omitempty"`
	System  System   `xml:"SYSTEM,omitempty" json:",omitempty"`
}

// System ...
type System struct {
	ID             string `xml:"SystemID,attr,omitempty" json:"SystemID,omitempty"`
	Description    string `xml:"SystemDescription,attr,omitempty" json:"SystemDescription,omitempty"`
	TimeZoneOffset string `xml:",attr,omitempty" json:",omitempty"`
	MetadataID     string `xml:",attr,omitempty" json:",omitempty"`

	Comments          string   `xml:"COMMENTS,omitempty" json:"COMMENTS,omitempty"`
	ResourceVersion   Version  `xml:",omitempty" json:",omitempty"`
	ResourceDate      DateTime `xml:",omitempty" json:",omitempty"`
	ForeignKeyVersion Version  `xml:",omitempty" json:",omitempty"`
	ForeignKeyDate    DateTime `xml:",omitempty" json:",omitempty"`
	FilterVersion     Version  `xml:",omitempty" json:",omitempty"`
	FilterDate        DateTime `xml:",omitempty" json:",omitempty"`

	MForeignKey MForeignKey `xml:"METADATA-FOREIGN_KEY,omitempty" json:"METADATA-FOREIGN_KEY,omitempty"`
	MResource   MResource   `xml:"METADATA-RESOURCE,omitempty" json:"METADATA-RESOURCE,omitempty"`
	MFilter     MFilter     `xml:"METADATA-FILTER,omitempty" json:"METADATA-FILTER,omitempty"`
}
