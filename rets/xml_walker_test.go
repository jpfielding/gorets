package rets

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"bufio"
	"bytes"

	test "github.move.com/Industry-Platforms/ghostwheel/lib/testutils"
)

var body = `<?xml version="1.0" encoding="UTF-8"?>
<Listings xmlns="http://rets.org/xsd/Syndication/2008-03" xmlns:commons="http://rets.org/xsd/RETSCommons/2007-08" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
<Listing>
    <Address>
        <commons:preference-order>1</commons:preference-order>
        <commons:address-preference-order>1</commons:address-preference-order>
        <commons:FullStreetAddress>200 Winston St.</commons:FullStreetAddress>
        <commons:City>Chunchula</commons:City>
        <commons:StateOrProvince>AL</commons:StateOrProvince>
        <commons:PostalCode>36521</commons:PostalCode>
        <commons:Country>United States</commons:Country>
        </Address>
        <ListPrice>210000.00</ListPrice>
        <ListingURL>http://example.com</ListingURL>
        <ProviderName>Assets-Sell, Inc.</ProviderName>
        <ProviderURL>www.Assetssell.com</ProviderURL>
        <Bedrooms>4</Bedrooms>
        <Bathrooms>3</Bathrooms>
        <PropertyType>Single Family</PropertyType>
        <ListingKey>347502</ListingKey>
        <ListingCategory>Purchase</ListingCategory>
        <ListingStatus>Active</ListingStatus>
        <ModificationTimestamp>2012-05-18T16:57:00</ModificationTimestamp>
        <Photos>
        <Photo><URL>http://example.com.com/photos/listing/original/824715.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824716.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824717.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824718.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824719.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824720.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824721.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824722.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824723.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824724.jpeg</URL></Photo>
        <Photo><URL>http://example.com.com/photos/listing/original/824725.jpeg</URL></Photo>
        </Photos>
        <ListingDescription>&amp;lt;br&amp;gt;Large country home with 5.75 cleared acres.  Country Kitchen with loads of counter space, Breakfast area, formal Dining room and Family room with wood-burning fireplace.  Master bedroom is downstairs with a walk-in closet, and there are three (3) large bedrooms upstairs.  Upstairs bedroom opens to large wooden deck.  Security bars on windows and doors on first floor.  A special feature of the home is a full basement...great for your workshop, special hobbies and a hurricane shelter.  Seller said basement has never flooded.  Seller will consider all offers.  Great country home location for travel either to Mobile or Washington County.</ListingDescription>
        <MlsNumber>525557</MlsNumber>
        <LivingArea>2880</LivingArea>
        <LotSize>5.750000</LotSize>
        <ListingDate>2015-10-09</ListingDate>
        <ListingTitle>COUNTRY LIVING</ListingTitle>
        <FullBathrooms>3</FullBathrooms>
        <ListingParticipants>
            <Participant>
            <ParticipantId>3448</ParticipantId>
            <FirstName>Richard</FirstName>
            <LastName>March</LastName>
            <OfficePhone>304-555-9040</OfficePhone>
            <Email>clownpants@gmail.com</Email>
            <Fax>251-665-4622</Fax>
            </Participant>
        </ListingParticipants>
        <Brokerage>
            <OfficeID>3H5</OfficeID>
            <Name>Buyers and Sellers Realty</Name>
            <Phone>304-555-9040</Phone>
            <Email>info@4mobileareahomes.com</Email>
            <WebsiteURL>4MobileAreaHomes.com</WebsiteURL>
            <Address>
            <commons:preference-order>1</commons:preference-order>
            <commons:address-preference-order>1</commons:address-preference-order>
            <commons:FullStreetAddress>3951 Burma Road</commons:FullStreetAddress>
            <commons:City>Mobile</commons:City>
            <commons:StateOrProvince>AL</commons:StateOrProvince>
            <commons:PostalCode>36693</commons:PostalCode>
            <commons:Country>United States</commons:Country>
            </Address>
        </Brokerage>
        <Franchise>
            <Name>Assets-2-Sell, Inc.</Name>
            <Phone>800 555-7816</Phone>
            <Email>Info@assets2sell.com</Email>
            <WebsiteURL>www.assets2sell.com</WebsiteURL>
            <Address>
                <commons:preference-order>1</commons:preference-order>
                <commons:address-preference-order>1</commons:address-preference-order>
                <commons:FullStreetAddress>1610 Meadow Wood Lane</commons:FullStreetAddress>
                <commons:City>Reno</commons:City>
                <commons:StateOrProvince>Nevada</commons:StateOrProvince>
                <commons:PostalCode>89502</commons:PostalCode>
                <commons:Country>United States</commons:Country>
            </Address>
        </Franchise>
        <Location>
        <Neighborhoods>
        <Neighborhood>
        <Description>Other</Description>
        </Neighborhood>
        </Neighborhoods>
        </Location>
        <Schools>
            <SchoolDistrict/>
        </Schools>
        <DetailedCharacteristics>
        	<NumFloors>2</NumFloors>
        	<NumParkingSpaces>0</NumParkingSpaces>
        </DetailedCharacteristics>
    </Listing>
    <Listing><Address><commons:preference-order>1</commons:preference-order><commons:address-preference-order>1</commons:address-preference-order><commons:FullStreetAddress>1313 Mockingbird Lane</commons:FullStreetAddress><commons:City>Daphne</commons:City><commons:StateOrProvince>AL</commons:StateOrProvince><commons:PostalCode>36526</commons:PostalCode><commons:Country>United States</commons:Country></Address><ListPrice>134900.00</ListPrice><ListingURL>http://assets2sell.com/listing/abasdfsadfa234</ListingURL><ProviderName>Assets-2-Sell, Inc.</ProviderName><ProviderURL>www.Assetssell.com</ProviderURL><Bedrooms>3</Bedrooms><Bathrooms>2</Bathrooms><PropertyType>Single Family</PropertyType><ListingKey>347707</ListingKey><ListingCategory>Purchase</ListingCategory><ListingStatus>Active</ListingStatus><ModificationTimestamp>2015-10-21T09:59:00</ModificationTimestamp><Photos><Photo><URL>http://example.com.com/photos/listing/original/826831.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826832.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826833.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826834.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826835.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826836.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826837.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826838.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826839.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826840.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826841.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826842.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826843.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826844.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826845.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826846.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826847.jpeg</URL></Photo><Photo><URL>http://example.com.com/photos/listing/original/826848.jpeg</URL></Photo></Photos><ListingDescription>Contemporary design with 3 bedrooms, 2 baths, features include vaulted ceiling and woodburning fireplace in the living room, separate dining, marble vanities in both bathrooms, large deck and spacious backyard- all at the end of a cul-de-sac. Just off Main Street in Daphne, so you can enjoy all the Lake Forest amenities without driving all the way through Lake Forest.</ListingDescription><LivingArea>1452</LivingArea><LotSize>0.367309</LotSize><YearBuilt>1980</YearBuilt><ListingDate>2015-10-19</ListingDate><ListingTitle>Great Outdoor Living Space!</ListingTitle><FullBathrooms>2</FullBathrooms><ListingParticipants><Participant><ParticipantId>7761</ParticipantId><FirstName>Helter</FirstName><LastName>Skelter</LastName><OfficePhone>(304) 555-9090</OfficePhone><Email>helterskelter@chimichangas.com</Email><Fax>(866) 555-8696</Fax></Participant></ListingParticipants><Brokerage><OfficeID>7EE</OfficeID><Name>Balding Realty</Name><Phone>(304) 555-9090</Phone><Email>info@4BaldwinAreaHomes.com</Email><WebsiteURL>4BaldwinAreaHomes.com</WebsiteURL><Address><commons:preference-order>1</commons:preference-order><commons:address-preference-order>1</commons:address-preference-order><commons:FullStreetAddress>666 Nohope Avenue</commons:FullStreetAddress><commons:City>Fairhope</commons:City><commons:StateOrProvince>AL</commons:StateOrProvince><commons:PostalCode>36532</commons:PostalCode><commons:Country>United States</commons:Country></Address></Brokerage><Franchise><Name>Assets-2-Sell, Inc.</Name><Phone>800 555-7816</Phone><Email>Info@assets2sell.com</Email><WebsiteURL>www.assets2sell.com</WebsiteURL><Address><commons:preference-order>1</commons:preference-order><commons:address-preference-order>1</commons:address-preference-order><commons:FullStreetAddress>1610 Meadow Wood Lane</commons:FullStreetAddress><commons:City>Reno</commons:City><commons:StateOrProvince>Nevada</commons:StateOrProvince><commons:PostalCode>89502</commons:PostalCode><commons:Country>United States</commons:Country></Address></Franchise><Location><Neighborhoods><Neighborhood><Description>Daphne</Description></Neighborhood></Neighborhoods></Location><Schools><SchoolDistrict/></Schools><Taxes><Tax><Amount>388.00</Amount></Tax></Taxes><DetailedCharacteristics><NumFloors>1</NumFloors><NumParkingSpaces>0</NumParkingSpaces></DetailedCharacteristics></Listing>
</Listings>
`

func TestSimpleXMLListings(t *testing.T) {
	var content = ioutil.NopCloser(bytes.NewReader([]byte(body)))
	parser := DefaultXMLDecoder(content, false)
	repeatable := []string{
		"Photos/Photo",
		"ListingParticipants/Participant",
	}
	xds := XMLData{
		Prefix:      "Listing",
		EndFunc:     XMLDataExit("Listings"),
		RepeatElems: repeatable,
	}
	var datas []map[string]string
	err := xds.Walk(parser, func(data map[string]string, err error) error {
		if err != nil {
			return err
		}
		datas = append(datas, data)
		return nil
	})
	test.Ok(t, err)
	test.Equals(t, datas[0]["Address/preference-order"], "1")
	test.Equals(t, datas[0]["Address/FullStreetAddress"], "200 Winston St.")
	test.Equals(t, datas[0]["ListingParticipants/Participant[1]/Email"], "clownpants@gmail.com")
	test.Equals(t, datas[0]["Photos/Photo[11]/URL"], "http://example.com.com/photos/listing/original/824725.jpeg")
	test.Equals(t, datas[1]["Address/preference-order"], "1")
	test.Equals(t, datas[1]["Photos/Photo[18]/URL"], "http://example.com.com/photos/listing/original/826848.jpeg")
	test.Equals(t, datas[1]["ListingParticipants/Participant[1]/OfficePhone"], "(304) 555-9090")
}

func TestSimpleXMLListingsError(t *testing.T) {
	var content = ioutil.NopCloser(bytes.NewReader([]byte(body)))
	parser := DefaultXMLDecoder(content, false)

	repeatable := []string{}
	xds := XMLData{
		Prefix:      "Listing",
		EndFunc:     XMLDataExit("Listings"),
		RepeatElems: repeatable,
	}
	xds.Walk(parser, func(data map[string]string, err error) error {
		test.NotOk(t, err)
		test.Assert(t, len(data) > 0, "should still have some data")
		return err
	})
}

func TestSimpleXMLBrokerage(t *testing.T) {
	var content = ioutil.NopCloser(bytes.NewReader([]byte(body)))
	parser := DefaultXMLDecoder(content, false)

	repeatable := []string{}
	xds := XMLData{
		Prefix:      "Brokerage",
		EndFunc:     XMLDataExit("Listings"),
		RepeatElems: repeatable,
	}
	var datas []map[string]string
	err := xds.Walk(parser, func(data map[string]string, err error) error {
		if err != nil {
			return err
		}
		datas = append(datas, data)
		return nil
	})
	test.Ok(t, err)

	test.Equals(t, datas[0]["Address/preference-order"], "1")
	test.Equals(t, datas[0]["OfficeID"], "3H5")
	test.Equals(t, datas[0]["Phone"], "304-555-9040")

	test.Equals(t, datas[1]["Address/preference-order"], "1")
	test.Equals(t, datas[1]["OfficeID"], "7EE")
	test.Equals(t, datas[1]["Phone"], "(304) 555-9090")

}

func TestReplaceXMLEncoding(t *testing.T) {
	content := func() io.ReadCloser {
		content := bytes.NewReader([]byte(body))
		// read out the bad xml header
		buffered := bufio.NewReaderSize(content, 1024)
		buffered.ReadLine()
		return ioutil.NopCloser(buffered)
	}()
	parser := DefaultXMLDecoder(content, false)
	repeatable := []string{
		"Photos/Photo",
		"ListingParticipants/Participant",
	}
	xds := XMLData{
		Prefix:      "Listing",
		EndFunc:     XMLDataExit("Listings"),
		RepeatElems: repeatable,
	}
	err := xds.Walk(parser, func(data map[string]string, err error) error {
		test.Assert(t, len(data) != 0, "should not be empty")
		return nil
	})
	test.Ok(t, err)
}

func dump(data map[string]string) {
	row, _ := json.Marshal(data)
	var out bytes.Buffer
	json.Indent(&out, row, "", "\t")
	fmt.Println(out.String())
}
