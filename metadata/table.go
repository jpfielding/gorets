package metadata

// MTable ...
type MTable struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource string   `xml:"Resource,attr"`
	Class    string   `xml:"Class,attr"`
	Table    []Table  `xml:"Table"`
}

// Table ...
type Table struct {
	MetadataEntryID    string     `xml:"MetadataEntryID"`
	SystemName         string     `xml:"SystemName"`
	StandardName       string     `xml:"StandardName"`
	LongName           string     `xml:"LongName"`
	DBName             string     `xml:"DBName"`
	ShortName          string     `xml:"ShortName"`
	MaximumLength      Number     `xml:"MaximumLength"`
	DataType           string     `xml:"DataType"`
	Precision          Number     `xml:"Precision"`
	Searchable         Boolean    `xml:"Searchable"`
	Interpretation     string     `xml:"Interpretation"`
	Alignment          string     `xml:"Alignment"`
	UseSeparator       Boolean    `xml:"UseSeparator"`
	EditMaskID         StringList `xml:"EditMaskID"`
	LookupName         string     `xml:"LookupName"`
	MaxSelect          Number     `xml:"MaxSelect"`
	Units              string     `xml:"Units"`
	Index              Boolean    `xml:"Index"`
	Minimum            string     `xml:"Minimum"`
	Maximum            string     `xml:"Maxiumum"`
	Default            Number     `xml:"Default"`
	Required           Number     `xml:"Required"`
	SearchHelpID       string     `xml:"SearchHelpID"`
	Unique             Boolean    `xml:"Unique"`
	ModTimeStamp       DateTime   `xml:"ModTimeStamp"`
	ForeignKeyName     string     `xml:"ForeignKeyName"`
	ForeignField       string     `xml:"ForeignField"`
	KeyQuery           string     `xml:"KeyQuery"`  // deprecated
	KeySelect          string     `xml:"KeySelect"` // deprecated
	InKeyIndex         Boolean    `xml:"InKeyIndex"`
	FilterParentField  string     `xml:"FilterParentField"`
	DefaultSearchOrder Number     `xml:"DefaultSearchOrder"`
	Case               string     `xml:"Case"`
}
