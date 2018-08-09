package metadata

// MTable ...
type MTable struct {
	Date     DateTime `xml:",attr,omitempty"`
	Version  Version  `xml:",attr,omitempty"`
	Resource RETSID   `xml:",attr,omitempty"`
	Class    RETSID   `xml:",attr,omitempty"`
	Field    []Field  `xml:",omitempty"`
}

// Field ...
type Field struct {
	MetadataEntryID    RETSID    `xml:",omitempty" json:",omitempty"`
	SystemName         RETSName  `xml:",omitempty" json:",omitempty"`
	StandardName       RETSName  `xml:",omitempty" json:",omitempty"`
	LongName           Text      `xml:",omitempty" json:",omitempty"`
	DBName             AlphaNum  `xml:",omitempty" json:",omitempty"`
	ShortName          Text      `xml:",omitempty" json:",omitempty"`
	MaximumLength      Numeric   `xml:",omitempty" json:",omitempty"` // TODO limit to positive
	DataType           string    `xml:",omitempty" json:",omitempty"` // TODO limit to options
	Precision          Numeric   `xml:",omitempty" json:",omitempty"` // TODO limit to positive
	Searchable         Boolean   `xml:",omitempty" json:",omitempty"`
	Interpretation     string    `xml:",omitempty" json:",omitempty"` // TODO limit to options
	Alignment          string    `xml:",omitempty" json:",omitempty"` // TODO limit to options
	UseSeparator       Boolean   `xml:",omitempty" json:",omitempty"` // TODO limit to options
	EditMaskID         RETSNames `xml:",omitempty" json:",omitempty"`
	LookupName         RETSName  `xml:",omitempty" json:",omitempty"`
	MaxSelect          Numeric   `xml:",omitempty" json:",omitempty"`
	Units              string    `xml:",omitempty" json:",omitempty"`
	Index              Boolean   `xml:",omitempty" json:",omitempty"`
	Minimum            Numeric   `xml:",omitempty" json:",omitempty"`
	Maximum            Numeric   `xml:",omitempty" json:",omitempty"`
	Default            Numeric   `xml:",omitempty" json:",omitempty"`
	Required           Numeric   `xml:",omitempty" json:",omitempty"`
	SearchHelpID       RETSName  `xml:",omitempty" json:",omitempty"`
	Unique             Boolean   `xml:",omitempty" json:",omitempty"`
	ModTimeStamp       DateTime  `xml:",omitempty" json:",omitempty"`
	ForeignKeyName     RETSID    `xml:",omitempty" json:",omitempty"`
	ForeignField       RETSName  `xml:",omitempty" json:",omitempty"`
	KeyQuery           Boolean   `xml:",omitempty" json:",omitempty"` // deprecated
	KeySelect          Boolean   `xml:",omitempty" json:",omitempty"` // deprecated
	InKeyIndex         Boolean   `xml:",omitempty" json:",omitempty"`
	FilterParentField  RETSName  `xml:",omitempty" json:",omitempty"`
	DefaultSearchOrder Numeric   `xml:",omitempty" json:",omitempty"`
	Case               string    `xml:",omitempty" json:",omitempty"` // TODO limit to known options
}
