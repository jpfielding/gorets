package metadata

import "time"

// MResource ...
type MResource struct {
	Date     time.Time  `xml:"Date,attr"`
	Version  Version    `xml:"Version,attr"`
	Resource []Resource `xml:"Resource"`
}

// Resource ...
type Resource struct {
	ResourceID                  string    `xml:"ResourceID"`
	StandardName                string    `xml:"StandardName"`
	VisibleName                 string    `xml:"VisibleName"`
	Description                 string    `xml:"Description"`
	KeyField                    string    `xml:"KeyField"`
	ClassCount                  int       `xml:"ClassCount"`
	ClassVersion                Version   `xml:"ClassVersion"`
	ClassDate                   time.Time `xml:"ClassDate"`
	ObjectVersion               Version   `xml:"ObjectVersion"`
	ObjectDate                  time.Time `xml:"ObjectDate"`
	SearchHelpVersion           Version   `xml:"SearchHelpVersion"`
	SearchHelpDate              time.Time `xml:"SearchHelpDate"`
	EditMaskVersion             Version   `xml:"EditMaskVersion"`
	EditMaskDate                time.Time `xml:"EditMaskDate"`
	LookupVersion               Version   `xml:"LookupVersion"`
	LookupDate                  time.Time `xml:"LookupDate"`
	UpdateHelpVersion           Version   `xml:"UpdateHelpVersion"`
	UpdateHelpDate              time.Time `xml:"UpdateHelpDate"`
	ValidationExpressionVersion Version   `xml:"ValidationExpressionVersion"`
	ValidationExpressionDate    time.Time `xml:"ValidationExpressionDate"`
	ValidationLookupVersion     Version   `xml:"ValidationLookupVersion"`
	ValidationLookupDate        time.Time `xml:"ValidationLookupDate"`
	ValidationExternalVersion   Version   `xml:"ValidationExternalVersion"`
	ValidationExternalDate      time.Time `xml:"ValidationExternalDate"`

	// the resource children
	MClass                MClass                `xml:"METADATA-CLASS"`
	MObject               MObject               `xml:"METADATA-OBJECT"`
	MLookup               MLookup               `xml:"METADATA-LOOKUP"`
	MSearchHelp           MSearchHelp           `xml:"METADATA-SEARCH_HELP"`
	MEditMask             MEditMask             `xml:"METADATA-EDIT_MASK"`
	MUpdateHelp           MUpdateHelp           `xml:"METADATA-UPDATE_HELP"`
	MValidationLookup     MValidationLookup     `xml:"METADATA-VALIDATION_LOOKUP"`
	MValidationExternal   MValidationExternal   `xml:"METADATA-VALIDATION_EXTERNAL"`
	MValidationExpression MValidationExpression `xml:"METADATA-VALIDATION_EXPRESSION"`
}
