/**
provides the searching core
*/
package rets

import (
	"fmt"
	"strings"
	"testing"

	testutils "github.com/jpfielding/gotest/testutils"
)

var payloadlist = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
	<DELIMITER value = "09"/>
	<RETSPayloadList
		Resource="RESOURCE"
		Class="CLASS_1"
		Version="Version"
		Date="Tue, 3 Sep 2013 00:00:00 GMT">
		<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
	</RETSPayloadList>
	<RETSPayloadList
		Resource="RESOURCE"
		Class="CLASS_2"
		Version="Version"
		Date="Tue, 3 Sep 2013 00:00:00 GMT">
		<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
		<DATA>	1	2	3	4	5	6	</DATA>
	</RETSPayloadList>
</RETS>
`

/*func TestParseGetPayloadList(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewReader([]byte(payloadlist)))

	pl, err := parseGetPayloadList(body)
	if err != nil {
		t.Error("error parsing body: "+ err.Error())
	}
	AssertEqualsInt(t, "bad code", 0, pl.Rets.ReplyCode)
	AssertEquals(t, "bad text", "V2.7.0 2315: Success", pl.Rets.ReplyText)

	verifyCompactData(t, pl,"RESOURCE:CLASS_1")
	verifyCompactData(t, pl,"RESOURCE:CLASS_2")

}*/

func verifyCompactData(t *testing.T, pl *PayloadList, id string) {
	payload := <-pl.Payloads
	testutils.Equals(t, 6, len(payload.Columns()))

	testutils.Equals(t, id, payload.ID)
	testutils.Equals(t, "A,B,C,D,E,F", strings.Join(payload.Columns(), ","))

	counter := 0
	payload.Rows(func(i int, r Row) {
		testutils.Assert(t, strings.Join(r, ",") == "1,2,3,4,5,6", fmt.Sprintf("bad row %d: %s", counter, r))

		testutils.Ok(t, pl.Error)
		counter = counter + 1
	})
	testutils.Equals(t, 8, counter)
}
