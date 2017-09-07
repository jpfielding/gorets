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

var payloadlist = `<RETS ReplyCode="0" ReplyText="V2.7.0 2315: Success">
	<DELIMITER value = "09"/>
	<RETSPayloadList
		Resource="RESOURCE"
		Class="CLASS_1"
		Version="Version"
		Date="Tue, 3 Sep 2013 00:00:00 GMT">
		<COLUMNS>	A	B	C	D	E	F	</COLUMNS>
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
	</RETSPayloadList>
</RETS>
`

func TestParseGetPayloadList(t *testing.T) {
	body := ioutil.NopCloser(strings.NewReader(payloadlist))

	pl, err := NewPayloadList(body)
	if err != nil {
		t.Error("error parsing body: " + err.Error())
	}
	assert.Equal(t, StatusOK, pl.Response.Code)
	assert.Equal(t, "V2.7.0 2315: Success", pl.Response.Text)

	var payload []CompactData
	err = pl.ForEach(func(cd CompactData, err error) error {
		payload = append(payload, cd)
		return err
	})
	assert.Equal(t, "CLASS_1", payload[0].Attr["Class"])
	assert.Equal(t, 1, len(payload[0].Entries()))
	assert.Equal(t, "CLASS_2", payload[1].Attr["Class"])
	assert.Equal(t, 2, len(payload[1].Entries()))
}
