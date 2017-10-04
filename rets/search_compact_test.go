/**
provides the searching core
*/
package rets

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// GET /platinum/search?
// Class=ALL&
// Count=1&
// Format=COMPACT-DECODED&
// Limit=10&
// Offset=50&
// Query=%28%28LocaleListingStatus%3D%7CACTIVE-CORE%2CCNTG%2FKO-CORE%2CCNTG%2FNO+KO-CORE%2CAPP+REG-CORE%29%2C%7E%28VOWList%3D0%29%29&
// QueryType=DMQL2&
// SearchType=Property

var compactDecoded = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
<COUNT Records="10" />
<DELIMITER value = "09"/>
<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<DATA>	1	2	3	4		6	</DATA>
<MAXROWS/>
</RETS>
`

func TestSearchCompactEof(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(""))

	_, err := NewCompactSearchResult(body)
	assert.NotNil(t, err)
}

func TestSearchCompactBadChar(t *testing.T) {
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<COUNT Records="1" />
			<DELIMITER value = "09"/>
			<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
			<DATA>	1` + "\x0b" + `1	2	3	4		6	</DATA>

			</RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))

	cr, err := NewCompactSearchResult(body)
	assert.Nil(t, err)
	assert.Equal(t, StatusOK, cr.Response.Code)
	counter := 0
	cr.ForEach(func(row Row, err error) error {
		assert.Nil(t, err)
		assert.Equal(t, "1 1,2,3,4,,6", strings.Join(row, ","))
		counter = counter + 1
		return nil
	})
}

func TestSearchCompactNoEof(t *testing.T) {
	rets := `<RETS ReplyCode="20201" ReplyText="No Records Found." ></RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))

	cr, err := NewCompactSearchResult(body)
	defer cr.Close()
	assert.Nil(t, err)
	assert.Equal(t, StatusNoRecords, cr.Response.Code)
	found := 0
	max, err := cr.ForEach(func(data Row, err error) error {
		if err != nil {
			assert.Nil(t, err)
			return err
		}
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, max, false)
	assert.Equal(t, 0, found)
}

func TestSearchCompactWithDelimNoEof(t *testing.T) {
	rets := `<?xml version="1.0" encoding="iso-8859-1"?>
			<RETS ReplyCode="0" ReplyText="Operation successful" >
			<DELIMITER value="09"/>
			</RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))

	cr, err := NewCompactSearchResult(body)
	defer cr.Close()
	assert.Equal(t, cr.Delimiter, "	")
	assert.Nil(t, err)
	assert.Equal(t, StatusOK, cr.Response.Code)
	found := 0
	max, err := cr.ForEach(func(data Row, err error) error {
		if err != nil {
			assert.Nil(t, err)
			return err
		}
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, max, false)
	assert.Equal(t, 0, found)
}

func TestSearchCompactEmbeddedRetsStatus(t *testing.T) {
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<RETS-STATUS ReplyCode="20201" ReplyText="No matching records were found" />
			</RETS>`
	body := ioutil.NopCloser(strings.NewReader(rets))
	cr, err := NewCompactSearchResult(body)
	assert.Nil(t, err)
	assert.Equal(t, StatusNoRecords, cr.Response.Code)
}

func TestSearchCompactParseSearchQuit(t *testing.T) {
	noEnd := strings.Split(compactDecoded, "<MAXROWS/>")[0]
	body := ioutil.NopCloser(strings.NewReader(noEnd))

	cr, err := NewCompactSearchResult(body)
	assert.Nil(t, err)

	rowsFound := 0
	cr.ForEach(func(data Row, err error) error {
		if err != nil {
			assert.Equal(t, strings.Contains(err.Error(), "EOF"), "found something not eof")
			return err
		}
		assert.Equal(t, "1,2,3,4,,6", strings.Join(data, ","))
		rowsFound++
		return nil
	})
	assert.Equal(t, 8, rowsFound)
}

func TestSearchCompactParseCompact(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(compactDecoded))

	cr, err := NewCompactSearchResult(body)
	assert.Nil(t, err)

	assert.Equal(t, StatusOK == cr.Response.Code, "bad code")
	assert.Equal(t, "V2.7.0 2315: Success" == cr.Response.Text, "bad text")

	assert.Equal(t, 10 == int(cr.Count), "bad count")
	assert.Equal(t, 6 == len(cr.Columns), "bad header count")

	assert.Equal(t, "A,B,C,D,E,F" == strings.Join(cr.Columns, ","), "bad headers")

	counter := 0
	maxRows, err := cr.ForEach(func(row Row, err error) error {
		if strings.Join(row, ",") != "1,2,3,4,,6" {
			t.Errorf("bad row %d: %s", counter, row)
		}

		if err != nil {
			t.Errorf("error parsing body at row %d: %s", counter, err.Error())
		}
		counter = counter + 1
		return nil
	})
	assert.Nil(t, err)

	assert.Equal(t, 8 == counter, "bad count")
	assert.Equal(t, maxRows, "bad max rows")

}
