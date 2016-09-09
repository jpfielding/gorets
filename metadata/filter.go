package metadata

// MFilter ...
type MFilter struct {
	Date    DateTime `xml:"Date,attr,omitempty"`
	Version Version  `xml:"Version,attr,omitempty"`
	Filter  []Filter `xml:"Filter,omitempty"`
}

// Filter ...
type Filter struct {
	FilterID          RETSID   `xml:"FilterID,omitempty"`
	ParentResource    RETSID   `xml:"ParentResource,omitempty"`
	ParentLookupName  RETSName `xml:"ParentLookupName,omitempty"`
	ChildResource     RETSID   `xml:"ChildResource,omitempty"`
	ChildLookupName   RETSName `xml:"ChildLookupName,omitempty"`
	NotShownByDefault Boolean  `xml:"NotShownByDefault,omitempty"`
}

// MFilterType ...
type MFilterType struct {
	Date    DateTime `xml:"Date,attr,omitempty"`
	Version Version  `xml:"Version,attr,omitempty"`

	Filter []FilterType `xml:"FilterType,omitempty"`
}

// FilterType ...
type FilterType struct {
	FilterTypeID RETSID    `xml:"FilterTypeID,omitempty"`
	ParentValue  PlainText `xml:"ParentValue,omitempty"`
	ChildValue   PlainText `xml:"ChildValue,omitempty"`
}
