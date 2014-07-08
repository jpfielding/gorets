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

	_, err := parseCompactResult(body, 1)
	if err != io.EOF {
		t.Error("error parsing body: " + err.Error())
	}
}

func TestParseCompact(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(compactDecoded)))

	cr, err := parseCompactResult(body, 1)
	if err != nil {
		t.Error("error parsing body: " + err.Error())
	}
	AssertEqualsInt(t, "bad code", 0, cr.RetsResponse.ReplyCode)
	AssertEquals(t, "bad text", "V2.7.0 2315: Success", cr.RetsResponse.ReplyText)

	AssertEqualsInt(t, "bad count", 10, int(cr.Count))
	AssertEqualsInt(t, "bad header count", 6, len(cr.Columns))

	AssertEquals(t, "bad headers", "A,B,C,D,E,F", strings.Join(cr.Columns, ","))

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

	AssertEqualsInt(t, "bad count", 8, counter)
	AssertEqualsTrue(t, "bad max rows", cr.MaxRows)

}

func TestParseStandardXml(t *testing.T) {
}
