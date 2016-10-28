package dictionary

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// lookupValue is ...
type lookupValue []string

func (lv *lookupValue) UnmarshalText(data []byte) (err error) {
	if string(data) == "" {
		lv = nil
		return
	}

	*lv = strings.Split(string(data), ",")
	return
}

// dateTime
type dateTime time.Time

func (dt *dateTime) UnmarshalText(data []byte) (err error) {
	if string(data) == "" {
		dt = nil
		return
	}

	if len(string(data)) == 10 {
		ddt, err := time.Parse("01/02/2006", string(data))
		if err != nil {
			return err
		}

		*dt = (dateTime)(ddt)
		return err
	}

	ddt, err := time.Parse("01/02/2006 15:04:05", string(data))
	if err != nil {
		return err
	}

	*dt = (dateTime)(ddt)
	return
}

func (dt *dateTime) MarshalJSON() (data []byte, err error) {
	if (time.Time)(*dt).IsZero() {
		return json.Marshal(nil)
	}

	return json.Marshal((time.Time)(*dt))
}

// boolYN
type boolYN bool

func (b *boolYN) UnmarshalText(data []byte) (err error) {
	if string(data) == "" {
		b = nil
		return
	}

	if string(data) == "Yes" {
		*b = true
		return
	}

	if string(data) == "No" {
		*b = false
		return
	}

	bb, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}

	*b = (boolYN)(bb)
	return
}

// AgentOffice is ...
type AgentOffice struct {
	BuyerAgent    *BuyerAgent    `xml:"BuyerAgent" json:",omitempty"`
	BuyerOffice   *BuyerOffice   `xml:"BuyerOffice" json:",omitempty"`
	CoBuyerAgent  *CoBuyerAgent  `xml:"CoBuyerAgent" json:",omitempty"`
	CoBuyerOffice *CoBuyerOffice `xml:"CoBuyerOffice" json:",omitempty"`
	CoListAgent   *CoListAgent   `xml:"CoListAgent" json:",omitempty"`
	CoListOffice  *CoListOffice  `xml:"CoListOffice" json:",omitempty"`
	ListAgent     *ListAgent     `xml:"ListAgent" json:",omitempty"`
	ListOffice    *ListOffice    `xml:"ListOffice" json:",omitempty"`
	Team          *Team          `xml:"Team" json:",omitempty"`
}

// Price is ...
type Price struct {
	ClosePrice        json.Number `xml:"ClosePrice" json:",omitempty"`
	ListPrice         json.Number `xml:"ListPrice" json:",omitempty"`
	ListPriceLow      json.Number `xml:"ListPriceLow" json:",omitempty"`
	OriginalListPrice json.Number `xml:"OriginalListPrice" json:",omitempty"`
	PreviousListPrice json.Number `xml:"PreviousListPrice" json:",omitempty"`
}

// Area is ...
type Area struct {
	MLSAreaMajor    string `xml:"MLSAreaMajor" json:",omitempty"`
	MLSAreaMinor    string `xml:"MLSAreaMinor" json:",omitempty"`
	SubdivisionName string `xml:"SubdivisionName" json:",omitempty"`
}

// CoBuyerOffice is ...
type CoBuyerOffice struct {
	CoBuyerOfficeAOR      string      `xml:"CoBuyerOfficeAOR" json:",omitempty"`
	CoBuyerOfficeEmail    string      `xml:"CoBuyerOfficeEmail" json:",omitempty"`
	CoBuyerOfficeFax      string      `xml:"CoBuyerOfficeFax" json:",omitempty"`
	CoBuyerOfficeKey      string      `xml:"CoBuyerOfficeKey" json:",omitempty"`
	CoBuyerOfficeMlsId    string      `xml:"CoBuyerOfficeMlsId" json:",omitempty"`
	CoBuyerOfficeName     string      `xml:"CoBuyerOfficeName" json:",omitempty"`
	CoBuyerOfficePhone    string      `xml:"CoBuyerOfficePhone" json:",omitempty"`
	CoBuyerOfficePhoneExt json.Number `xml:"CoBuyerOfficePhoneExt" json:",omitempty"`
	CoBuyerOfficeURL      string      `xml:"CoBuyerOfficeURL" json:",omitempty"`
}

// CoListAgent is ...
type CoListAgent struct {
	CoListAgentAOR               string      `xml:"CoListAgentAOR" json:",omitempty"`
	CoListAgentCellPhone         string      `xml:"CoListAgentCellPhone" json:",omitempty"`
	CoListAgentDesignation       lookupValue `xml:"CoListAgentDesignation" json:",omitempty"`
	CoListAgentDirectPhone       string      `xml:"CoListAgentDirectPhone" json:",omitempty"`
	CoListAgentEmail             string      `xml:"CoListAgentEmail" json:",omitempty"`
	CoListAgentFax               string      `xml:"CoListAgentFax" json:",omitempty"`
	CoListAgentFirstName         string      `xml:"CoListAgentFirstName" json:",omitempty"`
	CoListAgentFullName          string      `xml:"CoListAgentFullName" json:",omitempty"`
	CoListAgentHomePhone         string      `xml:"CoListAgentHomePhone" json:",omitempty"`
	CoListAgentKey               string      `xml:"CoListAgentKey" json:",omitempty"`
	CoListAgentLastName          string      `xml:"CoListAgentLastName" json:",omitempty"`
	CoListAgentMiddleName        string      `xml:"CoListAgentMiddleName" json:",omitempty"`
	CoListAgentMlsId             string      `xml:"CoListAgentMlsId" json:",omitempty"`
	CoListAgentNamePrefix        string      `xml:"CoListAgentNamePrefix" json:",omitempty"`
	CoListAgentNameSuffix        string      `xml:"CoListAgentNameSuffix" json:",omitempty"`
	CoListAgentOfficePhone       string      `xml:"CoListAgentOfficePhone" json:",omitempty"`
	CoListAgentOfficePhoneExt    json.Number `xml:"CoListAgentOfficePhoneExt" json:",omitempty"`
	CoListAgentPager             string      `xml:"CoListAgentPager" json:",omitempty"`
	CoListAgentPreferredPhone    string      `xml:"CoListAgentPreferredPhone" json:",omitempty"`
	CoListAgentPreferredPhoneExt json.Number `xml:"CoListAgentPreferredPhoneExt" json:",omitempty"`
	CoListAgentStateLicense      string      `xml:"CoListAgentStateLicense" json:",omitempty"`
	CoListAgentTollFreePhone     string      `xml:"CoListAgentTollFreePhone" json:",omitempty"`
	CoListAgentURL               string      `xml:"CoListAgentURL" json:",omitempty"`
	CoListAgentVoiceMail         string      `xml:"CoListAgentVoiceMail" json:",omitempty"`
	CoListAgentVoiceMailExt      json.Number `xml:"CoListAgentVoiceMailExt" json:",omitempty"`
}

// ListAgent is ...
type ListAgent struct {
	ListAgentAOR               string      `xml:"ListAgentAOR" json:",omitempty"`
	ListAgentCellPhone         string      `xml:"ListAgentCellPhone" json:",omitempty"`
	ListAgentDesignation       lookupValue `xml:"ListAgentDesignation" json:",omitempty"`
	ListAgentDirectPhone       string      `xml:"ListAgentDirectPhone" json:",omitempty"`
	ListAgentEmail             string      `xml:"ListAgentEmail" json:",omitempty"`
	ListAgentFax               string      `xml:"ListAgentFax" json:",omitempty"`
	ListAgentFirstName         string      `xml:"ListAgentFirstName" json:",omitempty"`
	ListAgentFullName          string      `xml:"ListAgentFullName" json:",omitempty"`
	ListAgentHomePhone         string      `xml:"ListAgentHomePhone" json:",omitempty"`
	ListAgentKey               string      `xml:"ListAgentKey" json:",omitempty"`
	ListAgentLastName          string      `xml:"ListAgentLastName" json:",omitempty"`
	ListAgentMiddleName        string      `xml:"ListAgentMiddleName" json:",omitempty"`
	ListAgentMlsId             string      `xml:"ListAgentMlsId" json:",omitempty"`
	ListAgentNamePrefix        string      `xml:"ListAgentNamePrefix" json:",omitempty"`
	ListAgentNameSuffix        string      `xml:"ListAgentNameSuffix" json:",omitempty"`
	ListAgentOfficePhone       string      `xml:"ListAgentOfficePhone" json:",omitempty"`
	ListAgentOfficePhoneExt    json.Number `xml:"ListAgentOfficePhoneExt" json:",omitempty"`
	ListAgentPager             string      `xml:"ListAgentPager" json:",omitempty"`
	ListAgentPreferredPhone    string      `xml:"ListAgentPreferredPhone" json:",omitempty"`
	ListAgentPreferredPhoneExt json.Number `xml:"ListAgentPreferredPhoneExt" json:",omitempty"`
	ListAgentStateLicense      string      `xml:"ListAgentStateLicense" json:",omitempty"`
	ListAgentTollFreePhone     string      `xml:"ListAgentTollFreePhone" json:",omitempty"`
	ListAgentURL               string      `xml:"ListAgentURL" json:",omitempty"`
	ListAgentVoiceMail         string      `xml:"ListAgentVoiceMail" json:",omitempty"`
	ListAgentVoiceMailExt      json.Number `xml:"ListAgentVoiceMailExt" json:",omitempty"`
}

// HOA is ...
type HOA struct {
	AssociationAmenities     lookupValue `xml:"AssociationAmenities" json:",omitempty"`
	AssociationFee           json.Number `xml:"AssociationFee" json:",omitempty"`
	AssociationFee2          json.Number `xml:"AssociationFee2" json:",omitempty"`
	AssociationFee2Frequency string      `xml:"AssociationFee2Frequency" json:",omitempty"`
	AssociationFeeFrequency  string      `xml:"AssociationFeeFrequency" json:",omitempty"`
	AssociationFeeIncludes   lookupValue `xml:"AssociationFeeIncludes" json:",omitempty"`
	AssociationName          string      `xml:"AssociationName" json:",omitempty"`
	AssociationName2         string      `xml:"AssociationName2" json:",omitempty"`
	AssociationPhone         string      `xml:"AssociationPhone" json:",omitempty"`
	AssociationPhone2        string      `xml:"AssociationPhone2" json:",omitempty"`
	AssociationYN            boolYN      `xml:"AssociationYN" json:",omitempty"`
}

// Room is ...
type Room struct {
	Area              json.Number `xml:"Area" json:",omitempty"`
	AreaSource        string      `xml:"AreaSource" json:",omitempty"`
	AreaUnits         string      `xml:"AreaUnits" json:",omitempty"`
	Description       string      `xml:"Description" json:",omitempty"`
	Dimensions        string      `xml:"Dimensions" json:",omitempty"`
	Features          lookupValue `xml:"Features" json:",omitempty"`
	Length            json.Number `xml:"Length" json:",omitempty"`
	LengthWidthSource string      `xml:"LengthWidthSource" json:",omitempty"`
	LengthWidthUnits  string      `xml:"LengthWidthUnits" json:",omitempty"`
	Level             string      `xml:"Level" json:",omitempty"`
	Width             json.Number `xml:"Width" json:",omitempty"`
	Type              lookupValue `xml:"Type" json:",omitempty"`
}

// BuyerOffice is ...
type BuyerOffice struct {
	BuyerOfficeAOR      string      `xml:"BuyerOfficeAOR" json:",omitempty"`
	BuyerOfficeEmail    string      `xml:"BuyerOfficeEmail" json:",omitempty"`
	BuyerOfficeFax      string      `xml:"BuyerOfficeFax" json:",omitempty"`
	BuyerOfficeKey      string      `xml:"BuyerOfficeKey" json:",omitempty"`
	BuyerOfficeMlsId    string      `xml:"BuyerOfficeMlsId" json:",omitempty"`
	BuyerOfficeName     string      `xml:"BuyerOfficeName" json:",omitempty"`
	BuyerOfficePhone    string      `xml:"BuyerOfficePhone" json:",omitempty"`
	BuyerOfficePhoneExt json.Number `xml:"BuyerOfficePhoneExt" json:",omitempty"`
	BuyerOfficeURL      string      `xml:"BuyerOfficeURL" json:",omitempty"`
}

// Marketing is ...
type Marketing struct {
	IDXAddressDisplayYN            boolYN `xml:"IDXAddressDisplayYN" json:",omitempty"`
	IDXAutomatedValuationDisplayYN boolYN `xml:"IDXAutomatedValuationDisplayYN" json:",omitempty"`
	IDXConsumerCommentYN           boolYN `xml:"IDXConsumerCommentYN" json:",omitempty"`
	IDXEntireListingDisplayYN      boolYN `xml:"IDXEntireListingDisplayYN" json:",omitempty"`
	SignOnPropertyYN               boolYN `xml:"SignOnPropertyYN" json:",omitempty"`
	VirtualTourURLBranded          string `xml:"VirtualTourURLBranded" json:",omitempty"`
	VirtualTourURLUnbranded        string `xml:"VirtualTourURLUnbranded" json:",omitempty"`
	VOWAddressDisplayYN            boolYN `xml:"VOWAddressDisplayYN" json:",omitempty"`
	VOWAutomatedValuationDisplayYN boolYN `xml:"VOWAutomatedValuationDisplayYN" json:",omitempty"`
	VOWConsumerCommentYN           boolYN `xml:"VOWConsumerCommentYN" json:",omitempty"`
	VOWEntireListingDisplayYN      boolYN `xml:"VOWEntireListingDisplayYN" json:",omitempty"`
}

// Characteristics is ...
type Characteristics struct {
	AnchorsCoTenants         string      `xml:"AnchorsCoTenants" json:",omitempty"`
	CommunityFeatures        lookupValue `xml:"CommunityFeatures" json:",omitempty"`
	CurrentUse               lookupValue `xml:"CurrentUse" json:",omitempty"`
	DevelopmentStatus        lookupValue `xml:"DevelopmentStatus" json:",omitempty"`
	Fencing                  lookupValue `xml:"Fencing" json:",omitempty"`
	FrontageLength           string      `xml:"FrontageLength" json:",omitempty"`
	FrontageType             lookupValue `xml:"FrontageType" json:",omitempty"`
	HorseAmenities           lookupValue `xml:"HorseAmenities" json:",omitempty"`
	HorseYN                  boolYN      `xml:"HorseYN" json:",omitempty"`
	LandLeaseAmount          json.Number `xml:"LandLeaseAmount" json:",omitempty"`
	LandLeaseAmountFrequency string      `xml:"LandLeaseAmountFrequency" json:",omitempty"`
	LandLeaseExpirationDate  *dateTime   `xml:"LandLeaseExpirationDate" json:",omitempty"`
	LandLeaseYN              boolYN      `xml:"LandLeaseYN" json:",omitempty"`
	LaundryFeatures          lookupValue `xml:"LaundryFeatures" json:",omitempty"`
	LeaseTerm                string      `xml:"LeaseTerm" json:",omitempty"`
	LotDimensionsSource      string      `xml:"LotDimensionsSource" json:",omitempty"`
	LotFeatures              lookupValue `xml:"LotFeatures" json:",omitempty"`
	LotSizeAcres             json.Number `xml:"LotSizeAcres" json:",omitempty"`
	LotSizeArea              json.Number `xml:"LotSizeArea" json:",omitempty"`
	LotSizeDimensions        string      `xml:"LotSizeDimensions" json:",omitempty"`
	LotSizeSource            string      `xml:"LotSizeSource" json:",omitempty"`
	LotSizeSquareFeet        json.Number `xml:"LotSizeSquareFeet" json:",omitempty"`
	LotSizeUnits             string      `xml:"LotSizeUnits" json:",omitempty"`
	MobileHomeRemainsYN      boolYN      `xml:"MobileHomeRemainsYN" json:",omitempty"`
	NumberOfBuildings        json.Number `xml:"NumberOfBuildings" json:",omitempty"`
	NumberofLots             json.Number `xml:"NumberofLots" json:",omitempty"`
	NumberOfPads             json.Number `xml:"NumberOfPads" json:",omitempty"`
	NumberOfUnitsTotal       json.Number `xml:"NumberOfUnitsTotal" json:",omitempty"`
	ParkManagerName          string      `xml:"ParkManagerName" json:",omitempty"`
	ParkManagerPhone         string      `xml:"ParkManagerPhone" json:",omitempty"`
	ParkName                 string      `xml:"ParkName" json:",omitempty"`
	PoolFeatures             lookupValue `xml:"PoolFeatures" json:",omitempty"`
	PoolPrivateYN            boolYN      `xml:"PoolPrivateYN" json:",omitempty"`
	PossibleUse              lookupValue `xml:"PossibleUse" json:",omitempty"`
	RoadFrontageType         lookupValue `xml:"RoadFrontageType" json:",omitempty"`
	RoadResponsibility       lookupValue `xml:"RoadResponsibility" json:",omitempty"`
	RoadSurfaceType          lookupValue `xml:"RoadSurfaceType" json:",omitempty"`
	SeniorCommunityYN        boolYN      `xml:"SeniorCommunityYN" json:",omitempty"`
	SpaFeatures              lookupValue `xml:"SpaFeatures" json:",omitempty"`
	SpaYN                    boolYN      `xml:"SpaYN" json:",omitempty"`
	Topography               string      `xml:"Topography" json:",omitempty"`
	View                     lookupValue `xml:"View" json:",omitempty"`
	ViewYN                   boolYN      `xml:"ViewYN" json:",omitempty"`
	WaterBodyName            string      `xml:"WaterBodyName" json:",omitempty"`
	WaterfrontYN             boolYN      `xml:"WaterfrontYN" json:",omitempty"`
}

// Closing is ...
type Closing struct {
	BuyerFinancing lookupValue `xml:"BuyerFinancing" json:",omitempty"`
	ClosingTerms   lookupValue `xml:"ClosingTerms" json:",omitempty"`
	Contingency    string      `xml:"Contingency" json:",omitempty"`
	Possession     lookupValue `xml:"Possession" json:",omitempty"`
}

// Address is ...
type Address struct {
	CarrierRoute         string `xml:"CarrierRoute" json:",omitempty"`
	City                 string `xml:"City" json:",omitempty"`
	Country              string `xml:"Country" json:",omitempty"`
	CountyOrParish       string `xml:"CountyOrParish" json:",omitempty"`
	PostalCity           string `xml:"PostalCity" json:",omitempty"`
	PostalCode           string `xml:"PostalCode" json:",omitempty"`
	PostalCodePlus4      string `xml:"PostalCodePlus4" json:",omitempty"`
	StateOrProvince      string `xml:"StateOrProvince" json:",omitempty"`
	StreetAdditionalInfo string `xml:"StreetAdditionalInfo" json:",omitempty"`
	StreetDirPrefix      string `xml:"StreetDirPrefix" json:",omitempty"`
	StreetDirSuffix      string `xml:"StreetDirSuffix" json:",omitempty"`
	StreetName           string `xml:"StreetName" json:",omitempty"`
	StreetNumber         string `xml:"StreetNumberNumeric" json:",omitempty"`
	StreetSuffix         string `xml:"StreetSuffix" json:",omitempty"`
	StreetSuffixModifier string `xml:"StreetSuffixModifier" json:",omitempty"`
	Township             string `xml:"Township" json:",omitempty"`
	UnitNumber           string `xml:"UnitNumber" json:",omitempty"`
	UnparsedAddress      string `xml:"UnparsedAddress" json:",omitempty"`
}

// Contract is ...
type Contract struct {
	Disclosures              lookupValue `xml:"Disclosures" json:",omitempty"`
	Exclusions               string      `xml:"Exclusions" json:",omitempty"`
	Inclusions               string      `xml:"Inclusions" json:",omitempty"`
	ListingFinancing         lookupValue `xml:"ListingFinancing" json:",omitempty"`
	ListingTerms             lookupValue `xml:"ListingTerms" json:",omitempty"`
	Ownership                string      `xml:"Ownership" json:",omitempty"`
	SpecialListingConditions lookupValue `xml:"SpecialListingConditions" json:",omitempty"`
}

// Showing is ...
type Showing struct {
	AccessCode                   string      `xml:"AccessCode" json:",omitempty"`
	LockBoxLocation              string      `xml:"LockBoxLocation" json:",omitempty"`
	LockBoxSerialNumber          string      `xml:"LockBoxSerialNumber" json:",omitempty"`
	LockBoxType                  string      `xml:"LockBoxType" json:",omitempty"`
	ShowingContactPhoneExtension json.Number `xml:"ShowingContactPhoneExtension" json:",omitempty"`
	ShowingContactPhoneNumber    string      `xml:"ShowingContactPhoneNumber" json:",omitempty"`
	ShowingInstructions          string      `xml:"ShowingInstructions" json:",omitempty"`
}

// GIS is ...
type GIS struct {
	CrossStreet         string      `xml:"CrossStreet" json:",omitempty"`
	Directions          string      `xml:"Directions" json:",omitempty"`
	Elevation           json.Number `xml:"Elevation" json:",omitempty"`
	ElevationUnits      string      `xml:"ElevationUnits" json:",omitempty"`
	Latitude            json.Number `xml:"Latitude" json:",omitempty"`
	Longitude           json.Number `xml:"Longitude" json:",omitempty"`
	MapCoordinate       string      `xml:"MapCoordinate" json:",omitempty"`
	MapCoordinateSource string      `xml:"MapCoordinateSource" json:",omitempty"`
	MapURL              string      `xml:"MapURL" json:",omitempty"`
}

// Location is ...
type Location struct {
	Address                  *Address `xml:"Address" json:",omitempty"`
	Area                     *Area    `xml:"Area" json:",omitempty"`
	DistanceFromSchoolBus    string   `xml:"DistanceFromSchoolBus" json:",omitempty"`
	DistanceFromShopping     string   `xml:"DistanceFromShopping" json:",omitempty"`
	DistanceToBus            string   `xml:"DistanceToBus" json:",omitempty"`
	DistanceToFreeway        string   `xml:"DistanceToFreeway" json:",omitempty"`
	DistanceToPlaceofWorship string   `xml:"DistanceToPlaceofWorship" json:",omitempty"`
	DistanceToSchools        string   `xml:"DistanceToSchools" json:",omitempty"`
	GIS                      *GIS     `xml:"GIS" json:",omitempty"`
	School                   *School  `xml:"School" json:",omitempty"`
}

// OccupantOwner is ...
type OccupantOwner struct {
	OccupantName  string `xml:"OccupantName" json:",omitempty"`
	OccupantPhone string `xml:"OccupantPhone" json:",omitempty"`
	OccupantType  string `xml:"OccupantType" json:",omitempty"`
	OwnerName     string `xml:"OwnerName" json:",omitempty"`
	OwnerPhone    string `xml:"OwnerPhone" json:",omitempty"`
}

// Business is ...
type Business struct {
	BusinessType              string      `xml:"BusinessType" json:",omitempty"`
	HoursDaysofOperation      string      `xml:"HoursDaysofOperation" json:",omitempty"`
	LaborInformation          string      `xml:"LaborInformation" json:",omitempty"`
	LeaseAmount               json.Number `xml:"LeaseAmount" json:",omitempty"`
	LeaseAmountFrequency      string      `xml:"LeaseAmountFrequency" json:",omitempty"`
	LeaseAssignableYN         boolYN      `xml:"LeaseAssignableYN" json:",omitempty"`
	LeaseExpiration           *dateTime   `xml:"LeaseExpiration" json:",omitempty"`
	LeaseRenewalOptionYN      boolYN      `xml:"LeaseRenewalOptionYN" json:",omitempty"`
	NumberOfFullTimeEmployees json.Number `xml:"NumberOfFullTimeEmployees" json:",omitempty"`
	NumberOfPartTimeEmployees json.Number `xml:"NumberOfPartTimeEmployees" json:",omitempty"`
	OwnershipType             string      `xml:"OwnershipType" json:",omitempty"`
	SeatingCapacity           json.Number `xml:"SeatingCapacity" json:",omitempty"`
	SpecialLicenses           string      `xml:"SpecialLicenses" json:",omitempty"`
	YearEstablished           json.Number `xml:"YearEstablished" json:",omitempty"`
	YearsCurrentOwner         json.Number `xml:"YearsCurrentOwner" json:",omitempty"`
}

// ListOffice is ...
type ListOffice struct {
	ListOfficeAOR      string      `xml:"ListOfficeAOR" json:",omitempty"`
	ListOfficeEmail    string      `xml:"ListOfficeEmail" json:",omitempty"`
	ListOfficeFax      string      `xml:"ListOfficeFax" json:",omitempty"`
	ListOfficeKey      string      `xml:"ListOfficeKey" json:",omitempty"`
	ListOfficeMlsId    string      `xml:"ListOfficeMlsId" json:",omitempty"`
	ListOfficeName     string      `xml:"ListOfficeName" json:",omitempty"`
	ListOfficePhone    string      `xml:"ListOfficePhone" json:",omitempty"`
	ListOfficePhoneExt json.Number `xml:"ListOfficePhoneExt" json:",omitempty"`
	ListOfficeURL      string      `xml:"ListOfficeURL" json:",omitempty"`
}

// Structure is ...
type Structure struct {
	AboveGradeFinishedArea       json.Number  `xml:"AboveGradeFinishedArea" json:",omitempty"`
	AboveGradeFinishedAreaSource string       `xml:"AboveGradeFinishedAreaSource" json:",omitempty"`
	AboveGradeFinishedAreaUnits  string       `xml:"AboveGradeFinishedAreaUnits" json:",omitempty"`
	AccessibilityFeatures        lookupValue  `xml:"AccessibilityFeatures" json:",omitempty"`
	ArchitecturalStyle           lookupValue  `xml:"ArchitecturalStyle" json:",omitempty"`
	AttachedGarageYN             boolYN       `xml:"AttachedGarageYN" json:",omitempty"`
	Basement                     lookupValue  `xml:"Basement" json:",omitempty"`
	BathroomsFull                json.Number  `xml:"BathroomsFull" json:",omitempty"`
	BathroomsHalf                json.Number  `xml:"BathroomsHalf" json:",omitempty"`
	BathroomsOneQuarter          json.Number  `xml:"BathroomsOneQuarter" json:",omitempty"`
	BathroomsThreeQuarter        json.Number  `xml:"BathroomsThreeQuarter" json:",omitempty"`
	BathroomsTotal               string       `xml:"BathroomsTotal" json:",omitempty"`
	BedroomsPossible             json.Number  `xml:"BedroomsPossible" json:",omitempty"`
	BedroomsTotal                json.Number  `xml:"BedroomsTotal" json:",omitempty"`
	BelowGradeFinishedArea       json.Number  `xml:"BelowGradeFinishedArea" json:",omitempty"`
	BelowGradeFinishedAreaSource string       `xml:"BelowGradeFinishedAreaSource" json:",omitempty"`
	BelowGradeFinishedAreaUnits  string       `xml:"BelowGradeFinishedAreaUnits" json:",omitempty"`
	BodyType                     lookupValue  `xml:"BodyType" json:",omitempty"`
	BuilderModel                 string       `xml:"BuilderModel" json:",omitempty"`
	BuilderName                  string       `xml:"BuilderName" json:",omitempty"`
	BuildingAreaSource           string       `xml:"BuildingAreaSource" json:",omitempty"`
	BuildingAreaTotal            json.Number  `xml:"BuildingAreaTotal" json:",omitempty"`
	BuildingAreaUnits            string       `xml:"BuildingAreaUnits" json:",omitempty"`
	BuildingFeatures             lookupValue  `xml:"BuildingFeatures" json:",omitempty"`
	BuildingName                 string       `xml:"BuildingName" json:",omitempty"`
	CarportSpaces                json.Number  `xml:"CarportSpaces" json:",omitempty"`
	CarportYN                    boolYN       `xml:"CarportYN" json:",omitempty"`
	CommonWalls                  lookupValue  `xml:"CommonWalls" json:",omitempty"`
	ConstructionMaterials        lookupValue  `xml:"ConstructionMaterials" json:",omitempty"`
	Cooling                      lookupValue  `xml:"Cooling" json:",omitempty"`
	CoolingYN                    boolYN       `xml:"CoolingYN" json:",omitempty"`
	CoveredSpaces                json.Number  `xml:"CoveredSpaces" json:",omitempty"`
	DirectionFaces               string       `xml:"DirectionFaces" json:",omitempty"`
	DOH1                         string       `xml:"DOH1" json:",omitempty"`
	DOH2                         string       `xml:"DOH2" json:",omitempty"`
	DOH3                         string       `xml:"DOH3" json:",omitempty"`
	DoorFeatures                 lookupValue  `xml:"DoorFeatures" json:",omitempty"`
	EntryLevel                   json.Number  `xml:"EntryLevel" json:",omitempty"`
	EntryLocation                string       `xml:"EntryLocation" json:",omitempty"`
	ExteriorFeatures             lookupValue  `xml:"ExteriorFeatures" json:",omitempty"`
	FireplaceFeatures            lookupValue  `xml:"FireplaceFeatures" json:",omitempty"`
	FireplacesTotal              json.Number  `xml:"FireplacesTotal" json:",omitempty"`
	FireplaceYN                  boolYN       `xml:"FireplaceYN" json:",omitempty"`
	Flooring                     lookupValue  `xml:"Flooring" json:",omitempty"`
	FoundationArea               json.Number  `xml:"FoundationArea" json:",omitempty"`
	FoundationDetails            lookupValue  `xml:"FoundationDetails" json:",omitempty"`
	GarageSpaces                 json.Number  `xml:"GarageSpaces" json:",omitempty"`
	GarageYN                     boolYN       `xml:"GarageYN" json:",omitempty"`
	HabitableResidenceYN         boolYN       `xml:"HabitableResidenceYN" json:",omitempty"`
	Heating                      lookupValue  `xml:"Heating" json:",omitempty"`
	HeatingYN                    boolYN       `xml:"HeatingYN" json:",omitempty"`
	InteriorFeatures             lookupValue  `xml:"InteriorFeatures" json:",omitempty"`
	Levels                       lookupValue  `xml:"Levels" json:",omitempty"`
	License1                     string       `xml:"License1" json:",omitempty"`
	License2                     string       `xml:"License2" json:",omitempty"`
	License3                     string       `xml:"License3" json:",omitempty"`
	LivingArea                   json.Number  `xml:"LivingArea" json:",omitempty"`
	LivingAreaSource             string       `xml:"LivingAreaSource" json:",omitempty"`
	LivingAreaUnits              string       `xml:"LivingAreaUnits" json:",omitempty"`
	Make                         string       `xml:"Make" json:",omitempty"`
	MobileDimUnits               string       `xml:"MobileDimUnits" json:",omitempty"`
	MobileLength                 json.Number  `xml:"MobileLength" json:",omitempty"`
	MobileWidth                  json.Number  `xml:"MobileWidth" json:",omitempty"`
	Model                        string       `xml:"Model" json:",omitempty"`
	NewConstructionYN            boolYN       `xml:"NewConstructionYN" json:",omitempty"`
	OpenParkingSpaces            json.Number  `xml:"OpenParkingSpaces" json:",omitempty"`
	OpenParkingYN                boolYN       `xml:"OpenParkingYN" json:",omitempty"`
	OtherParking                 string       `xml:"OtherParking" json:",omitempty"`
	OtherStructures              lookupValue  `xml:"OtherStructures" json:",omitempty"`
	ParkingFeatures              lookupValue  `xml:"ParkingFeatures" json:",omitempty"`
	ParkingTotal                 json.Number  `xml:"ParkingTotal" json:",omitempty"`
	PatioAndPorchFeatures        lookupValue  `xml:"PatioAndPorchFeatures" json:",omitempty"`
	Performance                  *Performance `xml:"Performance" json:",omitempty"`
	PropertyAttachedYN           boolYN       `xml:"PropertyAttachedYN" json:",omitempty"`
	Roof                         lookupValue  `xml:"Roof" json:",omitempty"`
	Rooms                        *Rooms       `xml:"Rooms" json:",omitempty"`
	RVParkingDimensions          string       `xml:"RVParkingDimensions" json:",omitempty"`
	SerialU                      string       `xml:"SerialU" json:",omitempty"`
	SerialX                      string       `xml:"SerialX" json:",omitempty"`
	SerialXX                     string       `xml:"SerialXX" json:",omitempty"`
	Skirt                        lookupValue  `xml:"Skirt" json:",omitempty"`
	Stories                      string       `xml:"Stories" json:",omitempty"`
	StoriesTotal                 string       `xml:"StoriesTotal" json:",omitempty"`
	WindowFeatures               lookupValue  `xml:"WindowFeatures" json:",omitempty"`
	YearBuilt                    json.Number  `xml:"YearBuilt" json:",omitempty"`
	YearBuiltDetails             string       `xml:"YearBuiltDetails" json:",omitempty"`
	YearBuiltEffective           json.Number  `xml:"YearBuiltEffective" json:",omitempty"`
	YearBuiltSource              string       `xml:"YearBuiltSource" json:",omitempty"`
}

// Tax is ...
type Tax struct {
	AdditionalParcelsDescription      string      `xml:"AdditionalParcelsDescription" json:",omitempty"`
	AdditionalParcelsYN               boolYN      `xml:"AdditionalParcelsYN" json:",omitempty"`
	ParcelNumber                      string      `xml:"ParcelNumber" json:",omitempty"`
	PublicSurveyRange                 string      `xml:"PublicSurveyRange" json:",omitempty"`
	PublicSurveySection               string      `xml:"PublicSurveySection" json:",omitempty"`
	PublicSurveyTownship              string      `xml:"PublicSurveyTownship" json:",omitempty"`
	TaxAmountFrequency                string      `xml:"TaxAmountFrequency" json:",omitempty"`
	TaxAnnualAmount                   json.Number `xml:"TaxAnnualAmount" json:",omitempty"`
	TaxAssessedValue                  json.Number `xml:"TaxAssessedValue" json:",omitempty"`
	TaxBlock                          string      `xml:"TaxBlock" json:",omitempty"`
	TaxBookNumber                     string      `xml:"TaxBookNumber" json:",omitempty"`
	TaxExemptions                     lookupValue `xml:"TaxExemptions" json:",omitempty"`
	TaxLegalDescription               string      `xml:"TaxLegalDescription" json:",omitempty"`
	TaxLot                            string      `xml:"TaxLot" json:",omitempty"`
	TaxMapNumber                      string      `xml:"TaxMapNumber" json:",omitempty"`
	TaxOtherAnnualAssessmentAmount    json.Number `xml:"TaxOtherAnnualAssessmentAmount" json:",omitempty"`
	TaxOtherAssessmentAmountFrequency string      `xml:"TaxOtherAssessmentAmountFrequency" json:",omitempty"`
	TaxParcelLetter                   string      `xml:"TaxParcelLetter" json:",omitempty"`
	TaxStatusCurrent                  lookupValue `xml:"TaxStatusCurrent" json:",omitempty"`
	TaxTract                          string      `xml:"TaxTract" json:",omitempty"`
	TaxYear                           json.Number `xml:"TaxYear" json:",omitempty"`
	Zoning                            string      `xml:"Zoning" json:",omitempty"`
	ZoningDescription                 string      `xml:"ZoningDescription" json:",omitempty"`
}

// Utilities is ...
type Utilities struct {
	DistanceToElectric             string      `xml:"DistanceToElectric" json:",omitempty"`
	DistanceToGas                  string      `xml:"DistanceToGas" json:",omitempty"`
	DistanceToPhoneService         string      `xml:"DistanceToPhoneService" json:",omitempty"`
	DistanceToSewer                string      `xml:"DistanceToSewer" json:",omitempty"`
	DistanceToStreet               string      `xml:"DistanceToStreet" json:",omitempty"`
	DistanceToWater                string      `xml:"DistanceToWater" json:",omitempty"`
	Electric                       lookupValue `xml:"Electric" json:",omitempty"`
	ElectricOnPropertyYN           boolYN      `xml:"ElectricOnPropertyYN" json:",omitempty"`
	Gas                            string      `xml:"Gas" json:",omitempty"`
	IrrigationSource               lookupValue `xml:"IrrigationSource" json:",omitempty"`
	IrrigationWaterRightsAcres     json.Number `xml:"IrrigationWaterRightsAcres" json:",omitempty"`
	IrrigationWaterRightsYN        boolYN      `xml:"IrrigationWaterRightsYN" json:",omitempty"`
	NumberOfSeparateElectricMeters json.Number `xml:"NumberOfSeparateElectricMeters" json:",omitempty"`
	NumberOfSeparateGasMeters      json.Number `xml:"NumberOfSeparateGasMeters" json:",omitempty"`
	NumberOfSeparateWaterMeters    json.Number `xml:"NumberOfSeparateWaterMeters" json:",omitempty"`
	Sewer                          lookupValue `xml:"Sewer" json:",omitempty"`
	Telephone                      string      `xml:"Telephone" json:",omitempty"`
	WaterSource                    lookupValue `xml:"WaterSource" json:",omitempty"`
}

// CoBuyerAgent is ...
type CoBuyerAgent struct {
	CoBuyerAgentAOR               string      `xml:"CoBuyerAgentAOR" json:",omitempty"`
	CoBuyerAgentCellPhone         string      `xml:"CoBuyerAgentCellPhone" json:",omitempty"`
	CoBuyerAgentDesignation       lookupValue `xml:"CoBuyerAgentDesignation" json:",omitempty"`
	CoBuyerAgentDirectPhone       string      `xml:"CoBuyerAgentDirectPhone" json:",omitempty"`
	CoBuyerAgentEmail             string      `xml:"CoBuyerAgentEmail" json:",omitempty"`
	CoBuyerAgentFax               string      `xml:"CoBuyerAgentFax" json:",omitempty"`
	CoBuyerAgentFirstName         string      `xml:"CoBuyerAgentFirstName" json:",omitempty"`
	CoBuyerAgentFullName          string      `xml:"CoBuyerAgentFullName" json:",omitempty"`
	CoBuyerAgentHomePhone         string      `xml:"CoBuyerAgentHomePhone" json:",omitempty"`
	CoBuyerAgentKey               string      `xml:"CoBuyerAgentKey" json:",omitempty"`
	CoBuyerAgentLastName          string      `xml:"CoBuyerAgentLastName" json:",omitempty"`
	CoBuyerAgentMiddleName        string      `xml:"CoBuyerAgentMiddleName" json:",omitempty"`
	CoBuyerAgentMlsId             string      `xml:"CoBuyerAgentMlsId" json:",omitempty"`
	CoBuyerAgentNamePrefix        string      `xml:"CoBuyerAgentNamePrefix" json:",omitempty"`
	CoBuyerAgentNameSuffix        string      `xml:"CoBuyerAgentNameSuffix" json:",omitempty"`
	CoBuyerAgentOfficePhone       string      `xml:"CoBuyerAgentOfficePhone" json:",omitempty"`
	CoBuyerAgentOfficePhoneExt    json.Number `xml:"CoBuyerAgentOfficePhoneExt" json:",omitempty"`
	CoBuyerAgentPager             string      `xml:"CoBuyerAgentPager" json:",omitempty"`
	CoBuyerAgentPreferredPhone    string      `xml:"CoBuyerAgentPreferredPhone" json:",omitempty"`
	CoBuyerAgentPreferredPhoneExt json.Number `xml:"CoBuyerAgentPreferredPhoneExt" json:",omitempty"`
	CoBuyerAgentStateLicense      string      `xml:"CoBuyerAgentStateLicense" json:",omitempty"`
	CoBuyerAgentTollFreePhone     string      `xml:"CoBuyerAgentTollFreePhone" json:",omitempty"`
	CoBuyerAgentURL               string      `xml:"CoBuyerAgentURL" json:",omitempty"`
	CoBuyerAgentVoiceMail         string      `xml:"CoBuyerAgentVoiceMail" json:",omitempty"`
	CoBuyerAgentVoiceMailExt      json.Number `xml:"CoBuyerAgentVoiceMailExt" json:",omitempty"`
}

// Remarks is ...
type Remarks struct {
	PrivateOfficeRemarks string `xml:"PrivateOfficeRemarks" json:",omitempty"`
	PrivateRemarks       string `xml:"PrivateRemarks" json:",omitempty"`
	PublicRemarks        string `xml:"PublicRemarks" json:",omitempty"`
	SyndicationRemarks   string `xml:"SyndicationRemarks" json:",omitempty"`
}

// GreenMarketing is ...
type GreenMarketing struct {
	GreenEnergyEfficient   lookupValue `xml:"GreenEnergyEfficient" json:",omitempty"`
	GreenEnergyGeneration  lookupValue `xml:"GreenEnergyGeneration" json:",omitempty"`
	GreenIndoorAirQuality  lookupValue `xml:"GreenIndoorAirQuality" json:",omitempty"`
	GreenLocation          lookupValue `xml:"GreenLocation" json:",omitempty"`
	GreenSustainability    lookupValue `xml:"GreenSustainability" json:",omitempty"`
	GreenWaterConservation lookupValue `xml:"GreenWaterConservation" json:",omitempty"`
}

// GreenCertification is ...
type GreenCertification struct {
	Rating        string      `xml:"Rating" json:",omitempty"`
	VerifyingBody string      `xml:"VerifyingBody" json:",omitempty"`
	YearCertified string      `xml:"YearCertified" json:",omitempty"`
	Type          lookupValue `xml:"Type" json:",omitempty"`
}

// Rooms is ...
type Rooms struct {
	Room       *[]Room     `xml:"Room" json:",omitempty"`
	RoomsTotal json.Number `xml:"RoomsTotal" json:",omitempty"`
}

// Team is ...
type Team struct {
	BuyerTeamDisplayName string `xml:"BuyerTeamDisplayName" json:",omitempty"`
	ListTeamDisplayName  string `xml:"ListTeamDisplayName" json:",omitempty"`
}

// Compensation is ...
type Compensation struct {
	BuyerAgencyCompensation           string      `xml:"BuyerAgencyCompensation" json:",omitempty"`
	BuyerAgencyCompensationType       string      `xml:"BuyerAgencyCompensationType" json:",omitempty"`
	DualVariableCompensationYN        boolYN      `xml:"DualVariableCompensationYN" json:",omitempty"`
	LeaseRenewalCompensation          lookupValue `xml:"LeaseRenewalCompensation" json:",omitempty"`
	SubAgencyCompensation             string      `xml:"SubAgencyCompensation" json:",omitempty"`
	SubAgencyCompensationType         string      `xml:"SubAgencyCompensationType" json:",omitempty"`
	TransactionBrokerCompensation     string      `xml:"TransactionBrokerCompensation" json:",omitempty"`
	TransactionBrokerCompensationType string      `xml:"TransactionBrokerCompensationType" json:",omitempty"`
}

// ListingMedia is ...
type ListingMedia struct {
	DocumentsAvailable       lookupValue `xml:"DocumentsAvailable" json:",omitempty"`
	DocumentsChangeTimestamp *dateTime   `xml:"DocumentsChangeTimestamp" json:",omitempty"`
	DocumentsCount           json.Number `xml:"DocumentsCount" json:",omitempty"`
	PhotosChangeTimestamp    *dateTime   `xml:"PhotosChangeTimestamp" json:",omitempty"`
	PhotosCount              json.Number `xml:"PhotosCount" json:",omitempty"`
	VideosChangeTimestamp    *dateTime   `xml:"VideosChangeTimestamp" json:",omitempty"`
	VideosCount              json.Number `xml:"VideosCount" json:",omitempty"`
}

// Equipment is ...
type Equipment struct {
	Appliances       lookupValue `xml:"Appliances" json:",omitempty"`
	OtherEquipment   lookupValue `xml:"OtherEquipment" json:",omitempty"`
	SecurityFeatures lookupValue `xml:"SecurityFeatures" json:",omitempty"`
}

// School is ...
type School struct {
	ElementarySchool             string `xml:"ElementarySchool" json:",omitempty"`
	ElementarySchoolDistrict     string `xml:"ElementarySchoolDistrict" json:",omitempty"`
	HighSchool                   string `xml:"HighSchool" json:",omitempty"`
	HighSchoolDistrict           string `xml:"HighSchoolDistrict" json:",omitempty"`
	MiddleOrJuniorSchool         string `xml:"MiddleOrJuniorSchool" json:",omitempty"`
	MiddleOrJuniorSchoolDistrict string `xml:"MiddleOrJuniorSchoolDistrict" json:",omitempty"`
}

// GreenCertificaiton is ...
type GreenCertificaiton struct {
	Metric json.Number `xml:"Metric" json:",omitempty"`
	URL    string      `xml:"URL" json:",omitempty"`
}

// Performance is ...
type Performance struct {
	GreenCertification *GreenCertification `xml:"GreenCertification" json:",omitempty"`
	GreenMarketing     *GreenMarketing     `xml:"GreenMarketing" json:",omitempty"`
}

// Property is ...
type Property struct {
	Characteristics *Characteristics `xml:"Characteristics" json:",omitempty"`
	Equipment       *Equipment       `xml:"Equipment" json:",omitempty"`
	Farming         *Farming         `xml:"Farming" json:",omitempty"`
	Financial       *Financial       `xml:"Financial" json:",omitempty"`
	HOA             *HOA             `xml:"HOA" json:",omitempty"`
	Location        *Location        `xml:"Location" json:",omitempty"`
	OccupantOwner   *OccupantOwner   `xml:"OccupantOwner" json:",omitempty"`
	PropertySubType string           `xml:"PropertySubType" json:",omitempty"`
	PropertyType    string           `xml:"PropertyType" json:",omitempty"`
	Structure       *Structure       `xml:"Structure" json:",omitempty"`
	Tax             *Tax             `xml:"Tax" json:",omitempty"`
	UnitType        *UnitType        `xml:"UnitType" json:",omitempty"`
	Utilities       *Utilities       `xml:"Utilities" json:",omitempty"`
	Listing         *Listing         `xml:"Listing" json:",omitempty"`
}

// BuyerAgent is ...
type BuyerAgent struct {
	BuyerAgentAOR               string      `xml:"BuyerAgentAOR" json:",omitempty"`
	BuyerAgentCellPhone         string      `xml:"BuyerAgentCellPhone" json:",omitempty"`
	BuyerAgentDesignation       lookupValue `xml:"BuyerAgentDesignation" json:",omitempty"`
	BuyerAgentDirectPhone       string      `xml:"BuyerAgentDirectPhone" json:",omitempty"`
	BuyerAgentEmail             string      `xml:"BuyerAgentEmail" json:",omitempty"`
	BuyerAgentFax               string      `xml:"BuyerAgentFax" json:",omitempty"`
	BuyerAgentFirstName         string      `xml:"BuyerAgentFirstName" json:",omitempty"`
	BuyerAgentFullName          string      `xml:"BuyerAgentFullName" json:",omitempty"`
	BuyerAgentHomePhone         string      `xml:"BuyerAgentHomePhone" json:",omitempty"`
	BuyerAgentKey               string      `xml:"BuyerAgentKey" json:",omitempty"`
	BuyerAgentLastName          string      `xml:"BuyerAgentLastName" json:",omitempty"`
	BuyerAgentMiddleName        string      `xml:"BuyerAgentMiddleName" json:",omitempty"`
	BuyerAgentMlsId             string      `xml:"BuyerAgentMlsId" json:",omitempty"`
	BuyerAgentNamePrefix        string      `xml:"BuyerAgentNamePrefix" json:",omitempty"`
	BuyerAgentNameSuffix        string      `xml:"BuyerAgentNameSuffix" json:",omitempty"`
	BuyerAgentOfficePhone       string      `xml:"BuyerAgentOfficePhone" json:",omitempty"`
	BuyerAgentOfficePhoneExt    json.Number `xml:"BuyerAgentOfficePhoneExt" json:",omitempty"`
	BuyerAgentPager             string      `xml:"BuyerAgentPager" json:",omitempty"`
	BuyerAgentPreferredPhone    string      `xml:"BuyerAgentPreferredPhone" json:",omitempty"`
	BuyerAgentPreferredPhoneExt json.Number `xml:"BuyerAgentPreferredPhoneExt" json:",omitempty"`
	BuyerAgentStateLicense      string      `xml:"BuyerAgentStateLicense" json:",omitempty"`
	BuyerAgentTollFreePhone     string      `xml:"BuyerAgentTollFreePhone" json:",omitempty"`
	BuyerAgentURL               string      `xml:"BuyerAgentURL" json:",omitempty"`
	BuyerAgentVoiceMail         string      `xml:"BuyerAgentVoiceMail" json:",omitempty"`
	BuyerAgentVoiceMailExt      json.Number `xml:"BuyerAgentVoiceMailExt" json:",omitempty"`
}

// CoListOffice is ...
type CoListOffice struct {
	CoListOfficeAOR      string      `xml:"CoListOfficeAOR" json:",omitempty"`
	CoListOfficeEmail    string      `xml:"CoListOfficeEmail" json:",omitempty"`
	CoListOfficeFax      string      `xml:"CoListOfficeFax" json:",omitempty"`
	CoListOfficeKey      string      `xml:"CoListOfficeKey" json:",omitempty"`
	CoListOfficeMlsId    string      `xml:"CoListOfficeMlsId" json:",omitempty"`
	CoListOfficeName     string      `xml:"CoListOfficeName" json:",omitempty"`
	CoListOfficePhone    string      `xml:"CoListOfficePhone" json:",omitempty"`
	CoListOfficePhoneExt json.Number `xml:"CoListOfficePhoneExt" json:",omitempty"`
	CoListOfficeURL      string      `xml:"CoListOfficeURL" json:",omitempty"`
}

// Listing is ...
type Listing struct {
	AgentOffice           *AgentOffice  `xml:"AgentOffice" json:",omitempty"`
	ApprovalStatus        string        `xml:"ApprovalStatus" json:",omitempty"`
	Closing               *Closing      `xml:"Closing" json:",omitempty"`
	Compensation          *Compensation `xml:"Compensation" json:",omitempty"`
	Contract              *Contract     `xml:"Contract" json:",omitempty"`
	CopyrightNotice       string        `xml:"CopyrightNotice" json:",omitempty"`
	Dates                 *Dates        `xml:"Dates" json:",omitempty"`
	Disclaimer            string        `xml:"Disclaimer" json:",omitempty"`
	HomeWarrantyYN        boolYN        `xml:"HomeWarrantyYN" json:",omitempty"`
	LeaseConsideredYN     boolYN        `xml:"LeaseConsideredYN" json:",omitempty"`
	ListAOR               string        `xml:"ListAOR" json:",omitempty"`
	ListingAgreement      string        `xml:"ListingAgreement" json:",omitempty"`
	ListingId             string        `xml:"ListingId" json:",omitempty"`
	ListingKey            string        `xml:"ListingKey" json:",omitempty"`
	ListingService        string        `xml:"ListingService" json:",omitempty"`
	Marketing             *Marketing    `xml:"Marketing" json:",omitempty"`
	Media                 *ListingMedia `xml:"Media" json:",omitempty"`
	MlsStatus             string        `xml:"MlsStatus" json:",omitempty"`
	OriginatingSystemKey  string        `xml:"OriginatingSystemKey" json:",omitempty"`
	OriginatingSystemName string        `xml:"OriginatingSystemName" json:",omitempty"`
	Price                 *Price        `xml:"Price" json:",omitempty"`
	Remarks               *Remarks      `xml:"Remarks" json:",omitempty"`
	Showing               *Showing      `xml:"Showing" json:",omitempty"`
	StandardStatus        string        `xml:"StandardStatus" json:",omitempty"`
}

// UnitType is ...
type UnitType struct {
	ActualRent       json.Number `xml:"ActualRent" json:",omitempty"`
	BathsTotal       json.Number `xml:"BathsTotal" json:",omitempty"`
	BedsTotal        json.Number `xml:"BedsTotal" json:",omitempty"`
	Description      string      `xml:"Description" json:",omitempty"`
	Furnished        string      `xml:"Furnished" json:",omitempty"`
	GarageAttachedYN boolYN      `xml:"GarageAttachedYN" json:",omitempty"`
	GarageSpaces     json.Number `xml:"GarageSpaces" json:",omitempty"`
	ProForma         json.Number `xml:"ProForma" json:",omitempty"`
	TotalRent        json.Number `xml:"TotalRent" json:",omitempty"`
	UnitsTotal       json.Number `xml:"UnitsTotal" json:",omitempty"`
	Type             lookupValue `xml:"Type" json:",omitempty"`
}

// Dates is ...
type Dates struct {
	CancelationDate          *dateTime   `xml:"CancelationDate" json:",omitempty"`
	CloseDate                *dateTime   `xml:"CloseDate" json:",omitempty"`
	ContingentDate           *dateTime   `xml:"ContingentDate" json:",omitempty"`
	ContractStatusChangeDate *dateTime   `xml:"ContractStatusChangeDate" json:",omitempty"`
	CumulativeDaysOnMarket   json.Number `xml:"CumulativeDaysOnMarket" json:",omitempty"`
	DaysOnMarket             json.Number `xml:"DaysOnMarket" json:",omitempty"`
	ExpirationDate           *dateTime   `xml:"ExpirationDate" json:",omitempty"`
	ListingContractDate      *dateTime   `xml:"ListingContractDate" json:",omitempty"`
	MajorChangeTimestamp     *dateTime   `xml:"MajorChangeTimestamp" json:",omitempty"`
	MajorChangeType          string      `xml:"MajorChangeType" json:",omitempty"`
	ModificationTimestamp    *dateTime   `xml:"ModificationTimestamp" json:",omitempty"`
	OffMarketDate            *dateTime   `xml:"OffMarketDate" json:",omitempty"`
	OffMarketTimestamp       *dateTime   `xml:"OffMarketTimestamp" json:",omitempty"`
	OnMarketDate             *dateTime   `xml:"OnMarketDate" json:",omitempty"`
	OnMarketTimestamp        *dateTime   `xml:"OnMarketTimestamp" json:",omitempty"`
	OriginalEntryTimestamp   *dateTime   `xml:"OriginalEntryTimestamp" json:",omitempty"`
	PendingTimestamp         *dateTime   `xml:"PendingTimestamp" json:",omitempty"`
	PriceChangeTimestamp     *dateTime   `xml:"PriceChangeTimestamp" json:",omitempty"`
	PurchaseContractDate     *dateTime   `xml:"PurchaseContractDate" json:",omitempty"`
	StatusChangeTimestamp    *dateTime   `xml:"StatusChangeTimestamp" json:",omitempty"`
	WithdrawnDate            *dateTime   `xml:"WithdrawnDate" json:",omitempty"`
}

// Farming is ...
type Farming struct {
	CropsIncludedYN               boolYN      `xml:"CropsIncludedYN" json:",omitempty"`
	CultivatedArea                json.Number `xml:"CultivatedArea" json:",omitempty"`
	FarmCreditServiceInclYN       boolYN      `xml:"FarmCreditServiceInclYN" json:",omitempty"`
	FarmLandAreaSource            string      `xml:"FarmLandAreaSource" json:",omitempty"`
	FarmLandAreaUnits             string      `xml:"FarmLandAreaUnits" json:",omitempty"`
	GrazingPermitsBlmYN           boolYN      `xml:"GrazingPermitsBlmYN" json:",omitempty"`
	GrazingPermitsForestServiceYN boolYN      `xml:"GrazingPermitsForestServiceYN" json:",omitempty"`
	GrazingPermitsPrivateYN       boolYN      `xml:"GrazingPermitsPrivateYN" json:",omitempty"`
	PastureArea                   json.Number `xml:"PastureArea" json:",omitempty"`
	RangeArea                     json.Number `xml:"RangeArea" json:",omitempty"`
	Vegetation                    lookupValue `xml:"Vegetation" json:",omitempty"`
	WoodedArea                    json.Number `xml:"WoodedArea" json:",omitempty"`
}

// Financial is ...
type Financial struct {
	CableTvExpense                json.Number `xml:"CableTvExpense" json:",omitempty"`
	CapRate                       json.Number `xml:"CapRate" json:",omitempty"`
	ElectricExpense               json.Number `xml:"ElectricExpense" json:",omitempty"`
	ExistingLeaseType             lookupValue `xml:"ExistingLeaseType" json:",omitempty"`
	FinancialDataSource           lookupValue `xml:"FinancialDataSource" json:",omitempty"`
	FuelExpense                   json.Number `xml:"FuelExpense" json:",omitempty"`
	FurnitureReplacementExpense   json.Number `xml:"FurnitureReplacementExpense" json:",omitempty"`
	GardnerExpense                json.Number `xml:"GardnerExpense" json:",omitempty"`
	GrossIncome                   json.Number `xml:"GrossIncome" json:",omitempty"`
	GrossScheduledIncome          json.Number `xml:"GrossScheduledIncome" json:",omitempty"`
	IncomeIncludes                lookupValue `xml:"IncomeIncludes" json:",omitempty"`
	InsuranceExpense              json.Number `xml:"InsuranceExpense" json:",omitempty"`
	LicensesExpense               json.Number `xml:"LicensesExpense" json:",omitempty"`
	MaintenanceExpense            json.Number `xml:"MaintenanceExpense" json:",omitempty"`
	ManagerExpense                json.Number `xml:"ManagerExpense" json:",omitempty"`
	NetOperatingIncome            json.Number `xml:"NetOperatingIncome" json:",omitempty"`
	NewTaxesExpense               json.Number `xml:"NewTaxesExpense" json:",omitempty"`
	NumberOfUnitsLeased           json.Number `xml:"NumberOfUnitsLeased" json:",omitempty"`
	NumberOfUnitsMoMo             json.Number `xml:"NumberOfUnitsMoMo" json:",omitempty"`
	NumberOfUnitsVacant           json.Number `xml:"NumberOfUnitsVacant" json:",omitempty"`
	OperatingExpense              json.Number `xml:"OperatingExpense" json:",omitempty"`
	OperatingExpenseIncludes      lookupValue `xml:"OperatingExpenseIncludes" json:",omitempty"`
	OtherExpense                  json.Number `xml:"OtherExpense" json:",omitempty"`
	OwnerPays                     lookupValue `xml:"OwnerPays" json:",omitempty"`
	PestControlExpense            json.Number `xml:"PestControlExpense" json:",omitempty"`
	PoolExpense                   json.Number `xml:"PoolExpense" json:",omitempty"`
	ProfessionalManagementExpense json.Number `xml:"ProfessionalManagementExpense" json:",omitempty"`
	RentControlYN                 boolYN      `xml:"RentControlYN" json:",omitempty"`
	RentIncludes                  lookupValue `xml:"RentIncludes" json:",omitempty"`
	SuppliesExpense               json.Number `xml:"SuppliesExpense" json:",omitempty"`
	TenantPays                    lookupValue `xml:"TenantPays" json:",omitempty"`
	TotalActualRent               json.Number `xml:"TotalActualRent" json:",omitempty"`
	TrashExpense                  json.Number `xml:"TrashExpense" json:",omitempty"`
	UnitsFurnished                string      `xml:"UnitsFurnished" json:",omitempty"`
	VacancyAllowance              json.Number `xml:"VacancyAllowance" json:",omitempty"`
	VacancyAllowanceRate          json.Number `xml:"VacancyAllowanceRate" json:",omitempty"`
	WaterSewerExpense             json.Number `xml:"WaterSewerExpense" json:",omitempty"`
	WorkmansCompensationExpense   json.Number `xml:"WorkmansCompensationExpense" json:",omitempty"`
}
