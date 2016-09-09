package metadata

// MTable ...
type MTable struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`
	Class    RETSID   `xml:"Class,attr,omitempty"`
	Field    []Field  `xml:"Field,omitempty"`
}

// Field ...
type Field struct {
	MetadataEntryID    RETSID    `xml:"MetadataEntryID,omitempty"`
	SystemName         RETSName  `xml:"SystemName,omitempty"`
	StandardName       RETSName  `xml:"StandardName,omitempty"`
	LongName           Text      `xml:"LongName,omitempty"`
	DBName             AlphaNum  `xml:"DBName,omitempty"`
	ShortName          Text      `xml:"ShortName,omitempty"`
	MaximumLength      Numeric   `xml:"MaximumLength,omitempty"` // TODO limit to postive
	DataType           string    `xml:"DataType,omitempty"`      // TODO limit to options
	Precision          Numeric   `xml:"Precision,omitempty"`     // TODO limit to positive
	Searchable         Boolean   `xml:"Searchable,omitempty"`
	Interpretation     string    `xml:"Interpretation,omitempty"` // TODO limit to options
	Alignment          string    `xml:"Alignment,omitempty"`      // TODO limit to options
	UseSeparator       Boolean   `xml:"UseSeparator,omitempty"`   // TODO limit to options
	EditMaskID         RETSNames `xml:"EditMaskID,omitempty"`
	LookupName         RETSName  `xml:"LookupName,omitempty"`
	MaxSelect          Numeric   `xml:"MaxSelect,omitempty"`
	Units              string    `xml:"Units,omitempty"`
	Index              Boolean   `xml:"Index,omitempty"`
	Minimum            Numeric   `xml:"Minimum,omitempty"`
	Maximum            Numeric   `xml:"Maxiumum,omitempty"`
	Default            Numeric   `xml:"Default,omitempty"`
	Required           Numeric   `xml:"Required,omitempty"`
	SearchHelpID       RETSName  `xml:"SearchHelpID,omitempty"`
	Unique             Boolean   `xml:"Unique,omitempty"`
	ModTimeStamp       DateTime  `xml:"ModTimeStamp,omitempty"`
	ForeignKeyName     RETSID    `xml:"ForeignKeyName,omitempty"`
	ForeignField       RETSName  `xml:"ForeignField,omitempty"`
	KeyQuery           Boolean   `xml:"KeyQuery,omitempty"`  // deprecated
	KeySelect          Boolean   `xml:"KeySelect,omitempty"` // deprecated
	InKeyIndex         Boolean   `xml:"InKeyIndex,omitempty"`
	FilterParentField  RETSName  `xml:"FilterParentField,omitempty"`
	DefaultSearchOrder Numeric   `xml:"DefaultSearchOrder,omitempty"`
	Case               string    `xml:"Case,omitempty"` // TODO limit to known options
}
