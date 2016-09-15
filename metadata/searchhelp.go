package metadata

// MSearchHelp ...
type MSearchHelp struct {
	Date       DateTime     `xml:",attr,omitempty" json:",omitempty"`
	Version    Version      `xml:",attr,omitempty" json:",omitempty"`
	Resource   RETSID       `xml:",attr,omitempty" json:",omitempty"`
	SearchHelp []SearchHelp `xml:",omitempty" json:",omitempty"`
}

// SearchHelp ...
type SearchHelp struct {
	MetadataEntryID RETSID   `xml:",omitempty" json:",omitempty"`
	SearchHelpID    RETSName `xml:",omitempty" json:",omitempty"`
	Value           Text     `xml:",omitempty" json:",omitempty"`
}
