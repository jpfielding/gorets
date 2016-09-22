package metadata

// MColumnGroup ...
type MColumnGroup struct {
	Date        DateTime      `xml:",attr,omitempty"`
	Version     Version       `xml:",attr,omitempty"`
	Resource    RETSID        `xml:",attr,omitempty"`
	Class       RETSID        `xml:",attr,omitempty"`
	ColumnGroup []ColumnGroup `xml:",omitempty"`
}

// ColumnGroup ...
type ColumnGroup struct {
	MetadataEntryID   RETSID   `xml:",omitempty" json:",omitempty"`
	ColumnGroupName   RETSID   `xml:",omitempty" json:",omitempty"`
	ControlSystemName RETSID   `xml:",omitempty" json:",omitempty"`
	LongName          RETSName `xml:"omitempty" json:",omitempty"`
	ShortName         RETSName `xml:"omitempty" json:",omitempty"`
	Description       Text     `xml:"omitempty" json:",omitempty"`

	MColumnGroupControl       MColumnGroupControl       `xml:"METADATA-COLUMN_GROUP_CONTROL,omitempty" json:"METADATA-COLUMN_GROUP_CONTROL,omitempty"`
	MColumnGroupTable         MColumnGroupTable         `xml:"METADATA-COLUMN_GROUP_TABLE,omitempty" json:"METADATA-COLUMN_GROUP_TABLE,omitempty"`
	MColumnGroupNormalization MColumnGroupNormalization `xml:"METADATA-COLUMN_GROUP_NORMALIZATION,omitempty" json:"METADATA-COLUMN_GROUP_NORMALIZATION,omitempty"`
}

// MColumnGroupControl ...
type MColumnGroupControl struct {
	Date               DateTime             `xml:",attr,omitempty"`
	Version            Version              `xml:",attr,omitempty"`
	Resource           RETSID               `xml:",attr,omitempty"`
	Class              RETSID               `xml:",attr,omitempty"`
	ColumnGroup        RETSID               `xml:",attr,omitempty"`
	ColumnGroupControl []ColumnGroupControl `xml:",omitempty"`
}

// ColumnGroupControl ...
type ColumnGroupControl struct {
	MetadataEntryID RETSID  `xml:",omitempty" json:",omitempty"`
	LowValue        Numeric `xml:",omitempty" json:",omitempty"`
	HighValue       Numeric `xml:",omitempty" json:",omitempty"`
}

// MColumnGroupTable ...
type MColumnGroupTable struct {
	Date             DateTime           `xml:",attr,omitempty"`
	Version          Version            `xml:",attr,omitempty"`
	Resource         RETSID             `xml:",attr,omitempty"`
	Class            RETSID             `xml:",attr,omitempty"`
	ColumnGroup      RETSID             `xml:",attr,omitempty"`
	ColumnGroupTable []ColumnGroupTable `xml:",omitempty"`
}

// ColumnGroupTable ...
type ColumnGroupTable struct {
	MetadataEntryID    RETSID   `xml:",omitempty" json:",omitempty"`
	SystemName         RETSID   `xml:"omitempty" json:",omitempty"`
	ColumnGroupSetName RETSID   `xml:"omitempty" json:",omitempty"`
	LongName           RETSName `xml:"omitempty" json:",omitempty"`
	ShortName          RETSName `xml:"omitempty" json:",omitempty"`
	DisplayOrder       Numeric  `xml:"omitempty" json:",omitempty"`
	DisplayLength      Numeric  `xml:"omitempty" json:",omitempty"`
	DisplayHeight      Numeric  `xml:"omitempty" json:",omitempty"`
	ImmediateRefresh   Boolean  `xml:"omitempty" json:",omitempty"`
}

// MColumnGroupNormalization ...
type MColumnGroupNormalization struct {
	Date                     DateTime                   `xml:",attr,omitempty"`
	Version                  Version                    `xml:",attr,omitempty"`
	Resource                 RETSID                     `xml:",attr,omitempty"`
	Class                    RETSID                     `xml:",attr,omitempty"`
	ColumnGroup              RETSID                     `xml:",attr,omitempty"`
	ColumnGroupNormalization []ColumnGroupNormalization `xml:",omitempty"`
}

// ColumnGroupNormalization ...
type ColumnGroupNormalization struct {
	MetadataEntryID RETSID  `xml:",omitempty" json:",omitempty"`
	TypeIdentifier  RETSID  `xml:"omitempty" json:",omitempty"`
	Sequence        Numeric `xml:"omitempty" json:",omitempty"`
	ColumnLabel     RETSID  `xml:"omitempty" json:",omitempty"`
	SystemName      RETSID  `xml:"omitempty" json:",omitempty"`
}
