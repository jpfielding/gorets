package metadata

// MSearchHelp ...
type MSearchHelp struct {
	Date       DateTime     `xml:"Date,attr,omitempty"`
	Version    Version      `xml:"Version,attr,omitempty"`
	Resource   RETSID       `xml:"Resource,attr,omitempty"`
	SearchHelp []SearchHelp `xml:"SearchHelp,omitempty"`
}

// SearchHelp ...
type SearchHelp struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID,omitempty"`
	SearchHelpID    RETSName `xml:"SearchHelpID,omitempty"`
	Value           Text     `xml:"Value,omitempty"`
}
