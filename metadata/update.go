package metadata

// MUpdate ...
type MUpdate struct {
	Date     DateTime `xml:",attr,omitempty" json:",omitempty"`
	Version  Version  `xml:",attr,omitempty" json:",omitempty"`
	Resource RETSID   `xml:",attr,omitempty" json:",omitempty"`
	Class    RETSID   `xml:",attr,omitempty" json:",omitempty"`
	Update   []Update `xml:",omitempty" json:",omitempty"`
}

// Update ...
type Update struct {
	MetadataEntryID   RETSID    `xml:",omitempty" json:",omitempty"`
	UpdateAction      AlphaNum  `xml:",omitempty" json:",omitempty"` // some standardish names add,clone,change,delete,beginupdate
	Description       PlainText `xml:",omitempty" json:",omitempty"`
	KeyField          RETSName  `xml:",omitempty" json:",omitempty"`
	UpdateTypeVersion Version   `xml:",omitempty" json:",omitempty"`
	UpdateTypeDate    DateTime  `xml:",omitempty" json:",omitempty"`
	RequiresBegin     Boolean   `xml:",omitempty" json:",omitempty"`

	MUpdateType MUpdateType `xml:"METADATA-UPDATE_TYPE,omitempty" json:"METADATA-UPDATE_TYPE,,omitempty"`
}

// MUpdateType ...
type MUpdateType struct {
	Date       DateTime     `xml:",attr,omitempty" json:",omitempty"`
	Version    Version      `xml:",attr,omitempty" json:",omitempty"`
	Resource   RETSID       `xml:",attr,omitempty" json:",omitempty"`
	Update     RETSID       `xml:",attr,omitempty" json:",omitempty"`
	UpdateType []UpdateType `xml:",omitempty" json:",omitempty"`
}

// UpdateType ...
type UpdateType struct {
	MetadataEntryID        RETSID      `xml:",omitempty" json:",omitempty"`
	SystemName             RETSName    `xml:",omitempty" json:",omitempty"`
	Sequence               Numeric     `xml:",omitempty" json:",omitempty"`
	Attributes             NumericList `xml:",omitempty" json:",omitempty"` // TODO limit to 1-7
	Default                PlainText   `xml:",omitempty" json:",omitempty"`
	ValidationExpressionID RETSNames   `xml:",omitempty" json:",omitempty"`
	UpdateHelpID           RETSName    `xml:",omitempty" json:",omitempty"`
	ValidationLookupName   RETSName    `xml:",omitempty" json:",omitempty"` // deprecated
	ValidationExternalName RETSName    `xml:",omitempty" json:",omitempty"`
	MaxUpdate              Numeric     `xml:",omitempty" json:",omitempty"`
	SearchResultOrder      Numeric     `xml:",omitempty" json:",omitempty"`
	SearchQueryOrder       Numeric     `xml:",omitempty" json:",omitempty"`
}
