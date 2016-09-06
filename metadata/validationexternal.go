package metadata

import "time"

// MValidationExternal ...
type MValidationExternal struct {
	Date     time.Time `xml:"Date,attr"`
	Version  Version   `xml:"Version,attr"`
	Resource string    `xml:"Resource,attr"`

	ValidationExternal []ValidationExternal `xml:"ValidationExternal"`
}

// ValidationExternal ...
type ValidationExternal struct {
	Date                   time.Time `xml:"Date,attr"`
	Version                Version   `xml:"Version,attr"`
	MetadataEntryID        string    `xml:"MetadataEntryID"`
	ValidationExternalName string    `xml:"ValidationExternalName"`
	SearchResource         string    `xml:"SearchResource"`
	SearchClass            string    `xml:"SearchClass"`
}
