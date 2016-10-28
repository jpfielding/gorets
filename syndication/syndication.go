package syndication

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

// EachListing listens to a function that walks listings
type EachListing func(Listing, error) error

// ToListing creates an adapter to be used with something that walks a large stream
// and segments it into smaller doms
func ToListing(each EachListing) func(io.ReadCloser, error) error {
	return func(body io.ReadCloser, err error) error {
		if err != nil {
			return err
		}
		listing := Listing{}
		err = xml.NewDecoder(body).Decode(&listing)
		return each(listing, err)
	}
}

// initial Syndication mappings courtesy of http://github.com/cridenour

// Listings is a top level element for sharing a
// collection of listings.
type Listings struct {
	// The unique id, or key for this collection of listings.
	ListingsKey string `xml:"listingsKey,attr"`
	// The version number for this schema.
	Version string `xml:"version,attr"`
	// The fully formatted version number.
	VersionTimestamp string `xml:"versionTimestamp,attr"`
	// The language used. Defaults to US English.
	Language string `xml:"lang,attr"`

	// Zero or more listings is contained within the item type.
	Listings []Listing `xml:"Listing"`
	// The disclaimer string for the collection of listings
	Disclaimer SecureString `xml:"Disclaimer"`
}

// OtherChoice is to be used within enumerations requiring
// the presence of an attribute descriptor. This
// attribute is called otherDescription and should only
// be used in the case of the enum having the value of
// "Other".
type OtherChoice struct {
	// Description only to be used when Value is "Other"
	Description *string `xml:"otherDescription,attr" json:",omitempty"`
	// Value of the other choices
	Value string `xml:",chardata"`
}

// GetValue ...
func (oc *OtherChoice) GetValue() string {
	if oc.Value == "Other" {
		if oc.Description != nil {
			return *oc.Description
		}
		return ""
	}
	return oc.Value
}

// MarshalJSON ...
func (oc *OtherChoice) MarshalJSON() ([]byte, error) {
	return json.Marshal(oc.GetValue())
}

// SecureString ...
// Disclaimer is text that serves as a negation or limitation of the rights under
// a warranty given by a seller to a buyer.
type SecureString struct {
	// The NAR Information Security Guidelines class.
	// Possible values: Public, Confidential, Secret
	// Default: Confidential
	SecurityClass string `xml:"http://rets.org/xsd/RETSCommons isgSecurityClass,attr"`
	// Inherits from xs:string
	Value string `xml:",chardata"`
}

// Address for a property
type Address struct {
	// Indicates the preference order within all
	// the ContactMethods. The highest preference
	// is 0.
	PreferenceOrder int `xml:"http://rets.org/xsd/RETSCommons preference-order"`
	// Indicates the preference order within all
	// the Addresses. The highest preference is 0.
	AddressPreferenceOrder int `xml:"http://rets.org/xsd/RETSCommons address-preference-order"`
	// Provide a category for the address.
	//
	// Possible values:
	//
	// Display
	// Mailing
	// Shipping
	// Billing
	// Legal
	// Tax
	// Other
	//
	Category OtherChoice `xml:"http://rets.org/xsd/RETSCommons category"`
	// The FullStreetAddress is a text representation of
	// the address with the full civic location as a single entity.
	//
	// It may optionally include any of City, StateOrProvince,
	// PostalCode and Country.
	FullStreetAddress string `xml:"http://rets.org/xsd/RETSCommons FullStreetAddress"`
	// Civic Address Fields
	//
	// May only contain civic fields or BoxNumber element
	// Text field that uniquely locates a building on a given
	// street. House numbers may have fractional or alphabetic
	// modifiers. This is the first component in a street
	// address.
	StreetNumber *string `xml:"http://rets.org/xsd/RETSCommons StreetNumber,omitempty" json:",omitempty"`
	// Text field containing the direction that follows the
	// house number and precedes the street name in an address.
	// The format may be either an abbreviation, such as "NE" or
	// "N.E" or the full direction, "Northeast".
	StreetDirPrefix *string `xml:"http://rets.org/xsd/RETSCommons StreetDirPrefix,omitempty" json:",omitempty"`
	// Text field containing the name of the street in an
	// address.  This may follow the house number or, if
	// applicable, the street direction prefix. It precedes the
	// street suffix. for example, in the address
	// "123 Main St.", "Main" is the street name.
	StreetName *string `xml:"http://rets.org/xsd/RETSCommons StreetName,omitempty" json:",omitempty"`
	// Text field describing the street type in an address. This
	// field follows the street name and precedes the street
	// direction suffix. A street suffix may be formatted as
	// either an  abbreviation or full name. Examples include:
	// Road, Rd., Avenue, Ave., etc.
	StreetSuffix *string `xml:"http://rets.org/xsd/RETSCommons StreetSuffix,omitempty" json:",omitempty"`
	// Text field containing the direction that follows the
	// street suffix in an address. The format may be either an
	// abbreviation, such as "NE" or "N.E" or the full
	// direction, "Northeast".
	StreetDirSuffix *string `xml:"http://rets.org/xsd/RETSCommons StreetDirSuffix,omitempty" json:",omitempty"`
	// Any additional elements needed to form the address
	StreetAdditionalInfo *string `xml:"http://rets.org/xsd/RETSCommons StreetAdditionalInfo,omitempty" json:",omitempty"`
	// Use the BoxNumber element to contain address location information
	// not covered by the Civic Address.
	//
	// May only contain if without the civic fields above
	// A container at a central mailing location, where the
	// incoming mail of a person or legal entity is held until
	// picked up by the person or legal entity. Also known as
	// a post office box.
	BoxNumber *string `xml:"http://rets.org/xsd/RETSCommons BoxNumber,omitempty" json:",omitempty"`
	// Text field containing the number or portion of a larger
	// building or complex. Unit Number should appear following
	// the street suffix or, if it exists, the street suffix
	// direction, in the street address. Examples are:
	// "APT G", "55", etc.
	UnitNumber *string `xml:"http://rets.org/xsd/RETSCommons UnitNumber,omitempty" json:",omitempty"`
	// The city, township, municipality, etc. portion of the
	// physical, legal or mailing address for a property,
	// person, etc.
	City *string `xml:"http://rets.org/xsd/RETSCommons City,omitempty" json:",omitempty"`
	// Text field containing either the accepted postal
	// abbreviation or the full name for one of the 50 U.S.
	// states or 13 Canadian provinces/territories.
	State *string `xml:"http://rets.org/xsd/RETSCommons StateOrProvince,omitempty" json:",omitempty"`
	// In the United states, the postal code (ZIP code) the
	// basic postal code format consists of five numerical
	// digits and may include a five digit ZIP+4 code that
	// allows delivery of a piece of mail to be even more
	// accurately defined. In Canada, the postal code is a six
	// character alpha-numerical code defined and maintained by
	// Canada Post Corporation for mail processing
	// (sorting and delivery).
	PostalCode *string `xml:"http://rets.org/xsd/RETSCommons PostalCode,omitempty" json:",omitempty"`
	// The group of addresses to which the USPS assigns the
	// same code to aid in mail delivery. For the USPS, these
	// codes are 9 digits: 5 numbers for the ZIP Code, one
	// letter for the carrier route type, and 3 numbers for the
	// carrier route number.
	CarrierRoute *string `xml:"http://rets.org/xsd/RETSCommons CarrierRoute,omitempty" json:",omitempty"`
	// The territory of nation or state included in a person or
	// property's legal or mailing address.
	Country *string `xml:"http://rets.org/xsd/RETSCommons Country,omitempty" json:",omitempty"`
	// Indicates the level of privacy for this information. Creation of
	// this attribute inspired by the need to provide this level of
	// privacy information for contact nformation: phone, email,
	// address, and the like.
	//
	// Possible values:
	//
	// Public
	// Agent/Member
	// MLS
	//
	Privacy *string `xml:"http://rets.org/xsd/RETSCommons privacyType,omitempty" json:",omitempty"`
}

// PriceWithFrequency ...
// For attributes that start with a XML namespace, use the URL of the
// namespace to avoid compilation errors
//
// eg. xmlns:commons="http://rets.org/xsd/RETSCommons"
// commons:isgSecurityClass
// will become
// `xml:"http://rets.org/xsd/RETSCommons isgSecurityClass,attr"`
type PriceWithFrequency struct {
	// The currencyPeriod attribute indicates that the price is repeated
	// at the frequency indicated. The abscence of the attribute indicates
	// a one-time payment.
	//
	// Possible values:
	//
	// Daily
	// Week
	// Bi-Weekly
	// Month
	// Bi-Monthly
	// Quarterly
	// Semi-Annually
	// Annually
	// Seasonal
	//
	Frequency string `xml:"http://rets.org/xsd/RETSCommons currencyPeriod,attr,omitempty"`
	// The NAR Information Security Guidelines class.
	// Possible values: Public, Confidential, Secret
	// Default: Confidential
	SecurityClass string `xml:"http://rets.org/xsd/RETSCommons isgSecurityClass,attr,omitempty"`
	// The currency code.

	// This document uses the three character ASCII currency code
	// values defined in ISO 4217.
	// Default: USD
	Currency string `xml:"http://rets.org/xsd/RETSCommons currencyCode,attr,omitempty"`
	// Price inherits from a nullable decimal in XSD.
	Value string `xml:",chardata"`
}

// AlternatePrice is a secondary price type in a different currency
type AlternatePrice struct {
	// The list price for the property.
	// Where range pricing is used, the higher of the two prices.
	ListPrice PriceWithFrequency `xml:"AlternateListPrice"`
	// Where range pricing is used, the lower price of the range
	ListPriceLow *PriceWithFrequency `xml:"AlternateListPriceLow,omitempty" json:",omitempty"`
}

// SecureBoolean is a boolean element with an isgSecurityClass attribute.
type SecureBoolean struct {
	// The NAR Information Security Guidelines class.
	// Possible values: Public, Confidential, Secret
	// Default: Confidential
	SecurityClass string `xml:"http://rets.org/xsd/RETSCommons isgSecurityClass,attr"`
	// Boolean value
	Value bool `xml:",chardata"`
}

// SecureDateTime is a datetime element with an isgSecurityClass attribute.
type SecureDateTime struct {
	// The NAR Information Security Guidelines class.
	// Possible values: Public, Confidential, Secret
	// Default: Confidential
	SecurityClass string `xml:"http://rets.org/xsd/RETSCommons isgSecurityClass,attr"`
	// DateTime inherits from a nullable datetime in XSD
	Value string `xml:",chardata"`
}

// MarketingInformation contains items related
// to the contract between the
// selling agent and the owner.
// These indicators are used to determine the
// visibility of the listing on the internet,
// visibility of address on the internet, visibility
// of the photo on the internet, and whether the
// property has a sign. Additional elements
// may be others discovered in the future.
type MarketingInformation struct {
	// The seller agreed to permit the listing
	// to be marketed on the internet.
	PermitInternet *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons PermitInternet,omitempty" json:",omitempty"`
	// The seller agreed to permit the property address
	// to be displayed on the internet.
	PermitAddress *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons PermitAddressOnInternet,omitempty" json:",omitempty"`
	// The seller agreed to permit the
	// display of image(s) of the property
	// on the internet.
	PermitPicture *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons PermitPictureOnInternet,omitempty" json:",omitempty"`
	// The seller agreed to permit a for-sale
	// sign on the property and asserts the
	// right to provide that permission.
	// This may be constrained by local
	// rules or home-owner rules.
	PermitSign *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons PermitSignOnProperty,omitempty" json:",omitempty"`
	// A for-sale sign is on the property.
	HasSign *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons HasSignOnProperty,omitempty" json:",omitempty"`
	// Indicates whether the listing may be displayed.
	PermitVOW *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons VOWEntireListingDisplay,omitempty" json:",omitempty"`
	// Indicates whether the listing address may be displayed.
	PermitAddressVOW *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons VOWAddressDisplay,omitempty" json:",omitempty"`
	// Indicates whether the listing approximated valuation may be displayed.
	PermitValuationVOW *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons VOWAutomatedValuationDisplay,omitempty" json:",omitempty"`
	// Indicates whether consumer comments may be displayed.
	PermitCommentsVOW *SecureBoolean `xml:"http://rets.org/xsd/RETSCommons VOWConsumerComment,omitempty" json:",omitempty"`
}

// Media represents either a photo, video or other type of media file
type Media struct {
	// The last time that the media was modified.
	ModificationTimestamp *SecureDateTime `xml:"MediaModificationTimestamp,omitempty" json:",omitempty"`
	// Url of the media file. Required.
	URL string `xml:"MediaURL"`
	// The display caption for this photo. It is intended
	// to be a short title for the media.
	Caption *string `xml:"MediaCaption,omitempty" json:",omitempty"`
	// A narrative description of the media element.
	Description *string `xml:"MediaDescription,omitempty" json:",omitempty"`
	// Zeroth (0) element is the preferred photo by convention
	OrderNumber *int `xml:"MediaOrderNumber,omitempty" json:",omitempty"`
	// A placeholder for a classification of the media item
	// Examples: Elevation (exterior), Interior, Community,
	// View, Plan, Plat
	// The purpose is to allow the media items (photos)
	// to be grouped.
	// Optional because it is not widely available.
	Classification *string `xml:"MediaClassification,omitempty" json:",omitempty"`
}

// Area is a measurement in a two (of five) dimensions
type Area struct {
	// Area measurement units.
	//
	// Possible values:
	//
	// sqaureFoot
	// squareYard
	// acre
	// squareCentimeter
	// squareMeter
	// hectare
	// unknown
	//
	// Defaults to squareFoot.
	Units string `xml:"http://rets.org/xsd/RETSCommons areaUnits,attr"`
	// The source of the measurement.
	//
	// Possible values:
	//
	// Appraisal
	// Builder
	// Measured
	// Public Records
	// Unknown
	// Other
	//
	Source string `xml:"http://rets.org/xsd/RETSCommons measurementSource,attr"`
	// Area inherits from a nullable decimal in XSD.
	Value string `xml:",chardata"`
}

// ParticipantCode is an application defined coding for the participant
type ParticipantCode struct {
	// The identifier for the participant code
	ID string `xml:"ParticipantCodeId"`
	// A description of the code. This could identify the
	// system that the code comes from or any other
	// purpose to assist in describing the code identifer
	Description *string `xml:"ParticipantCodeDescription,omitempty" json:",omitempty"`
}

// ValueDescription is a common pair of elements
type ValueDescription struct {
	// Value is an optional string
	Value *string `xml:"Value,omitempty" json:",omitempty"`
	// Description is an optional string
	Description *string `xml:"Description,omitempty" json:",omitempty"`
}

// License is a professional license
type License struct {
	// The values that a license may take. These licenses
	// are typically issued by the State or Province where
	// the person is practising. See Licensing.xsd for more.
	//
	// Possible values:
	//
	// Real Estate Broker
	// Salesperson
	// Auctioneer
	// Certified General Appraiser
	// Certified Residential Appraiser
	// Licensed Appraiser
	// Registered Appraiser
	// Leasing Agent
	// Mortgage Broker
	// Apprentice Inspector
	// Licensed Inspector
	// Registered Inspector
	// Property Management
	// Rental Property Management
	// Unknown
	// Other
	Category OtherChoice `xml:"LicenseCategory"`
	// The License Number.
	LicenseNumber SecureString `xml:"LicenseNumber"`
	// A string representing the jurisdiction
	// that the license issuing body is located.
	Jurisdiction SecureString `xml:"Jurisdiction"`
	// The license issuing body State.
	State *string `xml:"StateOrProvince,omitempty" json:",omitempty"`
	// DateTime when the license began being valid
	Start *SecureDateTime `xml:"LicenseStartDateTime,omitempty" json:",omitempty"`
	// DateTime when the license expires
	Expiration *SecureDateTime `xml:"LicenseExpirationDateTime,omitempty" json:",omitempty"`
	// DateTime when the license was transferred
	Transfer *SecureDateTime `xml:"LicenseTransferDateTime,omitempty" json:",omitempty"`
}

// Participant represents a party in the deal.
// Different from Participants.xsd/ParticipantType
type Participant struct {
	// A unique identifier for the participant.
	Key string `xml:"ParticipantKey"`
	// The well-known identifier for the
	// participant. It is not necessarily the
	// key field for the participant.
	ID *string `xml:"ParticipantId,omitempty" json:",omitempty"`
	// List of participant codes
	Codes []ParticipantCode `xml:"ParticipantCode,omitempty" json:",omitempty"`
	// First Name
	FirstName *string `xml:"FirstName,omitempty" json:",omitempty"`
	// Last Name
	LastName *string `xml:"LastName,omitempty" json:",omitempty"`
	// Role
	Role *string `xml:"Role,omitempty" json:",omitempty"`
	// Primary phone number
	PrimaryPhone *string `xml:"PrimaryContactPhone,omitempty" json:",omitempty"`
	// Office Phone
	OfficePhone *string `xml:"OfficePhone,omitempty" json:",omitempty"`
	// Mobile Phone
	MobilePhone *string `xml:"MobilePhone,omitempty" json:",omitempty"`
	// Email
	Email *string `xml:"Email,omitempty" json:",omitempty"`
	// Fax
	Fax *string `xml:"Fax,omitempty" json:",omitempty"`
	// Website
	Website *string `xml:"WebsiteURL,omitempty" json:",omitempty"`
	// Photo
	Photo *string `xml:"PhotoURL,omitempty" json:",omitempty"`
	// Address
	Address *Address `xml:"Address,omitempty" json:",omitempty"`
	// A collection of licenses for the participant
	Licenses []License `xml:"Licenses>License,omitempty" json:",omitempty"`
	// Allow for value/description pair as an open bucket for extensibility.
	AdditionalInfo []ValueDescription `xml:"ParticipantAdditionalInformation,omitempty" json:",omitempty"`
}

// OfficeCode is an application defined encoding for the office.
type OfficeCode struct {
	// The identifier for the office.
	// It should permit unique identification of
	// office within the system
	ID string `xml:"OfficeCodeId"`
	// A description of the code. This could identify the
	// system that the code comes from or any other
	// purpose to assist in describing the code identifer
	Desription *string `xml:"OfficeCodeDescription,omitempty" json:",omitempty"`
}

// Office ...
type Office struct {
	// A unique identifier for the office.
	Key string `xml:"OfficeKey"`
	// The well-known identifier for the
	// office. It is not necessarily the
	// key field for the office.
	ID          string       `xml:"OfficeId"`
	Level       *string      `xml:"Level,omitempty" json:",omitempty"`
	OfficeCodes []OfficeCode `xml:"OfficeCode,omitempty" json:",omitempty"`
	// The name of the office.
	// For example, "Le Page Westside"
	Name          *string `xml:"Name,omitempty" json:",omitempty"`
	CorporateName *string `xml:"CorporateName,omitempty" json:",omitempty"`
	BrokerID      *string `xml:"BrokerId,omitempty" json:",omitempty"`
	// The identifier for the main office of any type of group of offices
	MainOfficeID   *string            `xml:"MainOfficeId,omitempty" json:",omitempty"`
	Phone          *string            `xml:"PhoneNumber,omitempty" json:",omitempty"`
	Fax            *string            `xml:"Fax,omitempty" json:",omitempty"`
	Addresss       *Address           `xml:"Address,omitempty" json:",omitempty"`
	Email          *string            `xml:"OfficeEmail,omitempty" json:",omitempty"`
	Website        *string            `xml:"Website,omitempty" json:",omitempty"`
	Logo           *string            `xml:"OfficeLogoURL,omitempty" json:",omitempty"`
	AdditionalInfo []ValueDescription `xml:"OfficeAdditionalInformation,omitempty" json:",omitempty"`
}

// Business ...
type Business struct {
	// Name of the business. Required.
	Name string `xml:"Name"`
	// An identifier from the providing system.
	ID *string `xml:"BusinessId,omitempty" json:",omitempty"`
	// Phone number as a string, no specific format enforced.
	Phone *string `xml:"Phone,omitempty" json:",omitempty"`
	// Fax number as a string, no specific format enforced.
	Fax *string `xml:"Fax,omitempty" json:",omitempty"`
	// Email address as a string no specific format enforced.
	Email *string `xml:"Email,omitempty" json:",omitempty"`
	// URL for the company website. Format is enforced.
	Website *string `xml:"WebsiteURL,omitempty" json:",omitempty"`
	// URL for the business logo. Format is enforced.
	Logo *string `xml:"LogoURL,omitempty" json:",omitempty"`
	// Address for the business.
	Address        *Address `xml:"Address,omitempty" json:",omitempty"`
	AdditionalInfo *string  `xml:"BusinessAdditionalInformation,omitempty" json:",omitempty"`
}

// Structure ... TODO: Implement structure from RETSCommon.xsd
type Structure struct {
}

// School in the given property's area.
type School struct {
	// The name of the school.
	Name *string `xml:"http://rets.org/xsd/RETSCommons Name,omitempty" json:",omitempty"`
	// The type of school in question.
	// Examples include Middle, Junior High,
	// etc.
	Category *OtherChoice `xml:"http://rets.org/xsd/RETSCommons SchoolCategory,omitempty" json:",omitempty"`
	// The district that a school is in.
	// A school may only belong to a single
	// district.
	District *SecureString `xml:"http://rets.org/xsd/RETSCommons District,omitempty" json:",omitempty"`
	// Further information about the school.
	Description *string `xml:"http://rets.org/xsd/RETSCommons Description,omitempty" json:",omitempty"`
}

// Community describes the area around the house, including schools
type Community struct {
	// Text field containing the name of a particular area of
	// land laid out and divided into lots, blocks, and building
	// sites, and in which public facilities are laid out, such as
	// streets, alleys, parks, and easements for public utilities.
	// Types of subdivisions include common interests (condominiums),
	// planned developments, time-share projects,
	Subdivision *SecureString `xml:"http://rets.org/xsd/RETSCommons Subdivision,omitempty" json:",omitempty"`
	// The collection of schools for a given property.
	Schools []School `xml:"http://rets.org/xsd/RETSCommons Schools>School,omitempty" json:",omitempty"`
	// The name of the development, neighborhood or
	// association in which the property is located.
	Name *SecureString `xml:"http://rets.org/xsd/RETSCommons CommunityName,omitempty" json:",omitempty"`
	// Text description of the common amenities offered to residents
	// in a common interest where a major percentage of the residents
	// in each household are 55 or older. May include items such as:
	// assisted living, senior center, etc.
	SeniorCommunity *SecureString `xml:"http://rets.org/xsd/RETSCommons SeniorCommunity,omitempty" json:",omitempty"`
	Structures      []Structure   `xml:"http://rets.org/xsd/RETSCommons ExistingStructures>ExistingStructure,omitempty" json:",omitempty"`
}

// Neighborhood has a name and optional description
type Neighborhood struct {
	Name        *string `xml:"Name,omitempty" json:",omitempty"`
	Description *string `xml:"Description,omitempty" json:",omitempty"`
}

// Location ...
type Location struct {
	Latitude           *string    `xml:"Latitude,omitempty" json:",omitempty"`
	Longitude          *string    `xml:"Longitude,omitempty" json:",omitempty"`
	Elevation          *string    `xml:"Elevation,omitempty" json:",omitempty"`
	MapCoordinate      *string    `xml:"MapCoordinate,omitempty" json:",omitempty"`
	Directions         *string    `xml:"Directions,omitempty" json:",omitempty"`
	GeocodeOptions     *string    `xml:"GeocodeOptions,omitempty" json:",omitempty"`
	County             *string    `xml:"County,omitempty" json:",omitempty"`
	StreetIntersection *string    `xml:"StreetIntersection,omitempty" json:",omitempty"`
	ParcelID           *string    `xml:"ParcelId,omitempty" json:",omitempty"`
	Community          *Community `xml:"Community,omitempty" json:",omitempty"`
	// A string of the amenites found in the community
	CommunityAmenities *string `xml:"CommunityAmenities,omitempty" json:",omitempty"`
	// The mailing address of the community.
	CommunityAddress *Address `xml:"CommunityAddress,omitempty" json:",omitempty"`
	// The total number of floors in the building.
	TotalNumFloors *int `xml:"TotalNumFloors,omitempty" json:",omitempty"`
	// The civic building code zoning type
	Zoning *string `xml:"Zoning,omitempty" json:",omitempty"`
	// A string of the amenites found in the building of the listing
	BuildingAmenities *string `xml:"BuildingAmenities,omitempty" json:",omitempty"`
	// The total number of units in the building of the listing
	BuildingUnitCount *int `xml:"BuildingUnitCount,omitempty" json:",omitempty"`
	// Neighborhood the property is located in
	Neighborhoods []Neighborhood `xml:"Neighborhoods>Neighborhood,omitempty" json:",omitempty"`
}

// OpenHouse ...
type OpenHouse struct {
	Date        string  `xml:"Date"`
	StartTime   *string `xml:"StartTime,omitempty" json:",omitempty"`
	EndTime     *string `xml:"EndTime,omitempty" json:",omitempty"`
	Description *string `xml:"Description,omitempty" json:",omitempty"`
	Appointment *bool   `xml:"AppointmentRequiredYN,omitempty" json:",omitempty"`
}

// Tax ...
type Tax struct {
	Year        *int   `xml:"Year,omitempty" json:",omitempty"`
	Amount      string `xml:"Amount"`
	Description string `xml:"TaxDescription"`
}

// Expense ...
type Expense struct {
	Category string              `xml:"ExpenseCategory,omitempty" json:",omitempty"`
	Value    *PriceWithFrequency `xml:"ExpenseValue,omitempty" json:",omitempty"`
}

// Characteristics ...
type Characteristics struct {
	//
	//
	// Possible values:
	//
	// Barbeque or Grill
	// Coffee System
	// Coffee System - Rough in
	// Cooktop
	// Cooktop - Electric
	// Cooktop - Electric 2 burner
	// Cooktop - Electric 6 burner
	// Cooktop - Gas
	// Cooktop - Gas 2 burner
	// Cooktop - Gas 5 burner
	// Cooktop - Gas 6 burner
	// Cooktop - Gas Custom
	// Cooktop - Induction
	// Cooktop - Induction 2 burner
	// Cooktop - Induction 6 burner
	// Dishwasher
	// Dishwasher - Drawer
	// Dishwasher - Two or more
	// Dryer
	// Dryer - Dual fuel
	// Dryer - Electric 110V
	// Dryer - Electric 220V
	// Dryer - Gas
	// Dryer - Gas rough in
	// Freezer
	// Freezer - Compact
	// Freezer - Upright
	// Garbage Disposer
	// Ice Maker
	// Microwave
	// Oven
	// Oven - Convection
	// Oven - Double
	// Oven - Double Electric
	// Oven - Double Gas
	// Oven - Gas
	// Oven - Gas 3 wide
	// Oven - Self-Cleaning
	// Oven - Steam
	// Oven - Twin
	// Oven - Twin Electric
	// Oven - Twin Gas
	// Oven - Twin Gas 3 wide
	// Oven - Twin Mixed
	// Range
	// Range - Built In
	// Range - Dual
	// Range - Dual 6 burner
	// Range - Dual 8 burner
	// Range - Dual 10 burner
	// Range - Electric
	// Range - Gas
	// Range - Gas 6 burner
	// Range - Gas 8 burner
	// Range - Gas 10 burner
	// Range - Induction
	// Range - Other
	// Rangetop - Electric
	// Rangetop - Electric 2 burner
	// Rangetop - Electric 6 burner
	// Rangetop - Gas
	// Rangetop - Gas 2 burner
	// Rangetop - Gas 4 burner compact
	// Rangetop - Gas 6 burner
	// Rangetop - Gas 8 burner
	// Rangetop - Gas 10 burner
	// Rangetop - Gas Custom
	// Rangetop - Induction
	// Rangetop - Induction 2 burner
	// Rangetop - Induction 6 burner
	// Refrigerator
	// Refrigerator - Bar
	// Refrigerator - Built-in
	// Refrigerator - Built-in With Plumbing
	// Refrigerator - Drawer
	// Refrigerator - Side by Side
	// Refrigerator - Undercounter
	// Refrigerator - Wine Storage
	// Refrigerator - With Plumbing
	// Trash Compactor
	// Vacuum System
	// Vacuum System - Rough in
	// Vent Hood
	// Vent Hood 6 burner
	// Vent Hood 8 burner
	// Vent Hood 10 burner
	// Warming Drawer
	// Washer
	// Washer - Front load
	// Washer - Steamer
	// Washer - Top load
	// Washer/Dryer Combo
	// Washer/Dryer Stack
	// Water - Filter
	// Water - Instant Hot
	// Water - Purifier
	// Water - Softener
	// None
	// Other
	//
	Appliances []OtherChoice `xml:"Appliances>Appliance,omitempty" json:",omitempty"`
	// Description of the architectural design of the property listed.
	//
	// Possibile values:
	//
	// A Frame
	// Art Deco
	// Bungalow
	// Cape Cod
	// Colonial
	// Contemporary
	// Conventional
	// Cottage
	// Craftsman
	// Creole
	// Dome
	// Dutch Colonial
	// English
	// Federal
	// French
	// French Provincial
	// Georgian
	// Gothic Revival
	// Greek Revival
	// High Rise
	// Historical
	// International
	// Italianate
	// Loft
	// Mansion
	// Mediterranean
	// Modern
	// Monterey
	// Mountain
	// National
	// Neoclassical
	// New Traditional
	// Prairie
	// Pueblo
	// Queen Anne
	// Rambler
	// Ranch
	// Regency
	// Rustic
	// Saltbox
	// Santa Fe
	// Second Empire
	// Shed
	// Shingle
	// Shotgun
	// Spanish
	// Spanish Eclectic
	// Split Level
	// Stick
	// Tudor
	// Victorian
	// Other
	//
	ArchitectureStyle *OtherChoice `xml:"ArchitectureStyle,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has an attic.
	HasAttic *bool `xml:"HasAttic,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has a barbecue area.
	HasBarbecueArea *bool `xml:"HasBarbecueArea,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has a basement.
	HasBasement *bool `xml:"HasBasement,omitempty" json:",omitempty"`
	// The number of units in the building listed.
	BuildingUnitCount *int `xml:"BuildingUnitCount,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed is ready for cable.
	IsCableReady *bool `xml:"IsCableReady,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has one or more ceiling fans.
	HasCeilingFan *bool `xml:"HasCeilingFan,omitempty" json:",omitempty"`
	// Number of the floor the listed condominium is located on.
	CondoFloorNum *int `xml:"CondoFloorNum,omitempty" json:",omitempty"`
	// Collection of all the types of cooling system the listed property has.
	//
	// Possible values:
	//
	// Attic Fan
	// Ceiling Fan(s)
	// Central A/C
	// Central Evaporative
	// Central Fan
	// Chilled Water
	// Dehumidifiers
	// Dessicant Cooler
	// Evaporative
	// Heat Pumps
	// Partial
	// Radiant Floor
	// Radiant Floor Ground Loop
	// Refrigerator/Evaporative
	// Solar A/C-Active
	// Solar A/C-Passive
	// Wall Unit(s) A/C
	// Wall Unit(s) Evaporative
	// Window Unit(s) A/C
	// Window Unit(s) Evaporative
	// Zoned A/C
	// Unknown
	// Other
	// None
	//
	CoolingSystems []OtherChoice `xml:"CoolingSystems>CoolingSystem,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has one or more decks.
	HasDeck *bool `xml:"HasDeck,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has disabled access ramps,
	// elevators, or the like.
	HasDisabledAccess *bool `xml:"HasDisabledAccess,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has a dock.
	HasDock *bool `xml:"HasDock,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has a doorman.
	HasDoorman *bool `xml:"HasDoorman,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has double pane windows.
	HasDoublePaneWindows *bool `xml:"HasDoublePaneWindows,omitempty" json:",omitempty"`
	// Indicates whether or not the property listed has an elevator.
	HasElevator *bool `xml:"HasElevator,omitempty" json:",omitempty"`
	// Collection of types of exterior covering or adornment on the home.
	//
	// Possible values:
	//
	// Adobe
	// Aluminum Siding
	// Asbestos
	// Asphalt
	// Block
	// Board and Batten
	// Brick
	// Brick Veneer
	// Brick and Wood
	// Cedar Siding
	// Comb
	// Composition
	// Composition Shingles
	// Concrete
	// Concrete Block
	// EIFS
	// Fiberglass
	// Glass
	// Hardboard
	// Log
	// Log Siding
	// Masonite
	// Masonry
	// Metal
	// Metal Siding
	// Poured Concrete
	// Shingles (Not Wood)
	// Stone
	// Stone Veneer
	// Stucco
	// Stucco - Synthetic
	// Tile
	// Tilt-up (Pre-Cast Concrete)
	// Vinyl Siding
	// Wood
	// Wood Shingle
	// Wood Siding
	// Unknown
	// Other
	// None
	//
	ExteriorTypes []OtherChoice `xml:"ExteriorTypes>ExteriorType,omitempty" json:",omitempty"`
	// Indicates whether or not the listed property has a fireplace.
	HasFireplace *bool `xml:"HasFireplace,omitempty" json:",omitempty"`
	// Collection of floor coverings.
	//
	// Possible values:
	//
	// Bamboo
	// Brick
	// Carpet
	// Carpet - Full
	// Carpet - Partial
	// Concrete
	// Concrete - Bare
	// Concrete - Painted
	// Cork
	// Drainage
	// Engineered Wood
	// Glass
	// Granite
	// Hardwood
	// Laminate
	// Linoleum
	// Load Restriction
	// Marble
	// Parquet Wood
	// Rough-in
	// Slate
	// Soft Wood
	// Solid Wood
	// Specialty
	// Specialty Concrete
	// Tile
	// Tile - Ceramic
	// Tile - Porcelain
	// Tile - Stone
	// Tile or Stone
	// Vinyl
	// Wood
	// Unknown
	// Other
	// None
	//
	FloorCoverings []OtherChoice `xml:"FloorCoverings>FloorCovering,omitempty" json:",omitempty"`
	// Indicates whether or not a garden is located on the listed property.
	HasGarden *bool `xml:"HasGarden,omitempty" json:",omitempty"`
	// Indicates whether the listed property has gated entry.
	HasGatedEntry *bool `xml:"HasGatedEntry,omitempty" json:",omitempty"`
	// Indicates whether or not the listed property has a greenhouse.
	HasGreenhouse *bool `xml:"HasGreenhouse,omitempty" json:",omitempty"`
	// All the types of heating in use.
	//
	// Possible values:
	//
	// Butane Gas
	// Coal
	// Electric
	// Geothermal
	// Kerosene
	// Natural Gas
	// Oil
	// Passive Heat Pump
	// Passive Solar
	// Pellet
	// Propane Gas
	// Solar
	// Solar Panel
	// Wood
	// Unknown
	// Other
	// None
	//
	HeatingFuels []OtherChoice `xml:"HeatingFuels>HeatingFuel,omitempty" json:",omitempty"`
	// Types of heating system.
	//
	// Possible values:
	//
	// Central Furnace
	// Electric Air Filter
	// Fireplace
	// Fireplace - Insert
	// Floor Furnace
	// Floor Wall
	// Forced Air
	// Geothermal
	// Gravity Air
	// Gravity Hot Water
	// Heat Pump
	// Hot Water
	// Hot Water Radiant Floor
	// Humidifier
	// Pellet Stove
	// Radiant
	// Radiant Ceiling
	// Radiant Floor
	// Radiator
	// Solar Active
	// Solar Passive
	// Solar Active and Passive
	// Space Heater
	// Steam
	// Stove
	// S-W Changeover
	// Wall Unit
	// Zoned
	// Unknown
	// Other
	// None
	HeatingSystems []OtherChoice `xml:"HeatingSystems>HeatingSystem,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more hot tubs or spas.
	HasHotTubSpa *bool `xml:"HasHotTubSpa,omitempty" json:",omitempty"`
	// Indicates whether the property has an intercom.
	Intercom *bool `xml:"Intercom,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more jetted bath tubs.
	HasJettedBathTub bool `xml:"HasJettedBathTub,omitempty" json:",omitempty"`
	// Indicates whether the property has a lawn.
	HasLawn *bool `xml:"HasLawn,omitempty" json:",omitempty"`
	// Legal description.
	LegalDescription *string `xml:"LegalDescription,omitempty" json:",omitempty"`
	// Indicates that the property has a secondary suite or Mother-in-law suite.
	HasMotherInLaw *bool `xml:"HasMotherInLaw,omitempty" json:",omitempty"`
	// Indicates whether the property is new construction.
	IsNewConstruction *bool `xml:"IsNewConstruction,omitempty" json:",omitempty"`
	// Indicates the number of floors for the property.
	NumFloors *string `xml:"NumFloors,omitempty" json:",omitempty"`
	// Indicates the number of parking spaces for the property.
	NumParkingSpaces *int `xml:"NumParkingSpaces,omitempty" json:",omitempty"`
	// Collection of all the types of parking available for the property.
	//
	// Possible values:
	//
	// Alley
	// Assigned
	// Boat
	// Built-in
	// Carport
	// Commercial
	// Covered
	// Driveway
	// Fee
	// Fenced
	// Garage
	// Garage - Attached
	// Garage - Detached
	// Gated
	// Golf Cart
	// Guest
	// Heated
	// Leased
	// Mechanics
	// Mixed
	// Meter
	// Off Alley
	// Offsite
	// Off Street
	// On Street
	// Open
	// Oversized
	// Owned
	// Parking Lot
	// Parking Structure
	// Paved or Surfaced
	// Pole
	// Porte-Cochere
	// Pull-through
	// Ramp
	// RV
	// Secured
	// Side Apron
	// Side by Side
	// Special Needs
	// Stacked
	// Tandem
	// Tuck-Under
	// Unassigned
	// Underground/Basement
	// Unimproved
	// Valet
	// Workshop
	// Zoned Permit
	// Unknown
	// Other
	// None
	//
	ParkingTypes []OtherChoice `xml:"ParkingTypes>ParkingType,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more patios.
	HasPatio *bool `xml:"HasPatio,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more ponds.
	HasPond *bool `xml:"HasPond,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more pools.
	HasPool *bool `xml:"HasPool,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more porches.
	HasPorch *bool `xml:"HasPorch,omitempty" json:",omitempty"`
	// Collection of the type of roofing materials at the property.
	//
	// Possible values:
	//
	// Aluminum
	// Asbestos
	// Asphalt
	// Built-up
	// Clay Tile
	// Composition Shingle
	// ?Concrete?
	// Concrete Tile
	// Copper
	// Corrugated Metal
	// Green
	// ?gypsum?
	// Masonite or Cement Shake
	// Membrane
	// Metal
	// Shingle (Not wood)
	// Slate
	// Solar Panel
	// Standing Seam Steel
	// Steel
	// Tar and Gravel
	// Thatched
	// Tile
	// Urethane
	// Wood Shake
	// Wood Shingle
	// Unknown
	// Other
	//
	RoofTypes []OtherChoice `xml:"RoofTypes>RoofType,omitempty" json:",omitempty"`
	// Indicates the number of rooms.
	RoomCount *int `xml:"RoomCount,omitempty" json:",omitempty"`
	// Collection of rooms in the property.
	//
	// Possible values:
	//
	// Atrium
	// Attic-Finished
	// Attic-Unfinished
	// Basement
	// Basement-Finished
	// Basement-Unfinished
	// Bedroom
	// Bonus Room
	// Breakfast Nook
	// Breakfast Room
	// Crafts Room
	// Den
	// Dining Room
	// Eat-In Kitchen
	// Efficiency
	// Enclosed Patio
	// Family Room
	// Florida Room
	// Formal Dining Room
	// Foyer
	// Full Bath
	// Game Room
	// Great Room
	// Guest House
	// Guest Room
	// Half Bath
	// In-Law Suite
	// Kitchen
	// Kitchenette
	// Laundry Closet
	// Laundry Room
	// Library
	// Living Room
	// Loft
	// Master Bathroom
	// Media Room
	// Mudroom
	// Music Room
	// Office
	// One-Quarter Bath
	// Patio
	// Photo Lab
	// Recreational Room
	// Sauna
	// Servant Quarters
	// Sitting Room
	// Solarium
	// Storage
	// Studio
	// Study
	// Sunroom
	// Theatre
	// Three-Quarter Bath
	// Utility
	// Walk-In Closet
	// Walk-In Pantry
	// Wok Kitchen
	// Workshop
	// Unknown
	// Other
	//
	Rooms []OtherChoice `xml:"Rooms>Room,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more
	// RV Parking spot or area.
	HasRVParking *bool `xml:"HasRVParking,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more saunas.
	HasSauna *bool `xml:"HasSauna,omitempty" json:",omitempty"`
	// Indicates whether the property has a security system.
	HasSecuritySystem *bool `xml:"HasSecuritySystem,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more skylights.
	HasSkylight *bool `xml:"HasSkylight,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more sports court.
	HasSportsCourt *bool `xml:"HasSportsCourt,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more sprinkler system.
	HasSprinklerSystem *bool `xml:"HasSprinklerSystem,omitempty" json:",omitempty"`
	// Indicates whether the property a vaulted ceiling.
	HasVaultedCeiling *bool `xml:"HasVaultedCeiling,omitempty" json:",omitempty"`
	// A collection of the various view types from the property.
	//
	// Possible values:
	//
	// Airport
	// Average
	// Bluff
	// Bridge
	// Canyon
	// City
	// Desert
	// Forest
	// Golf Course
	// Harbor
	// Hills
	// Lake
	// Marina
	// Mountain
	// None
	// Ocean
	// Panorama
	// Park
	// Ravine
	// River
	// Territorial
	// Valley
	// Vista
	// Unknown
	// Water
	// Other
	//
	ViewTypes []OtherChoice `xml:"ViewTypes>ViewType,omitempty" json:",omitempty"`
	// Indicates whether the property is on the waterfront: ocean or lake.
	IsWaterfront *bool `xml:"IsWaterfront,omitempty" json:",omitempty"`
	// Indicates whether the property has one or more wet bars.
	HasWetBar *bool `xml:"HasWetBar,omitempty" json:",omitempty"`
	// A description of what the owner loves about the property.
	WhatOwnerLoves *string `xml:"WhatOwnerLoves,omitempty" json:",omitempty"`
	// Is the property wired for high tech purposes:
	// home network, speaker system, etc.
	IsWired *bool `xml:"IsWired,omitempty" json:",omitempty"`
	// Indicates the year the property received updates.
	YearUpdated *int `xml:"YearUpdated,omitempty" json:",omitempty"`
	// Allow for value/description pair as an open bucket
	// for extensibility. This allows addition of features,
	// remarks, etc.
	AdditionalInfo []ValueDescription `xml:"AdditionInformation,omitempty" json:",omitempty"`
}

// Listing is a top level element for sharing a single listing.
type Listing struct {
	// The following are required fields in each listing

	// The address for the property.
	Address Address `xml:"Address"`
	// The listed price for the property.
	ListPrice PriceWithFrequency `xml:"ListPrice"`
	// When the listing uses range pricing, the lower of the two list
	// prices for the property. Not required.
	ListPriceLow *PriceWithFrequency `xml:"ListPriceLow,omitempty" json:",omitempty"`
	// A collection of alternate list prices for the property. Each listed
	// price for the property is provided in a different currency.
	// The currency type is in the currency code attribute.
	//
	// Wherever possible, the ISO 4217 standard currency code values
	// should be used.
	AlternatePrices []AlternatePrice `xml:"AlernatePrices,omitempty" json:",omitempty"`
	// The URL for the original listing.
	ListingURL string `xml:"ListingURL"`
	// Name of the listing provider
	ProviderName string `xml:"ProviderName"`
	// URL of the listing provider
	ProviderURL string `xml:"ProviderURL"`
	// The source of the listing information
	//
	// Possible values:
	//
	// Aggregator
	// Broker
	// Franchiser
	// HomeBuilder
	// Member
	// MLS
	// Owner
	// Publisher
	// Unknown
	// Other
	//
	ProviderCategory OtherChoice `xml:"ProviderCategory"`
	// Email address to use for lead management.
	LeadRoutingEmail string `xml:"LeadRoutingEmail"`
	// The total number of bedrooms. For listings with no home on the
	// property, such as a lot or land, this element should have the
	// value of 0. When no information is available on Bedrooms, the
	// value should be empty.
	Bedrooms int `xml:"Bedrooms"`
	// The total number of bathrooms. For listings with no home on the
	// property, such as a lot or land, this element should have the
	// value of 0. When no information is available on Bedrooms, the
	// value should be empty.
	Bathrooms int `xml:"Bathrooms"`
	// Primary type of the listed property.
	//
	// Possible values:
	//
	// Residential
	// Lots And Land
	// Farm And Agriculture
	// MultiFamily
	// Commercial
	// Common Interest
	// Rental
	// Other
	//
	PropertyType OtherChoice `xml:"PropertyType"`
	// Secondary type of the listed property.
	//
	// Possible values:
	//
	// Apartment
	// Boatslip
	// Cabin
	// Condominium
	// Deeded Parking
	// Duplex
	// Farm
	// Manufactured Home
	// Mobile Home
	// Own Your Own
	// Quadruplex
	// Single Family Attached
	// Single Family Detached
	// Stock Cooperative
	// Townhouse
	// Timeshare
	// Triplex
	// Other
	//
	PropertySubType OtherChoice `xml:"PropertySubType"`
	// The identifier for the listing generated by the feed
	ListingKey string `xml:"ListingKey"`

	// The following fields are recommended

	// The category of the listing.
	//
	// The values are:
	//
	// Purchase
	// Lease
	// Rent
	//
	ListingCategory *string `xml:"ListingCategory,omitempty" json:",omitempty"`
	// The status of the listing.
	//
	// The values are:
	//
	// Active
	// Cancelled
	// Closed
	// Expired
	// Pending
	// Withdrawn
	//
	// This element is optional. If absent, the listing is assumed to be "Active".
	ListingStatus *string `xml:"ListingStatus,omitempty" json:",omitempty"`
	// Extensible:  Information related to the visibility
	// of the listing on the internet, visibility
	// of certain items in the listing on the
	// internet, and the permission to place
	// signage on the property. Information
	// related to how the listing is allowed
	// to be marketed as a part of the
	// agreement between the seller and the
	// agent. VOW information is included in this common
	// type.
	MarketingInformation *MarketingInformation `xml:"MarketingInformation,omitempty" json:",omitempty"`
	// A collection of photos for the property.
	Photos []Media `xml:"Photos>Photo,omitempty" json:",omitempty"`
	// Indicates whether or not the address may be publicly disclosed.
	DiscloseAddress *bool `xml:"DiscloseAddress,omitempty" json:",omitempty"`
	// Indicates if the listing is a short sale.
	ShortSale *bool `xml:"ShortSale,omitempty" json:",omitempty"`
	// This is a longer description of the listing.
	Description *string `xml:"ListingDescription,omitempty" json:",omitempty"`
	// The MLSId is the identifier for the MLS.
	MLSId *string `xml:"MlsId,omitempty" json:",omitempty"`
	// The string name for the MLS.
	MLSName *string `xml:"MlsName,omitempty" json:",omitempty"`
	// The MlsNumber is the identifier for the property within the MLS.
	MLSNumber *string `xml:"MlsNumber,omitempty" json:",omitempty"`
	// Total livable square feet of the listed property.
	LivingArea *Area `xml:"LivingArea,omitempty" json:",omitempty"`
	// Size of the lot. Attributes contain the units and measurement source
	LotSize *Area `xml:"LotSize,omitempty" json:",omitempty"`
	// Year the property was constructed.
	YearBuilt *int `xml:"YearBuilt,omitempty" json:",omitempty"`
	// Date the property was listed.
	ListingDate *string `xml:"ListingDate,omitempty" json:",omitempty"`
	// Container for the miscellaneous details about the listed property.
	Characteristics Characteristics `xml:"DetailedCharacteristics"`
	// When it was last modified
	// Example: 2007-03-11T12:00:00-05:00
	ModificationTimestamp *SecureDateTime `xml:"ModificationTimestamp,omitempty" json:",omitempty"`

	// The following fields are optional

	// An item for tracking the listing. May be a URI, for example.
	TrackingItem *string `xml:"TrackingItem,omitempty" json:",omitempty"`
	// A short title for the listing. Examples may be "Lovely Cape Cod in Downtown".
	Title *string `xml:"ListingTitle,omitempty" json:",omitempty"`
	// A total count for  the full bathrooms on the property. Full bath
	// generally includes sink, toilet, shower and bath, but may have
	// other local definitions
	FullBathrooms *int `xml:"FullBathrooms,omitempty" json:",omitempty"`
	// A total count of the three-quarter bathrooms on the property.
	// Three-quarter bathrooms contain a sink, toilet and either a shower
	// (most common definition) or bath but not both.
	ThreeQuarterBathrooms *int `xml:"ThreeQuarterBathrooms,omitempty" json:",omitempty"`
	// A total count of the half bathrooms on the property.
	// Half bathrooms contain a sink and toilet.
	HalfBathrooms *int `xml:"HalfBathrooms,omitempty" json:",omitempty"`
	// A total count of the one-quarter bathrooms on the property.
	// One-quarter bathrooms contain a sink.
	OneQuarterBathrooms *int `xml:"OneQuarterBathrooms,omitempty" json:",omitempty"`
	// A total count for all the partial bathrooms on the property. Partial bath
	// includes one quarter, one half, three quarter baths.
	PartialBathrooms *int `xml:"PartialBathrooms,omitempty" json:",omitempty"`
	// If the property is in some  form of foreclosure, this field provides a
	// value for the type of foreclosure.
	//
	// Possible values:
	//
	// Notice of Default (Pre-Foreclosure)
	// Lis Pendens (Pre-Foreclosure)
	// Notice of Trustee Sale (Auction)
	// Notice of Foreclosure Sale (Auction)
	// REO - Bank Owned
	// Foreclosure - Other
	// Other
	//
	ForeclosureStatus *OtherChoice `xml:"ForeclosureStatus,omitempty" json:",omitempty"`
	// A collection of all the participants in a listing.
	Participants []Participant `xml:"ListingParticipants>Participant,omitempty" json:",omitempty"`
	// A collection of all the virtual tours for a listing.
	VirtualTours []Media `xml:"VirtualTours>VirtualTour,omitempty" json:",omitempty"`
	// A collection of all the videos for a listing.
	// Does not include virtual tours or photos.
	Videos []Media `xml:"Videos>Video,omitempty" json:",omitempty"`
	// The collection of offices associated with the listing.
	Offices []Office `xml:"Offices>Office,omitempty" json:",omitempty"`
	// The brokerage for the listing.
	Brokerage *Business `xml:"Brokerage,omitempty" json:",omitempty"`
	// The franchise for the listing.
	Franchise *Business `xml:"Franchise,omitempty" json:",omitempty"`
	// The builder of the listing.
	Builder *Business `xml:"Builder,omitempty" json:",omitempty"`
	// The property manager for the listing.
	PropertyManager *Business `xml:"PropertyManager,omitempty" json:",omitempty"`
	// The information about the geographic location of the property.
	Location *Location `xml:"Location,omitempty" json:",omitempty"`
	// A collection of all the open houses for a property.
	OpenHouses []OpenHouse `xml:"OpenHouses>OpenHouse,omitempty" json:",omitempty"`
	// A collection of all the taxes reported for a property.
	Taxes []Tax `xml:"Taxes>Tax,omitempty" json:",omitempty"`
	// A collection of all the additional fees for a property.
	//
	// Each element has an enumerated description and a
	// fee (price) with optional frequency (example:monthly)
	//
	// Possible enumeration values:
	//
	// Annual Operating Expenses
	// Boat Fee
	// Community/Master Home Owner Fee
	// Condo/Coop Fee
	// Club Fee
	// Dock Fee
	// Elevator Use Fee
	// Equestrian Fee
	// Front Foot Fee
	// Ground Maintenance Fee
	// Home Owner Assessments Fee
	// Home Owner Transfer Fee
	// Land Assessment Fee
	// Move in Fee
	// Pet Deposit
	// Pool/Spa Fee
	// Processing Fee
	// Refuse Fee
	// Repair Deductible
	// Security Deposit
	// Security Guard Fee
	// Security Gate Fee
	// Special Assessment Fee
	// Water/Sewer Hookup Fee
	// Tenant Pays
	// Owner Pays
	// Other

	// In the case of the Other annotation, the value (description)
	// is contained in the attribute of the element
	Expenses []Expense `xml:"Expenses>Expense,omitempty" json:",omitempty"`
	// The disclaimer string for a specific listing
	Disclaimer SecureString `xml:"Disclaimer"`
}
