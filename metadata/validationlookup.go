package metadata

// MValidationLookup ...
type MValidationLookup struct {
	Date             DateTime           `xml:",attr,omitempty" json:",omitempty"`
	Version          Version            `xml:",attr,omitempty" json:",omitempty"`
	Resource         RETSID             `xml:",attr,omitempty" json:",omitempty"`
	ValidationLookup []ValidationLookup `xml:",omitempty" json:",omitempty"`
}

// ValidationLookup ...
type ValidationLookup struct {
	MetadataEntryID      RETSID   `xml:",omitempty" json:",omitempty"`
	ValidationLookupName RETSName `xml:",omitempty" json:",omitempty"`
	Parent1Field         RETSName `xml:",omitempty" json:",omitempty"`
	Parent2Field         RETSName `xml:",omitempty" json:",omitempty"`
	Date                 DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version              Version  `xml:",attr,omitempty" json:",omitempty"`

	MLookupType *MLookupType `xml:"METADATA-LOOKUP_TYPE,omitempty" json:"METADATA-LOOKUP_TYPE,omitempty"`
}

// MValidationLookupType ...
type MValidationLookupType struct {
	Date             DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version          Version  `xml:",attr,omitempty" json:",omitempty"`
	Resource         RETSID   `xml:",attr,omitempty" json:",omitempty"`
	ValidationLookup RETSName `xml:",attr,omitempty" json:",omitempty"`

	ValidationLookupType []ValidationLookupType `xml:",omitempty" json:",omitempty"`
}

// ValidationLookupType ...
type ValidationLookupType struct {
	MetadataEntryID RETSID   `xml:",omitempty" json:",omitempty"`
	ValidText       RETSName `xml:",omitempty" json:",omitempty"`
	Parent1Value    RETSName `xml:",omitempty" json:",omitempty"`
	Parent2Value    RETSName `xml:",omitempty" json:",omitempty"`
}
