package metadata

// MLookup ...
type MLookup struct {
	Date     DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version  Version  `xml:",attr,omitempty" json:",omitempty"`
	Resource RETSID   `xml:",attr,omitempty" json:",omitempty"`
	Lookup   []Lookup `xml:",omitempty" json:",omitempty"`
}

// Lookup ...
type Lookup struct {
	MetadataEntryID   RETSID    `xml:",omitempty" json:",omitempty"`
	LookupName        RETSName  `xml:",omitempty" json:",omitempty"`
	VisibleName       PlainText `xml:",omitempty" json:",omitempty"`
	LookupTypeVersion Version   `xml:",omitempty" json:",omitempty"`
	LookupTypeDate    DateTime  `xml:",omitempty" json:",omitempty"`
	FilterID          RETSID    `xml:",omitempty" json:",omitempty"`
	NotShowByDefault  Boolean   `xml:",omitempty" json:",omitempty"`

	MLookupType *MLookupType `xml:"METADATA-LOOKUP_TYPE,omitempty" json:"METADATA-LOOKUP_TYPE,omitempty"`
}

// MLookupType ...
type MLookupType struct {
	Date       DateTime     `xml:",attr,omitempty" json:",omitempty"`
	Version    Version      `xml:",attr,omitempty" json:",omitempty"`
	Resource   RETSID       `xml:",attr,omitempty" json:",omitempty"`
	Lookup     RETSID       `xml:",attr,omitempty" json:",omitempty"`
	LookupType []LookupType `xml:",omitempty" json:",omitempty"`
}

// LookupType ...
type LookupType struct {
	MetadataEntryID RETSID    `xml:",omitempty" json:",omitempty"`
	LongValue       PlainText `xml:",omitempty" json:",omitempty"`
	ShortValue      PlainText `xml:",omitempty" json:",omitempty"`
	Value           PlainText `xml:",omitempty" json:",omitempty"`
}
