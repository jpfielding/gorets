package metadata

import "time"

// MUpdateHelp ...
type MUpdateHelp struct {
	Date       time.Time    `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   string       `xml:"Resource,attr"`
	UpdateHelp []UpdateHelp `xml:"UpdateHelp"`
}

// UpdateHelp ...
type UpdateHelp struct {
	MetadataEntryID string `xml:"MetadataEntryID"`
	UpdateHelpID    string `xml:"UpdateHelpID"`
	Value           string `xml:"Value"`
}
