/**
provides the searching core
*/
package rets

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
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

func TestEof(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte("")))

	_, err := NewCompactSearchResult(body)
	testutils.NotOk(t, err)
}

func TestBadChar(t *testing.T) {
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<COUNT Records="1" />
			<DELIMITER value = "09"/>
			<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
			<DATA>	1` + "\x0b" + `1	2	3	4		6	</DATA>

			</RETS>`
	body := ioutil.NopCloser(bytes.NewReader([]byte(rets)))

	cr, err := NewCompactSearchResult(body)
	testutils.Ok(t, err)
	testutils.Equals(t, 0, cr.RetsResponse.ReplyCode)
	counter := 0
	cr.ForEach(func(row []string, err error) error {
		testutils.Ok(t, err)
		testutils.Equals(t, "1 1,2,3,4,,6", strings.Join(row, ","))
		counter = counter + 1
		return nil
	})
}

func TestNoEof(t *testing.T) {
	rets := `<RETS ReplyCode="20201" ReplyText="No Records Found." ></RETS>`
	body := ioutil.NopCloser(bytes.NewReader([]byte(rets)))

	cr, err := NewCompactSearchResult(body)
	testutils.Ok(t, err)
	testutils.Equals(t, 20201, cr.RetsResponse.ReplyCode)
}

func TestEmbeddedRetsStatus(t *testing.T) {
	rets := `<?xml version="1.0" encoding="UTF-8" ?>
			<RETS ReplyCode="0" ReplyText="Operation Successful">
			<RETS-STATUS ReplyCode="20201" ReplyText="No matching records were found" />
			</RETS>`
	body := ioutil.NopCloser(bytes.NewReader([]byte(rets)))
	cr, err := NewCompactSearchResult(body)
	testutils.Ok(t, err)
	testutils.Equals(t, 20201, cr.RetsResponse.ReplyCode)
}

func TestParseSearchQuit(t *testing.T) {
	noEnd := strings.Split(compactDecoded, "<MAXROWS/>")[0]
	body := ioutil.NopCloser(bytes.NewReader([]byte(noEnd)))

	cr, err := NewCompactSearchResult(body)
	testutils.Ok(t, err)

	rowsFound := 0
	cr.ForEach(func(data []string, err error) error {
		if err != nil {
			testutils.Assert(t, strings.Contains(err.Error(), "EOF"), "found something not eof")
			return err
		}
		testutils.Equals(t, "1,2,3,4,,6", strings.Join(data, ","))
		rowsFound++
		return nil
	})
	testutils.Equals(t, 8, rowsFound)
}

func TestParseCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	cr, err := NewCompactSearchResult(body)
	testutils.Ok(t, err)

	testutils.Assert(t, 0 == cr.RetsResponse.ReplyCode, "bad code")
	testutils.Assert(t, "V2.7.0 2315: Success" == cr.RetsResponse.ReplyText, "bad text")

	testutils.Assert(t, 10 == int(cr.Count), "bad count")
	testutils.Assert(t, 6 == len(cr.Columns), "bad header count")

	testutils.Assert(t, "A,B,C,D,E,F" == strings.Join(cr.Columns, ","), "bad headers")

	counter := 0
	maxRows, err := cr.ForEach(func(row []string, err error) error {
		if strings.Join(row, ",") != "1,2,3,4,,6" {
			t.Errorf("bad row %d: %s", counter, row)
		}

		if err != nil {
			t.Errorf("error parsing body at row %d: %s", counter, err.Error())
		}
		counter = counter + 1
		return nil
	})

	testutils.Assert(t, 8 == counter, "bad count")
	testutils.Assert(t, maxRows, "bad max rows")

}
