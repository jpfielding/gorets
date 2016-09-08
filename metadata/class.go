package metadata

// MClass ...
type MClass struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`
	Class    []Class  `xml:"Class"`
}

// Class ...
type Class struct {
	ClassName             RETSName  `xml:"ClassName"`
	StandardName          PlainText `xml:"StandardName"`
	VisibleName           PlainText `xml:"VisibleName"`
	Description           PlainText `xml:"Description"`
	TableVersion          Version   `xml:"TableVersion"`
	TableDate             DateTime  `xml:"TableDate"`
	UpdateVersion         Version   `xml:"UpdateVersion"`
	UpdateDate            DateTime  `xml:"UpdateDate"`
	ClassTimeStamp        RETSName  `xml:"ClassTimeStamp"`
	DeletedFlagField      RETSName  `xml:"DeletedFlagField"`
	DeletedFlagValue      AlphaNum  `xml:"DeletedFlagValue"`
	HasKeyIndex           Boolean   `xml:"HasKeyIndex"`
	OffsetSupport         Boolean   `xml:"OffsetSupport"`
	ColumnGroupVersion    Version   `xml:"ColumnGroupVersion"`
	ColumnGroupDate       DateTime  `xml:"ColumnGroupDate"`
	ColumnGroupSetVersion Version   `xml:"ColumnGroupSetVersion"`
	ColumnGroupSetDate    DateTime  `xml:"ColumnGroupSetDate"`

	MTable  MTable  `xml:"METADATA-TABLE"`
	MUpdate MUpdate `xml:"METADATA-UPDATE"`
}
