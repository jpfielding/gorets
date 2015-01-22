package gorets_metadata

import (
	"time"
)

type MTable struct {
	Version         string
	Date            time.Time
	Resource, Class string
}

type Table struct {
	MetadataEntryID string
	SystemName      string
	StandardName    string
	LongName        string
	DBName          string
	ShortName       string
	MaximumLength   int
	// Boolean, Character, Date, DateTime, Time, Tiny, Small, Int, Long, Decimal
	DataType   string
	Precision  string
	Searchable *bool
	// Number, Currency, Lookup, LookupMulti, URI
	Interpretation string
	// Left, Right, Center, Justify
	Alignment    string
	UseSeparator *bool
	EditMaskID   string
	LookupName   string
	MaxSelect    *int
	// Feet | Meters | SqFt | SqMeters | Acres | Hectares
	Units              string
	Index              *bool
	Minimum            *int
	Maximum            *int
	Default            *int
	Required           *int
	SearchHelpID       string
	Unique             *bool
	ForeignKeyName     string
	ForeignKeyField    string
	InKeyIndex         *bool
	FilterParentField  string
	DefaultSearchOrder *int
	// UPPER, LOWER, EXACT, MIXED
	Case string
}
