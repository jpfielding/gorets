package metadata

// MTable ...
type MTable struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`
	Class    RETSID   `xml:"Class,attr"`
	Field    []Field  `xml:"Field"`
}

// Field ...
type Field struct {
	MetadataEntryID    RETSID    `xml:"MetadataEntryID"`
	SystemName         RETSName  `xml:"SystemName"`
	StandardName       RETSName  `xml:"StandardName"`
	LongName           Text      `xml:"LongName"`
	DBName             AlphaNum  `xml:"DBName"`
	ShortName          Text      `xml:"ShortName"`
	MaximumLength      Numeric   `xml:"MaximumLength"` // TODO limit to postive
	DataType           string    `xml:"DataType"`      // TODO limit to options
	Precision          Numeric   `xml:"Precision"`     // TODO limit to positive
	Searchable         Boolean   `xml:"Searchable"`
	Interpretation     string    `xml:"Interpretation"` // TODO limit to options
	Alignment          string    `xml:"Alignment"`      // TODO limit to options
	UseSeparator       Boolean   `xml:"UseSeparator"`   // TODO limit to options
	EditMaskID         RETSNames `xml:"EditMaskID"`
	LookupName         RETSName  `xml:"LookupName"`
	MaxSelect          Numeric   `xml:"MaxSelect"`
	Units              string    `xml:"Units"`
	Index              Boolean   `xml:"Index"`
	Minimum            Numeric   `xml:"Minimum"`
	Maximum            Numeric   `xml:"Maxiumum"`
	Default            Numeric   `xml:"Default"`
	Required           Numeric   `xml:"Required"`
	SearchHelpID       RETSName  `xml:"SearchHelpID"`
	Unique             Boolean   `xml:"Unique"`
	ModTimeStamp       DateTime  `xml:"ModTimeStamp"`
	ForeignKeyName     RETSID    `xml:"ForeignKeyName"`
	ForeignField       RETSName  `xml:"ForeignField"`
	KeyQuery           Boolean   `xml:"KeyQuery"`  // deprecated
	KeySelect          Boolean   `xml:"KeySelect"` // deprecated
	InKeyIndex         Boolean   `xml:"InKeyIndex"`
	FilterParentField  RETSName  `xml:"FilterParentField"`
	DefaultSearchOrder Numeric   `xml:"DefaultSearchOrder"`
	Case               string    `xml:"Case"` // TODO limit to known options
}
