/**
provides the searching core
*/
package rets

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jpfielding/gominidom/minidom"
	"github.com/stretchr/testify/assert"
)

func TestSearchXMLEof(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(""))
	_, err := NewCompactSearchResult(body)
	assert.NotNil(t, err)
}

func TestSearchXMLBadChar(t *testing.T) {
	type Listing struct {
		Row string
	}
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<COUNT Records="5" />
			<Listings>
			<PropertyListing>
				<Row>bad` + "\x0b" + `row</Row>
			</PropertyListing>
			<PropertyListing>
				<Row>good row</Row>
			</PropertyListing>
			</Listings>
			<MAXROWS/>
			</RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))

	cr, err := NewStandardXMLSearchResult(body)
	assert.Nil(t, err)
	assert.Equal(t, StatusOK, cr.Response.Code)
	var listings []Listing
	count, maxRows, err := cr.ForEach(minidom.ByName("PropertyListing"), func(elem io.ReadCloser, err error) error {
		listing := Listing{}
		err = xml.NewDecoder(elem).Decode(&listing)
		listings = append(listings, listing)
		return err
	})
	assert.Nil(t, err)
	assert.Equal(t, true, maxRows)
	assert.Equal(t, 5, count)
	assert.Equal(t, 2, len(listings))
	assert.Equal(t, "bad row", listings[0].Row)
	assert.Equal(t, "good row", listings[1].Row)
}

func TestSearchXMLParseSearchQuit(t *testing.T) {
	noEnd := strings.Split(standardXML, "Commercial")[0]
	body := ioutil.NopCloser(strings.NewReader(noEnd))

	cr, err := NewStandardXMLSearchResult(body)
	assert.Nil(t, err)

	var listings [][]byte
	count, maxRows, err := cr.ForEach(minidom.ByName("PropertyListing"), func(elem io.ReadCloser, err error) error {
		tmp, _ := ioutil.ReadAll(elem)
		listings = append(listings, tmp)
		return err
	})
	assert.NotNil(t, err)
	assert.Equal(t, false, maxRows)
	assert.Equal(t, 10, count)
	assert.Equal(t, 1, len(listings))
}

func TestSearchXML(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(standardXML))

	cr, err := NewStandardXMLSearchResult(body)
	assert.Nil(t, err)

	var listings []io.ReadCloser
	count, maxRows, err := cr.ForEach(minidom.ByName("PropertyListing"), func(elem io.ReadCloser, err error) error {
		listings = append(listings, elem)
		return err
	})
	assert.Nil(t, err)
	assert.Equal(t, true, maxRows)
	assert.Equal(t, 10, count)
	assert.Equal(t, 2, len(listings))
}

func TestSearchXMLComplex(t *testing.T) {
	type Listing struct {
		Business      string  `xml:"Business>RESIOWNS"`
		Approved      bool    `xml:"Listing>Approved"`
		MLS           string  `xml:"Listing>MLS"`
		Disclaimer    string  `xml:"Listing>Disclaimer"`
		Status        string  `xml:"Listing>Status"`
		ListingPrice  float64 `xml:"Listing>Price>ListingPrice"`
		OriginalPrice float64 `xml:"Listing>Price>OriginalPrice"`
		SellPrice     float64 `xml:"Listing>Price>SellingPrice"`
	}
	body := ioutil.NopCloser(strings.NewReader(standardXML))

	cr, err := NewStandardXMLSearchResult(body)
	assert.Nil(t, err)

	var listings []Listing
	count, maxRows, err := cr.ForEach(minidom.ByName("PropertyListing"), func(elem io.ReadCloser, err error) error {
		if err != nil {
			return err
		}
		listing := Listing{}
		xml.NewDecoder(elem).Decode(&listing)
		listings = append(listings, listing)
		return err
	})
	assert.Nil(t, err)
	assert.Equal(t, true, maxRows)
	assert.Equal(t, 10, count)
	assert.Equal(t, 2, len(listings))
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
		  </PropertyListing>
	  </Residential>
	  <Commercial>
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
		  </PropertyListing>
		</Commercial>
      </REProperties>
    </REData>
	<MAXROWS/>
  </RETS>
`
