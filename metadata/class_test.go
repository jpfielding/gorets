package metadata

var example = `
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
