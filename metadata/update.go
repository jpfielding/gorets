package metadata

// MUpdate ...
type MUpdate struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`
	Class    RETSID   `xml:"Class,attr"`
	Update   []Update `xml:"Update"`
}

// Update ...
type Update struct {
	MetadataEntryID   RETSID    `xml:"MetadataEntryID"`
	UpdateAction      AlphaNum  `xml:"UpdateAction"` // some standardish names add,clone,change,delete,beginupdate
	Description       PlainText `xml:"Description"`
	KeyField          RETSName  `xml:"KeyField"`
	UpdateTypeVersion Version   `xml:"UpdateTypeVersion"`
	UpdateTypeDate    DateTime  `xml:"UpdateTypeDate"`
	RequiresBegin     Boolean   `xml:"RequiresBegin"`

	MUpdateType MUpdateType `xml:"METADATA-UPDATE_TYPE"`
}

// MUpdateType ...
type MUpdateType struct {
	Date       DateTime     `xml:"Date,attr"`
	Version    Version      `xml:"Version,attr"`
	Resource   RETSID       `xml:"Resource,attr"`
	Update     RETSID       `xml:"Update,attr"`
	UpdateType []UpdateType `xml:"Update"`
}

// UpdateType ...
type UpdateType struct {
	MetadataEntryID        RETSID      `xml:"MetadataEntryID"`
	SystemName             RETSName    `xml:"SystemName"`
	Sequence               Numeric     `xml:"Sequence"`
	Attributes             NumericList `xml:"Attributes"` // TODO limit to 1-7
	Default                PlainText   `xml:"Default"`
	ValidationExpressionID RETSNames   `xml:"ValidationExpressionID"`
	UpdateHelpID           RETSName    `xml:"UpdateHelpID"`
	ValidationLookupName   RETSName    `xml:"ValidationLookupName"` // deprecated
	ValidationExternalName RETSName    `xml:"ValidationExternalName"`
	MaxUpdate              Numeric     `xml:"MaxUpdate"`
	SearchResultOrder      Numeric     `xml:"SearchResultOrder"`
	SearchQueryOrder       Numeric     `xml:"SearchQueryOrder"`
}
