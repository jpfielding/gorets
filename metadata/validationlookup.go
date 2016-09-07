package metadata

// MValidationLookup ...
type MValidationLookup struct {
	Date             DateTime           `xml:"Date,attr"`
	Version          Version            `xml:"Version,attr"`
	Resource         RETSID             `xml:"Resource,attr"`
	ValidationLookup []ValidationLookup `xml:"ValidationLookup"`
}

// ValidationLookup ...
type ValidationLookup struct {
	MetadataEntryID      RETSID   `xml:"MetadataEntryID"`
	ValidationLookupName RETSName `xml:"ValidationLookupName"`
	Parent1Field         RETSName `xml:"Parent1Field"`
	Parent2Field         RETSName `xml:"Parent2Field"`
	Date                 DateTime `xml:"Date,attr"`
	Version              Version  `xml:"Version,attr"`

	MLookupType MLookupType `xml:"METADATA-LOOKUP_TYPE"`
}

// MValidationLookupType ...
type MValidationLookupType struct {
	Date             DateTime `xml:"Date,attr"`
	Version          Version  `xml:"Version,attr"`
	Resource         RETSID   `xml:"Resource,attr"`
	ValidationLookup RETSName `xml:"ValidationLookup,attr"`

	ValidationLookupType []ValidationLookupType `xml:"ValidationLookup"`
}

// ValidationLookupType ...
type ValidationLookupType struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID"`
	ValidText       RETSName `xml:"ValidText"`
	Parent1Value    RETSName `xml:"Parent1Value"`
	Parent2Value    RETSName `xml:"Parent2Value"`
}
