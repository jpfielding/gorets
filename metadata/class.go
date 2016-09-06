package metadata

import "time"

// MClass ...
type MClass struct {
	Date     time.Time `xml:"Date,attr"`
	Version  Version   `xml:"Version,attr"`
	Resource string    `xml:"Resource,attr"`
	Class    []Class   `xml:"Class"`
}

// Class ...
type Class struct {
	ClassName             string    `xml:"ClassName"`
	StandardName          string    `xml:"StandardName"`
	VisibleName           string    `xml:"VisibleName"`
	Description           string    `xml:"Description"`
	TableVersion          Version   `xml:"TableVersion"`
	TableDate             time.Time `xml:"TableDate"`
	UpdateVersion         Version   `xml:"UpdateVersion"`
	UpdateDate            time.Time `xml:"UpdateDate"`
	ClassTimeStamp        time.Time `xml:"ClassTimeStamp"`
	DeletedFlagField      string    `xml:"DeletedFlagField"`
	DeletedFlagValue      string    `xml:"DeletedFlagValue"`
	HasKeyIndex           *bool     `xml:"HasKeyIndex"`
	OffsetSupport         *bool     `xml:"OffsetSupport"`
	ColumnGroupVersion    Version   `xml:"ColumnGroupVersion"`
	ColumnGroupDate       time.Time `xml:"ColumnGroupDate"`
	ColumnGroupSetVersion Version   `xml:"ColumnGroupSetVersion"`
	ColumnGroupSetDate    time.Time `xml:"ColumnGroupSetDate"`

	MUpdate MUpdate `xml:"METADATA-UPDATE"`
}
