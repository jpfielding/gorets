package metadata

// MValidationLookup ...
type MValidationLookup struct {
	Date             DateTime           `xml:"Date,attr"`
	Version          Version            `xml:"Version,attr"`
	Resource         string             `xml:"Resource,attr"`
	ValidationLookup []ValidationLookup `xml:"ValidationLookup"`
}

// ValidationLookup ...
type ValidationLookup struct {
	Date         DateTime `xml:"Date,attr"`
	Version      Version  `xml:"Version,attr"`
	Parent1Field string   `xml:"Parent1Field"`
	Parent2Field string   `xml:"Parent2Field"`
	Value        string   `xml:"Value"`

	MLookupType MLookupType `xml:"METADATA-LOOKUP_TYPE"`
}

// MValidationLookupType ...
type MValidationLookupType struct {
	Date             DateTime `xml:"Date,attr"`
	Version          Version  `xml:"Version,attr"`
	Resource         string   `xml:"Resource,attr"`
	ValidationLookup string   `xml:"ValidationLookup,attr"`

	ValidationLookupType []ValidationLookupType `xml:"ValidationLookup"`
}

// ValidationLookupType ...
type ValidationLookupType struct {
	MetadataEntryID string `xml:"MetadataEntryID"`
	ValidText       string `xml:"ValidText"`
	Parent1Value    string `xml:"Parent1Value"`
	Parent2Value    string `xml:"Parent2Value"`
}
