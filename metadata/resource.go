package gorets_metadata

import (
	"time"
)

type MResource struct {
	Version   string
	Date      time.Time
	Resources []Resource
}

type Resource struct {
	ResourceID                  string
	StandardName                string
	VisibleName                 string
	Description                 string
	KeyField                    string
	ClassCount                  int
	ClassVersion                string
	ClassDate                   *time.Time
	ObjectVersion               string
	ObjectDate                  *time.Time
	SearchHelpVersion           string
	SearchHelpDate              *time.Time
	EditMaskVersion             string
	EditMaskDate                *time.Time
	LookupVersion               string
	LookupDate                  *time.Time
	UpdateHelpVersion           string
	UpdateHelpDate              *time.Time
	ValidationExpressionVersion string
	ValidationExpressionDate    *time.Time
	ValidationLookupVersion     string
	ValidationLookupDate        *time.Time
	ValidationExternalVersion   string
	ValidationExternalDate      *time.Time
}
