const dataA = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "Metadata":{
      	"Date": "2000-01-01T01:00:00",
      	"Version": "01.00",
      	"System": {
      		"SystemID": "TEST",
      		"SystemDescription": "TEST METADATA",
      		"TimeZoneOffset": "-00:00",
      		"METADATA-FOREIGN_KEY": {},
      		"METADATA-RESOURCE": {
      			"Date": "2000-01-01T01:00:00",
      			"Version": "01.00",
      			"Resource": [
      				{
      					"ResourceID": "A",
      					"StandardName": "A",
      					"VisibleName": "A-Test",
      					"Description": "A test resource",
      					"KeyField": "AKey",
      					"ClassCount": "1",
      					"ClassVersion": "01.00",
      					"ClassDate": "2000-01-01T01:00:00",
      					"ObjectVersion": "01.00",
      					"ObjectDate": "2000-01-01T01:00:00",
      					"SearchHelpVersion": "01.00",
      					"SearchHelpDate": "2000-01-01T01:00:00",
      					"LookupVersion": "01.00",
      					"LookupDate": "2000-01-01T01:00:00",
      					"METADATA-CLASS": {
      						"Date": "2000-01-01T01:00:00",
      						"Version": "01.00",
      						"Resource": "A",
      						"Class": [
      							{
      								"ClassName": "B",
      								"StandardName": "B",
      								"VisibleName": "BTest",
      								"Description": "B a test class of Resource A",
      								"TableVersion": "01.00",
      								"TableDate": "2000-01-01T01:00:00",
      								"ClassTimeStamp": "LastModifiedDateTime",
      								"HasKeyIndex": "0",
      								"METADATA-TABLE": {
      									"Date": "2000-01-01T01:00:00",
      									"Version": "01.00",
      									"Resource": "A",
      									"Class": "B",
      									"Field": [
      										{
      											"MetadataEntryID": "1",
      											"SystemName": "AKey",
      											"StandardName": "AKey",
      											"LongName": "AKey",
      											"DBName": "AKey",
      											"ShortName": "AKey",
      											"MaximumLength": "10",
      											"DataType": "Int",
      											"Precision": "0",
      											"Searchable": "1",
      											"Interpretation": "Number",
      											"Alignment": "Right",
      											"UseSeparator": "0",
      											"MaxSelect": "0",
      											"Index": "0",
      											"Minimum": "0",
      											"Maximum": "2147483647",
      											"Default": "0",
      											"Required": "0",
      											"SearchHelpID": "AKey",
      											"Unique": "1",
      											"ModTimeStamp": "0",
      											"ForeignField": "AKey",
      											"InKeyIndex": "0"
      										},
                          {
      											"MetadataEntryID": "2",
      											"SystemName": "Lookup",
      											"LongName": "Lookup",
      											"DBName": "Lookup",
      											"ShortName": "Lookup",
      											"MaximumLength": "1",
      											"DataType": "Character",
      											"Precision": "0",
      											"Searchable": "1",
      											"Interpretation": "Lookup",
      											"Alignment": "Left",
      											"UseSeparator": "0",
      											"LookupName": "Lookup",
      											"MaxSelect": "0",
      											"Index": "0",
      											"Minimum": "0",
      											"Maximum": "0",
      											"Default": "0",
      											"Required": "0",
      											"SearchHelpID": "Lookup",
      											"Unique": "0",
      											"ModTimeStamp": "0",
      											"ForeignField": "Lookup",
      											"InKeyIndex": "0"
      										},
                          {
      											"MetadataEntryID": "3",
      											"SystemName": "LastModifiedDateTime",
      											"StandardName": "ModificationTimestamp",
      											"LongName": "LastModified Date Time",
      											"DBName": "LstModDtTm",
      											"ShortName": "LastModified Date Time",
      											"MaximumLength": "19",
      											"DataType": "DateTime",
      											"Precision": "0",
      											"Searchable": "1",
      											"Alignment": "Center",
      											"UseSeparator": "0",
      											"MaxSelect": "0",
      											"Index": "0",
      											"Minimum": "0",
      											"Maximum": "0",
      											"Default": "0",
      											"Required": "0",
      											"SearchHelpID": "LastModifiedDateTime",
      											"Unique": "0",
      											"ModTimeStamp": "0",
      											"InKeyIndex": "0"
      										}
      									]
      								},
      								"METADATA-UPDATE": {},
      								"METADATA-COLUMN_GROUP_SET": {
      									"Date": "",
      									"Version": "",
      									"Resource": "",
      									"Class": "",
      									"ColumnGroupSet": null
      								},
      								"METADATA-COLUMN_GROUP": {
      									"Date": "",
      									"Version": "",
      									"Resource": "",
      									"Class": "",
      									"ColumnGroup": null
      								}
      							}
      						]
      					},
      					"METADATA-OBJECT": {
      						"Date": "2000-01-01T01:00:00",
      						"Version": "01.00",
      						"Resource": "A",
      						"Object": [
      							{
      								"MetadataEntryID": "1",
      								"ObjectType": "Photo",
      								"MIMEType": "image/jpeg",
      								"VisibleName": "Photo",
      								"Description": "Photo"
      							}
      						]
      					},
      					"METADATA-LOOKUP": {
      						"Date": "2000-01-01T01:00:00",
      						"Version": "01.00",
      						"Resource": "A",
      						"Lookup": [
      							{
      								"MetadataEntryID": "1",
      								"LookupName": "Lookup",
      								"VisibleName": "Lookup",
      								"METADATA-LOOKUP_TYPE": {
      									"Date": "2000-01-01T01:00:00",
      									"Version": "01.00",
      									"Resource": "A",
      									"Lookup": "Lookup",
      									"LookupType": [
      										{
      											"MetadataEntryID": "1",
      											"LongValue": "A",
      											"ShortValue": "A",
      											"Value": "A"
      										},
      										{
      											"MetadataEntryID": "2",
      											"LongValue": "B",
      											"ShortValue": "B",
      											"Value": "B"
      										},
      										{
      											"MetadataEntryID": "3",
      											"LongValue": "C",
      											"ShortValue": "C",
      											"Value": "C"
      										}
      									]
      								}
      							}
      						]
      					},
      					"METADATA-SEARCH_HELP": {
      						"Date": "2012-10-11T15:01:18",
      						"Version": "01.72.10927",
      						"Resource": "ActiveAgent",
      						"SearchHelp": [
      							{
      								"MetadataEntryID": "1",
      								"SearchHelpID": "AKey",
      								"Value": "The key value"
      							},
      							{
      								"MetadataEntryID": "2",
      								"SearchHelpID": "Lookup",
      								"Value": "Lookup Test"
      							},
      							{
      								"MetadataEntryID": "3",
      								"SearchHelpID": "LastModifiedDateTime",
      								"Value": "LastModifiedDateTime"
      							}
      						]
      					},
      					"METADATA-EDIT_MASK": {},
      					"METADATA-UPDATE": {},
      					"METADATA-VALIDATION_EXTERNAL": {},
      					"METADATA-VALIDATION_EXPRESSION": {},
      					"METADATA-VALIDATION_LOOKUP": {}
      				}
            ]
      		},
      		"METADATA-FILTER": {}
      	}
      }
    }
  }`
);

const dataB = JSON.parse(
  `{
    "result":{
      "wirelog":"U2FtcGxlIFdpcmVsb2c=",
      "Metadata":{
      	"Date": "2000-01-01T01:00:00",
      	"Version": "01.00",
      	"System": {
      		"SystemID": "ALPHA",
      		"SystemDescription": "ALTERNATE TEST METADATA",
      		"TimeZoneOffset": "-00:00",
      		"METADATA-FOREIGN_KEY": {},
      		"METADATA-RESOURCE": {
      			"Date": "2000-01-01T01:00:00",
      			"Version": "01.00",
      			"Resource": [
      				{
      					"ResourceID": "A",
      					"StandardName": "A",
      					"VisibleName": "A-Test",
      					"Description": "A test resource",
      					"KeyField": "AKey",
      					"ClassCount": "1",
      					"ClassVersion": "01.00",
      					"ClassDate": "2000-01-01T01:00:00",
      					"ObjectVersion": "01.00",
      					"ObjectDate": "2000-01-01T01:00:00",
      					"SearchHelpVersion": "01.00",
      					"SearchHelpDate": "2000-01-01T01:00:00",
      					"LookupVersion": "01.00",
      					"LookupDate": "2000-01-01T01:00:00",
      					"METADATA-CLASS": {
      						"Date": "2000-01-01T01:00:00",
      						"Version": "01.00",
      						"Resource": "A",
      						"Class": [
      							{
      								"ClassName": "B",
      								"StandardName": "B",
      								"VisibleName": "BTest",
      								"Description": "B a test class of Resource A",
      								"TableVersion": "01.00",
      								"TableDate": "2000-01-01T01:00:00",
      								"ClassTimeStamp": "LastModifiedDateTime",
      								"HasKeyIndex": "0",
      								"METADATA-TABLE": {
      									"Date": "2000-01-01T01:00:00",
      									"Version": "01.00",
      									"Resource": "A",
      									"Class": "B",
      									"Field": [
      										{
      											"MetadataEntryID": "1",
      											"SystemName": "AKey",
      											"StandardName": "AKey",
      											"LongName": "AKey",
      											"DBName": "AKey",
      											"ShortName": "AKey",
      											"MaximumLength": "10",
      											"DataType": "Int",
      											"Precision": "0",
      											"Searchable": "1",
      											"Interpretation": "Number",
      											"Alignment": "Right",
      											"UseSeparator": "0",
      											"MaxSelect": "0",
      											"Index": "0",
      											"Minimum": "0",
      											"Maximum": "2147483647",
      											"Default": "0",
      											"Required": "0",
      											"SearchHelpID": "AKey",
      											"Unique": "1",
      											"ModTimeStamp": "0",
      											"ForeignField": "AKey",
      											"InKeyIndex": "0"
      										},
                          {
      											"MetadataEntryID": "2",
      											"SystemName": "Lookup",
      											"LongName": "Lookup",
      											"DBName": "Lookup",
      											"ShortName": "Lookup",
      											"MaximumLength": "1",
      											"DataType": "Character",
      											"Precision": "0",
      											"Searchable": "1",
      											"Interpretation": "Lookup",
      											"Alignment": "Left",
      											"UseSeparator": "0",
      											"LookupName": "Lookup",
      											"MaxSelect": "0",
      											"Index": "0",
      											"Minimum": "0",
      											"Maximum": "0",
      											"Default": "0",
      											"Required": "0",
      											"SearchHelpID": "Lookup",
      											"Unique": "0",
      											"ModTimeStamp": "0",
      											"ForeignField": "Lookup",
      											"InKeyIndex": "0"
      										},
                          {
      											"MetadataEntryID": "3",
      											"SystemName": "LastModifiedDateTime",
      											"StandardName": "ModificationTimestamp",
      											"LongName": "LastModified Date Time",
      											"DBName": "LstModDtTm",
      											"ShortName": "LastModified Date Time",
      											"MaximumLength": "19",
      											"DataType": "DateTime",
      											"Precision": "0",
      											"Searchable": "1",
      											"Alignment": "Center",
      											"UseSeparator": "0",
      											"MaxSelect": "0",
      											"Index": "0",
      											"Minimum": "0",
      											"Maximum": "0",
      											"Default": "0",
      											"Required": "0",
      											"SearchHelpID": "LastModifiedDateTime",
      											"Unique": "0",
      											"ModTimeStamp": "0",
      											"InKeyIndex": "0"
      										}
      									]
      								},
      								"METADATA-UPDATE": {},
      								"METADATA-COLUMN_GROUP_SET": {
      									"Date": "",
      									"Version": "",
      									"Resource": "",
      									"Class": "",
      									"ColumnGroupSet": null
      								},
      								"METADATA-COLUMN_GROUP": {
      									"Date": "",
      									"Version": "",
      									"Resource": "",
      									"Class": "",
      									"ColumnGroup": null
      								}
      							}
      						]
      					},
      					"METADATA-OBJECT": {
      						"Date": "2000-01-01T01:00:00",
      						"Version": "01.00",
      						"Resource": "A",
      						"Object": [
      							{
      								"MetadataEntryID": "1",
      								"ObjectType": "Photo",
      								"MIMEType": "image/jpeg",
      								"VisibleName": "Photo",
      								"Description": "Photo"
      							}
      						]
      					},
      					"METADATA-LOOKUP": {
      						"Date": "2000-01-01T01:00:00",
      						"Version": "01.00",
      						"Resource": "A",
      						"Lookup": [
      							{
      								"MetadataEntryID": "1",
      								"LookupName": "Lookup",
      								"VisibleName": "Lookup",
      								"METADATA-LOOKUP_TYPE": {
      									"Date": "2000-01-01T01:00:00",
      									"Version": "01.00",
      									"Resource": "A",
      									"Lookup": "Lookup",
      									"LookupType": [
      										{
      											"MetadataEntryID": "1",
      											"LongValue": "A",
      											"ShortValue": "A",
      											"Value": "A"
      										},
      										{
      											"MetadataEntryID": "2",
      											"LongValue": "B",
      											"ShortValue": "B",
      											"Value": "B"
      										},
      										{
      											"MetadataEntryID": "3",
      											"LongValue": "C",
      											"ShortValue": "C",
      											"Value": "C"
      										}
      									]
      								}
      							}
      						]
      					},
      					"METADATA-SEARCH_HELP": {
      						"Date": "2012-10-11T15:01:18",
      						"Version": "01.72.10927",
      						"Resource": "ActiveAgent",
      						"SearchHelp": [
      							{
      								"MetadataEntryID": "1",
      								"SearchHelpID": "AKey",
      								"Value": "The key value"
      							},
      							{
      								"MetadataEntryID": "2",
      								"SearchHelpID": "Lookup",
      								"Value": "Lookup Test"
      							},
      							{
      								"MetadataEntryID": "3",
      								"SearchHelpID": "LastModifiedDateTime",
      								"Value": "LastModifiedDateTime"
      							}
      						]
      					},
      					"METADATA-EDIT_MASK": {},
      					"METADATA-UPDATE": {},
      					"METADATA-VALIDATION_EXTERNAL": {},
      					"METADATA-VALIDATION_EXPRESSION": {},
      					"METADATA-VALIDATION_LOOKUP": {}
      				}
            ]
      		},
      		"METADATA-FILTER": {}
      	}
      }
    }
  }`
);

const dataE = JSON.parse(
  `{
    "error":"Unknown Config"
  }`
);

module.exports = {
  getData: (config) => {
    console.log(config);
    if (config.id === 'testp:johndoe') {
      return { json: () => dataA };
    }
    if (config.id === 'zzzzp:janedoe') {
      return { json: () => dataB };
    }
    return { json: () => dataE };
  },
};
