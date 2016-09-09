package metadata

// MUpdate ...
type MUpdate struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`
	Class    RETSID   `xml:"Class,attr,omitempty"`
	Update   []Update `xml:"Update,omitempty"`
}

// Update ...
type Update struct {
	MetadataEntryID   RETSID    `xml:"MetadataEntryID,omitempty"`
	UpdateAction      AlphaNum  `xml:"UpdateAction,omitempty"` // some standardish names add,clone,change,delete,beginupdate
	Description       PlainText `xml:"Description,omitempty"`
	KeyField          RETSName  `xml:"KeyField,omitempty"`
	UpdateTypeVersion Version   `xml:"UpdateTypeVersion,omitempty"`
	UpdateTypeDate    DateTime  `xml:"UpdateTypeDate,omitempty"`
	RequiresBegin     Boolean   `xml:"RequiresBegin,omitempty"`

	MUpdateType MUpdateType `xml:"METADATA-UPDATE_TYPE,omitempty"`
}

// MUpdateType ...
type MUpdateType struct {
	Date       DateTime     `xml:"Date,attr,omitempty"`
	Version    Version      `xml:"Version,attr,omitempty"`
	Resource   RETSID       `xml:"Resource,attr,omitempty"`
	Update     RETSID       `xml:"Update,attr,omitempty"`
	UpdateType []UpdateType `xml:"Update,omitempty"`
}

// UpdateType ...
type UpdateType struct {
	MetadataEntryID        RETSID      `xml:"MetadataEntryID,omitempty"`
	SystemName             RETSName    `xml:"SystemName,omitempty"`
	Sequence               Numeric     `xml:"Sequence,omitempty"`
	Attributes             NumericList `xml:"Attributes,omitempty"` // TODO limit to 1-7
	Default                PlainText   `xml:"Default,omitempty"`
	ValidationExpressionID RETSNames   `xml:"ValidationExpressionID,omitempty"`
	UpdateHelpID           RETSName    `xml:"UpdateHelpID,omitempty"`
	ValidationLookupName   RETSName    `xml:"ValidationLookupName,omitempty"` // deprecated
	ValidationExternalName RETSName    `xml:"ValidationExternalName,omitempty"`
	MaxUpdate              Numeric     `xml:"MaxUpdate,omitempty"`
	SearchResultOrder      Numeric     `xml:"SearchResultOrder,omitempty"`
	SearchQueryOrder       Numeric     `xml:"SearchQueryOrder,omitempty"`
}
