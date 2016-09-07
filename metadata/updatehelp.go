package metadata

// MUpdateHelp ...
type MUpdateHelp struct {
	Date       DateTime     `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   RETSID       `xml:"Resource,attr"`
	UpdateHelp []UpdateHelp `xml:"UpdateHelp"`
}

// UpdateHelp ...
type UpdateHelp struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID"`
	UpdateHelpID    RETSName `xml:"UpdateHelpID"`
	Value           Text     `xml:"Value"`
}
