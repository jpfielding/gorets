package metadata

// MFilter ...
type MFilter struct {
	Date    DateTime `xml:"Date,attr"`
	Version Version  `xml:"Version,attr"`
	Filter  []Filter `xml:"Filter"`
}

// Filter ...
type Filter struct {
	FilterID          string `xml:"FilterID"`
	ParentResource    string `xml:"ParentResource"`
	ParentLookupName  string `xml:"ParentLookupName"`
	ChildResource     string `xml:"ChildResource"`
	ChildLookupName   string `xml:"ChildLookupName"`
	NotShownByDefault *bool  `xml:"NotShownByDefault"`
}

// MFilterType ...
type MFilterType struct {
	Date    DateTime `xml:"Date,attr"`
	Version Version  `xml:"Version,attr"`

	Filter []FilterType `"xml:FilterType"`
}

// FilterType ...
type FilterType struct {
	FilterTypeID string `xml:FilterTypeID`
	ParentValue  string `xml:ParentValue`
	ChildValue   string `xml:ChildValue`
}
