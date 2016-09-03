package metadata

import "time"

// MEditMask ...
type MEditMask struct {
	Date     time.Time  `xml:"Date,attr"`
	Version  Version    `xml:"Version,attr"`
	Resource string     `xml:"Resource,attr"`
	EditMask []EditMask `xml:"EditMask"`
}

// EditMask ...
type EditMask struct {
	MetadataEntryID string `xml:"MetadataEntryID"`
	EditMaskID      string `xml:"EditMaskID"`
	Value           string `xml:"Value"`
}
