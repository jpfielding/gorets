package metadata

import "time"

// MTable ...
type MTable struct {
	Date     time.Time `xml:"Date,attr"`
	Version  Version   `xml:"Version,attr"`
	Resource string    `xml:"Resource,attr"`
	Class    string    `xml:"Class,attr"`
	Table    []Table   `xml:"Table"`
}

// Table ...
type Table struct {
	SystemName     string `xml:"SystemName"`
	StandardName   string `xml:"StandardName"`
	LongName       string `xml:"LongName"`
	DBName         string `xml:"DBName"`
	ShortName      string `xml:"ShortName"`
	MaximumLength  int    `xml:"MaximumLength"`
	DataType       string `xml:"DataType"`
	Precision      string `xml:"Precision"`
	Searchable     int    `xml:"Searchable"`
	Interpretation string `xml:"Interpretation"`
	Alignment      string `xml:"Alignment"`
	UseSeparator   int    `xml:"UseSeparator"`
	EditMaskID     string `xml:"EditMaskID"`
	LookupName     string `xml:"LookupName"`
	MaxSelect      int    `xml:"MaxSelect"`
	Units          string `xml:"Units"`
	Index          int    `xml:"Index"`
	Minimum        string `xml:"Minimum"`
	Maximum        string `xml:"Maxiumum"`
	Default        string `xml:"Default"`
	Required       int    `xml:"Required"`
	SearchHelpID   string `xml:"SearchHelpID"`
	Unique         int    `xml:"Unique"`
}
