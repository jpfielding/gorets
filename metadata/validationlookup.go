package metadata

// MValidationLookup ...
type MValidationLookup struct {
	Date             DateTime           `xml:"Date,attr,omitempty"`
	Version          Version            `xml:"Version,attr,omitempty"`
	Resource         RETSID             `xml:"Resource,attr,omitempty"`
	ValidationLookup []ValidationLookup `xml:"ValidationLookup,omitempty"`
}

// ValidationLookup ...
type ValidationLookup struct {
	MetadataEntryID      RETSID   `xml:"MetadataEntryID,omitempty"`
	ValidationLookupName RETSName `xml:"ValidationLookupName,omitempty"`
	Parent1Field         RETSName `xml:"Parent1Field,omitempty"`
	Parent2Field         RETSName `xml:"Parent2Field,omitempty"`
	Date                 DateTime `xml:"Date,attr,omitempty"`
	Version              Version  `xml:"Version,attr,omitempty"`

	MLookupType MLookupType `xml:"METADATA-LOOKUP_TYPE,omitempty"`
}

// MValidationLookupType ...
type MValidationLookupType struct {
	Date             DateTime `xml:"Date,attr,omitempty"`
	Version          Version  `xml:"Version,attr,omitempty"`
	Resource         RETSID   `xml:"Resource,attr,omitempty"`
	ValidationLookup RETSName `xml:"ValidationLookup,attr,omitempty"`

	ValidationLookupType []ValidationLookupType `xml:"ValidationLookup,omitempty"`
}

// ValidationLookupType ...
type ValidationLookupType struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID,omitempty"`
	ValidText       RETSName `xml:"ValidText,omitempty"`
	Parent1Value    RETSName `xml:"Parent1Value,omitempty"`
	Parent2Value    RETSName `xml:"Parent2Value,omitempty"`
}
