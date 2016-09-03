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
	TableVersion          string    `xml:"TableVersion"`
	TableDate             time.Time `xml:"TableDate"`
	UpdateVersion         string    `xml:"UpdateVersion"`
	UpdateDate            time.Time `xml:"UpdateDate"`
	ClassTimeStamp        time.Time `xml:"ClassTimeStamp"`
	DeletedFlagField      string    `xml:"DeletedFlagField"`
	DeletedFlagValue      string    `xml:"DeletedFlagValue"`
	HasKeyIndex           *bool     `xml:"HasKeyIndex"`
	OffsetSupport         *bool     `xml:"OffsetSupport"`
	ColumnGroupVersion    string    `xml:"ColumnGroupVersion"`
	ColumnGroupDate       time.Time `xml:"ColumnGroupDate"`
	ColumnGroupSetVersion string    `xml:"ColumnGroupSetVersion"`
	ColumnGroupSetDate    time.Time `xml:"ColumnGroupSetDate"`

	// maybe put a placeholder in for unmapped fields ??
	// XDisplayOrder         int       `xml:"X-DisplayOrder"`
}
