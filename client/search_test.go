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

	testutils "github.com/jpfielding/gorets/testutils"
)

var compactDecoded string = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
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

	data := make(chan []string)
	errs := make(chan error)
	quit := make(chan struct{})
	defer close(quit)
	_, err := parseCompactResult(body, data, errs, quit)
	if err == io.EOF {
		t.Error("eof should not surface: " + err.Error())
	}
}

func TestParseSearchQuit(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	data := make(chan []string)
	errs := make(chan error)
	quit := make(chan struct{})
	defer close(quit)
	cr, err := parseCompactResult(body, data, errs, quit)
	testutils.Ok(t, err)

	row1 := <-cr.Data
	testutils.Equals(t, "1,2,3,4,,6", strings.Join(row1, ","))

	quit <- struct{}{}

	// the closed channel will emit a zero'd value of the proper type
	row2 := <-cr.Data
	testutils.Equals(t, 0, len(row2))

}

func TestParseCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	data := make(chan []string)
	errs := make(chan error)
	quit := make(chan struct{})
	defer close(quit)
	cr, err := parseCompactResult(body, data, errs, quit)
	testutils.Ok(t, err)

	testutils.Assert(t, 0 == cr.RetsResponse.ReplyCode, "bad code")
	testutils.Assert(t, "V2.7.0 2315: Success" == cr.RetsResponse.ReplyText, "bad text")

	testutils.Assert(t, 10 == int(cr.Count), "bad count")
	testutils.Assert(t, 6 == len(cr.Columns), "bad header count")

	testutils.Assert(t, "A,B,C,D,E,F" == strings.Join(cr.Columns, ","), "bad headers")

	filterTo := cr.FilterTo([]string{"A", "C", "E"})

	counter := 0
	for row := range cr.Data {
		if strings.Join(row, ",") != "1,2,3,4,,6" {
			t.Errorf("bad row %d: %s", counter, row)
		}
		filtered := filterTo(row)
		if strings.Join(filtered, ",") != "1,3," {
			t.Errorf("bad filtered row %d: %s", counter, filtered)
		}

		select {
		case err := <-cr.Errors:
			t.Errorf("error parsing body at row %d: %s", counter, err.Error())
			close(quit)
		default:
		}
		counter = counter + 1
	}

	testutils.Assert(t, 8 == counter, "bad count")
	testutils.Assert(t, cr.MaxRows, "bad max rows")

}

func TestParseStandardXml(t *testing.T) {
}
