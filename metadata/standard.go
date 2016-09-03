package metadata

import "encoding/xml"

var example = `
<RETS ReplyCode="0" ReplyText="Operation Successful">
<METADATA>
	<METADATA-SYSTEM Date="2016-08-15T14:47:13Z" Version="1.00.04717">
		<System/>
	</METADATA-SYSTEM>
	<METADATA-RESOURCE Date="2016-08-15T14:47:13Z" Version="1.00.03519">
		<Resource>
		<METADATA-CLASS Resource="Agent" Date="2016-07-11T18:22:38Z" Version="1.00.00045">
			<Class>
			<METADATA-TABLE Class="Agent" Resource="Agent" Date="2016-07-11T18:22:38Z" Version="1.00.00040">
				<Field/>
			</METADATA-TABLE>
			</Class>
		</METADATA-CLASS>
		<METADATA-OBJECT Resource="Agent" Date="2014-06-04T15:55:22Z" Version="1.00.00001">
			<Object/>
		</METADATA-OBJECT>
		<METADATA-LOOKUP Resource="Agent" Date="2016-08-08T16:00:15Z" Version="1.00.00110">
			<LookupType>
			<METADATA-LOOKUP_TYPE Lookup="AgentStatus" Resource="Agent" Date="2014-02-24T15:55:08Z" Version="1.00.00004">
				<Lookup/>
			</METADATA-LOOKUP_TYPE>
			</LookupType>
		</METADATA-LOOKUP>
		</Resource>
	</METADATA-RESOURCE>
</METADATA>
</RETS>`

// StandardXML ...
type StandardXML struct {
	XMLName   xml.Name    `xml:"RETS"`
	ReplyCode int         `xml:"ReplyCode"`
	ReplyText string      `xml:"ReplyText"`
	Metadata  XMLMetadata `xml:"METADATA"`
}

// XMLMetadata ...
type XMLMetadata struct {
	System MSystem `xml:"METADATA-SYSTEM"`
}
