package metadata

// MUpdateHelp ...
type MUpdateHelp struct {
	Date       DateTime     `xml:",attr,omitempty json:",omitempty"`
	Version    Version      `xml:",attr,omitempty json:",omitempty"`
	Resource   RETSID       `xml:",attr,omitempty json:",omitempty"`
	UpdateHelp []UpdateHelp `xml:",omitempty json:",omitempty"`
}

// UpdateHelp ...
type UpdateHelp struct {
	MetadataEntryID RETSID   `xml:",omitempty json:",omitempty"`
	UpdateHelpID    RETSName `xml:",omitempty json:",omitempty"`
	Value           Text     `xml:",omitempty json:",omitempty"`
}
