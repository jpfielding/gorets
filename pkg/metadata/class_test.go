package metadata

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClass(t *testing.T) {
	var raw = `<?xml version="1.0" encoding="utf-8"?>
    <RETS ReplyCode="0" ReplyText="Operation Successful">
    <METADATA>
    <METADATA-CLASS Version="01.72.11588" Date="2016-06-01T16:05:01" Resource="Property">
      <Class>
        <ClassName>COMM</ClassName>
        <StandardName>CommercialSale</StandardName>
        <VisibleName>Commercial</VisibleName>
        <Description>Contains data for Commercial searches.</Description>
        <TableVersion>01.72.11581</TableVersion>
        <TableDate>2016-02-09T06:02:17</TableDate>
        <UpdateVersion>01.72.10221</UpdateVersion>
        <UpdateDate/>
        <ClassTimeStamp>LastModifiedDateTime</ClassTimeStamp>
        <DeletedFlagField/>
        <DeletedFlagValue/>
        <HasKeyIndex>0</HasKeyIndex>
        <ColumnGroupVersion>01.72.11581</ColumnGroupVersion>
        <ColumnGroupDate>2016-02-09T06:02:17</ColumnGroupDate>
        <ColumnGroupSetVersion>01.72.11581</ColumnGroupSetVersion>
        <ColumnGroupSetDate>2016-02-09T06:02:17</ColumnGroupSetDate>
        <METADATA-TABLE Version="01.72.11581" Date="2016-02-09T06:02:17" System="ANNA" Resource="Property" Class="COMM">
        </METADATA-TABLE>
      </Class>
      <Class>
          <ClassName>LOTL</ClassName>
          <StandardName>Land</StandardName>
      </Class>
      </METADATA-CLASS>
    </METADATA>
    </RETS>`

	body := ioutil.NopCloser(strings.NewReader(raw))
	defer body.Close()

	extractor := &Extractor{Body: body}
	response, err := extractor.Open()

	assert.Nil(t, err)
	assert.Equal(t, "Operation Successful", response.ReplyText)

	mclass := &MClass{}
	err = extractor.DecodeNext("METADATA-CLASS", mclass)
	assert.Nil(t, err)
	assert.Equal(t, "Property", string(mclass.Resource))
	assert.Equal(t, "01.72.11588", string(mclass.Version))
	assert.Equal(t, "2016-06-01T16:05:01", string(mclass.Date))
	assert.Equal(t, 2, len(mclass.Class))

	comm := mclass.Class[0]
	assert.Equal(t, "COMM", string(comm.ClassName))
	assert.Equal(t, "CommercialSale", string(comm.StandardName))
	assert.Equal(t, "Commercial", string(comm.VisibleName))
	assert.Equal(t, "Contains data for Commercial searches.", string(comm.Description))
	assert.Equal(t, "01.72.11581", string(comm.TableVersion))
	assert.Equal(t, "2016-02-09T06:02:17", string(comm.TableDate))
	// TODO fill in the rest when time permits

	lotl := mclass.Class[1]
	assert.Equal(t, "LOTL", string(lotl.ClassName))
	assert.Equal(t, "Land", string(lotl.StandardName))
}
