/**
provides the searching core
*/
package rets

import (
	"io/ioutil"
	"strings"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

func TestSearchXMLEof(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(""))
	_, err := NewCompactSearchResult(body)
	testutils.NotOk(t, err)
}

func TestSearchXMLBadChar(t *testing.T) {
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<COUNT Records="1" />
			<Listings>
			<PropertyListing>
			<Row>bad` + "\x0b" + `row</Row>
			</PropertyListing>
			</Listings>
			<MAXROWS/>
			</RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))

	cr, err := NewStandardXMLSearchResult(body, "PropertyListing")
	testutils.Ok(t, err)
	testutils.Equals(t, StatusOK, cr.Response.Code)
	testutils.Equals(t, 1, cr.Count)
	counter := 0
	maxRows, err := cr.ForEach(func(row map[string]string, err error) error {
		testutils.Equals(t, "bad row", row["Row"])
		counter = counter + 1
		return err
	})
	testutils.Ok(t, err)
	testutils.Equals(t, true, maxRows)
	testutils.Equals(t, 1, counter)
}

func TestSearchXMLNoEof(t *testing.T) {
	rets := `<RETS ReplyCode="20201" ReplyText="No Records Found." ></RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))

	cr, err := NewStandardXMLSearchResult(body, "PropertyListing")
	testutils.Ok(t, err)
	testutils.Equals(t, StatusNoRecords, cr.Response.Code)
}

func TestSearchXMLEmbeddedRetsStatus(t *testing.T) {
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<RETS-STATUS ReplyCode="20201" ReplyText="No matching records were found" />
			</RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))
	cr, err := NewStandardXMLSearchResult(body, "PropertyListing")
	testutils.Ok(t, err)
	testutils.Equals(t, StatusNoRecords, cr.Response.Code)
}

func TestSearchXMLParseSearchQuit(t *testing.T) {
	noEnd := strings.Split(standardXML, "<Commerical>")[0]
	body := ioutil.NopCloser(strings.NewReader(noEnd))

	cr, err := NewStandardXMLSearchResult(body, "PropertyListing")
	testutils.Ok(t, err)

	rowsFound := 0
	cr.ForEach(func(data map[string]string, err error) error {
		if err != nil {
			testutils.Assert(t, strings.Contains(err.Error(), "EOF"), "found something not eof")
			return err
		}
		testutils.Assert(t, len(data) > 0, "should have something")
		rowsFound++
		return nil
	})
	testutils.Equals(t, 1, rowsFound)
}

func TestSearchXMLParseCompact(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(standardXML))

	cr, err := NewStandardXMLSearchResult(body, "PropertyListing")
	testutils.Ok(t, err)

	testutils.Equals(t, StatusOK, cr.Response.Code)
	testutils.Equals(t, "Operation successful.", cr.Response.Text)

	testutils.Equals(t, 10, int(cr.Count))

	counter := 0
	var datas []map[string]string
	maxRows, err := cr.ForEach(func(row map[string]string, err error) error {
		datas = append(datas, row)
		counter++
		return err
	})
	testutils.Ok(t, err)

	testutils.Equals(t, 2, counter)
	testutils.Equals(t, datas[0]["Business/RESIOWNS"], "Private Owned")
	testutils.Equals(t, datas[0]["Listing/ListingRid"], "6798")
	testutils.Equals(t, datas[1]["Business/RESIOWNS"], "Business Owned")
	testutils.Equals(t, datas[1]["Listing/ListingRid"], "1234")
	testutils.Equals(t, true, maxRows)
}

var standardXML = `<?xml version="1.0" encoding="utf-8"?>
<RETS ReplyCode="0" ReplyText="Operation successful.">
  <COUNT Records="10" />
  <REData>
    <REProperties>
      <Residential>
        <PropertyListing>
          <Business>
            <RESIOWNS>Private Owned</RESIOWNS>
          </Business>
          <Listing>
            <Approved>TRUE</Approved>
            <MLS>Timbuktu</MLS>
            <PropertyDisclaimer>Information should be deemed reliable but not guaranteed, all representations are approximate, and individual verification is recommended. Copyright 2016 Rapattoni Corporation. All rights reserved. U.S. Patent 6,910,045</PropertyDisclaimer>
            <Status>Active</Status>
            <ListingRid>6798</ListingRid>
            <ListingNumberDisplay>3240009</ListingNumberDisplay>
            <MLSOrigin>Timbuktu</MLSOrigin>
            <RESIAGRR>Exclusive Rt to Sell</RESIAGRR>
            <RESIASSA>Timbuktu MLS</RESIASSA>
            <Price>
              <ListingPrice>1200000.00</ListingPrice>
              <OriginalPrice>1200000.00</OriginalPrice>
              <SellingPrice>0.00</SellingPrice>
            </Price>
          </Listing>
          <Property>
            <PropertySubtype1>Single Family Residence</PropertySubtype1>
            <PropertyType>Residential</PropertyType>
            <Characteristics>
              <Acres>34.0000</Acres>
              <LotMeasurement>Acres</LotMeasurement>
              <LotSizeSource>(Owner)</LotSizeSource>
              <LotSquareFootage>1481040.00</LotSquareFootage>
              <RESIROAD>Public</RESIROAD>
            </Characteristics>
          </Property>
		  </PropertyListing>
		<PropertyListing>
          <Business>
            <RESIOWNS>Business Owned</RESIOWNS>
          </Business>
          <Listing>
            <Approved>TRUE</Approved>
            <MLS>Timbuktu</MLS>
            <PropertyDisclaimer>Information should be deemed reliable but not guaranteed, all representations are approximate, and individual verification is recommended. Copyright 2016 Rapattoni Corporation. All rights reserved. U.S. Patent 6,910,045</PropertyDisclaimer>
            <Status>Active</Status>
            <ListingRid>1234</ListingRid>
            <ListingNumberDisplay>3240009</ListingNumberDisplay>
            <MLSOrigin>Timbuktu</MLSOrigin>
            <RESIAGRR>Exclusive Rt to Sell</RESIAGRR>
            <RESIASSA>Timbuktu MLS</RESIASSA>
            <Price>
              <ListingPrice>1200000.00</ListingPrice>
              <OriginalPrice>1200000.00</OriginalPrice>
              <SellingPrice>0.00</SellingPrice>
            </Price>
          </Listing>
          <Property>
            <PropertySubtype1>Single Family Residence</PropertySubtype1>
            <PropertyType>Residential</PropertyType>
            <Characteristics>
              <Acres>34.0000</Acres>
              <LotMeasurement>Acres</LotMeasurement>
              <LotSizeSource>(Owner)</LotSizeSource>
              <LotSquareFootage>1481040.00</LotSquareFootage>
              <RESIROAD>Public</RESIROAD>
            </Characteristics>
          </Property>
		  </PropertyListing>
		  </Residential>
      </REProperties>
    </REData>
	<MAXROWS/>
  </RETS>
`
