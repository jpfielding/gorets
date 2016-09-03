package metadata

import "time"

// MLookup ...
type MLookup struct {
	Date     time.Time `xml:"Date,attr"`
	Version  Version   `xml:"Version,attr"`
	Resource string    `xml:"Resource,attr"`
	Lookup   []Lookup  `xml:"Lookup"`
}

// Lookup ...
type Lookup struct {
	MetadataEntryID   string    `xml:"MetadataEntryID"`
	LookupName        string    `xml:"LookupName"`
	VisibleName       string    `xml:"VisibleName"`
	LookupTypeVersion Version   `xml:"LookupTypeVersion"`
	LookupTypeDate    time.Time `xml:"LookupTypeDate"`
	FilterID          string    `xml:"FilterID"`
	NotShowByDefault  *bool     `xml:"NotShowByDefault"`

	MLookupType MLookupType `xml:"METADATA-LOOKUP_TYPE"`
}

// MLookupType ...
type MLookupType struct {
	Date       time.Time    `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   string       `xml:"Resource,attr"`
	Lookup     string       `xml:"Lookup,attr"`
	LookupType []LookupType `xml:"Lookup"`
}

// LookupType ...
type LookupType struct {
	MetadataEntryID string `xml:"MetadataEntryID"`
	LongValue       string `xml:"LongValue"`
	ShortValue      string `xml:"ShortValue"`
	Value           string `xml:"Value"`
	// XDisplayOrder int    `xml:"X-DisplayOrder"`
}
