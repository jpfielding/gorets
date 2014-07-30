/**
provides the searching core
*/
package gorets_client

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
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

	done := make(chan struct{})
	defer close(done)
	_, err := parseCompactResult(done, body, 1)
	if err != io.EOF {
		t.Error("error parsing body: " + err.Error())
	}
}

func TestParseSearchQuit(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	quit := make(chan struct{})
	defer close(quit)
	cr, err := parseCompactResult(quit, body, 1)
	ok(t, err)

	row1 := <- cr.Data
	equals(t, "1,2,3,4,,6", strings.Join(row1, ","))

	quit <- struct{}{}

	// the closed channel will emit a zero'd value of the proper type
	row2 := <- cr.Data
	equals(t, "", strings.Join(row2, ","))

}

func TestParseCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	done := make(chan struct{})
	defer close(done)
	cr, err := parseCompactResult(done, body, 1)
	ok(t, err)

	assert(t, 0 == cr.RetsResponse.ReplyCode, "bad code")
	assert(t, "V2.7.0 2315: Success" == cr.RetsResponse.ReplyText, "bad text")

	assert(t, 10 == int(cr.Count), "bad count")
	assert(t, 6 == len(cr.Columns), "bad header count")

	assert(t, "A,B,C,D,E,F" == strings.Join(cr.Columns, ","), "bad headers")

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

		if cr.ProcessingFailure != nil {
			t.Errorf("error parsing body at row %d: %s", counter, err.Error())
		}

		counter = counter + 1
	}

	assert(t, 8 == counter, "bad count")
	assert(t, cr.MaxRows, "bad max rows")

}

func TestParseStandardXml(t *testing.T) {
}
