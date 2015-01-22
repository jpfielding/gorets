package gorets_metadata

import (
	"time"
)

type MClass struct {
	Version  string
	Date     time.Time
	Resource string
	Classes  []Class
}
type Class struct {
	ClassName        string
	StandardName     string
	VisibleName      string
	Description      string
	TableVersion     string
	TableDate        string
	UpdateVersion    string
	UpdateDate       string
	ClassTimeStamp   *time.Time
	DeletedFlagField string
	DeletedFlagValue string
	HasKeyIndex      *bool
	OffsetSupport    *bool
}
