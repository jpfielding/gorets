package gorets_metadata

import (
	"time"
)

type MUpdate struct {
	Version         string
	Date            time.Time
	Resource, Class string
	Updates         []Update
}

type Update struct {
	MetadataEntryID string
	// Add, Clone, Change, BeginUpdate, CancelUpdate, ShowLocks
	UpdateAction      string
	Description       string
	KeyField          string
	UpdateTypeVersion string
	UpdateTypeDate    *time.Time
	RequiresBeing     *bool
}
