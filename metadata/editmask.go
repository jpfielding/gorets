package metadata

// MEditMask ...
type MEditMask struct {
	Date     DateTime   `xml:"Date,attr"`
	Version  Version    `xml:"Version,attr"`
	Resource RETSID     `xml:"Resource,attr"`
	EditMask []EditMask `xml:"EditMask"`
}

// EditMask ...
type EditMask struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID"`
	EditMaskID      RETSName `xml:"EditMaskID"`
	Value           Text     `xml:"Value"`
}
