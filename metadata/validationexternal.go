package metadata

// MValidationExternal ...
type MValidationExternal struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource string   `xml:"Resource,attr"`

	ValidationExternal []ValidationExternal `xml:"ValidationExternal"`
}

// ValidationExternal ...
type ValidationExternal struct {
	MetadataEntryID        string   `xml:"MetadataEntryID"`
	ValidationExternalName string   `xml:"ValidationExternalName"`
	SearchResource         string   `xml:"SearchResource"`
	SearchClass            string   `xml:"SearchClass"`
	Date                   DateTime `xml:"Date,attr"`
	Version                Version  `xml:"Version,attr"`
}
