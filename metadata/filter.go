package metadata

// MFilter ...
type MFilter struct {
	Date    DateTime `xml:"Date,attr"`
	Version Version  `xml:"Version,attr"`
	Filter  []Filter `xml:"Filter"`
}

// Filter ...
type Filter struct {
	FilterID          RETSID   `xml:"FilterID"`
	ParentResource    RETSID   `xml:"ParentResource"`
	ParentLookupName  RETSName `xml:"ParentLookupName"`
	ChildResource     RETSID   `xml:"ChildResource"`
	ChildLookupName   RETSName `xml:"ChildLookupName"`
	NotShownByDefault Boolean  `xml:"NotShownByDefault"`
}

// MFilterType ...
type MFilterType struct {
	Date    DateTime `xml:"Date,attr"`
	Version Version  `xml:"Version,attr"`

	Filter []FilterType `xml:"FilterType"`
}

// FilterType ...
type FilterType struct {
	FilterTypeID RETSID    `xml:"FilterTypeID"`
	ParentValue  PlainText `xml:"ParentValue"`
	ChildValue   PlainText `xml:"ChildValue"`
}
