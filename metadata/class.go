package metadata

// MClass ...
type MClass struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource string   `xml:"Resource,attr"`
	Class    []Class  `xml:"Class"`
}

// Class ...
type Class struct {
	ClassName             string   `xml:"ClassName"`
	StandardName          string   `xml:"StandardName"`
	VisibleName           string   `xml:"VisibleName"`
	Description           string   `xml:"Description"`
	TableVersion          Version  `xml:"TableVersion"`
	TableDate             DateTime `xml:"TableDate"`
	UpdateVersion         Version  `xml:"UpdateVersion"`
	UpdateDate            DateTime `xml:"UpdateDate"`
	ClassTimeStamp        DateTime `xml:"ClassTimeStamp"`
	DeletedFlagField      string   `xml:"DeletedFlagField"`
	DeletedFlagValue      string   `xml:"DeletedFlagValue"`
	HasKeyIndex           Boolean  `xml:"HasKeyIndex"`
	OffsetSupport         Boolean  `xml:"OffsetSupport"`
	ColumnGroupVersion    Version  `xml:"ColumnGroupVersion"`
	ColumnGroupDate       DateTime `xml:"ColumnGroupDate"`
	ColumnGroupSetVersion Version  `xml:"ColumnGroupSetVersion"`
	ColumnGroupSetDate    DateTime `xml:"ColumnGroupSetDate"`

	MUpdate MUpdate `xml:"METADATA-UPDATE"`
}
