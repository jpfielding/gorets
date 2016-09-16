package metadata

// MFilter ...
type MFilter struct {
	Date    DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version Version  `xml:",attr,omitempty" json:",omitempty"`
	Filter  []Filter `xml:",omitempty" json:",omitempty"`
}

// Filter ...
type Filter struct {
	FilterID          RETSID   `xml:",omitempty" json:",omitempty"`
	ParentResource    RETSID   `xml:",omitempty" json:",omitempty"`
	ParentLookupName  RETSName `xml:",omitempty" json:",omitempty"`
	ChildResource     RETSID   `xml:",omitempty" json:",omitempty"`
	ChildLookupName   RETSName `xml:",omitempty" json:",omitempty"`
	NotShownByDefault Boolean  `xml:",omitempty" json:",omitempty"`

	MFilterType MFilterType `xml:"METADATA-FILTER_TYPE,omitempty" json:"METADATA-FILTER_TYPE,omitempty"`
}

// MFilterType ...
type MFilterType struct {
	Date    DateTime `xml:",attr,omitempty"`
	Version Version  `xml:",attr,omitempty"`

	FilterType []FilterType `xml:",omitempty" json:",omitempty"`
}

// FilterType ...
type FilterType struct {
	FilterTypeID RETSID    `xml:",omitempty" json:",omitempty"`
	ParentValue  PlainText `xml:",omitempty" json:",omitempty"`
	ChildValue   PlainText `xml:",omitempty" json:",omitempty"`
}
