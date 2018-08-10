package syndication

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jpfielding/gominidom/minidom"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	doms := ioutil.NopCloser(strings.NewReader(example))
	parser := xml.NewDecoder(doms)
	listings := Listings{}
	// minidom isnt necessary but its crazy useful for massive streams
	md := minidom.MiniDom{
		StartFunc: func(start xml.StartElement) {
			switch start.Name.Local {
			case "Listings":
				attrs := map[string]string{}
				for _, v := range start.Attr {
					attrs[v.Name.Local] = v.Value
				}
				listings.ListingsKey = attrs["listingsKey"]
				listings.Version = attrs["version"]
				listings.VersionTimestamp = attrs["versionTimestamp"]
				listings.Language = attrs["lang"]
			case "Disclaimer":
				parser.DecodeElement(listings.Disclaimer, &start)
			}
		},
		// quit on the the xml tag
		EndFunc: minidom.QuitAt("Listings"),
	}
	err := md.Walk(parser, minidom.ByName("Listing"), ToListing(func(l Listing, err error) error {
		listings.Listings = append(listings.Listings, l)
		return err
	}))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(listings.Listings))
	assert.Equal(t, "http://www.somemls.com/lisings/1234567890", listings.Listings[0].ListingURL)
	assert.Equal(t, "New Light Fixtures", *listings.Listings[0].Photos[1].Caption)
	assert.Equal(t, "1100.0", listings.Listings[0].Expenses[2].Value.Value)
}

var example = `<Listings xmlns="http://rets.org/xsd/Syndication/2012-03" xmlns:commons="http://rets.org/xsd/RETSCommons" xmlns:schemaLocation="http://rets.org/xsd/Syndication/2012-03/Syndication.xsd" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" listingsKey="2012-03-06T22:14:47" version="0.96" versionTimestamp="2012-02-07T03:00:00Z" xml:lang="en-us">
	<Listing>
		<Address>
			<commons:preference-order>1</commons:preference-order>
			<commons:address-preference-order>1</commons:address-preference-order>
			<commons:FullStreetAddress>2245 Don Knotts Blvd.</commons:FullStreetAddress>
			<commons:UnitNumber>2</commons:UnitNumber>
			<commons:City>Morgantown</commons:City>
			<commons:StateOrProvince>WV</commons:StateOrProvince>
			<commons:PostalCode>26501</commons:PostalCode>
			<commons:Country>true</commons:Country>
		</Address>
		<ListPrice commons:isgSecurityClass="Public">234000</ListPrice>
		<ListPriceLow commons:isgSecurityClass="Public">214000</ListPriceLow>
		<AlternatePrices>
			<AlternatePrice>
				<AlternateListPrice commons:currencyCode="EUR" commons:isgSecurityClass="Public">483999.0</AlternateListPrice>
				<AlternateListPriceLow commons:currencyCode="EUR" commons:isgSecurityClass="Public">470000.0</AlternateListPriceLow>
			</AlternatePrice>
		</AlternatePrices>
		<ListingURL>http://www.somemls.com/lisings/1234567890</ListingURL>
		<ProviderName>SomeMLS</ProviderName>
		<ProviderURL>http://www.somemls.com</ProviderURL>
		<ProviderCategory>MLS</ProviderCategory>
		<LeadRoutingEmail>agent.lead.email@listhub.net</LeadRoutingEmail>
		<Bedrooms>3</Bedrooms>
		<Bathrooms>8</Bathrooms>
		<PropertyType otherDescription="Ranch">Commercial</PropertyType>
		<PropertySubType otherDescription="Ranch">Apartment</PropertySubType>
		<ListingKey>3yd-SOMEMLS-1234567890</ListingKey>
		<ListingCategory>Purchase</ListingCategory>
		<ListingStatus>Active</ListingStatus>
		<MarketingInformation>
			<commons:PermitAddressOnInternet commons:isgSecurityClass="Public">true</commons:PermitAddressOnInternet>
			<commons:VOWAddressDisplay commons:isgSecurityClass="Public">true</commons:VOWAddressDisplay>
			<commons:VOWAutomatedValuationDisplay commons:isgSecurityClass="Public">true</commons:VOWAutomatedValuationDisplay>
			<commons:VOWConsumerComment commons:isgSecurityClass="Public">true</commons:VOWConsumerComment>
		</MarketingInformation>
		<Photos>
			<Photo>
				<MediaModificationTimestamp commons:isgSecurityClass="Public">2012-03-06T17:14:47-05:00</MediaModificationTimestamp>
				<MediaURL>http://photos.listhub.com/listing123/1</MediaURL>
				<MediaCaption>Awesome Kitchen</MediaCaption>
				<MediaDescription>Kitchen was recently remodeled</MediaDescription>
			</Photo>
			<Photo>
				<MediaModificationTimestamp commons:isgSecurityClass="Public">2012-03-06T17:14:47-05:00</MediaModificationTimestamp>
				<MediaURL>http://photos.listhub.net/listing123/1</MediaURL>
				<MediaCaption>New Light Fixtures</MediaCaption>
				<MediaDescription>All light fixtures have been replaced</MediaDescription>
			</Photo>
		</Photos>
		<DiscloseAddress>true</DiscloseAddress>
		<ListingDescription>Fabulous home in terrific condition. Beautiful hard wood floors. Great place to raise a family. etc...</ListingDescription>
		<MlsId>SOMEMLS</MlsId>
		<MlsName>Listing Exchange Group</MlsName>
		<MlsNumber>1234567890</MlsNumber>
		<LivingArea>2200</LivingArea>
		<LotSize>130680.000000</LotSize>
		<YearBuilt>1992</YearBuilt>
		<ListingDate>2012-01-06</ListingDate>
		<ListingTitle>Ranch, Ranch - Morgantown, WV</ListingTitle>
		<FullBathrooms>2</FullBathrooms>
		<ThreeQuarterBathrooms>3</ThreeQuarterBathrooms>
		<HalfBathrooms>2</HalfBathrooms>
		<OneQuarterBathrooms>1</OneQuarterBathrooms>
		<ForeclosureStatus>REO - Bank Owned</ForeclosureStatus>
		<ListingParticipants>
			<Participant>
				<ParticipantKey>3yd-A2SELL-12345</ParticipantKey>
				<ParticipantId>12345</ParticipantId>
				<FirstName>John</FirstName>
				<LastName>Doe</LastName>
				<Role>Listing</Role>
				<PrimaryContactPhone>555555lead</PrimaryContactPhone>
				<OfficePhone>555555555</OfficePhone>
				<Email>l.0.null.null.2@leads.listhub.net</Email>
				<Fax>555-555-5555</Fax>
				<WebsiteURL>http://www.somemls.com/agents/12345</WebsiteURL>
			</Participant>
		</ListingParticipants>
		<VirtualTours>
			<VirtualTour>
				<MediaModificationTimestamp commons:isgSecurityClass="Public">2012-03-04T17:14:47-05:00</MediaModificationTimestamp>
				<MediaURL>http://virtualtour.com/listing/10923921</MediaURL>
				<MediaCaption>my virtual tour</MediaCaption>
				<MediaDescription>come see this property</MediaDescription>
			</VirtualTour>
		</VirtualTours>
		<Videos>
			<Video>
				<MediaModificationTimestamp commons:isgSecurityClass="Public">2012-03-04T17:14:47-05:00</MediaModificationTimestamp>
				<MediaURL>http://videos.listhub.com/listing/u93422/1</MediaURL>
				<MediaCaption>Awesome View</MediaCaption>
				<MediaDescription>This property overlooks downtown Morgantown</MediaDescription>
			</Video>
		</Videos>
		<Offices>
			<Office>
				<OfficeKey>3yd-A2SELL-OC1</OfficeKey>
				<OfficeId>OC1</OfficeId>
				<Level>???</Level>
				<OfficeCode>
					<OfficeCodeId>OC1</OfficeCodeId>
				</OfficeCode>
				<Name>Preview Listing Office</Name>
				<CorporateName>Preview Listing Office</CorporateName>
				<BrokerId>br0ker1d</BrokerId>
				<PhoneNumber>555-555-555</PhoneNumber>
				<Address>
					<commons:preference-order>1</commons:preference-order>
					<commons:address-preference-order>1</commons:address-preference-order>
					<commons:FullStreetAddress>2245 Don Knotts Blvd.</commons:FullStreetAddress>
					<commons:UnitNumber>2</commons:UnitNumber>
					<commons:City>Morgantown</commons:City>
					<commons:StateOrProvince>WV</commons:StateOrProvince>
					<commons:PostalCode>26501</commons:PostalCode>
					<commons:Country>USA</commons:Country>
				</Address>
				<Website>http://www.listoffice.com</Website>
			</Office>
		</Offices>
		<Brokerage>
			<Name>John Doe</Name>
			<Phone>555-555-lead</Phone>
			<Email>l.0.null.null.0@leads.listhub.net</Email>
			<WebsiteURL>http://johndoebrokerage.com/</WebsiteURL>
			<LogoURL>http://johndoebrokerage.com/logo.png</LogoURL>
			<Address>
				<commons:preference-order>1</commons:preference-order>
				<commons:address-preference-order>1</commons:address-preference-order>
				<commons:FullStreetAddress>2245 Don Knotts Blvd.</commons:FullStreetAddress>
				<commons:UnitNumber>2</commons:UnitNumber>
				<commons:City>Morgantown</commons:City>
				<commons:StateOrProvince>WV</commons:StateOrProvince>
				<commons:PostalCode>26501</commons:PostalCode>
				<commons:Country>true</commons:Country>
			</Address>
		</Brokerage>
		<Franchise>
			<Name>Advanced</Name>
    </Franchise>
    <Builder>
			<Name>Building Builders</Name>
			<Phone>999-999-9998</Phone>
			<Fax>999-999-9999</Fax>
			<Email>email@building.com</Email>
			<WebsiteURL>http://www.buildingbuilders.com</WebsiteURL>
			<Address>
				<commons:preference-order>1</commons:preference-order>
				<commons:address-preference-order>1</commons:address-preference-order>
				<commons:FullStreetAddress>2245 Don Knotts Blvd.</commons:FullStreetAddress>
				<commons:City>Morgantown</commons:City>
				<commons:StateOrProvince>WV</commons:StateOrProvince>
				<commons:PostalCode>26501</commons:PostalCode>
			</Address>
		</Builder>
		<Location>
			<Latitude>39.231</Latitude>
			<Longitude>-89.9383</Longitude>
			<Elevation>1000ft</Elevation>
			<Directions>Down the road, take a left, ford the stream, left at roundabout</Directions>
			<GeocodeOptions>WGS84</GeocodeOptions>
			<County>Monongalia</County>
			<ParcelId>12321</ParcelId>
			<Community>
				<commons:Subdivision commons:isgSecurityClass="Public">Cheat Crossings</commons:Subdivision>
				<commons:Schools>
					<commons:School>
						<commons:Name>Valley View</commons:Name>
						<commons:SchoolCategory>Elementary</commons:SchoolCategory>
						<commons:District commons:isgSecurityClass="Public">Monongalia</commons:District>
						<commons:Description>true</commons:Description>
					</commons:School>
					<commons:School>
						<commons:Name>MHS</commons:Name>
						<commons:SchoolCategory>High</commons:SchoolCategory>
						<commons:District commons:isgSecurityClass="Public">Monongalia</commons:District>
						<commons:Description>true</commons:Description>
					</commons:School>
					<commons:School>
						<commons:Name>Morgantown Jr High School</commons:Name>
						<commons:SchoolCategory>JuniorHigh</commons:SchoolCategory>
						<commons:District commons:isgSecurityClass="Public">Monongalia</commons:District>
						<commons:Description>true</commons:Description>
					</commons:School>
					<commons:School>
						<commons:Name>South</commons:Name>
						<commons:SchoolCategory>Middle</commons:SchoolCategory>
						<commons:District commons:isgSecurityClass="Public">Monongalia</commons:District>
						<commons:Description>true</commons:Description>
					</commons:School>
				</commons:Schools>
			</Community>
			<Neighborhoods>
				<Neighborhood>
					<Name>Downtown</Name>
					<Description>Fabulous Downtown Morgantown</Description>
				</Neighborhood>
				<Neighborhood>
					<Name>Industrial</Name>
					<Description>Business opportunities abound</Description>
				</Neighborhood>
			</Neighborhoods>
		</Location>
		<OpenHouses>
			<OpenHouse>
				<Date>2012-03-16</Date>
				<StartTime>5:14 PM</StartTime>
				<EndTime>9:14 PM</EndTime>
				<Description>Come out and see this lovely property!</Description>
			</OpenHouse>
		</OpenHouses>
		<Taxes>
			<Tax>
				<Year>2011</Year>
				<Amount>3400.0</Amount>
				<TaxDescription>tax description1</TaxDescription>
			</Tax>
			<Tax>
				<Year>2010</Year>
				<Amount>3300.0</Amount>
				<TaxDescription>tax description2</TaxDescription>
			</Tax>
		</Taxes>
		<Expenses>
			<Expense>
				<commons:ExpenseCategory>Trash Fee</commons:ExpenseCategory>
				<commons:ExpenseValue commons:currencyPeriod="Quarterly" commons:isgSecurityClass="Public">2000.0</commons:ExpenseValue>
			</Expense>
			<Expense>
				<commons:ExpenseCategory>Yard Care Fee</commons:ExpenseCategory>
				<commons:ExpenseValue commons:currencyPeriod="Annually" commons:isgSecurityClass="Public">2000.0</commons:ExpenseValue>
			</Expense>
			<Expense>
				<commons:ExpenseCategory>Home Owner Assessments Fee</commons:ExpenseCategory>
				<commons:ExpenseValue commons:currencyPeriod="Annually" commons:isgSecurityClass="Public">1100.0</commons:ExpenseValue>
			</Expense>
		</Expenses>
		<DetailedCharacteristics>
			<Appliances>
				<Appliance>Dishwasher</Appliance>
				<Appliance>Refrigerator</Appliance>
			</Appliances>
			<ArchitectureStyle otherDescription="Ranch">Cape Cod</ArchitectureStyle>
			<HasAttic>true</HasAttic>
			<HasBarbecueArea>true</HasBarbecueArea>
			<HasBasement>true</HasBasement>
			<BuildingUnitCount>1</BuildingUnitCount>
			<IsCableReady>true</IsCableReady>
			<HasCeilingFan>true</HasCeilingFan>
			<CondoFloorNum>1</CondoFloorNum>
			<CoolingSystems>
				<CoolingSystem>Central A/C</CoolingSystem>
			</CoolingSystems>
			<HasDeck>true</HasDeck>
			<HasDisabledAccess>true</HasDisabledAccess>
			<HasDock>true</HasDock>
			<HasDoorman>true</HasDoorman>
			<HasDoublePaneWindows>true</HasDoublePaneWindows>
			<HasElevator>true</HasElevator>
			<ExteriorTypes>
				<ExteriorType>Brick</ExteriorType>
				<ExteriorType>Vinyl Siding</ExteriorType>
			</ExteriorTypes>
			<HasFireplace>true</HasFireplace>
			<FloorCoverings>
				<FloorCovering>Carpet</FloorCovering>
				<FloorCovering>Wood</FloorCovering>
			</FloorCoverings>
			<HasGarden>true</HasGarden>
			<HasGatedEntry>true</HasGatedEntry>
			<HasGreenhouse>true</HasGreenhouse>
			<HeatingFuels>
				<HeatingFuel>Natural Gas</HeatingFuel>
			</HeatingFuels>
			<HeatingSystems>
				<HeatingSystem>Forced Air</HeatingSystem>
			</HeatingSystems>
			<HasHotTubSpa>true</HasHotTubSpa>
			<Intercom>true</Intercom>
			<HasJettedBathTub>true</HasJettedBathTub>
			<HasLawn>true</HasLawn>
			<LegalDescription>Legal description</LegalDescription>
			<HasMotherInLaw>true</HasMotherInLaw>
			<IsNewConstruction>false</IsNewConstruction>
			<NumFloors>2.0</NumFloors>
			<NumParkingSpaces>23</NumParkingSpaces>
			<HasPatio>true</HasPatio>
			<HasPond>true</HasPond>
			<HasPool>true</HasPool>
			<HasPorch>true</HasPorch>
			<RoofTypes>
				<RoofType>Composition Shingle</RoofType>
			</RoofTypes>
			<RoomCount>28</RoomCount>
			<Rooms>
				<Room>Bedroom</Room>
				<Room>Bedroom</Room>
				<Room>Bedroom</Room>
				<Room>Full Bath</Room>
				<Room>Full Bath</Room>
				<Room>Half Bath</Room>
				<Room>Half Bath</Room>
				<Room>Theatre</Room>
			</Rooms>
			<HasRVParking>true</HasRVParking>
			<HasSauna>true</HasSauna>
			<HasSecuritySystem>true</HasSecuritySystem>
			<HasSkylight>true</HasSkylight>
			<HasSportsCourt>true</HasSportsCourt>
			<HasSprinklerSystem>true</HasSprinklerSystem>
			<HasVaultedCeiling>true</HasVaultedCeiling>
			<ViewTypes>
				<ViewType>Mountain</ViewType>
			</ViewTypes>
			<IsWaterfront>true</IsWaterfront>
			<HasWetBar>true</HasWetBar>
			<IsWired>true</IsWired>
			<YearUpdated>2008</YearUpdated>
		</DetailedCharacteristics>
		<ModificationTimestamp commons:isgSecurityClass="Public">2012-03-06T17:14:47-05:00</ModificationTimestamp>
		<Disclaimer commons:isgSecurityClass="Public">Copyright © 2014 Listing Exchange Group. All rights reserved. All information provided by the listing agent/broker is deemed reliable but is not guaranteed and should be independently verified.</Disclaimer>
	</Listing>
</Listings>
`
