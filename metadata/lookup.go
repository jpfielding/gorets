package metadata

// MLookup ...
type MLookup struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`
	Lookup   []Lookup `xml:"Lookup"`
}

// Lookup ...
type Lookup struct {
	MetadataEntryID   RETSID    `xml:"MetadataEntryID"`
	LookupName        RETSName  `xml:"LookupName"`
	VisibleName       PlainText `xml:"VisibleName"`
	LookupTypeVersion Version   `xml:"LookupTypeVersion"`
	LookupTypeDate    DateTime  `xml:"LookupTypeDate"`
	FilterID          RETSID    `xml:"FilterID"`
	NotShowByDefault  Boolean   `xml:"NotShowByDefault"`

	MLookupType MLookupType `xml:"METADATA-LOOKUP_TYPE"`
}

// MLookupType ...
type MLookupType struct {
	Date       DateTime     `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   RETSID       `xml:"Resource,attr"`
	Lookup     RETSID       `xml:"Lookup,attr"`
	LookupType []LookupType `xml:"Lookup"`
}

// LookupType ...
type LookupType struct {
	MetadataEntryID RETSID    `xml:"MetadataEntryID"`
	LongValue       PlainText `xml:"LongValue"`
	ShortValue      PlainText `xml:"ShortValue"`
	Value           PlainText `xml:"Value"`
}
