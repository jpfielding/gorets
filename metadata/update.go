package metadata

import "time"

// MUpdate ...
type MUpdate struct {
	Date     time.Time `xml:"Date,attr"`
	Version  Version   `xml:"Version,attr"`
	Resource string    `xml:"Resource,attr"`
	Class    string    `xml:"Class,attr"`
	Update   []Update  `xml:"Update"`
}

// Update ...
type Update struct {
	MetadataEntryID   string    `xml:"MetadataEntryID"`
	UpdateAction      string    `xml:"UpdateAction"`
	Description       string    `xml:"Description"`
	KeyField          string    `xml:"KeyField"`
	UpdateTypeVersion Version   `xml:"UpdateTypeVersion"`
	UpdateTypeDate    time.Time `xml:"UpdateTypeDate"`
	RequiresBegin     *int      `xml:"RequiresBegin"`

	MUpdateType MUpdateType `xml:"METADATA-UPDATE_TYPE"`
}

// MUpdateType ...
type MUpdateType struct {
	Date       time.Time    `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   string       `xml:"Resource,attr"`
	Lookup     string       `xml:"Update,attr"`
	UpdateType []UpdateType `xml:"Update"`
}

// UpdateType ...
type UpdateType struct {
	MetadataEntryID        string `xml:"MetadataEntryID"`
	SystemName             string `xml:"SystemName"`
	Sequence               int    `xml:"Sequence"`
	Attributes             string `xml:"Attributes"`
	Default                string `xml:"Defualt"`
	ValidationExpressionID string `xml:"ValidationExpressionID"`
	UpdateHelpID           string `xml:"UpdateHelpID"`
	ValidationLookupName   string `xml:"ValidationLookupName"` // deprecated
	ValidationExternalName string `xml:"ValidationExternalName"`
	MaxUpdate              *int   `xml:"MaxUpdate"`
	SearchResultOrder      *int   `xml:"SearchResultOrder"`
	SearchQueryOrder       *int   `xml:"SearchQueryOrder"`
}
