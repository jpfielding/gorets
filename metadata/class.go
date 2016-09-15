package metadata

// MClass ...
type MClass struct {
	Date     DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version  Version  `xml:",attr,omitempty" json:",omitempty"`
	Resource RETSID   `xml:",attr,omitempty" json:",omitempty"`
	Class    []Class  `xml:",omitempty" json:",omitempty"`
}

// Class ...
type Class struct {
	ClassName             RETSName  `xml:",omitempty" json:",omitempty"`
	StandardName          PlainText `xml:",omitempty" json:",omitempty"`
	VisibleName           PlainText `xml:",omitempty" json:",omitempty"`
	Description           PlainText `xml:",omitempty" json:",omitempty"`
	TableVersion          Version   `xml:",omitempty" json:",omitempty"`
	TableDate             DateTime  `xml:",omitempty" json:",omitempty"`
	UpdateVersion         Version   `xml:",omitempty" json:",omitempty"`
	UpdateDate            DateTime  `xml:",omitempty" json:",omitempty"`
	ClassTimeStamp        RETSName  `xml:",omitempty" json:",omitempty"`
	DeletedFlagField      RETSName  `xml:",omitempty" json:",omitempty"`
	DeletedFlagValue      AlphaNum  `xml:",omitempty" json:",omitempty"`
	HasKeyIndex           Boolean   `xml:",omitempty" json:",omitempty"`
	OffsetSupport         Boolean   `xml:",omitempty" json:",omitempty"`
	ColumnGroupVersion    Version   `xml:",omitempty" json:",omitempty"`
	ColumnGroupDate       DateTime  `xml:",omitempty" json:",omitempty"`
	ColumnGroupSetVersion Version   `xml:",omitempty" json:",omitempty"`
	ColumnGroupSetDate    DateTime  `xml:",omitempty" json:",omitempty"`

	MTable  MTable  `xml:"METADATA-TABLE,omitempty" json:"METADATA-TABLE,omitempty"`
	MUpdate MUpdate `xml:"METADATA-UPDATE,omitempty" json:"METADATA-UPDATE,omitempty"`
}
