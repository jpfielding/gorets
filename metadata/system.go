package metadata

import "time"

// MSystem ...
type MSystem struct {
	Date    time.Time `xml:"Date,attr"`
	Version Version   `xml:"Version,attr"`
	System  System    `xml:"SYSTEM"`
}

// System ...
type System struct {
	ID             string `xml:"SystemID,attr"`
	Description    string `xml:"SystemDescription,attr"`
	TimeZoneOffset string `xml:"TimeZoneOffset,attr"`
	MetadataID     string `xml:"MetadataID,attr"`

	Comments          string    `xml:"COMMENTS"`
	ResourceVersion   Version   `xml:"ResourceVersion"`
	ResourceDate      time.Time `xml:ResourceDate"`
	ForeignKeyVersion Version   `xml:"ForeignKeyVersion"`
	ForeignKeyDate    time.Time `xml:"ForeignKeyDate"`
	FilterVersion     Version   `xml:"FilterVerision"`
	FilterDate        time.Time `xml:"FilterDate"`

	MForeignKey MForeignKey `xml:"METADATA-FOREIGN_KEY"`
	MResource   MResource   `xml:"METADATA-RESOURCE"`
	MFilter     MFilter     `xml:"METADATA-FILTER"`
}
