package metadata

// MLookup ...
type MLookup struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`
	Lookup   []Lookup `xml:"Lookup,omitempty"`
}

// Lookup ...
type Lookup struct {
	MetadataEntryID   RETSID    `xml:"MetadataEntryID,omitempty"`
	LookupName        RETSName  `xml:"LookupName,omitempty"`
	VisibleName       PlainText `xml:"VisibleName,omitempty"`
	LookupTypeVersion Version   `xml:"LookupTypeVersion,omitempty"`
	LookupTypeDate    DateTime  `xml:"LookupTypeDate,omitempty"`
	FilterID          RETSID    `xml:"FilterID,omitempty"`
	NotShowByDefault  Boolean   `xml:"NotShowByDefault,omitempty"`

	MLookupType MLookupType `xml:"METADATA-LOOKUP_TYPE,omitempty"`
}

// MLookupType ...
type MLookupType struct {
	Date       DateTime     `xml:"Date,attr,omitempty"`
	Version    Version      `xml:"Version,attr,omitempty"`
	Resource   RETSID       `xml:"Resource,attr,omitempty"`
	Lookup     RETSID       `xml:"Lookup,attr,omitempty"`
	LookupType []LookupType `xml:"Lookup,omitempty"`
}

// LookupType ...
type LookupType struct {
	MetadataEntryID RETSID    `xml:"MetadataEntryID,omitempty"`
	LongValue       PlainText `xml:"LongValue,omitempty"`
	ShortValue      PlainText `xml:"ShortValue,omitempty"`
	Value           PlainText `xml:"Value,omitempty"`
}
