package metadata

// MColumnGroupSet ...
type MColumnGroupSet struct {
	Date           DateTime         `xml:",attr,omitempty"`
	Version        Version          `xml:",attr,omitempty"`
	Resource       RETSID           `xml:",attr,omitempty"`
	Class          RETSID           `xml:",attr,omitempty"`
	ColumnGroupSet []ColumnGroupSet `xml:",omitempty"`
}

// ColumnGroupSet ...
type ColumnGroupSet struct {
	MetadataEntryID      RETSID   `xml:",omitempty" json:",omitempty"`
	ColumnGroupSetName   RETSID   `xml:",omitempty" json:",omitempty"`
	ColumnGroupSetParent RETSID   `xml:",omitempty" json:",omitempty"`
	Sequence             Numeric  `xml:"omitempty" json:",omitempty"`
	LongName             RETSName `xml:"omitempty" json:",omitempty"`
	ShortName            RETSName `xml:"omitempty" json:",omitempty"`
	Description          Text     `xml:"omitempty" json:",omitempty"`
	ColumnGroupName      RETSID   `xml:"omitempty" json:",omitempty"`
	PresentationStyle    Text     `xml:"omitempty" json:",omitempty"`
	URL                  Text     `xml:"omitempty" json:",omitempty"`
	ForeignKeyID         RETSID   `xml:"omitempty" json:",omitempty"`
}
