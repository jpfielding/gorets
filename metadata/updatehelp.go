package metadata

// MUpdateHelp ...
type MUpdateHelp struct {
	Date       DateTime     `xml:"Date,attr,omitempty"`
	Version    Version      `xml:"Version,attr,omitempty"`
	Resource   RETSID       `xml:"Resource,attr,omitempty"`
	UpdateHelp []UpdateHelp `xml:"UpdateHelp,omitempty"`
}

// UpdateHelp ...
type UpdateHelp struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID,omitempty"`
	UpdateHelpID    RETSName `xml:"UpdateHelpID,omitempty"`
	Value           Text     `xml:"Value,omitempty"`
}
