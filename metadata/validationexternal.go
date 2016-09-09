package metadata

// MValidationExternal ...
type MValidationExternal struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`

	ValidationExternal []ValidationExternal `xml:"ValidationExternal,omitempty"`
}

// ValidationExternal ...
type ValidationExternal struct {
	MetadataEntryID        RETSID   `xml:"MetadataEntryID,omitempty"`
	ValidationExternalName RETSName `xml:"ValidationExternalName,omitempty"`
	SearchResource         RETSName `xml:"SearchResource,omitempty"`
	SearchClass            RETSName `xml:"SearchClass,omitempty"`
	Date                   DateTime `xml:"Date,attr,omitempty"`
	Version                Version  `xml:"Version,attr,omitempty"`
}

// MValidationExternalType ...
type MValidationExternalType struct {
	Date                   DateTime                 `xml:"Date,attr,omitempty"`
	Version                Version                  `xml:"Version,attr,omitempty"`
	Resource               RETSID                   `xml:"Resource,attr,omitempty"`
	ValidationExternalName RETSID                   `xml:"ValidationExternalName,attr,omitempty"`
	ValidationExternalType []ValidationExternalType `xml:"ValidationExternalType,omitempty"`
}

// ValidationExternalType ...
type ValidationExternalType struct {
	MetadataEntryID RETSID     `xml:"MetadataEntryID,omitempty"`
	SearchField     PlainText  `xml:"SearchField,omitempty"`
	DisplayField    PlainText  `xml:"DisplayField,omitempty"`
	ResultFields    StringList `xml:"ResultFields,omitempty"` // TODO plaintext list
}
