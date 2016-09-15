package metadata

// MValidationExternal ...
type MValidationExternal struct {
	Date     DateTime `xml:",attr,omitempty json:",omitempty"`
	Version  Version  `xml:",attr,omitempty json:",omitempty"`
	Resource RETSID   `xml:",attr,omitempty json:",omitempty"`

	ValidationExternal []ValidationExternal `xml:",omitempty json:",omitempty"`
}

// ValidationExternal ...
type ValidationExternal struct {
	MetadataEntryID        RETSID   `xml:",omitempty json:",omitempty"`
	ValidationExternalName RETSName `xml:",omitempty json:",omitempty"`
	SearchResource         RETSName `xml:",omitempty json:",omitempty"`
	SearchClass            RETSName `xml:",omitempty json:",omitempty"`
	Date                   DateTime `xml:",attr,omitempty json:",omitempty"`
	Version                Version  `xml:",attr,omitempty json:",omitempty"`

	MValidationExternalType MValidationExternalType `xml:"METADATA-VALIDATION_EXTERNAL_TYPE,omitempty json:"METADATA-VALIDATION_EXTERNAL_TYPE,omitempty"`
}

// MValidationExternalType ...
type MValidationExternalType struct {
	Date                   DateTime                 `xml:",attr,omitempty json:",omitempty"`
	Version                Version                  `xml:",attr,omitempty json:",omitempty"`
	Resource               RETSID                   `xml:",attr,omitempty json:",omitempty"`
	ValidationExternalName RETSID                   `xml:",attr,omitempty json:",omitempty"`
	ValidationExternalType []ValidationExternalType `xml:",omitempty json:",omitempty"`
}

// ValidationExternalType ...
type ValidationExternalType struct {
	MetadataEntryID RETSID     `xml:",omitempty json:",omitempty"`
	SearchField     PlainText  `xml:",omitempty json:",omitempty"`
	DisplayField    PlainText  `xml:",omitempty json:",omitempty"`
	ResultFields    StringList `xml:",omitempty json:",omitempty"` // TODO plaintext list
}
