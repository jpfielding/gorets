package metadata

// MValidationExpression ...
type MValidationExpression struct {
	Date     DateTime `xml:"Date,attr,omitempty"`
	Version  Version  `xml:"Version,attr,omitempty"`
	Resource RETSID   `xml:"Resource,attr,omitempty"`

	ValidationExpression []ValidationExpression `xml:"ValidationExpression,omitempty"`
}

// ValidationExpression ...
type ValidationExpression struct {
	MetadataEntryID          RETSID   `xml:"MetadataEntryID,omitempty"`
	ValidationExpressionID   RETSName `xml:"ValidationExpressionID,omitempty"`
	ValidationExpressionType AlphaNum `xml:"ValidationExpressionType,omitempty"`
	Message                  Text     `xml:"Message,omitempty"`
	IsCaseSensitive          Boolean  `xml:"IsCaseSensitive,omitempty"`
}
