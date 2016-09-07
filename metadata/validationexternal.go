package metadata

// MValidationExternal ...
type MValidationExternal struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`

	ValidationExternal []ValidationExternal `xml:"ValidationExternal"`
}

// ValidationExternal ...
type ValidationExternal struct {
	MetadataEntryID        RETSID   `xml:"MetadataEntryID"`
	ValidationExternalName RETSName `xml:"ValidationExternalName"`
	SearchResource         RETSName `xml:"SearchResource"`
	SearchClass            RETSName `xml:"SearchClass"`
	Date                   DateTime `xml:"Date,attr"`
	Version                Version  `xml:"Version,attr"`
}

// MValidationExternalType ...
type MValidationExternalType struct {
	Date                   DateTime                 `xml:"Date,attr"`
	Version                Version                  `xml:"Version,attr"`
	Resource               RETSID                   `xml:"Resource,attr"`
	ValidationExternalName RETSID                   `xml:"ValidationExternalName,attr"`
	ValidationExternalType []ValidationExternalType `xml:"ValidationExternalType"`
}

// ValidationExternalType ...
type ValidationExternalType struct {
	MetadataEntryID RETSID     `xml:"MetadataEntryID"`
	SearchField     PlainText  `xml:"SearchField"`
	DisplayField    PlainText  `xml:"DisplayField"`
	ResultFields    StringList `xml:"ResultFields"` // TODO plaintext list
}
