package rets

const (
	StatusOK          = 0
	StatusSystemError = 10000

	// -----------------------------------------
	// 4.11 Login Reply Codes

	StatusZeroBalance               = 20003
	StatusBrokerCodeRequired        = 20012
	StatusBrokerCodeInvalid         = 20013
	StatusMultipleLogins            = 20022
	StatusServerLoginError          = 20036
	StatusClientAuthFailed          = 20037
	StatusUserAgentAuthRequired     = 20041
	StatusServerTemporarilyDisabled = 20050

	StatusInsecurePassword         = 20140
	StatusRepeatPassword           = 20141
	StatusInvalidEncryptedUsername = 20142

	// -----------------------------------------
	// 7.8 Search Reply Codes

	StatusUnknownQueryField         = 20200
	StatusNoRecords                 = 20201
	StatusInvalidSelect             = 20202
	StatusSearchError               = 20203
	StatusInvalidQuerySyntax        = 20206
	StatusUnauthorizedQuery         = 20207
	StatusMaxRecordsExceeded        = 20208
	StatusTimeout                   = 20209
	StatusTooManyOutstandingQueries = 20210
	StatusQueryTooComplex           = 20211
	StatusInvalidKeyRequest         = 20212 // deprecated
	StatusInvalidKey                = 20213 // deprecated

	// -----------------------------------------
	// 10.7 Update Reply Codes

	StatusInvalidParameter = 20301
	StatusSaveError        = 20302
	StatusUpdateError      = 20303
	StatusWarningResponse2 = 20311
	StatusWarningResponse0 = 20312

	// -----------------------------------------
	// 5.12 GetObject Reply Codes

	StatusInvalidResource            = 20400
	StatusInvalidType                = 20401
	StatusInvalidIdentifier          = 20402
	StatusObjectNotFound             = 20403
	StatusUnsupportedMimeType        = 20406
	StatusUnauthorizedRetrieval      = 20407
	StatusResourceUnvailable         = 20408
	StatusObjectUnavailable          = 20409
	StatusRequestTooLarge            = 20410
	StatusExecutionTimeout           = 20411
	StatusTooManyOutstandingRequests = 20412
	StatusMiscGetObjectError         = 20413

	// -----------------------------------------
	// 12.8 Metadata Reply Codes

	StatusUnknwownMetadataResource           = 20500
	StatusUnknownMetadataType                = 20501
	StatusUnknownMetadataIdentifier          = 20502
	StatusNoMetadataFound                    = 20503
	StatusUnsupportedMetadataMimeType        = 20506
	StatusUnauthorizedMetadataRetrieval      = 20507
	StatusMetadataResourceUnvailable         = 20508
	StatusMetadataUnavailable                = 20509
	StatusMetadataRequestTooLarge            = 20510
	StatusMetadataExecutionTimeout           = 20511
	StatusTooManyOutstandingMetadataRequests = 20512

	StatusMiscMetadataError     = 20513
	StatusDTDVersionUnavailable = 20514

	// -----------------------------------------
	// 6.6 Logout Reply Codes

	StatusNotLoggedIn     = 20701
	StatusMiscLogoutError = 20702
)
