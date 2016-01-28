/**
provides the searching core
*/
package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/net/context"

	testutils "github.com/jpfielding/gorets/testutils"
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

	_, err := NewCompactSearchResult(context.Background(), body, 0)
	if err == io.EOF {
		t.Error("eof should not surface: " + err.Error())
	}
}

func TestParseSearchQuit(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	ctx, cancel := context.WithCancel(context.Background())
	cr, err := NewCompactSearchResult(ctx, body, 0)
	testutils.Ok(t, err)

	row1 := <-cr.Data
	testutils.Equals(t, "1,2,3,4,,6", strings.Join(row1, ","))

	cancel()
	testutils.Equals(t, "context canceled", ctx.Err().Error())

	// still one left before the cancel call is seen
	row2 := <-cr.Data
	testutils.Equals(t, "1,2,3,4,,6", strings.Join(row2, ","))

	// the closed channel will emit a zero'd value of the proper type
	row3 := <-cr.Data
	testutils.Equals(t, 0, len(row3))

}

func TestParseCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cr, err := NewCompactSearchResult(ctx, body, 0)
	testutils.Ok(t, err)

	testutils.Assert(t, 0 == cr.RetsResponse.ReplyCode, "bad code")
	testutils.Assert(t, "V2.7.0 2315: Success" == cr.RetsResponse.ReplyText, "bad text")

	testutils.Assert(t, 10 == int(cr.Count), "bad count")
	testutils.Assert(t, 6 == len(cr.Columns), "bad header count")

	testutils.Assert(t, "A,B,C,D,E,F" == strings.Join(cr.Columns, ","), "bad headers")

	counter := 0
	for row := range cr.Data {
		if strings.Join(row, ",") != "1,2,3,4,,6" {
			t.Errorf("bad row %d: %s", counter, row)
		}

		select {
		case err := <-cr.Errors:
			t.Errorf("error parsing body at row %d: %s", counter, err.Error())
			cancel()
		default:
		}
		counter = counter + 1
	}

	testutils.Assert(t, 8 == counter, "bad count")
	testutils.Assert(t, cr.MaxRows, "bad max rows")

}

func TestParseStandardXml(t *testing.T) {
}
