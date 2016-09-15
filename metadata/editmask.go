package metadata

// MEditMask ...
type MEditMask struct {
	Date     DateTime   `xml:",attr,omitempty" json:",omitempty"`
	Version  Version    `xml:",attr,omitempty" json:",omitempty"`
	Resource RETSID     `xml:",attr,omitempty" json:",omitempty"`
	EditMask []EditMask `xml:",omitempty" json:",omitempty"`
}

// EditMask ...
type EditMask struct {
	MetadataEntryID RETSID   `xml:",omitempty" json:",omitempty"`
	EditMaskID      RETSName `xml:",omitempty" json:",omitempty"`
	Value           Text     `xml:",omitempty" json:",omitempty"`
}
