package metadata

// MClass ...
type MClass struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`
	Class    []Class  `xml:"Class,omitempty"`
}

// Class ...
type Class struct {
	ClassName             RETSName  `xml:"ClassName,omitempty"`
	StandardName          PlainText `xml:"StandardName,omitempty"`
	VisibleName           PlainText `xml:"VisibleName,omitempty"`
	Description           PlainText `xml:"Description,omitempty"`
	TableVersion          Version   `xml:"TableVersion,omitempty"`
	TableDate             DateTime  `xml:"TableDate,omitempty"`
	UpdateVersion         Version   `xml:"UpdateVersion,omitempty"`
	UpdateDate            DateTime  `xml:"UpdateDate,omitempty"`
	ClassTimeStamp        RETSName  `xml:"ClassTimeStamp,omitempty"`
	DeletedFlagField      RETSName  `xml:"DeletedFlagField,omitempty"`
	DeletedFlagValue      AlphaNum  `xml:"DeletedFlagValue,omitempty"`
	HasKeyIndex           Boolean   `xml:"HasKeyIndex,omitempty"`
	OffsetSupport         Boolean   `xml:"OffsetSupport,omitempty"`
	ColumnGroupVersion    Version   `xml:"ColumnGroupVersion,omitempty"`
	ColumnGroupDate       DateTime  `xml:"ColumnGroupDate,omitempty"`
	ColumnGroupSetVersion Version   `xml:"ColumnGroupSetVersion,omitempty"`
	ColumnGroupSetDate    DateTime  `xml:"ColumnGroupSetDate,omitempty"`

	MTable  MTable  `xml:"METADATA-TABLE,omitempty"`
	MUpdate MUpdate `xml:"METADATA-UPDATE,omitempty"`
}
