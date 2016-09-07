package metadata

// MValidationExpression ...
type MValidationExpression struct {
	Date     DateTime `xml:"Date,attr"`
	Version  Version  `xml:"Version,attr"`
	Resource RETSID   `xml:"Resource,attr"`

	ValidationExpression []ValidationExpression `xml:"ValidationExpression"`
}

// ValidationExpression ...
type ValidationExpression struct {
	MetadataEntryID          RETSID   `xml:"MetadataEntryID"`
	ValidationExpressionID   RETSName `xml:"ValidationExpressionID"`
	ValidationExpressionType AlphaNum `xml:"ValidationExpressionType"`
	Message                  Text     `xml:"Message"`
	IsCaseSensitive          Boolean  `xml:"IsCaseSensitive"`
}
