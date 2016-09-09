package metadata

// MEditMask ...
type MEditMask struct {
	Date     DateTime   `xml:"Date,attr,omitempty"`
	Version  Version    `xml:"Version,attr,omitempty"`
	Resource RETSID     `xml:"Resource,attr,omitempty"`
	EditMask []EditMask `xml:"EditMask,omitempty"`
}

// EditMask ...
type EditMask struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID,omitempty"`
	EditMaskID      RETSName `xml:"EditMaskID,omitempty"`
	Value           Text     `xml:"Value,omitempty"`
}
