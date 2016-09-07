package metadata

// MSearchHelp ...
type MSearchHelp struct {
	Date       DateTime     `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   RETSID       `xml:"Resource,attr"`
	SearchHelp []SearchHelp `xml:"SearchHelp"`
}

// SearchHelp ...
type SearchHelp struct {
	MetadataEntryID RETSID   `xml:"MetadataEntryID"`
	SearchHelpID    RETSName `xml:"SearchHelpID"`
	Value           Text     `xml:"Value"`
}
