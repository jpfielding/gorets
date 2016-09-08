package metadata

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/jpfielding/gotest/testutils"
)

func TestReadClass(t *testing.T) {
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
        <ColumnGroupVersion>01.72.11581</TableVersion>
        <ColumnGroupDate>2016-02-09T06:02:17</TableDate>
        <ColumnGroupSetVersion>01.72.11581</TableVersion>
        <ColumnGroupSetDate>2016-02-09T06:02:17</TableDate>
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

	body := ioutil.NopCloser(bytes.NewReader([]byte(raw)))
	defer body.Close()

	extractor := &Extractor{Body: body}
	rets, err := extractor.Open()

	testutils.Ok(t, err)
	testutils.Equals(t, "Operation Successful", rets.ReplyText)
	testutils.Equals(t, 0, rets.ReplyCode)

	mclass := &MClass{}
	err = extractor.Next("METADATA-CLASS", mclass)
	testutils.Ok(t, err)
	testutils.Equals(t, "Property", string(mclass.Resource))
	testutils.Equals(t, "01.72.11588", string(mclass.Version))
	testutils.Equals(t, "2016-06-01T16:05:01", string(mclass.Date))
	testutils.Equals(t, 2, len(mclass.Class))

	comm := mclass.Class[0]
	testutils.Equals(t, "COMM", string(comm.ClassName))
	testutils.Equals(t, "CommercialSale", string(comm.StandardName))
	testutils.Equals(t, "Commercial", string(comm.VisibleName))
	testutils.Equals(t, "Contains data for Commercial searches.", string(comm.Description))
	testutils.Equals(t, "TableVersion", string(comm.TableVersion))
	testutils.Equals(t, "TableDate", string(comm.TableDate))

	lotl := mclass.Class[1]
	testutils.Equals(t, "LOTL", string(lotl.ClassName))
	testutils.Equals(t, "Land", string(lotl.StandardName))
}
